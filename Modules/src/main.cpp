#include "main.h"

// I2C from AIN Address Table
uint8_t convertToAddress(uint16_t addrVIn)
{
    // Top
    if (addrVIn >= 93 && addrVIn < 186)
    {
        return 0x30;
    }
    else if (addrVIn >= 186 && addrVIn < 279)
    {
        return 0x31;
    }
    else if (addrVIn >= 279 && addrVIn < 372)
    {
        return 0x32;
    }
    else if (addrVIn >= 372 && addrVIn < 465)
    {
        return 0x33;
    }
    else if (addrVIn >= 465 && addrVIn < 558)
    {
        return 0x34;
    }
    // Bottom
    else if (addrVIn >= 558 && addrVIn < 651)
    {
        return 0x40;
    }
    else if (addrVIn >= 651 && addrVIn < 744)
    {
        return 0x41;
    }
    else if (addrVIn >= 744 && addrVIn < 837)
    {
        return 0x42;
    }
    else if (addrVIn >= 837 && addrVIn < 930)
    {
        return 0x43;
    }
    else if (addrVIn >= 930 && addrVIn <= 1024)
    {
        return 0x44;
    }
    else
    {
        return 0x00;
    }
}

// This function is called when the module code determines that the module has been solved
void FlagModuleSolved()
{
    // Set the module solved flag to 1
    gameplayModuleSolved = 1;
    outgoingI2CData[0] = 0x1;
    // Turn on the success LED
    digitalWrite(SuccessLEDPin, HIGH);
    // Turn off the failure LED
    S2MInteruptCallTime = millis();
    // Turn off the failure LED
    digitalWrite(S2MInteruptPin, LOW);
    gameplayTimerRunning = false;
    mod.clearModule();
    gameplaySeed = 0;
    for (uint8_t i = 0; i < GAMEPLAYSERIALNUMBERLENGTH; i++)
    {
        gameplaySerialNumber[i] = '\0';
    }
    gameplayLitIndicatorCount = 0;
    gameplayStrikeReductionRate = 0.25;
    gameplayNumBattery = 0;
    gameplayPorts = 0;
}

// This function is called when the module code determines that the module has been failed
void FlagModuleFailed()
{
    // Set the module solved flag to -1
    gameplayModuleSolved -= 1;
    mod.setNumStrike(gameplayModuleSolved * -1);
    // Record the on time of the failure LED
    FailureLEDCallTime = millis();
    // Turn on the failure LED
    digitalWrite(FailureLEDPin, HIGH);
    // Record the start of the signal to the controller interupt
    S2MInteruptCallTime = millis();
    // Signal the main controller interupt
    digitalWrite(S2MInteruptPin, LOW);
    GamePlayLockout = true;
}

#ifdef TIMER_ENABLE
// Updates to defused time of the module
void decrementCounter()
{
    if (!gameplayTimerRunning || gameplayModuleSolved > 0)
    {
        return;
    }
    unsigned long currentTime = millis();
    // Calculate how long it has been since the last run
    unsigned long timesincelastrun = currentTime - timekeeperLastRun;
    // Subtrack it out
    gameplayCountdownTime -= timesincelastrun;

    // If we have a strike
    if (gameplayModuleSolved < 0)
    {
        float everyrate = (1 / gameplayStrikeReductionRate) / (gameplayModuleSolved * -1);
        if (textraCount >= everyrate)
        {
            unsigned long textra = everyrate;
            if (everyrate < 1)
            {
                textra += 1 / everyrate;
            }
            textraCount = 0;
            gameplayCountdownTime -= textra;
        }
        else
        {
            if (timesincelastrun >= 1)
            {
                textraCount += 1;
            }
        }
    }
    if (gameplayCountdownTime == 0 || gameplayModuleSolved >= 18000000)
    {
        gameplayModuleSolved = 0;
        gameplayTimerRunning = false;
        FlagModuleFailed();
    }
    timekeeperLastRun = currentTime;

#ifdef DEBUG_MODE_TIMER
    uint16_t hundrethsTime = gameplayCountdownTime / 10;
    uint16_t seconds = gameplayCountdownTime / 1000;
    uint16_t minutes = seconds / 60;

    String ctime = String(minutes) + ":" + String(seconds % 60);
    if (minutes < 1)
    {
        ctime += "." + String(hundrethsTime / 10 % 10) + String(hundrethsTime % 10);
    }
    Serial.print("Time: ");
    Serial.print(ctime);
    Serial.print(", timmil:");
    Serial.println(gameplayCountdownTime);

#endif
}
#endif

// Sends out our output buffer
void requestEvent()
{
    if (uint8_tsToSend > 0)
    {

        Serial.print("Sending: ");
        Serial.println(uint8_tsToSend);

        for (size_t i = 0; i < uint8_tsToSend; i++)
        {
            Wire1.write(outgoingI2CData, uint8_tsToSend);
        }
        uint8_tsToSend = 0;
    }
    else
    {

        Wire1.write(0xFF);
    }
}

// Copy the incoming data into our input buffer
void receiveEvent(int numuint8_ts)
{

    Serial.print("Received: ");
    Serial.println(numuint8_ts);
    for (int i = 0; i < numuint8_ts; i++)
    {
        if (i > 10)
        {
            Wire1.read();
        }
        else
        {
            uint8_t wirein = Wire1.read();
            Serial.print(wirein, HEX);
            Serial.print(" ");
            incomeingI2CData[i] = wirein;

            incomeingI2CData[i] = Wire1.read();
        }
    }
    uint8_tsReceived = numuint8_ts;

    Serial.println();
}

// Processe a command from the controller and if necessary copy data into the output buffer
void I2CCommandProcessor()
{

    Serial.print("I2C cmdgroup:");
    Serial.print(incomeingI2CData[0] >> 4, HEX);
    Serial.print(", cmd:");
    Serial.println(incomeingI2CData[0] & 0xF, HEX);
    switch (incomeingI2CData[0] >> 4)
    {
    case 0x4:
        switch (incomeingI2CData[0] & 0xF)
        {
        case 0x0: // Stop
            gameplayTimerRunning = false;

            Serial.println("Game Stoped");
            break;

        default:
            break;
        }
        break;
    case 0x3:
        switch (incomeingI2CData[0] & 0xF)
        {
        case 0x0: // Start
            gameplayTimerRunning = true;
            GamePlayLockout = false;
            if (gameplaySeed == 0)
            {
                gameplaySeed = random(1, 65535);
                mod.setSeed(gameplaySeed);
            }

            Serial.println("Game Started");
            break;

        default:
            break;
        }
        break;
    // this is a clear command
    case 0x2:
        switch (incomeingI2CData[0] & 0xF)
        {
        // Clear SerialNumber
        case 0x4:
            for (int i = 0; i < GAMEPLAYSERIALNUMBERLENGTH; i++)
            {
                gameplaySerialNumber[i] = '\0';
            }
            mod.setSerialNumber(gameplaySerialNumber);

            Serial.println("Serial Number Cleared");
            break;
        // Clear LitIndicators
        case 0x5:
            gameplayLitIndicatorCount = 0;
            for (int i = 0; i < GAMEPLAYMAXLITINDICATOR; i++)
            {
                for (int j = 0; j < 3; j++)
                {
                    gameplayLitIndicators[i][j] = '\0';
                }
            }
            mod.setIndicators(gameplayLitIndicators);

            Serial.println("Lit Indicators Cleared");
            break;
        // Clear Number Batteries
        case 0x6:
            gameplayNumBattery = 0;
            mod.setBatteries(gameplayNumBattery);

            Serial.println("Number Batteries Cleared");
            break;
        // Clear Port Identities
        case 0x7:
            gameplayPorts = 0x0;
            mod.setPorts(gameplayPorts);

            Serial.println("Port Identities Cleared");
            break;
        // Clear Seed
        case 0x8:
            gameplaySeed = 0;
            mod.setSeed(gameplaySeed);
            mod.clearModule();

            Serial.println("Seed Cleared");
            break;
        default:
            break;
        }
        break;
    // This is a command to the module for configuration
    case 0x1:
        switch (incomeingI2CData[0] & 0xF)
        {
        // Set solved status
        case 0x1:
            // greater than 1 uint8_ts because uint8_ts received includes the command uint8_t
            if (uint8_tsReceived > 1)
            {
                gameplayModuleSolved = incomeingI2CData[1];

                Serial.print("SolvedStat: ");
                Serial.println(gameplayModuleSolved);
                if (gameplayModuleSolved > 0)
                {
                    gameplayModuleSolved = 1;
                    outgoingI2CData[0] = 0x1;
                    // Turn on the success LED
                    digitalWrite(SuccessLEDPin, HIGH);
                    gameplayTimerRunning = false;
                    mod.clearModule();
                    gameplaySeed = 0;
                    for (uint8_t i = 0; i < GAMEPLAYSERIALNUMBERLENGTH; i++)
                    {
                        gameplaySerialNumber[i] = '\0';
                    }
                    gameplayLitIndicatorCount = 0;
                    gameplayStrikeReductionRate = 0.25;
                    gameplayNumBattery = 0;
                    gameplayPorts = 0;
                }
                else
                {
                    digitalWrite(SuccessLEDPin, LOW);
                    if (gameplayModuleSolved == 0)
                    {
                        mod.setNumStrike(0);
                    }
                    else
                    {
                        mod.setNumStrike(1);
                    }
                }
            }
            break;
        // Sync Time between the module and the device
        case 0x2:
            // Time should be a unsigned long so 4 uint8_ts
            // greater than 4 uint8_ts because uint8_ts received includes the command uint8_t
            if (uint8_tsReceived > 4)
            {
                uint16_t data34 = incomeingI2CData[3] << 8 | incomeingI2CData[4];
                unsigned long data12 = incomeingI2CData[1] << 8 | incomeingI2CData[2];
                gameplayCountdownTime = data12 << 16 | data34;

                Serial.print("SyncTime: ");
                Serial.println(gameplayCountdownTime);
            }
            break;
        // Set Strike Rate
        case 0x3:
            // Reduction Rate should be a float so 4 uint8_ts
            // greater than 4 uint8_ts because uint8_ts received includes the command uint8_t
            if (uint8_tsReceived > 4)
            {
                uint16_t data34 = incomeingI2CData[3] << 8 | incomeingI2CData[4];
                unsigned long data12 = incomeingI2CData[1] << 8 | incomeingI2CData[2];
                gameplayStrikeReductionRate = data12 << 16 | data34;

                Serial.print("StrikeRate: ");
                Serial.println(gameplayStrikeReductionRate);
            }
            break;
        // Set Serial Number
        case 0x4:
            // Serial Number should be a String max 8 char;
            // greater than 8 uint8_ts because uint8_ts received includes the command uint8_t
            if (uint8_tsReceived > GAMEPLAYSERIALNUMBERLENGTH)
            {
                for (int i = 0; i < uint8_tsReceived && i < GAMEPLAYSERIALNUMBERLENGTH; i++)
                {
                    gameplaySerialNumber[i] = (char)incomeingI2CData[i + 1];
                }
                mod.setSerialNumber(gameplaySerialNumber);

                Serial.print("SerialNumber: ");
                Serial.println(gameplaySerialNumber);
            }
            break;
        // Set LitIndicator
        case 0x5:
            // Lit Indicator should be a string of 3 char
            // greater than 3 uint8_ts because uint8_ts received includes the command uint8_t
            if (uint8_tsReceived > 3)
            {
                for (int i = 0; i < uint8_tsReceived && i < 3; i++)
                {
                    gameplayLitIndicators[gameplayLitIndicatorCount][i] = (char)incomeingI2CData[i + 1];
                }
                if (gameplayLitIndicatorCount < GAMEPLAYMAXLITINDICATOR - 1)
                {
                    gameplayLitIndicatorCount++;
                }
                mod.setIndicators(gameplayLitIndicators);

                Serial.print("LitIndicators: {");
                for (uint8_t i = 0; i < gameplayLitIndicatorCount; i++)
                {
                    for (uint8_t j = 0; j < 3; j++)
                    {
                        Serial.print(gameplayLitIndicators[i][j]);
                    }
                    Serial.print(", ");
                }
                Serial.println("}");
            }
            break;
        // Set Number of Batteries
        case 0x6:
            // Number of Batteries should be a uint8_t
            // greater than 1 uint8_ts because uint8_ts received includes the command uint8_t
            if (uint8_tsReceived > 1)
            {
                gameplayNumBattery = incomeingI2CData[1];
                mod.setBatteries(gameplayNumBattery);

                Serial.print("NumBattery: ");
                Serial.println(gameplayNumBattery);
            }
            break;
        // Set Active Ports
        case 0x7:
            // Port Identities should be a uint8_t
            // greater than 1 uint8_ts because uint8_ts received includes the command uint8_t
            if (uint8_tsReceived > 1)
            {
                gameplayPorts = incomeingI2CData[1];
                mod.setPorts(gameplayPorts);

                Serial.print("Ports: ");
                Serial.println(gameplayPorts);
            }
            break;
        // Set Seed
        case 0x8:
            // Seed should be 2 uint8_ts
            // greater than 2 uint8_ts because uint8_ts received includes the command uint8_t
            if (uint8_tsReceived > 2)
            {
                gameplaySeed = incomeingI2CData[1] << 8 | incomeingI2CData[2];
                mod.setSeed(gameplaySeed);

                Serial.print("Seed: ");
                Serial.println(gameplaySeed);
            }
            break;
        default:
            break;
        }
        break;
    // This is a command to the device for data
    case 0x0:
        switch (incomeingI2CData[0] & 0xF)
        {
        // Get the modules ID
        case 0x1:
            uint8_tsToSend = 3;
            for (int i = 0; i < uint8_tsToSend; i++)
            {
                outgoingI2CData[i] = mod.getModuleName()[i];
            }
            break;
        // Get the Solved Status
        case 0x2:
            uint8_tsToSend = 1;
            outgoingI2CData[0] = gameplayModuleSolved;
            break;
        }
    }
    uint8_tsReceived = 0;
}

void setup()
{

    Serial.begin(115200);
    // while (!Serial); // wait for serial port to connect.
    delay(500);
    Serial.println("Serial Connected");
    // Setup LED's
    pinMode(SuccessLEDPin, OUTPUT);
    pinMode(FailureLEDPin, OUTPUT);
    pinMode(S2MInteruptPin, OUTPUT);

    // Set output pins to inital state
    digitalWrite(S2MInteruptPin, HIGH);
    digitalWrite(SuccessLEDPin, HIGH);
    digitalWrite(FailureLEDPin, HIGH);

    randomSeed(analogRead(RandSorcePin));

    // Setup I2C
    // Start Listening for address
    Wire.onReceive(receiveEvent);
    Wire.onRequest(requestEvent);

    Serial.println("RP2040 Set SPI TX");
    SPI.setTX(23);

    Serial.println("RP2040 Set SPI CLK");
    SPI.setSCK(22);

    Serial.println("RP2040 Set Wire SCL");
    Wire1.setSCL(3);

    Serial.println("RP2040 Set Wire SDA");
    Wire1.setSDA(2);

    uint8_t address = convertToAddress(analogRead(AddressInPin)/4);
    if (address == 0 ){
        digitalWrite(SuccessLEDPin, LOW);
        digitalWrite(FailureLEDPin, HIGH);
        while(true);
    }
    
    Serial.print("I2C Listening on ADDR: ");
    Serial.println(address, HEX);
    Wire1.begin(address);

    Wire1.onReceive(receiveEvent);
    Wire1.onRequest(requestEvent);

    mod.setupModule();

    // Turn off LEDs to signal initialization is complete
    digitalWrite(FailureLEDPin, LOW);
    digitalWrite(SuccessLEDPin, LOW);
}

void loop()
{
    // Pull the S2M Interupt Pin back UP after a short delay
    if (millis() - S2MInteruptCallTime > 50 && S2MInteruptCallTime != 0)
    {
        digitalWrite(S2MInteruptPin, HIGH);
        S2MInteruptCallTime = 0;
    }
    // Turn off the failure LED after a short delay
    if (millis() - FailureLEDCallTime > 700 && FailureLEDCallTime != 0)
    {
        digitalWrite(FailureLEDPin, LOW);
        FailureLEDCallTime = 0;
    }
    if (millis() - FailureLEDCallTime > 400 && FailureLEDCallTime != 0)
    {
        GamePlayLockout = false;
    }

    // Module Specific Code
    if (gameplayModuleSolved <= 0 && gameplayTimerRunning == true && !GamePlayLockout)
    {
        // Success LED should be on for only module success
        digitalWrite(SuccessLEDPin, LOW);
        mod.tickModule(gameplayCountdownTime);

        if (mod.checkSuccess())
        {
            FlagModuleSolved();
        }
        else if (mod.checkFailure())
        {
            FlagModuleFailed();
        }
    }

    if (uint8_tsReceived != 0)
    {
        I2CCommandProcessor();
    }

#ifdef TIMER_ENABLE
    // Timekeeping
    decrementCounter();
#endif
}