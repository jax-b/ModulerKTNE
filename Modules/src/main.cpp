#include <Arduino.h>
#include <Wire.h>

// ********************************
// Include the type of module that we use
// ********************************
// #include "basemodule\basemodule.h"
// BaseModule mod = BaseModule();
#include "OTS_Button\OTS_Button.h"
OTS_Button mod = OTS_Button();

// *************** External Communications Controller/Player *****************
/// Pins
#define AddressInPin A0
#define SuccessLEDPin 8
#define FailureLEDPin 9
#define S2MInteruptPin 4
#define MinMaxStable 5

#define DEBUG_MODE true

/// Timekeeping variables for external pins for their reset
unsigned long S2MInteruptCallTime = 0;
unsigned long FailureLEDCallTime = 0;

/// Buffers and data tracking for I2C communication
byte incomeingI2CData[10];
byte outgoingI2CData[10];
uint8_t bytesToSend = 0;
uint8_t bytesReceived = 0;
// **********************************************************************

// In side of the desired module use the following to acctivate the timer
// #define TIMER_ENABLE True

// ************** Default Variable Sizes ******************
#define GAMEPLAYSERIALNUMBERLENGTH 8
#define GAMEPLAYMAXLITINDICATOR 6

// *************** Gameplay Variables *****************
// the specific module  might not use the following variables but they will get set
// by the master controller so we want to keep trac of them for the module specific code
/// each run will have a serial number that is used in that run
char gameplaySerialNumber[GAMEPLAYSERIALNUMBERLENGTH];
/// The device can have a number of lit indicators on the sides but gameplay only cares if they are lit
char gameplayLitIndicators[GAMEPLAYMAXLITINDICATOR][3];
uint8_t gameplayLitIndicatorCount = 0;
/// Most modules will have a seed that it can be set to in order to load a specific module configuration
uint16_t gameplaySeed = NULL;
/// Number of battery's that esists on the entire device
uint8_t gameplayNumBattery = 0;
/// There is only six possible ports in the game and they have a bit possition assigned to them.
/// See reference chart in typed notes for more information about which is which
byte gameplayPorts;
/// Timekeeping and gameplay variables
unsigned long gameplayCountdownTime;
bool gameplayTimerRunning = false;
float gameplayStrikeReductionRate = 0.25;
/// All most modules will have a state that they can be set to
/// if a module is failed it will be set to a negative number
/// it will decrement by 1 each time it is failed
int8_t gameplayModuleSolved = 0;
// ****************************************************

// I2C from AIN Address Table
uint8_t convertToAddress(uint16_t addrVIn)
{
    // Top
    if (addrVIn >= 93 && addrVIn < 186)
    {
        return 0x30;
    }
    else if (addrVIn >= 186 && addrVIn < 279)
    {
        return 0x31;
    }
    else if (addrVIn >= 279 && addrVIn < 372)
    {
        return 0x32;
    }
    else if (addrVIn >= 372 && addrVIn < 465)
    {
        return 0x33;
    }
    else if (addrVIn >= 465 && addrVIn < 558)
    {
        return 0x34;
    }
    // Bottom
    else if (addrVIn >= 558 && addrVIn < 651)
    {
        return 0x40;
    }
    else if (addrVIn >= 651 && addrVIn < 744)
    {
        return 0x41;
    }
    else if (addrVIn >= 744 && addrVIn < 837)
    {
        return 0x42;
    }
    else if (addrVIn >= 837 && addrVIn < 930)
    {
        return 0x43;
    }
    else if (addrVIn >= 930 && addrVIn <= 1024)
    {
        return 0x44;
    }
    else
    {
        return 0x00;
    }
}

// Checks to make sure the the voltage coming in is stable with in a certain range
uint16_t getStableVoltage(int pin)
{
    bool VoltageStable = false;
    uint16_t AnalogReading = analogRead(pin);
    while (!VoltageStable)
    {
        // Read the voltage
        uint16_t currentReading = analogRead(pin);
        // If the voltage is within the range of the voltage we are looking for
        if (currentReading - MinMaxStable <= AnalogReading || currentReading + MinMaxStable >= AnalogReading)
        {
            VoltageStable = true;
        }
        else
        {
            AnalogReading = currentReading;
        }
        delay(1);
    }
    return AnalogReading;
}

// This function is called when the module code determines that the module has been solved
void FlagModuleSolved()
{
    // Set the module solved flag to 1
    gameplayModuleSolved = 1;
    // Turn on the success LED
    digitalWrite(SuccessLEDPin, HIGH);
    // Turn off the failure LED
    S2MInteruptCallTime = millis();
    // Turn off the failure LED
    digitalWrite(S2MInteruptPin, LOW);
    gameplayTimerRunning = false;
}

// This function is called when the module code determines that the module has been failed
void FlagModuleFailed()
{
    // Set the module solved flag to -1
    gameplayModuleSolved -= 1;
    // Record the on time of the failure LED
    FailureLEDCallTime = millis();
    // Turn on the failure LED
    digitalWrite(FailureLEDPin, HIGH);
    // Record the start of the signal to the controller interupt
    S2MInteruptCallTime = millis();
    // Signal the main controller interupt
    digitalWrite(S2MInteruptPin, LOW);
}

#ifdef TIMER_ENABLE
unsigned long timekeeperLastRun = 0;
// Updates to defused time of the module
void decrementCounter()
{
    if (!gameplayTimerRunning || gameplayModuleSolved > 0)
    {
        return;
    }
    unsigned long currentTime = millis();
    // Calculate how long it has been since the last run
    unsigned long timesincelastrun = timekeeperLastRun - currentTime;
    // Subtrack it out
    gameplayCountdownTime -= timesincelastrun;

    // If we have a strike
    if (gameplayModuleSolved <= 0)
    {
        float everyrate = (1 / gameplayStrikeReductionRate) / (gameplayModuleSolved * -1);
        unsigned long textra;
        if (everyrate > 1)
        {
            textra = gameplayCountdownTime % (unsigned long)everyrate;
            if (everyrate < 1)
            {
                textra = currentTime - timekeeperLastRun;
                textra += textra * (1 - everyrate);
            }
            else
            {
                textra = (currentTime - timekeeperLastRun) / everyrate;
            }
            gameplayCountdownTime -= textra;
        }
    }
    if (gameplayCountdownTime = 0 || gameplayModuleSolved >= 18000000)
    {
        gameplayModuleSolved = 0;
        gameplayTimerRunning = false;
        FlagModuleFailed();
    }
    timekeeperLastRun = millis();
#ifdef DEBUG_MODE
    Serial.print("Time: ");
    Serial.println(gameplayCountdownTime);
#endif
}
#endif

// Sends out our output buffer
void requestEvent()
{
    if (bytesToSend > 0)
    {
        for (size_t i = 0; i < bytesToSend; i++)
        {
            Wire.write(outgoingI2CData, bytesToSend);
        }
        bytesToSend = 0;
    }
}

// Copy the incoming data into our input buffer
void receiveEvent(int numBytes)
{
#ifdef DEBUG_MODE
    Serial.print("Received: ");
    Serial.println(numBytes);
#endif
    bytesReceived = numBytes;
    for (int i = 0; i < numBytes; i++)
    {
        if (i > 10)
        {
            Wire.read();
        }
        else
        {
#ifdef DEBUG_MODE
            uint8_t wirein = Wire.read();
            Serial.print(wirein, HEX);
            Serial.print(" ");
            incomeingI2CData[i] = wirein;
#else
            incomeingI2CData[i] = Wire.read();
#endif
        }
    }
#ifdef DEBUG_MODE
    Serial.println();
#endif
}

// Processe a command from the controller and if necessary copy data into the output buffer
void I2CCommandProcessor()
{
#ifdef DEBUG_MODE
    Serial.print("I2C cmdgroup:");
    Serial.print(incomeingI2CData[0] >> 4, HEX);
    Serial.print(", cmd:");
    Serial.println(incomeingI2CData[0] & 0xff, HEX);
#endif
    switch (incomeingI2CData[0] >> 4)
    {
    case 0x4:
        switch (incomeingI2CData[0] & 0xF)
        {
        case 0x0: // Stop
            gameplayTimerRunning = false;
#ifdef DEBUG_MODE
            Serial.println("Game Stoped");
#endif
            break;

        default:
            break;
        }
        break;
    case 0x3:
        switch (incomeingI2CData[0] & 0xF)
        {
        case 0x0: // Start
            gameplayTimerRunning = true;
            if (gameplaySeed == NULL)
            {
                gameplaySeed = random(1, 65535);
            }
#ifdef DEBUG_MODE
            Serial.println("Game Started");
#endif
            break;

        default:
            break;
        }
        break;
    // this is a clear command
    case 0x2:
        switch (incomeingI2CData[0] & 0xF)
        {
        // Clear SerialNumber
        case 0x4:
            for (int i = 0; i < GAMEPLAYSERIALNUMBERLENGTH; i++)
            {
                gameplaySerialNumber[i] = '\0';
            }
            mod.setSerialNumber(gameplaySerialNumber);
#ifdef DEBUG_MODE
            Serial.println("Serial Number Cleared");
#endif
            break;
        // Clear LitIndicators
        case 0x5:
            gameplayLitIndicatorCount = 0;
            for (int i = 0; i < GAMEPLAYMAXLITINDICATOR; i++)
            {
                for (int j = 0; j < 3; j++)
                {
                    gameplayLitIndicators[i][j] = '\0';
                }
            }
            mod.setIndicators(gameplayLitIndicators);
#ifdef DEBUG_MODE
            Serial.println("Lit Indicators Cleared");
#endif
            break;
        // Clear Number Batteries
        case 0x6:
            gameplayNumBattery = 0;
            mod.setBatteries(gameplayNumBattery);
#ifdef DEBUG_MODE
            Serial.println("Number Batteries Cleared");
#endif
            break;
        // Clear Port Identities
        case 0x7:
            gameplayPorts = 0x0;
            mod.setPorts(gameplayPorts);
#ifdef DEBUG_MODE
            Serial.println("Port Identities Cleared");
#endif
            break;
        // Clear Seed
        case 0x8:
            gameplaySeed = NULL;
            mod.setSeed(gameplaySeed);
#ifdef DEBUG_MODE
            Serial.println("Seed Cleared");
#endif
            break;
        default:
            break;
        }
        break;
    // This is a command to the module for configuration
    case 0x1:
        switch (incomeingI2CData[0] & 0xF)
        {
        // Set solved status
        case 0x1:
            // greater than 1 bytes because bytes received includes the command byte
            if (bytesReceived > 1)
            {
                gameplayModuleSolved = incomeingI2CData[1];
#ifdef DEBUG_MODE
                Serial.print("SolvedStat: ");
                Serial.println(gameplayModuleSolved);
#endif
                if (gameplayModuleSolved > 0)
                {
                    FlagModuleSolved();
                }
                else
                {
                    digitalWrite(SuccessLEDPin, LOW);
                    mod.setNumStrike(gameplayModuleSolved * -1);
                }
            }
            break;
        // Sync Time between the module and the device
        case 0x2:
            // Time should be a unsigned long so 4 bytes
            // greater than 4 bytes because bytes received includes the command byte
            if (bytesReceived > 4)
            {
                gameplayCountdownTime = incomeingI2CData[1] << 24 | incomeingI2CData[2] << 16 | incomeingI2CData[3] << 8 | incomeingI2CData[4];
#ifdef DEBUG_MODE
                Serial.print("SyncTime: ");
                Serial.println(gameplayCountdownTime);
#endif
            }
            break;
        // Set Strike Rate
        case 0x3:
            // Reduction Rate should be a float so 4 bytes
            // greater than 4 bytes because bytes received includes the command byte
            if (bytesReceived > 4)
            {
                gameplayStrikeReductionRate = incomeingI2CData[1] << 24 | incomeingI2CData[2] << 16 | incomeingI2CData[3] << 8 | incomeingI2CData[4];
#ifdef DEBUG_MODE
                Serial.print("StrikeRate: ");
                Serial.println(gameplayStrikeReductionRate);
#endif
            }
            break;
        // Set Serial Number
        case 0x4:
            // Serial Number should be a String max 8 char;
            // greater than 8 bytes because bytes received includes the command byte
            if (bytesReceived > GAMEPLAYSERIALNUMBERLENGTH)
            {
                for (int i = 0; i < bytesReceived && i < GAMEPLAYSERIALNUMBERLENGTH; i++)
                {
                    gameplaySerialNumber[i] = (char)incomeingI2CData[i + 1];
                }
                mod.setSerialNumber(gameplaySerialNumber);
#ifdef DEBUG_MODE
                Serial.print("SerialNumber: ");
                Serial.println(gameplaySerialNumber);
#endif
            }
            break;
        // Set LitIndicator
        case 0x5:
            // Lit Indicator should be a string of 3 char
            // greater than 3 bytes because bytes received includes the command byte
            if (bytesReceived > 3)
            {
                for (int i = 0; i < bytesReceived && i < 3; i++)
                {
                    gameplayLitIndicators[gameplayLitIndicatorCount][i] = (char)incomeingI2CData[i + 1];
                }
                if (gameplayLitIndicatorCount < GAMEPLAYMAXLITINDICATOR - 1)
                {
                    gameplayLitIndicatorCount++;
                }
                mod.setIndicators(gameplayLitIndicators);
#ifdef DEBUG_MODE
                Serial.print("LitIndicators: {");
                for (uint8_t i = 0; i < gameplayLitIndicatorCount; i++)
                {
                    Serial.print(gameplayLitIndicators[i]);
                    Serial.print(", ");
                }
                Serial.println();
#endif
            }
            break;
        // Set Number of Batteries
        case 0x6:
            // Number of Batteries should be a byte
            // greater than 1 bytes because bytes received includes the command byte
            if (bytesReceived > 1)
            {
                gameplayNumBattery = incomeingI2CData[1];
                mod.setBatteries(gameplayNumBattery);
#ifdef DEBUG_MODE
                Serial.print("NumBattery: ");
                Serial.println(gameplayNumBattery);
#endif
            }
            break;
        // Set Active Ports
        case 0x7:
            // Port Identities should be a byte
            // greater than 1 bytes because bytes received includes the command byte
            if (bytesReceived > 1)
            {
                gameplayPorts = incomeingI2CData[1];
                mod.setPorts(gameplayPorts);
#ifdef DEBUG_MODE
                Serial.print("Ports: ");
                Serial.println(gameplayPorts);
#endif
            }
            break;
        // Set Seed
        case 0x8:
            // Seed should be 2 bytes
            // greater than 2 bytes because bytes received includes the command byte
            if (bytesReceived > 2)
            {
                gameplaySeed = incomeingI2CData[1] << 8 | incomeingI2CData[2];
                mod.setSeed(gameplaySeed);
#ifdef DEBUG_MODE
                Serial.print("Seed: ");
                Serial.println(gameplaySeed);
#endif
            }
            break;
        default:
            break;
        }
        break;
    // This is a command to the device for data
    case 0x0:
        switch (incomeingI2CData[0] & 0xF)
        {
        // Get the modules ID
        case 0x0:
            bytesToSend = 4;
            for (int i = 0; i < bytesToSend; i++)
            {
                outgoingI2CData[i] = mod.getModuleName()[i];
            }
            break;
        // Get the Solved Status
        case 0x1:
            bytesToSend = 1;
            outgoingI2CData[0] = gameplayModuleSolved;
            break;
        }
    }
    bytesReceived = 0;
}

void setup()
{
#ifdef DEBUG_MODE
    Serial.begin(9600);
    while (!Serial) // wait for serial port to connect.
#endif
        // Setup LED's
        pinMode(SuccessLEDPin, OUTPUT);
    pinMode(FailureLEDPin, OUTPUT);
    pinMode(S2MInteruptPin, OUTPUT);
    // Set output pins to inital state
    digitalWrite(S2MInteruptPin, HIGH);
    digitalWrite(SuccessLEDPin, HIGH);
    digitalWrite(FailureLEDPin, HIGH);

    randomSeed(analogRead(A5));

    // Setup I2C
    // Start Listening for address
    Wire.begin(convertToAddress(getStableVoltage(AddressInPin)));
#ifdef DEBUG_MODE
    Serial.print("I2C Listening on ADDR: ");
    Serial.println(convertToAddress(getStableVoltage(AddressInPin)), HEX);
#endif
    Wire.onReceive(receiveEvent);
    Wire.onRequest(requestEvent);

    mod.setupModule();

    // Turn off LEDs to signal initialization is complete
    digitalWrite(FailureLEDPin, LOW);
    digitalWrite(SuccessLEDPin, LOW);
}

void loop()
{
    // Pull the S2M Interupt Pin back UP after a short delay
    if (millis() - S2MInteruptCallTime > 10 && S2MInteruptCallTime != 0)
    {
        digitalWrite(S2MInteruptPin, HIGH);
        S2MInteruptCallTime = 0;
    }
    // Turn off the failure LED after a short delay
    if (millis() - FailureLEDCallTime > 500 && FailureLEDCallTime != 0)
    {
        digitalWrite(FailureLEDPin, LOW);
        FailureLEDCallTime = 0;
    }

    // Module Specific Code
    if (gameplayModuleSolved != 1 && gameplayTimerRunning == true)
    {
        // Success LED should be on for only module success
        digitalWrite(SuccessLEDPin, LOW);
        if (mod.checkSuccess())
        {
            FlagModuleSolved();
        }
        if (mod.checkFailure())
        {
            FlagModuleFailed();
        }
        mod.tickModule(gameplayCountdownTime);
    }

    if (bytesReceived != 0)
    {
        I2CCommandProcessor();
    }

#ifdef TIMER_ENABLE
    // Timekeeping
    decrementCounter();
#endif
}