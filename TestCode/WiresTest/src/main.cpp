#include <Arduino.h>
#include <FastLED.h>

#define AddressInPin A0
#define SuccessLEDPin 8
#define FailureLEDPin 9
#define MinMaxStable 5
#define ALED_Pin 10

#define Wire1Button 14
#define Wire2Button 15
#define Wire3Button 16
#define Wire4Button 5
#define Wire5Button 6
#define Wire6Button 7

const uint8_t ButtonOrder[] = {Wire1Button, Wire2Button, Wire3Button, Wire4Button, Wire5Button, Wire6Button};
bool buttonStates[] = {0, 0, 0, 0, 0, 0};
bool buttonStatesFlicker[] = {0, 0, 0, 0, 0, 0};
unsigned long lastDebounceTime[] = {0, 0, 0, 0, 0, 0};
#define DEBOUNCE_DELAY 50

CRGB aleds[12];

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
  Serial.begin(9600);
  while (!Serial)
    ;

  pinMode(SuccessLEDPin, OUTPUT);
  pinMode(FailureLEDPin, OUTPUT);
  pinMode(AddressInPin, INPUT);

  for (int i = 0; i < 6; i++)
  {
    pinMode(ButtonOrder[i], INPUT_PULLUP);
  }

  Serial.println("Hello World!");
  Serial.println("Wires Test");
  Serial.println("Address Read Test");
  uint16_t stblVoltage = getStableVoltage(AddressInPin);
  Serial.print("Stable Voltage: ");
  Serial.print(stblVoltage);
  uint8_t address = convertToAddress(stblVoltage);
  Serial.print(", Address: ");
  Serial.println(address);

  Serial.println("Success Test");
  Serial.println("Green LED On");
  digitalWrite(SuccessLEDPin, HIGH);
  delay(500);
  Serial.println("Red LED On");
  digitalWrite(FailureLEDPin, HIGH);
  delay(500);
  Serial.println("Green LED Off");
  digitalWrite(SuccessLEDPin, LOW);
  delay(500);
  Serial.println("Red LED Off");
  digitalWrite(FailureLEDPin, LOW);
  Serial.println("Red LED Off");

  Serial.println("Aled Test");
  FastLED.addLeds<WS2812B, ALED_Pin, GRB>(aleds, 12);
  for (int i = 0; i < 12; i++)
  {
    aleds[i] = CRGB::Red;
  }
  FastLED.show();
  delay(500);
  for (int i = 0; i < 12; i++)
  {
    aleds[i] = CRGB::Green;
  }
  FastLED.show();
  delay(500);
  for (int i = 0; i < 12; i++)
  {
    aleds[i] = CRGB::Blue;
  }
  FastLED.show();
  delay(500);
  for (int i = 0; i < 12; i++)
  {
    aleds[i] = CRGB::Yellow;
  }
  FastLED.show();
  delay(500);
  for (int i = 0; i < 12; i++)
  {
    aleds[i] = CRGB::Purple;
  }
  FastLED.show();
  delay(500);
}

void loop()
{
  for (int i = 0; i < 6; i++)
  {
    debounce(i);
  }
  for (int i = 0; i < 6; i++)
  {
    if (buttonStates[i] == HIGH)
    {
      aleds[i] = CRGB::Purple;
      aleds[i+6] = CRGB::Orange;
    }
    else
    {
      aleds[i] = CRGB::Black;
      aleds[i+6] = CRGB::Black;
    }
  }
  FastLED.show();
}