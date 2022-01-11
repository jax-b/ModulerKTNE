#include "..\baseModule\baseModule.h"
#include <Adafruit_NeoPixel.h>

#define NEOPIXEL_PIN 10
#define WIRE_BUTTON_PIN_1 15
#define WIRE_BUTTON_PIN_2 14
#define WIRE_BUTTON_PIN_3 16
#define WIRE_BUTTON_PIN_4 5
#define WIRE_BUTTON_PIN_5 6
#define WIRE_BUTTON_PIN_6 7

class OTS_Wires : public baseModule
{
protected:
    Adafruit_NeoPixel *_pixels;
    uint8_t wireButtonPins[6] = {
        WIRE_BUTTON_PIN_1,
        WIRE_BUTTON_PIN_2,
        WIRE_BUTTON_PIN_3,
        WIRE_BUTTON_PIN_4,
        WIRE_BUTTON_PIN_5,
        WIRE_BUTTON_PIN_6
    };
    bool buttonStates[] = {0, 0, 0, 0, 0, 0};
    bool buttonStatesFlicker[] = {0, 0, 0, 0, 0, 0};
    unsigned long lastDebounceTime[] = {0, 0, 0, 0, 0, 0};
    
    void processButtons();
}