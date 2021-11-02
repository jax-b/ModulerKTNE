#include "basemodule.h"

void BaseModule::setupModule()
{
    Serial.print("HelloWorld");
}
void BaseModule::runModule()
{
    Serial.print(".");
}
bool BaseModule::checkFailure()
{
    return failuretriggered;
    failuretriggered = false;
}
bool BaseModule::checkSuccess()
{
    return successtriggered;
}
