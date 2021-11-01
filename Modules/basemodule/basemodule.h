namespace N
{
    class BaseModule
    {
    public:
        bool checkSuccess();
        bool checkFailure();
        void runModule();
        void setupModule();
    private:
        bool successtriggered;
        bool failuretriggered;
    };
}