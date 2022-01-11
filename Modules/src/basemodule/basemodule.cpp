#include "baseModule.h"

// Check if there is a failure in the module and if so, clear it
bool baseModule::checkFailure()
{   
    bool out = failureTriggered;
    failureTriggered = false;
    return out;
}
// Check if the module is successfully solved
bool baseModule::checkSuccess()
{
    return successTriggered;
}
char[] baseModule::getID()
{
    return modID;
}
void baseModule::setIndicators(char[][3] inIndicator) {
    for (uint8_t i = 0; i < sizeof(inIndicator); i++) {
        litIndicators[i] = inIndicator[i];
    }
}
void baseModule::setBatteries(uint8_t inNumBattery) {
    numBatteries = inNumBattery;
}
void baseModule::setNumStrike(uint8_t inNumStrike) {
    numStrike = inNumStrike;
}

bool baseModule::checkIndicator(char[3] inIndicator) {
    for (uint8_t i = 0; i < sizeof(litIndicators); i++) {
        if (litIndicators[i] == inIndicator) {
            return true;
        }
    }
    return false
}