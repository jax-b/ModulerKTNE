#include <Arduino.h>
#include <GxEPD2_BW.h>
#include <Fonts/FreeMonoBold12pt7b.h>
#include <Adafruit_I2CDevice.h>
#include <Adafruit_NeoPixel.h>

#define AddressInPin A0
#define SuccessLEDPin 8
#define FailureLEDPin 9
#define MinMaxStable 5
#define BUTTON_PIN 7
#define NEOPIXEL_PIN 10

bool buttonStates = 0;
bool buttonStatesFlicker = 0;
unsigned long lastDebounceTime = 0;
#define DEBOUNCE_DELAY 50

Adafruit_NeoPixel *pixels;

#define EPD_CS      6
#define EPD_DC      5
#define EPD_RESET   1 // can set to -1 and share with microcontroller Reset!
#define EPD_BUSY    A1 // can set to -1 to not use a pin (will wait a fixed delay)
GxEPD2_BW<GxEPD2_154_D67, 32> display(GxEPD2_154_D67(EPD_CS, EPD_DC, EPD_RESET, EPD_BUSY));
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
    if (currentReading - MinMaxStable <= AnalogReading && currentReading + MinMaxStable >= AnalogReading)
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

  Serial.println(analogRead(AddressInPin));

  pinMode(SuccessLEDPin, OUTPUT);
  pinMode(FailureLEDPin, OUTPUT);
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
  display.init();
  display.setTextColor(GxEPD_BLACK);
  display.setFont(&FreeMonoBold12pt7b);
  int16_t  x1, y1;
  uint16_t w, h;
  String stringtoprint = "Detonate";
  display.getTextBounds(stringtoprint, 0, 0, &x1, &y1, &w, &h);
  uint16_t x = ((display.width() - w) / 2) - x1;
  uint16_t y = ((display.height() - h) / 2) - y1;
  display.firstPage();
  do
  {
    display.fillScreen(GxEPD_WHITE);
    // comment out next line to have no or minimal Adafruit_GFX code
    display.setCursor(x, y);
    display.print(stringtoprint);
    
  }
  while (display.nextPage());
  
  pixels = new Adafruit_NeoPixel(4, NEOPIXEL_PIN, NEO_GRB + NEO_KHZ800);
  pixels->begin();
  for (int i = 0; i < 4; i++)
  {
    pixels->setPixelColor(i, 0xFF0000);
  }
  pixels->show();
  delay(500);
  for (int i = 0; i < 4; i++)
  {
    pixels->setPixelColor(i, 0x00FF00);
  }
  pixels->show();
  delay(500);
  for (int i = 0; i < 4; i++)
  {
    pixels->setPixelColor(i, 0x0000ff);
  }
  pixels->show();
  delay(500);
  for (int i = 0; i < 4; i++)
  {
    pixels->setPixelColor(i, 0xFFFF00);
  }
  pixels->show();
  delay(500);
  for (int i = 0; i < 4; i++)
  {
    pixels->setPixelColor(i, 0xFF00FF);
  }
  pixels->show();
  delay(500);
}

void loop()
{
  debounce();
  if (buttonStates)
  {
    for (int i = 0; i < 4; i++)
    {
      pixels->setPixelColor(i, 0);
    }
    pixels->show();
  }
  else
  {
    for (int i = 0; i < 4; i++)
    {
      pixels->setPixelColor(i, 0x00FF00);
    }
    pixels->show();
  }
}