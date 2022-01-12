namespace N
{
    class baseModule
    {
    public:
        bool checkSuccess(); 
        bool checkFailure();
        void tickModule(uint16_t); // Needs to be set by the child class to preform loop specific tasks
        void setupModule(); // Needs to be set by the child class to preform startup code 
        char* getModuleName();
        void setIndicators(char[][3]);
        void setBatteries(uint8_t);
        void setNumStrike(uint8_t);
        void setSeed(uint16_t); // Needs to be set by the child class as it might set up viewer related stuff
        void setSerialNumber(char[8]);
    protected:
        bool successTriggered;
        bool failureTriggered;
        char[] modID; // Needs to be set by the child class for identification
        char[][3] litIndicators;
        uint8_t numBatteries;
        uint8_t numStrike;
        char[8] serialNumber;
        uint16_t seed;
        bool checkIndicator(char[3]);
    };
}