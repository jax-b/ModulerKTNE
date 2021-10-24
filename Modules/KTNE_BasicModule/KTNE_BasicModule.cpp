#include <Arduino.h>
#include <Wire.h>

#define AddressInPin A0
#define SuccessLEDPin 2
#define FailureLEDPin 3
#define S2MInteruptPin  4
#define MinMaxStable 5

const String ModuleID = "Blnk";

// Most modules will have a scenario that it can be set to
uint8_t ScenarioID = 0;

// All most modules will have a state that they can be set to
// if a module is failed it will be set to a negative number
// it will decrement by 1 each time it is failed
int8_t moduleSolved = 0;

unsigned long S2MInteruptCallTime = 0;
unsigned long FailureLEDCallTime = 0;
byte I2C_command = 0;

unsigned long deviceCountdownTime;
float StrikeReductionRate = 0.25;

byte incomeingI2CData[10];
byte outgoingI2CData[10];
uint8_t bytesToSend = 0;
uint8_t bytesReceived = 0;

// Address Table
uint8_t convertToAddress(uint16_t addrVIn){
    // Top
    if (addrVIn >= 93 && addrVIn < 186){
        return 0x30;
    }
    else if (addrVIn >= 186 && addrVIn < 279){
        return 0x31;
    }
    else if (addrVIn >= 279 && addrVIn < 372){
        return 0x32;
    }
    else if (addrVIn >= 372 && addrVIn < 465){
        return 0x33;
    }
    else if (addrVIn >= 465 && addrVIn < 558){
        return 0x34;
    }
    // Bottom
    else if (addrVIn >= 558 && addrVIn < 651){
        return 0x40;
    }
    else if (addrVIn >= 651 && addrVIn < 744){
        return 0x41;
    }
    else if (addrVIn >= 744 && addrVIn < 837){
        return 0x42;
    }
    else if (addrVIn >= 837 && addrVIn < 930){
        return 0x43;
    }
    else if (addrVIn >= 930 && addrVIn <= 1024){
        return 0x44;
    }
    else {
        return 0x00;
    }
}

// Checks to make sure the the voltage coming in is stable with in a certain range
uint16_t getStableVoltage(int pin) {
    bool VoltageStable = false;
    uint_16t AnalogReading = 0;
    while (!VoltageStable){
        // Read the voltage
        uint16_t currentReading = analogRead(pin);
        // If the voltage is within the range of the voltage we are looking for
        if (currentReading - MinMaxStable <= AnalogReading || currentReading + MinMaxStable >= AnalogReading) {
            VoltageStable = true;
        }
        else{
            AnalogReading = currentReading;
        }
        delay(1);
    }
    return analogReading;
}

// This function is called when the module code determines that the module has been solved
void FlagModuleSolved() {
    // Set the module solved flag to 1
    moduleSolved = 1;
    // Turn on the success LED
    digitalWrite(SuccessLEDPin, HIGH);
    // Turn off the failure LED
    S2MInteruptCallTime = millis();
    // Turn off the failure LED
    digitalWrite(S2MInteruptPin, LOW);
}

// This function is called when the module code determines that the module has been failed
void FlagModuleFailed() {
    // Set the module solved flag to -1
    moduleSolved -= 1;
    // Record the on time of the failure LED
    FailureLEDCallTime = millis();
    // Turn on the failure LED
    digitalWrite(FailureLEDPin, HIGH);
    // Record the start of the signal to the controller interupt
    S2MInteruptCallTime = millis();
    // Signal the main controller interupt
    digitalWrite(S2MInteruptPin, LOW);
}

float timekeeperLastRun = 0;
// Updates to defused time of the module
void decrementCounter() {
    // Calculate our reduction rate based on current number of strikes
    int BaseReductionRate = 10;
    int reductionRate;
    if (moduleSolved == 0) {
        reductionRate = BaseReductionRate;
    } else if (moduleSolved < 0) {
        reductionRate = (moduleSolved * -1 * StrikeReductionRate + 1) * BaseReductionRate;
    } else {
        return;
    }
    
    // Calculate the current time minus how fast the clock is running
    if (millis() - reductionRate > timekeeperLastRun) {
        timekeeperLastRun = millis();
        // Reduce the clock
        deviceCountdownTime -= reductionRate;
    }
}

// Sends out our output buffer
void requestEvent() {
    if ( bytesToSend > 0 ) {
        Wire.write(outgoingI2CData, bytesToSend);
    }
    bytesToSend = 0;
}

// Copy the incoming data into our input buffer
void receiveEvent(int numBytes) {
    bytesReceived = numBytes;
    for (int i = 0; i < numBytes; i++) {
        if (i > 10) {
            Wire.read();
        } else {
            incomeingI2CData[i] = Wire.read();
        }
    }
    processCommands();
}

// Processe a command from the controller and if necessary copy data into the output buffer
void processCommands(){
    switch (incomeingI2CData[0] >> 4) {
    // This is a command to the module for configuration
    case 0x1:
        switch (incomeingI2CData[0] & 0xF) {
            // Set solved status
            case 0x1:
                if (bytesReceived > 2) {
                    moduleSolved = incomeingI2CData[1];
                }
                break;
            // Sync Time between the module and the device
            case 0x2:
                // Time should be a unsigned long so 4 bytes
                if (bytesReceived > 5){
                    deviceCountdownTime = incomeingI2CData[1] << 24 | incomeingI2CData[2] << 16 | incomeingI2CData[3] << 8 | incomeingI2CData[4];
                }
                break;
            // Set Strike Rate
            case 0x3:
                // Reduction Rate should be a float so 4 bytes
                if (bytesReceived > 5){
                    StrikeReductionRate = incomeingI2CData[1] << 24 | incomeingI2CData[2] << 16 | incomeingI2CData[3] << 8 | incomeingI2CData[4];
                }
            default:
                break;
        }
        break;
    // This is a command to the device for data
    case 0x0:
        switch (incomeingI2CData[0] & 0xF) {
            case 0x0:
                bytesToSend = ModuleID.length();
                for (int i = 0; i < bytesToSend; i++) {
                    outgoingI2CData[i] = ModuleID[i];
                }
            // Get the Solved Status
            case 0x1:
                bytesToSend = 1;
                outgoingI2CData[0] = moduleSolved;
                break;
        }
    }
}

void setup() {
    // Setup LED's
    pinMode(SuccessLEDPin, OUTPUT);
    pinMode(FailureLEDPin, OUTPUT);
    pinMode(S2MInteruptPin, OUTPUT);
    digitalWrite(S2MInteruptPin, HIGH);
    // Setup I2C
    // Start Listening for address
    Wire.begin(convertToAddress(getStableVoltage(AddressInPin)));
    Wire.onReceive(receiveEvent);
    Wire.onRequest(requestEvent);
}

void loop(){
    // Pull the S2M Interupt Pin back UP after a short delay
    if (millis() - S2MInteruptCallTime > 10) {
        digitalWrite(S2MInteruptPin, HIGH);
        S2MInteruptCallTime = 0;
    }
    // Turn off the failure LED after a short delay
    if (millis() - FailureLEDCallTime > 500) {
        digitalWrite(FailureLEDPin, LOW);
        FailureLEDCallTime = 0;
    }
    // Module Code
    if (moduleSolved != 1) {
        //run the module code
    }
    decrementCounter();
}