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
        char[][] indicators;
        int seed;
        bool successtriggered;
        bool failuretriggered;
        char[] modID = "base";
    };
}
