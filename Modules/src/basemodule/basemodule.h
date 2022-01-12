#include <Arduino.h>

#define GAMEPLAYMAXLITINDICATOR 6
#define GAMEPLAYSERIALNUMBERLENGTH 8


class baseModule
{
public:
    bool checkSuccess(); 
    bool checkFailure();
    void tickModule(uint16_t); // Needs to be set by the child class to preform loop specific tasks
    void setupModule(); // Needs to be set by the child class to preform startup code 
    char* getModuleName();
    void setIndicators(char[GAMEPLAYMAXLITINDICATOR][3]);
    void setBatteries(uint8_t);
    void setNumStrike(uint8_t);
    void setSeed(uint16_t); // Needs to be set by the child class as it might set up viewer related stuff
    void setSerialNumber(char[8]);
    void setPorts(uint8_t);
protected:
    bool successTriggered;
    bool failureTriggered;
    char modID[4]; // Needs to be set by the child class for identification
    char litIndicators[GAMEPLAYMAXLITINDICATOR][3];
    uint8_t numBatteries;
    uint8_t numStrike;
    uint8_t activePorts;
    char serialNumber[GAMEPLAYSERIALNUMBERLENGTH];
    uint16_t seed;
    bool checkIndicator(const char[3]);
};
