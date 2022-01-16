#include <Arduino.h>

#define GAMEPLAYMAXLITINDICATOR 6
#define GAMEPLAYSERIALNUMBERLENGTH 8

class baseModule
{
public:
    bool checkSuccess();
    bool checkFailure();
    void tickModule(uint16_t); // Needs to be set by the child class to preform loop specific tasks
    void setupModule();        // Needs to be set by the child class to preform startup code
    char *getModuleName();
    void setIndicators(char[GAMEPLAYMAXLITINDICATOR][3]);
    void setBatteries(uint8_t);
    void setNumStrike(uint8_t);
    void setSeed(uint16_t); // Needs to be set by the child class as it might set up viewer related stuff
    void setSerialNumber(char[8]);
    void setPorts(uint8_t);
    void clearModule();
protected:
    bool successTriggered = false;
    bool failureTriggered = false;
    char modID[4] = {0, 0, 0, 0}; // Needs to be set by the child class for identification
    char litIndicators[GAMEPLAYMAXLITINDICATOR][3];
    uint8_t numBatteries = 0;
    uint8_t numStrike = 0;
    uint8_t activePorts = 0;
    char serialNumber[GAMEPLAYSERIALNUMBERLENGTH] = {0, 0, 0, 0, 0, 0, 0, 0};
    uint16_t seed = 0;
    bool checkIndicator(const char[3]);
};
