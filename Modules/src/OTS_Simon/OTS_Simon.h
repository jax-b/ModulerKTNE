#include "../baseModule/baseModule.h"
#include <Adafruit_NeoPixel.h>

#define GREEN_LED_PIN 10
#define GREEN_BUTTON_PIN 7 
#define GREEN_TONE_HZ 765
#define RED_LED_PIN 16 
#define RED_BUTTON_PIN 14
#define RED_TONE_HZ 535
#define BLUE_LED_PIN 15
#define BLUE_BUTTON_PIN A1
#define BLUE_TONE_HZ 645
#define YELLOW_LED_PIN 5
#define YELLOW_BUTTON_PIN 6 
#define YELLOW_TONE_HZ 850
#define TONE_PIN 3

// vowel, strike, position
const uint8_t PROGMEM PlayerSimonMap[2][3][4] = {
    {
        // Zero -- False -- No Vowel
        {
            // Zero -- No Strikes
            3, // Green Flash = Yellow
            2, // Red Flash = Blue
            1, // Blue Flash = Red
            0  // Yellow Flash = Green
        },
        {
            // One -- One Strike
            2, // Green Flash = Blue
            3, // Red Flash = Yellow
            0, // Blue Flash = Green
            1  // Yellow Flash = Red
        },
        {
            // Two -- Two Strikes
            3, // Green Flash = Yellow
            0, // Red Flash = Green
            1, // Blue Flash = Red
            2  // Yellow Flash = Blue
        },
    },
    {
        // One -- True -- Vowel
        {
            // Zero -- No Strikes
            0, // Green Flash = Green
            2, // Red Flash = Blue
            3, // Blue Flash = Yellow
            1, // Yellow Flash = Red
        },
        {
            // One -- One Strike
            3, // Green Flash = Yellow
            1, // Red Flash = Red
            2, // Blue Flash = Blue
            0, // Yellow Flash = Green
        },
        {
            // Two -- Two Strikes
            2, // Green Flash = Blue
            3, // Red Flash = Yellow
            0, // Blue Flash = Green
            1, // Yellow Flash = Red
        },
    }};
const uint8_t PROGMEM ButtonOrder[4] = {GREEN_BUTTON_PIN, RED_BUTTON_PIN, BLUE_BUTTON_PIN, YELLOW_BUTTON_PIN};
const uint8_t PROGMEM LEDOrder[4] = {GREEN_LED_PIN, RED_LED_PIN, BLUE_LED_PIN, YELLOW_LED_PIN};
const uint16_t PROGMEM TONEOrder[4] = {GREEN_TONE_HZ, RED_TONE_HZ, BLUE_TONE_HZ, YELLOW_TONE_HZ};
const char PROGMEM vowels[5] = {'a', 'e', 'i', 'o', 'u'};
class OTS_Simon : public baseModule
{
public:
    OTS_Simon();
    // inherated from baseModule
    virtual void setupModule();
    virtual void setSeed(uint16_t);
    virtual void tickModule(uint16_t);
    virtual void clearModule();

private:
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

protected:
    void processButtons();
    bool checkSerialNumVowel();
    void resetLEDs();
    void playSequence();
};