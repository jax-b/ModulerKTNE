namespace N
{
    class BaseModule
    {
    public:
        bool checkSuccess();
        bool checkFailure();
        void runModule();
        void setupModule();
        char[] modID();
    protected:
        char[][3] indicators;
        int seed;
        int batteries;
        bool successtriggered;
        bool failuretriggered;
        char[] modID = "base";
    };
}
