#include "OTS_Simon.h"

OTS_Simon::OTS_Simon()
{
    OTS_Simon::modID[0] = "s";
    OTS_Simon::modID[1] = "m";
    OTS_Simon::modID[2] = "o";
    OTS_Simon::modID[3] = "n";
}

const uint8_t PROGMEM ButtonOrder[4] = {GREEN_BUTTON_PIN, RED_BUTTON_PIN, BLUE_BUTTON_PIN, YELLOW_BUTTON_PIN};
const uint8_t PROGMEM LEDOrder[4] = {GREEN_LED_PIN, RED_LED_PIN, BLUE_LED_PIN, YELLOW_LED_PIN};

void OTS_Simon::setupModule()
{
    // Setup Button Pin States
    for (uint8_t i = 0; i < 4; i++)
    {
        pinMode(ButtonOrder[i], INPUT_PULLUP);
        pinMode(LEDOrder[i], OUTPUT);
    }
}

void OTS_Simon::processButtons()
{
    for (uint8_t ButtonNumber = 0; ButtonNumber < 4; ButtonNumber++)
    {
        // read the state of the switch/button:
        bool currentState = !digitalRead(ButtonOrder[ButtonNumber]);

        if (currentState != OTS_Simon::buttonStatesFlicker[ButtonNumber])
        {
            OTS_Simon::lastDebounceTime[ButtonNumber] = millis();
            OTS_Simon::buttonStatesFlicker[ButtonNumber] = currentState;
        }

        if ((millis() - OTS_Simon::lastDebounceTime[ButtonNumber]) > 5)
        {
            // save the the last state
            OTS_Simon::buttonStates[ButtonNumber] = currentState;
        }
    }
}

const char PROGMEM vowels[5] = {'a', 'e', 'i', 'o', 'u'};
bool OTS_Simon::checkSerialNumVowel()
{
    for (uint8_t i = 0; i < 5; i++)
    {
        for (uint8_t j = 0; j < 8; j++)
        {
            if (vowels[i] == OTS_Simon::serialNumber[j])
            {
                return true;
            }
        }
    }
    return false;
}

void OTS_Simon::setSeed(uint16_t inSeed)
{
    OTS_Simon::seed = inSeed;
    // Seed to set the sequence
    // sequence = ;
}

void OTS_Simon::resetLEDs()
{
    for (uint8_t i = 0; i < 4; i++)
    {
        if (millis() - OTS_Simon::lastLEDOnTime[i] > 100 || !OTS_Simon::playerTurn)
        {
            digitalWrite(LEDOrder[i], LOW);
        }
    }
}

void OTS_Simon::playSequance()
{
    if (OTS_Simon::sequencePosition >= OTS_Simon::sequenceLength)
    {
        OTS_Simon::sequencePosition = 0;
        OTS_Simon::lastFullSequenceTime = millis();
    }
    if (millis() - OTS_Simon::lastFullSequenceTime > 1000 && millis() - OTS_Simon::lastSequenceTime > 150)
    {
        digitalWrite(LEDOrder[OTS_Simon::sequencePosition], HIGH);
        OTS_Simon::lastLEDOnTime[OTS_Simon::sequencePosition] = millis();
        OTS_Simon::lastSequenceTime = millis();
        OTS_Simon::sequencePosition++;
    }
}

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
void OTS_Simon::tickModule(uint16_t currentGameTime)
{
    this->processButtons();
    this->resetLEDs();
    // We need to wait if the master controller has not cleared out the failure flag
    if (OTS_Simon::failureTriggered)
    {
        return;
    }

    for (uint8_t i = 0; i < 4 && !OTS_Simon::playerTurn; i++)
    {
        if (OTS_Simon::buttonStates[i] == true)
        {
            OTS_Simon::playerTurn = true;
        }
    }

    if (!playerTurn)
    { // Play the sequence to the player
        this->playSequance();
    }
    else
    {
        uint8_t btnpressed;
        for (uint8_t i = 0; i < 4; i++)
        {
            if (OTS_Simon::buttonStates[i] == true)
            {
                btnpressed = i;
            }
        }
        uint8_t actualBTNPressed = PlayerSimonMap[OTS_Simon::checkSerialNumVowel()][OTS_Simon::numStrike][btnpressed];
        OTS_Simon::PlayerEntrys[OTS_Simon::PlayerEntrysPosition] = actualBTNPressed;

        // Check for an invalid press
        for (uint8_t i = 0; i <= OTS_Simon::sequencePosition; i++)
        {
            if (OTS_Simon::PlayerEntrys[i] != OTS_Simon::sequence[i])
            {
                OTS_Simon::failureTriggered = true;
                return;
            }
        }
        if (OTS_Simon::PlayerEntrysPosition == OTS_Simon::sequencePosition)
        { // Check for a completed round of simon says
            for (uint8_t i = 0; i < 4; i++)
            {
                OTS_Simon::PlayerEntrys[i] = 0;
            }
            OTS_Simon::PlayerEntrysPosition = 0;
            OTS_Simon::sequenceLength++;
        }
    }
}