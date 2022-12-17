using namespace std;
#include <Arduino.h>
#include <Wire.h>
#include <SPI.h>

//#define DEBUG_MODE_TIMER
#define TIMER_ENABLE

// ********************************
// Include the type of module that we use
// ********************************
// #include "basemodule/basemodule.h"
// BaseModule mod = BaseModule();
#include "OTS_Button/OTS_Button.h"
OTS_Button mod = OTS_Button();
// #include "OTS_Wires/OTS_Wires.h"
// OTS_Wires mod = OTS_Wires();


// *************** External Communications Controller/Player *****************
#define AddressInPin A0
#define RandSorcePin 27
#define SuccessLEDPin 8
#define FailureLEDPin 9
#define S2MInteruptPin 4


/// Timekeeping variables for external pins for their reset
unsigned long S2MInteruptCallTime = 0;
unsigned long FailureLEDCallTime = 0;
bool GamePlayLockout = 0;

/// Buffers and data tracking for I2C communication
uint8_t incomeingI2CData[10];
uint8_t outgoingI2CData[10];
uint8_t uint8_tsToSend = 0;
uint8_t uint8_tsReceived = 0;
// **********************************************************************

// In side of the desired module use the following to acctivate the timer
// #define TIMER_ENABLE True

// ************** Default Variable Sizes ******************
#define GAMEPLAYSERIALNUMBERLENGTH 8
#define GAMEPLAYMAXLITINDICATOR 6

// *************** Gameplay Variables *****************
// the specific module  might not use the following variables but they will get set
// by the master controller so we want to keep trac of them for the module specific code
/// each run will have a serial number that is used in that run
char gameplaySerialNumber[GAMEPLAYSERIALNUMBERLENGTH];
/// The device can have a number of lit indicators on the sides but gameplay only cares if they are lit
char gameplayLitIndicators[GAMEPLAYMAXLITINDICATOR][3];
uint8_t gameplayLitIndicatorCount = 0;
/// Most modules will have a seed that it can be set to in order to load a specific module configuration
uint16_t gameplaySeed = 0;
/// Number of battery's that esists on the entire device
uint8_t gameplayNumBattery = 0;
/// There is only six possible ports in the game and they have a bit possition assigned to them.
/// See reference chart in typed notes for more information about which is which
uint8_t gameplayPorts = 0;
/// Timekeeping and gameplay variables
unsigned long gameplayCountdownTime = 0;
bool gameplayTimerRunning = false;
float gameplayStrikeReductionRate = 0.25;
/// All most modules will have a state that they can be set to
/// if a module is failed it will be set to a negative number
/// it will decrement by 1 each time it is failed
int8_t gameplayModuleSolved = 0;
// ****************************************************

uint8_t convertToAddress(uint16_t addrVIn);
void FlagModuleSolved();
void FlagModuleFailed();
#ifdef TIMER_ENABLE
unsigned long timekeeperLastRun = 0;
uint8_t textraCount = 0;
void decrementCounter();
#endif
void requestEvent();
void receiveEvent(int numuint8_ts);
void I2CCommandProcessor();
void setup();
void loop();