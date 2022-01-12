#include "OTS_Simon.h"

OTS_Simon::OTS_Simon()
{
    OTS_Simon::modID = "smon"
}

void OTS_Simon::setupModule() 
{
    // Setup Button Pin States
    for (uint8_t i = 0; i < 4; i++)
    {
        pinMode(OTS_Simon::ButtonOrder[i], INPUT_PULLUP);
        pinMode(OTS_Simon::LEDOrder[i], OUTPUT);
    }
}

void OTS_Simon::processButtons()
{
    for (uint8_t ButtonNumber = 0; ButtonNumber < 4; ButtonNumber++)
    {
        // read the state of the switch/button:
        bool currentState = !digitalRead(OTS_Simon::ButtonOrder[ButtonNumber]);

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

bool OTS_Simon::checkSerialNumVowel() {
    const char PROGMEM vowels[5] = {'a', 'e', 'i', 'o', 'u'};
    for (uint8_t i = 0; i < 5; i++) {
        for (uint8_t j = 0; j < 8; j++) {
            if (vowels[i] == OTS_Simon::serialNumber[j]) {
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
            digitalWrite(OTS_Simon::LEDOrder[i], LOW);
        }
        {
            digitalWrite(OTS_Simon::LEDOrder[i], 0);
        }
       
    }
}

void OTS_Simon::playSequance() {
    if (OTS_Simon::sequencePosition >= OTS_Simon::sequenceLength) {
        OTS_Simon::sequencePosition = 0;
        OTS_Simon::lastFullSequenceTime = millis();
    }
    if (millis() - OTS_Simon::lastFullSequenceTime > 1000 && millis() - OTS_Simon::lastSequenceTime > 150) {
        digitalWrite(OTS_Simon::LEDOrder[OTS_Simon::sequencePosition], HIGH);
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

    for (uint8_t i = 0; i < 4 && !OTS_Simon::playerTurn; i++)
    {
        if (OTS_Simon::buttonStates[i] == true)
        {
            OTS_Simon::playerTurn = true;
        }
    }
    
    if (!playerTurn) { // Play the sequence to the player
        this->playSequance();
    } else {
        uint8_t btnpressed;
        for (uint8_t i = 0; i < 4; i++)
        {
            if (OTS_Simon::buttonStates[i] == true)
            {
                btnpressed = i;
                
            }
        }
        uint8_t actualBTNPressed = OTS_Simon::PlayerSimonMap[OTS_Simon::checkSerialNumVowel()][OTS_Simon::numStrike][btnpressed];
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
        if (OTS_Simon::PlayerEntrysPosition == OTS_Simon::sequencePosition) { // Check for a completed round of simon says
            for (uint8_t i = 0; i < 4; i++)
            {
                OTS_Simon::PlayerEntrys[i] = 0;
            }
            OTS_Simon::PlayerEntrysPosition = 0;
            OTS_Simon::sequenceLength++;
        }

    }
}