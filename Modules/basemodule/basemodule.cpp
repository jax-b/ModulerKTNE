#include "basemodule.h"

void BaseModule::setupModule()
{
    Serial.print("HellWoWorld");
}
void BaseModule::runModule()
{
    Serial.print(".");
}
bool BaseModule::checkFailure()
{
    return true;
}
bool BaseModule::checkSuccess()
{
    return true;
}