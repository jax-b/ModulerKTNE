#include <unistd.h>				//Needed for I2C port
#include <fcntl.h>				//Needed for I2C port
#include <sys/ioctl.h>			//Needed for I2C port
#include <linux/i2c-dev.h>		//Needed for I2C port

class ModControl {
    public:
        // Front Face
        // Main | 1 | 2
        // 3  | 4 | 5
        const uint8_t FRONT_MOD_1 = 0x30;
        const uint8_t FRONT_MOD_2 = 0x31;
        const uint8_t FRONT_MOD_3 = 0x32;
        const uint8_t FRONT_MOD_4 = 0x33;
        const uint8_t FRONT_MOD_5 = 0x34;
        // Back Face
        // 1 | 2 | Main
        // 3 | 4 | 5
        const uint8_t BACK_MOD_1 = 0x40;
        const uint8_t BACK_MOD_2 = 0x41;
        const uint8_t BACK_MOD_3 = 0x42;
        const uint8_t BACK_MOD_4 = 0x43;
        const uint8_t BACK_MOD_5 = 0x44;
        
        ModControl(int *busname);

        // Game Status
        /// Stop all gameplay functions
        /// Returns sucess of command
        int8_t stopGame(uint8_t address);
        /// Start all gameplay functions
        /// Tells the module to start the game and start its internal timer
        /// if no game seed is set a random seed will be generated
        /// Returns sucess of command
        int8_t startGame(uint8_t address);

        // Clear
        /// Clears out the gameplay serial number 
        /// Returns sucess of command
        int8_t clearGameSerialNumber(uint8_t address);
        /// Clears out all gameplay lit indicators
        /// Returns sucess of command
        int8_t clearGameLitIndicator(uint8_t address);
        /// Clears out the number of batteries
        /// Returns sucess of command
        int8_t clearGameNumBatteries(uint8_t address);
        /// Clears out all ports from the module
        /// Returns sucess of command
        int8_t clearGamePortIDS(uint8_t address);
        /// Clears out the game seed from the module
        /// Returns sucess of command
        int8_t clearGameSeed(uint8_t address);

        // Set 

        /// Sets the solved status of the module
        /// 0 = Unsolved
        /// 1 = Solved
        /// anything negative is the number of strikes upto -128
        /// Returns sucess of command
        int8_t setSolvedStatus(uint8_t address, int8_t status);
        /// Sync Game Time
        /// Sets the game time to the current time
        /// Returns sucess of command
        int8_t syncGameTime(uint8_t address, int value);
        /// Sets how fast the module should accelerate per strike
        /// should be 0 <= x < 1
        /// defaults to 0.25
        /// Returns sucess of command
        int8_t setStrikeReductionRate(uint8_t address, float rate);
        /// Set Game Serial Number
        /// Returns sucess of command
        int8_t setGameSerialNumber(uint8_t address, char serialnumber[8]);
        /// Sets a lit indicator, only provide indicators that are lit
        /// Each time a lit indicator is sent, the module will append it to the list of indicators
        /// There is a max of 6 indicators that can be lit at a time
        /// once all are sent the last indicator will be overwritten if more data is sent
        /// use clearGameLitIndicators to clear the list
        /// the indicator label is exactly 3 characters long
        /// Returns sucess of command
        int8_t setGameLitIndicator(uint8_t address, char indlabel[3]);
        /// Set Game Num Batteries
        /// 0 - 255
        /// Returns sucess of command
        int8_t setGameNumBatteries(uint8_t address, uint8_t num);
        /// Set Game Port IDS
        /// 0x1: DVI-D, 0x2: Parallel, 0x3: PS2, 0x4: RJ-45, 0x5: Serial, 0x6: Stereo RCA
        /// There is a max of 6 ports that can be set at a time
        /// Once all ports are set the last port will be overwritten if more ports are set
        /// use clearGamePortIDS to clear the list
        /// Only send that specific port ID once, thats all that matters for the game logic
        /// Returns sucess of command
        int8_t setGamePortID(uint8_t address, uint8_t port);
        /// Set Game Seed
        /// The seed is a 2 byte number, 1-65535
        /// Returns sucess of command
        int8_t setGameSeed(uint8_t address, uint16_t seed);

        // Get
        /// Get Module Type
        char * getModuleType(uint8_t address);
        /// Gets the solved status of the module
        int8_t getSolvedStatus(uint8_t address);

        // User Automation Functions
        // Clear All Game Data from the specified module
        /// Returns sucess of command
        int8_t clearAllGameData(uint8_t address);
        // Setup All Game Data from the specified module
        /// Returns sucess of command
        int8_t setupAllGameData(uint8_t address, char serialnumber[], char litIndicators[][3], uint8_t numBatteries, uint8_t portIDs[],uint16_t seed);
    protected:
        int file_i2c;
        
}