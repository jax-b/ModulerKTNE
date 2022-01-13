#include "..\baseModule\baseModule.h"
#include <Adafruit_NeoPixel.h>

#define GREEN_LED_PIN 10
#define GREEN_BUTTON_PIN 7
#define RED_LED_PIN 16
#define RED_BUTTON_PIN 14
#define BLUE_LED_PIN 15
#define BLUE_BUTTON_PIN A1
#define YELLOW_LED_PIN 5
#define YELLOW_BUTTON_PIN 6
#define TONE_LED_PIN 3
class OTS_Simon : public baseModule
{
public:
    OTS_Simon();
    // inherated from baseModule
    virtual void setupModule();
    virtual void setSeed(uint16_t);
    virtual void tickModule(uint16_t);

protected:
    bool buttonStates[4] = {0, 0, 0, 0};
    bool buttonStatesFlicker[4] = {0, 0, 0, 0};
    unsigned long lastDebounceTime[4] = {0, 0, 0, 0};
    unsigned long lastLEDOnTime[4] = {0, 0, 0, 0};
    bool playerTurn = false;
    uint8_t sequencePosition = 0;
    uint8_t sequenceLength = 1;
    uint8_t sequence[4] = {0, 0, 0, 0};
    unsigned long lastFullSequenceTime = 0;
    unsigned long lastSequenceTime = 0;
    uint8_t PlayerEntrys[4] = {0, 0, 0, 0};
    uint8_t PlayerEntrysPosition = 0;

private:
    void processButtons();
    bool checkSerialNumVowel();
    void resetLEDs();
    void playSequance();
};