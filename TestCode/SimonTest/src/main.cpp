#define PCBVERSION_1_01
// #define PCBVERSION_1_1


#include <arduino.h>
#ifdef PCBVERSION_1_01
#define AddressInPin 26
#define SuccessLEDPin 9
#define FailureLEDPin 8
#define BlueLEDPin 10
#define MinMaxStable 5
#define M2CInteruptPin 4
#define AudioOut 7
#define GreenLED 12
#define GreenButton 15
#define RedLED 11
#define RedButton 16
#define BlueLED 14
#define BlueButton 17
#define YellowLED 13
#define YellowButton 18
#endif
#ifdef PCBVERSION_1_1
#define AddressInPin 26
#define SuccessLEDPin 9
#define FailureLEDPin 8
#define BlueLEDPin 10
#define MinMaxStable 5
#define M2CInteruptPin 4
#define AudioOut 7
#define GreenLED 12
#define GreenButton 15
#define RedLED 11
#define RedButton 16
#define BlueLED 14
#define BlueButton 17
#define YellowLED 13
#define YellowButton 18
// Audio Shutdown is a inverted signal active low
#define AudioShutdown 6
#endif


const uint8_t ButtonOrder[] = {GreenButton, RedButton, BlueButton, YellowButton};
const uint8_t LEDOrder[] = {GreenLED, RedLED, BlueLED, YellowLED};
bool buttonStates[] = {0, 0, 0, 0};
bool buttonStatesFlicker[] = {0, 0, 0, 0};
unsigned long lastDebounceTime[] = {0, 0, 0, 0};
#define DEBOUNCE_DELAY 50

// Address Table
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

void debounce(int ButtonNumber)
{
    // read the state of the switch/button:
    bool currentState = digitalRead(ButtonOrder[ButtonNumber]);

    if (currentState != buttonStatesFlicker[ButtonNumber])
    {
        lastDebounceTime[ButtonNumber] = millis();
        buttonStatesFlicker[ButtonNumber] = currentState;
    }

    if ((millis() - lastDebounceTime[ButtonNumber]) > DEBOUNCE_DELAY)
    {
        // save the the last state
        buttonStates[ButtonNumber] = currentState;
    }
}

void setup()
{
    
    pinMode(SuccessLEDPin, OUTPUT);
    pinMode(FailureLEDPin, OUTPUT);
    pinMode(BlueLEDPin, OUTPUT);
    pinMode(AddressInPin, INPUT);
    pinMode(AudioOut, OUTPUT);
    digitalWrite(AudioOut, false);
    #ifdef PCBVERSION_1_1
    pinMode(AudioShutdown, OUTPUT);
    digitalWrite(AudioShutdown, true); 
    #endif
    Serial.begin(9600);
    while (!Serial);
    delay(10);

    

    for (int i = 0; i < 4; i++)
    {
        pinMode(ButtonOrder[i], INPUT_PULLUP);
        pinMode(LEDOrder[i], OUTPUT);
    }


    Serial.println("Hello World!");
    Serial.println("Simon Test");
    Serial.println("Address Read Test");
    uint16_t stblVoltage = getStableVoltage(AddressInPin);
    Serial.print("Stable Voltage: ");
    Serial.print(stblVoltage);
    uint8_t address = convertToAddress(stblVoltage);
    Serial.print(", Address: ");
    Serial.println(address);

    Serial.println("Green LED On");
    digitalWrite(SuccessLEDPin, HIGH);
    delay(500);
    Serial.println("Green LED Off");
    digitalWrite(SuccessLEDPin, LOW);
    delay(500);
    Serial.println("Red LED On");
    digitalWrite(FailureLEDPin, HIGH);
    delay(500);
    Serial.println("Red LED Off");
    digitalWrite(FailureLEDPin, LOW);
    delay(500);
    Serial.println("Blue LED On");
    digitalWrite(BlueLEDPin, HIGH);
    delay(500);
    Serial.println("Blue LED Off");
    digitalWrite(BlueLEDPin, LOW);

    for (int i = 0; i < 4; i++)
    {
        digitalWrite(LEDOrder[i], HIGH);
        Serial.print("LED ");
        Serial.println(i);
        delay(500);
    }
    for (int i = 0; i < 4; i++)
    {
        digitalWrite(LEDOrder[i], LOW);
        delay(500);
    }

    Serial.println("Audio Test");
    // analogWrite(AudioOut, 100);
    
    #ifdef PCBVERSION_1_1
    digitalWrite(AudioShutdown, false); 
    #endif
    float frq = 100;
    for (uint8_t i = 0; i < 8; i++)
    {
        Serial.print(frq);
        Serial.print(" Hz, ");
        tone(AudioOut, frq, 100);
        frq += 50;
        delay(500);
    }
    Serial.println("Audio Off");
    digitalWrite(AudioOut, false);
    #ifdef PCBVERSION_1_1
    digitalWrite(AudioShutdown, true); 
    #endif
}

void loop()
{
    for (int i = 0; i < 4; i++)
    {
        debounce(i);
    }
    for (int i = 0; i < 4; i++)
    {
        digitalWrite(LEDOrder[i], !buttonStates[i]);
    }
}