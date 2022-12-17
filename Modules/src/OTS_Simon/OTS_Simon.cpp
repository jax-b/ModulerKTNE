#include "OTS_Simon.h"

OTS_Simon::OTS_Simon()
{
    const char realModID[] = "smon";
    for (uint8_t i = 0; i < 4; i++)
    {
        modID[i] = realModID[i];
    }
}

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
    OTS_Simon::seed = inSeed; // This does not need to happen as simon says only cares about the vowel and number of total strikes
}

void OTS_Simon::resetLEDs()
{
    if (!OTS_Simon::playerTurn) {
        for (uint8_t i = 0; i < 4; i++) {
            if (millis() - OTS_Simon::lastLEDOnTime[i] > 200) {
                if (digitalRead(LEDOrder[i])) {
                    noTone(TONE_PIN);
                }
                digitalWrite(LEDOrder[i], LOW);
            }
        }
    }
}

void OTS_Simon::clearModule() {
    OTS_Simon::seed = 0;
    for (uint8_t i = 0; i < 4; i++)
    {
        digitalWrite(LEDOrder[i], LOW);
    }
    noTone(TONE_PIN);
}

void OTS_Simon::playSequence()
{
    if (OTS_Simon::sequencePosition >= OTS_Simon::sequenceLength)
    {
        OTS_Simon::sequencePosition = 0;
        OTS_Simon::lastFullSequenceTime = millis();
    }
    if (millis() - OTS_Simon::lastFullSequenceTime > 5000 && millis() - OTS_Simon::lastSequenceTime > 250)
    {
        digitalWrite(LEDOrder[OTS_Simon::sequencePosition], HIGH);
        noTone(TONE_PIN);
        tone(TONE_PIN, uint(TONEOrder[OTS_Simon::sequencePosition]));
        OTS_Simon::lastLEDOnTime[OTS_Simon::sequencePosition] = millis();
        OTS_Simon::lastSequenceTime = millis();
        OTS_Simon::sequencePosition++;
    }
}

void OTS_Simon::tickModule(uint16_t currentGameTime)
{
    this->processButtons();
    this->resetLEDs();
    // We need to wait if the master controller has not cleared out the failure flag
    if (OTS_Simon::failureTriggered)
    {
        return;
    }

    // Check to see if the player has pushed a button, then switch from playback mode to player mode
    for (uint8_t i = 0; i < 4 && !OTS_Simon::playerTurn; i++)
    {
        if (OTS_Simon::buttonStates[i] == true)
        {
            OTS_Simon::playerTurn = true;
            // OTS_Simon::lastFullSequenceTime = millis(); // Reusing this var to tell how long 
        }
    }

    if (!playerTurn)
    { // Play the sequence to the player
        this->playSequence();
    }
    else
    {   
        uint8_t btnPressed;
        bool btnFound = false;
        for (uint8_t i = 0; i < 4; i++)
        {
            if (OTS_Simon::buttonStates[i] == true)
            {
                // Set the found button
                btnPressed = i;
                btnFound = true;
                break;

            } 
        }
        if (btnFound) {
            noTone(TONE_PIN);
            tone(TONE_PIN, uint(TONEOrder[btnPressed]));
            uint8_t actualBTNPressed = PlayerSimonMap[OTS_Simon::checkSerialNumVowel()][OTS_Simon::numStrike][btnPressed];
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
        } else {
            noTone(TONE_PIN);
        }
    }
}