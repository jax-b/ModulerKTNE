#include "baseModule.h"

// Check if there is a failure in the module and if so, clear it
bool baseModule::checkFailure()
{
    bool out = baseModule::failureTriggered;
    baseModule::failureTriggered = false;
    return out;
}
// Check if the module is successfully solved
bool baseModule::checkSuccess()
{
    bool out = baseModule::successTriggered;
    baseModule::successTriggered = false;
    return out;
}
char *baseModule::getModuleName()
{
    #ifdef DEBUG_MODE
        Serial.println("baseModule::modID");
    #endif
    return baseModule::modID;
}
void baseModule::setIndicators(char inIndicators[GAMEPLAYMAXLITINDICATOR][3])
{
    for (uint8_t  i = 0; i < GAMEPLAYMAXLITINDICATOR; i++)
    {
        for (uint8_t  j = 0; j < 3; j++)
        {
            baseModule::litIndicators[i][j] = inIndicators[i][j];
        }
    }
}
void baseModule::setBatteries(uint8_t inNumBattery)
{
    baseModule::numBatteries = inNumBattery;
}
void baseModule::setNumStrike(uint8_t inNumStrike)
{
    baseModule::numStrike = inNumStrike;
}

void baseModule::setSerialNumber(char inSerialNumber[GAMEPLAYSERIALNUMBERLENGTH])
{
    for (uint8_t i = 0; i < GAMEPLAYSERIALNUMBERLENGTH; i++)
    {
        baseModule::serialNumber[i] = inSerialNumber[i];
    }
}

void baseModule::setPorts(uint8_t inPorts)
{
    baseModule::activePorts = inPorts;
}

bool baseModule::checkIndicator(const char inIndicator[3])
{
    for (uint8_t i = 0; i < sizeof(baseModule::litIndicators); i++)
    {
        if (baseModule::litIndicators[i] == inIndicator)
        {
            return true;
        }
    }
    return false;
}

void baseModule::setupModule(){}
void baseModule::tickModule(unsigned long inTime){}
void baseModule::setSeed(uint16_t inSeed){}