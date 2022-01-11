#include "..\baseModule\baseModule.h"
#include <Adafruit_NeoPixel.h>
#include <GxEPD2_BW.h>
#include <Fonts/FreeMonoBold12pt7b.h>

// To be changed bases off of selected microcontroller and pcb version
#define NEOPIXEL_PIN 10
#define BUTTON_PIN 7
#define EPD_CS     6
#define EPD_DC     5
#define EPD_BUSY   A1

class OTS_Button : public baseModule
{
protected:
    uint8_t buttonColor;
    uint8_t stripColor;
    uint32_t stripColorHex;
    uint8_t chosenWord; // index of the chosen word rather than having it be a whole character array. just makes things faster
    // only 4 possible words, placing them into an array then using the seed to choose which word will be displayed
    char[4][8] const possibleButtonWords = {
        {'A', 'b', 'o', 'r', 't', '\0', '\0', '\0'},
        {'D', 'e', 't', 'o', 'n', 'a', 't', 'e'},
        {'H', 'o', 'l', 'd', '\0', '\0', '\0', '\0'},
        { 'P', 'r', 'e', 's', 's', '\0', '\0', '\0'}};
    Adafruit_NeoPixel *_pixels;
    GxEPD2_BW<GxEPD2_154_D67, 32> _display(GxEPD2_154_D67(EPD_CS, EPD_DC, -1, EPD_BUSY));
private:
    uint16_t timeLastBtn;
    bool lastBTNState;
    void relHeldButton();
}