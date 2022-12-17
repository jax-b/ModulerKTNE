#include "OTS_Wires.h"

OTS_Wires::OTS_Wires()
{   
    const char realModID[] = "wire";
    for (uint8_t i = 0; i < 4; i++)
    {
        modID[i] = realModID[i];
    }
}

const uint8_t PROGMEM wireButtonPins[6] = {
    WIRE_BUTTON_PIN_1,
    WIRE_BUTTON_PIN_2,
    WIRE_BUTTON_PIN_3,
    WIRE_BUTTON_PIN_4,
    WIRE_BUTTON_PIN_5,
    WIRE_BUTTON_PIN_6
};

void OTS_Wires::clearModule()
{
    for (uint8_t i = 0; i < 12; i++) {
        OTS_Wires::_pixels->setPixelColor(i, 0);
    }
    OTS_Wires::_pixels->show();
    OTS_Wires::failureTriggered = false;
    OTS_Wires::successTriggered = false;
    OTS_Wires::numBatteries = 0;
    OTS_Wires::numStrike = 0;
    OTS_Wires::activePorts = 0;
    OTS_Wires::seed = 0;
    for (uint8_t i = 0; i < 6; i++) {
        OTS_Wires::wireCuts[i] = false;
    }
}
    
void OTS_Wires::setupModule()
{
    // Setup Button Pin States
    for (uint8_t i = 0; i < 6; i++)
    {
        pinMode(wireButtonPins[i], INPUT_PULLUP);
    }

    // Setup the NeoPixel strip
    _pixels = new Adafruit_NeoPixel(12, NEOPIXEL_PIN, NEO_GRB + NEO_KHZ800);
    _pixels->begin();
    for (uint8_t i = 0; i < 12; i++)
    {
        _pixels->setPixelColor(i, WIRE_COLOR[0]);
    }
    _pixels->show();
}

void OTS_Wires::processButtons()
{
    for (uint8_t ButtonNumber = 0; ButtonNumber < 6; ButtonNumber++)
    {
        // read the state of the switch/button:
        bool currentState = digitalRead(wireButtonPins[ButtonNumber]);

        if (currentState != OTS_Wires::buttonStatesFlicker[ButtonNumber])
        {
            OTS_Wires::lastDebounceTime[ButtonNumber] = millis();
            OTS_Wires::buttonStatesFlicker[ButtonNumber] = currentState;
        }

        if ((millis() - OTS_Wires::lastDebounceTime[ButtonNumber]) > 5)
        {
            // save the the last state
            OTS_Wires::buttonStates[ButtonNumber] = currentState;
        }
    }
}

void OTS_Wires::setSeed(uint16_t inSeed)
{
    OTS_Wires::seed = inSeed;
    if (seed == 0)
    {
        return;
    }
    OTS_Wires::numWires = seed; //fancy math
    
    for (uint8_t wire = 0; wire < 6; wire++) {
        // Need math to tell which wire is which color and if its connected
        // while making sure that there is at 3 wires up to all 6
        OTS_Wires::wireCuts[wire] = true;
        OTS_Wires::wireColors[wire] = seed % 5;
        if (wireCuts[wire]) {
            OTS_Wires::_pixels->setPixelColor(wire, WIRE_COLOR[wireColors[wire]]);
            OTS_Wires::_pixels->setPixelColor(wire+6, WIRE_COLOR[wireColors[wire]]);
        }
    }
    OTS_Wires::_pixels->show();
}

void OTS_Wires::cutWire(uint8_t buttonNumber){
    if (OTS_Wires::wireCuts[buttonNumber]) {
        
        bool cutGood = false;
        // GameLogic for wire cutting Set Cut good to true if the wire is the right color

        if (cutGood) {
            for (uint8_t i = 0; i < 12; i++)
            {
                OTS_Wires::_pixels->setPixelColor(i, WIRE_COLOR[0]);
            }
            OTS_Wires::_pixels->show();
            OTS_Wires::successTriggered = true;
            OTS_Wires::wireCuts[buttonNumber] = false;
        } else {
            OTS_Wires::failureTriggered = true;
            OTS_Wires::_pixels->setPixelColor(buttonNumber+6, WIRE_COLOR[0]);
            OTS_Wires::_pixels->show();
            OTS_Wires::wireCuts[buttonNumber] = false;
        }
    }
}

void OTS_Wires::tickModule(uint16_t currentGameTime)
{
    OTS_Wires::processButtons();
    for (uint8_t buttonNumber = 0; buttonNumber < 6; buttonNumber++){
        if (buttonStates[buttonNumber] == true) {
            OTS_Wires::cutWire(buttonNumber);
        }
    }

    // We need to wait if the master controller has not cleared out the failure flag
    if (failureTriggered)
    {
        return;
    }
}