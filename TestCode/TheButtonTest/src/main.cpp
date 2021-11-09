#include <Arduino.h>
#include <FastLED.h>
#include <SPI.h>
#include <EPD1in54.h>
#include <EPDPaint.h>

#define AddressInPin A0
#define SuccessLEDPin 8
#define FailureLEDPin 9
#define MinMaxStable 5
#define BUTTON_PIN 7
#define ALED_Pin 10

bool buttonStates = 0;
bool buttonStatesFlicker = 0;
unsigned long lastDebounceTime = 0;
#define DEBOUNCE_DELAY 50

CRGB aleds[4];


#define EPD_CS      6
#define EPD_DC      5
#define EPD_RESET   1 // can set to -1 and share with microcontroller Reset!
#define EPD_BUSY    A1 // can set to -1 to not use a pin (will wait a fixed delay)
#define COLORED     0
#define UNCOLORED   1
unsigned char image[1024];
EPDPaint paint(image, 0, 0); 
EPD1in54 epd(EPD_RESET,EPD_DC,EPD_CS,EPD_BUSY);

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
  uint16_t AnalogReading = 0;
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

void debounce()
{
  // read the state of the switch/button:
  bool currentState = digitalRead(BUTTON_PIN);

  if (currentState != buttonStatesFlicker)
  {
    lastDebounceTime = millis();
    buttonStatesFlicker = currentState;
  }

  if ((millis() - lastDebounceTime) > DEBOUNCE_DELAY)
  {
    // save the the last state
    buttonStates = currentState;
  }
}

void setup()
{
  Serial.begin(9600);
  while (!Serial);

  pinMode(SuccessLEDPin, OUTPUT);
  pinMode(FailureLEDPin, OUTPUT);
  pinMode(AddressInPin, INPUT);
  pinMode(BUTTON_PIN, INPUT_PULLUP);

  Serial.println("Hello World!");
  Serial.println("The Button Test");
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

  Serial.println("EPD Test");
  if (epd.init(lutPartialUpdate) != 0) {
    Serial.print("e-Paper init failed");
  } else {
    Serial.println("e-Paper init succeed");
    epd.clearFrameMemory(0xFF);   // bit set = white, bit reset = black
    epd.displayFrame();
    epd.clearFrameMemory(0xFF);   // bit set = white, bit reset = black
    epd.displayFrame();

    paint.setWidth(200);
    paint.setHeight(24);
    paint.setRotate(ROTATE_0);

    paint.clear(UNCOLORED);
    paint.drawStringAt(30, 4, "Hello world!", &Font16, UNCOLORED);
    epd.setFrameMemory(paint.getImage(), 0, 10, paint.getWidth(), paint.getHeight());

    paint.clear(UNCOLORED);
    paint.drawStringAt(30, 4, "e-Paper Demo", &Font16, UNCOLORED);
    epd.setFrameMemory(paint.getImage(), 0, 30, paint.getWidth(), paint.getHeight());

    paint.setWidth(64);
    paint.setHeight(64);

    paint.clear(UNCOLORED);
    paint.drawRectangle(0, 0, 40, 50, UNCOLORED);
    paint.drawLine(0, 0, 40, 50, UNCOLORED);
    paint.drawLine(40, 0, 0, 50, UNCOLORED);
    epd.setFrameMemory(paint.getImage(), 16, 60, paint.getWidth(), paint.getHeight());

    paint.clear(UNCOLORED);
    paint.drawCircle(32, 32, 30, UNCOLORED);
    epd.setFrameMemory(paint.getImage(), 120, 60, paint.getWidth(), paint.getHeight());

    paint.clear(UNCOLORED);
    paint.drawFilledRectangle(0, 0, 40, 50, UNCOLORED);
    epd.setFrameMemory(paint.getImage(), 16, 130, paint.getWidth(), paint.getHeight());

    paint.clear(UNCOLORED);
    paint.drawFilledCircle(32, 32, 30, UNCOLORED);
    epd.setFrameMemory(paint.getImage(), 120, 130, paint.getWidth(), paint.getHeight());
    epd.displayFrame();
  }

  Serial.println("ledtest");
  FastLED.addLeds<WS2812B, ALED_Pin, GRB>(aleds, 4);
  for (int i = 0; i < 4; i++)
  {
    aleds[i] = CRGB::Red;
  }
  FastLED.show();
  delay(500);
  for (int i = 0; i < 4; i++)
  {
    aleds[i] = CRGB::Green;
  }
  FastLED.show();
  delay(500);
  for (int i = 0; i < 4; i++)
  {
    aleds[i] = CRGB::Blue;
  }
  FastLED.show();
  delay(500);
  for (int i = 0; i < 4; i++)
  {
    aleds[i] = CRGB::Yellow;
  }
  FastLED.show();
  delay(500);
  for (int i = 0; i < 4; i++)
  {
    aleds[i] = CRGB::Purple;
  }
  FastLED.show();
  delay(500);
}

void loop()
{
  debounce();
  if (buttonStates)
  {
    for (int i = 0; i < 4; i++)
    {
      aleds[i] = CRGB::Black;
    }
    FastLED.show();
  }
  else
  {
    for (int i = 0; i < 4; i++)
    {
      aleds[i] = CRGB::Green;
    }
    FastLED.show();
  }
}