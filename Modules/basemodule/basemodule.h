namespace N
{
    class BaseModule
    {
    public:
        bool checkSuccess();
        bool checkFailure();
        void runModule();
        void setupModule();
    protected:
        bool successtriggered;
        bool failuretriggered;
        char[] modID;
    };
}