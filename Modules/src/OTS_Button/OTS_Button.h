#include "..\baseModule\baseModule.h"
#include <Adafruit_NeoPixel.h>
#include <GxEPD2_BW.h>
#include <Fonts/FreeMonoBold12pt7b.h>

// To be changed bases off of selected microcontroller and pcb version
#define NEOPIXEL_PIN 10
#define BUTTON_PIN 7
#define EPD_CS 6
#define EPD_DC 5
#define EPD_BUSY A1

// Enable the timer in the main software loop
#define TIMER_ENABLE True

class OTS_Button : public baseModule
{
public:
    // inherated from baseModule
    virtual void setupModule();
    virtual void setSeed(uint16_t);
    virtual void tickModule(uint16_t);

protected:
    uint8_t buttonColor = 0x0;
    uint8_t stripColor = 0x0;
    uint32_t stripColorHex = 0x0;
    uint8_t chosenWord = 0x0; // index of the chosen word rather than having it be a whole character array. just makes things faster
    // only 4 possible words, placing them into an array then using the seed to choose which word will be displayed
    const char possibleButtonWords[4][8] = {
        {'A', 'b', 'o', 'r', 't', '\0', '\0', '\0'},
        {'D', 'e', 't', 'o', 'n', 'a', 't', 'e'},
        {'H', 'o', 'l', 'd', '\0', '\0', '\0', '\0'},
        {'P', 'r', 'e', 's', 's', '\0', '\0', '\0'}};
    Adafruit_NeoPixel *_pixels;
    GxEPD2_BW<GxEPD2_154_D67, 32> _display(GxEPD2_154_D67(EPD_CS, EPD_DC, -1, EPD_BUSY));

private:
    uint16_t timeLastBtn = 0;
    bool lastBTNState = 0;
    void relHeldButton(uint16_t);
    uint16_t btnDebounce();
};