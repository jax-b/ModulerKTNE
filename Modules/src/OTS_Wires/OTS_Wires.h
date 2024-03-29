#include "../baseModule/baseModule.h"
#include <Adafruit_NeoPixel.h>
#include <Arduino.h>
#include <Wire.h>
#include <SPI.h>

#define NEOPIXEL_PIN 10
#define WIRE_BUTTON_PIN_1 15
#define WIRE_BUTTON_PIN_2 14
#define WIRE_BUTTON_PIN_3 16
#define WIRE_BUTTON_PIN_4 5
#define WIRE_BUTTON_PIN_5 6
#define WIRE_BUTTON_PIN_6 7

const uint32_t PROGMEM WIRE_COLOR[6] = {
    0x00000000, // OFF/No Wire
    0x00FFFFFF, // White
    0x00FF0000, // Red
    0x000000FF, // Blue
    0x00FFFF00, // Yellow
    0x00FF00FF  // Magenta (Black in the manual)
};
const uint8_t PROGMEM wireButtonPins[6] = {
    WIRE_BUTTON_PIN_1,
    WIRE_BUTTON_PIN_2,
    WIRE_BUTTON_PIN_3,
    WIRE_BUTTON_PIN_4,
    WIRE_BUTTON_PIN_5,
    WIRE_BUTTON_PIN_6
};

class OTS_Wires : public baseModule
{
public:
    OTS_Wires();
    // inherated from baseModule
    virtual void setupModule();
    virtual void setSeed(uint16_t);
    virtual void tickModule(uint16_t);
    virtual void clearModule();

private:
    Adafruit_NeoPixel *_pixels;
    bool buttonStates[6] = {0, 0, 0, 0, 0, 0};
    bool buttonStatesFlicker[6] = {0, 0, 0, 0, 0, 0};
    bool wireCuts[6] = {0, 0, 0, 0, 0, 0};
    uint8_t wireColors[6] = {0, 0, 0, 0, 0, 0};
    uint8_t numWires = 0;
    unsigned long lastDebounceTime[6] = {0, 0, 0, 0, 0, 0};
    
protected:
    void processButtons();
    void cutWire(uint8_t buttonNumber);
};