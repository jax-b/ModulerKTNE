#include "OTS_Wires.h"

OTS_Wires::OTS_Wires()
{
    modID = "wire"
}

void OTS_Wires::setupModule() override
{
    // Setup Button Pin States
    for (uint8_t i = 0; i < 6; i++)
    {
        pinMode(wireButtonPins[i], INPUT_PULLUP);
    }

    // Setup the NeoPixel strip
    _pixels = new Adafruit_NeoPixel(12, NEOPIXEL_PIN, NEO_GRB + NEO_KHZ800);
    _pixels->begin();
    for (uint8_t i = 0; i < 4; i++)
    {
        _pixels->setPixelColor(i, 0x0);
    }
    _pixels->show();
}

void OTS_Wires::processButtons()
{
    for (uint8_t ButtonNumber = 0; ButtonNumber < 6; ButtonNumber++)
    {
        // read the state of the switch/button:
        bool currentState = digitalRead(ButtonOrder[ButtonNumber]);

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

void OTS_Wires::setSeed(uint16_t inSeed) override
{
    seed = inSeed;
}

void OTS_Wires::tickModule(uint16_t currentGameTime) override
{
    processButtons();

    // We need to wait if the master controller has not cleared out the failure flag
    if (failureTriggered)
    {
        return;
    }
}