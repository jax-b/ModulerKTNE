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
        bool successtriggered;
        bool failuretriggered;
        char[] modID = "base";
    };
}