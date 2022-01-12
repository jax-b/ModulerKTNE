#include "OTS_Simon.h"

OTS_Simon::OTS_Simon()
{
    modID = "smon"
}

void OTS_Simon::setupModule() override
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

        if (currentState != buttonStatesFlicker[ButtonNumber])
        {
            lastDebounceTime[ButtonNumber] = millis();
            buttonStatesFlicker[ButtonNumber] = currentState;
        }

        if ((millis() - lastDebounceTime[ButtonNumber]) > 5)
        {
            // save the the last state
            buttonStates[ButtonNumber] = currentState;
        }
    }
}

bool OTS_Simon::checkSerialNumVowel() {
    char[] vowels = {'a', 'e', 'i', 'o', 'u'};
    for (uint8_t i = 0; i < 5; i++) {
        for (uint8_t j = 0; j < 8; j++) {
            if (vowels[i] == serialNum[j]) {
                return true;
            }
        }
    }
    return false;
}

void OTS_Simon::setSeed(uint16_t inSeed) override
{
    seed = inSeed;
    // Seed to set the sequence
    // sequence = ;
}

void OTS_Simon::resetLEDs()
{
    for (uint8_t i = 0; i < 4; i++)
    {
        if (millis() - lastLEDOnTime[i] > 100 || !playerTurn)
        {
            digitalWrite(LEDOrder[i], LOW);
        }
        {
            digitalWrite(LEDOrder[i], 0);
        }
       
    }
}

void OTS_Simon::playSequance() {
    if sequencePosition >= sequenceLength {
        sequencePosition = 0;
        lastFullSequenceTime = millis();
    }
    if (millis() - lastFullSequenceTime > 1000 && millis() - lastSequenceTime > 150) {
        digitalWrite(LEDOrder[sequencePosition], HIGH);
        lastLEDOnTime[sequencePosition] = millis();
        lastSequenceTime = millis();
        sequencePosition++;
    }
}

void OTS_Simon::tickModule(uint16_t currentGameTime) override
{ 
    this->processButtons();
    this->resetLEDs();
    // We need to wait if the master controller has not cleared out the failure flag
    if (failureTriggered)
    {
        return;
    }

    for (uint8_t i = 0; i < 4 && !playerTurn; i++)
    {
        if (buttonStates[i] == true)
        {
            playerTurn = true;
        }
    }
    
    if (!playerTurn) { // Play the sequence to the player
        this->playSequance();
    } else {
        uint8_t btnpressed;
        for (uint8_t i = 0; i < 4; i++)
        {
            if (buttonStates[i] == true)
            {
                btnpressed = i;
                
            }
        }
        actualBTNPressed = PlayerSimonMap[checkSerialNumVowel()][numStrike][btnpressed];
        PlayerEntrys[PlayerEntrysPosition] = actualBTNPressed;

        // Check for an invalid press
        for (uint8_t i = 0; i <= sequencePosition; i++)
        {
            if (PlayerEntrys[i] != sequence[i])
            {
                failureTriggered = true;
                return;
            }
        } 
        if PlayerEntrysPosition == sequencePosition { // Check for a completed round of simon says
            PlayerEntrys = {0,0,0,0};
            PlayerEntrysPosition = 0;
            sequenceLength++;
        }

    }
}