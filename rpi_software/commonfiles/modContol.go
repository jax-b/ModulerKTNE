package commonfiles

import (
	"fmt"
	"math"
	"time"

	"github.com/d2r2/go-i2c"
	"go.uber.org/zap"
)

const (
	// Front Face
	// Main | 1 | 2
	// 3  | 4 | 5
	FRONT_MOD_1 byte = 0x30
	FRONT_MOD_2 byte = 0x31
	FRONT_MOD_3 byte = 0x32
	FRONT_MOD_4 byte = 0x33
	FRONT_MOD_5 byte = 0x34
	// Back Face
	// 1 | 2 | Main
	// 3 | 4 | 5
	BACK_MOD_1 byte = 0x40
	BACK_MOD_2 byte = 0x41
	BACK_MOD_3 byte = 0x42
	BACK_MOD_4 byte = 0x43
	BACK_MOD_5 byte = 0x44
)

type ModControl struct {
	i2c *i2c.I2C
	log *zap.SugaredLogger
}

// Returns a new module controller
func NewModControl(logger *zap.SugaredLogger, address byte, bus int) *ModControl {
	logger = logger.Named("ModControl-" + fmt.Sprintf("%02X", address))
	i2c, err := i2c.NewI2C(address, bus)
	if err != nil {
		logger.Error("Error opening i2c bus", err)
	}
	mc := &ModControl{
		i2c: i2c,
		log: logger,
	}
	return mc
}

func (smc *ModControl) GetAddress() byte {
	return smc.i2c.GetAddr()
}

// Safely closes the i2c connection
func (smc *ModControl) Close() {
	smc.i2c.Close()
}

// Sends 1 byte of data to the module
// If the byte is not written then the module is not installed
func (smc *ModControl) TestIfPresent() bool {
	smc.log.Debug("Testing if Present")
	bytesWritten, err := smc.i2c.WriteBytes([]byte{0x0})
	if err != nil {
		smc.log.Debug("Module Not Connected")
		return false
	}
	if bytesWritten < 1 {
		smc.log.Debug("Module Not Connected")
		return false
	}
	smc.log.Debug("Module Connected")
	return true
}

// Game Status: Stop all gameplay functions
func (smc *ModControl) StopGame() error {
	_, err := smc.i2c.WriteBytes([]byte{0x40})
	if err != nil {
		smc.log.Error("Failed to write stop to bus: ", err)
		return err
	}
	smc.log.Debug("Stopped Game")
	return nil
}

// Game Status: Start all gameplay functions
/// Tells the module to start the game and start its internal timer
/// if no game seed is set a random seed will be generated
func (smc *ModControl) StartGame() error {
	_, err := smc.i2c.WriteBytes([]byte{0x30})
	if err != nil {
		smc.log.Error("Failed to write start to bus: ", err)
		return err
	}
	smc.log.Debug("Started Game: ")
	return nil
}

// Clears out the gameplay serial number
func (smc *ModControl) ClearGameSerialNumber() error {
	_, err := smc.i2c.WriteBytes([]byte{0x24})
	if err != nil {
		smc.log.Error("Failed to write Clear Serial Number to bus: ", err)
		return err
	}
	smc.log.Debug("Cleared Game Serial Number")
	return nil
}

// Clears out all gameplay lit indicators
func (smc *ModControl) ClearGameLitIndicator() error {
	_, err := smc.i2c.WriteBytes([]byte{0x25})
	if err != nil {
		smc.log.Error("Failed to write Clear Game Lit Indicator to bus: ", err)
		return err
	}
	smc.log.Debug("Cleared Game Lit Indicator")
	return nil
}

// Clears out the number of batteries
func (smc *ModControl) ClearGameNumBatteries() error {
	_, err := smc.i2c.WriteBytes([]byte{0x26})
	if err != nil {
		smc.log.Error("Failed to write Clear batteries to bus: ", err)
		return err
	}
	smc.log.Debug("Cleared Game Num Batteries")
	return nil
}

// Clears out all ports from the module
func (smc *ModControl) ClearGamePortIDS() error {
	_, err := smc.i2c.WriteBytes([]byte{0x27})
	if err != nil {
		smc.log.Error("Failed to write Clear Game ID to bus: ", err)
		return err
	}
	smc.log.Debug("Cleared Game Port IDs")
	return nil
}

// Clears out the game seed from the module
func (smc *ModControl) ClearGameSeed() error {
	_, err := smc.i2c.WriteBytes([]byte{0x28})
	if err != nil {
		smc.log.Error("Failed to write Clear Game Seed to bus: ", err)
		return err
	}
	smc.log.Debug("Cleared Game Seed")
	return nil
}

// Sets the solved status of the module
// 0 = Unsolved
// 1 = Solved
// anything negative is the number of strikes up to -128
func (smc *ModControl) SetSolvedStatus(status int8) error {
	_, err := smc.i2c.WriteBytes([]byte{0x11, byte(status)})
	if err != nil {
		smc.log.Error("Failed to write Solved Status to bus: ", err)
		return err
	}
	smc.log.Debug("Set Solved Status: ", status)
	return nil
}

// Set: Sync Game Time
// Sets the game time to the current time
func (smc *ModControl) SyncGameTime(time time.Duration) error {
	timeout := time.Milliseconds()
	_, err := smc.i2c.WriteBytes([]byte{0x12, byte(timeout >> 24), byte(timeout >> 16), byte(timeout >> 8), byte(timeout)})
	if err != nil {
		smc.log.Error("Failed to write Sync Game Time to bus: ", err)
		return err
	}
	smc.log.Debug("Set Game Time: ", timeout)
	return nil
}

// Sets how fast the module should accelerate per strike
// should be 0 <= x < 1
// defaults to 0.25
func (smc *ModControl) SetStrikeReductionRate(rate float32) error {
	n := math.Float32bits(rate)
	_, err := smc.i2c.WriteBytes([]byte{0x13, byte(n >> 24), byte(n >> 16), byte(n >> 8), byte(n)})
	if err != nil {
		smc.log.Error("Failed to write Strike Reduction Rate to bus: ", err)
		return err
	}
	smc.log.Debug("Set Strike Reduction Rate: ", rate)
	return nil
}

// Set Game Serial Number
func (smc *ModControl) SetGameSerialNumber(serialnumber [8]rune) error {
	buff := []byte{0x14}
	for i := 0; i < 8; i++ {
		buff = append(buff, byte(serialnumber[i]))
	}
	_, err := smc.i2c.WriteBytes(buff)
	if err != nil {
		smc.log.Error("Failed to write Game Serial Number to bus: ", err)
		return err
	}
	smc.log.Debug("Set Game Serial Number: ", serialnumber)
	return nil
}

// Sets a lit indicator, only provide indicators that are lit
// Each time a lit indicator is sent, the module will append it to the list of indicators
// There is a max of 6 indicators that can be lit at a time
// once all are sent the last indicator will be overwritten if more data is sent
// use clearGameLitIndicators to clear the list
// the indicator label is exactly 3 characters long
func (smc *ModControl) SetGameLitIndicator(indlabel [3]rune) error {
	buff := []byte{0x15}
	for i := 0; i < 3; i++ {
		buff = append(buff, byte(indlabel[i]))
	}
	_, err := smc.i2c.WriteBytes(buff)
	if err != nil {
		smc.log.Error("Failed to write Game Lit Indicator to bus: ", err)
		return err
	}
	smc.log.Debug("Set Game Lit Indicator: ", indlabel)
	return nil
}

// Set Game Num Batteries
// 0 - 255
func (smc *ModControl) SetGameNumBatteries(num uint8) error {
	_, err := smc.i2c.WriteBytes([]byte{0x16, byte(num)})
	if err != nil {
		smc.log.Error("Failed to write Game Num Batteries to bus: ", err)
		return err
	}
	smc.log.Debug("Set Game Num Batteries: ", num)
	return nil
}

// Set Game Port IDS
// 0x1: DVI-D, 0x2: Parallel, 0x3: PS2, 0x4: RJ-45, 0x5: Serial, 0x6: Stereo RCA
/// There is a max of 6 ports that can be set at a time
// Once all ports are set the last port will be overwritten if more ports are set
// use clearGamePortIDS to clear the list
// Only send that specific port ID once, thats all that matters for the game logic
func (smc *ModControl) SetGamePortID(id byte) error {
	_, err := smc.i2c.WriteBytes([]byte{0x17, id})
	if err != nil {
		smc.log.Error("Failed to write Game Port ID to bus: ", err)
		return err
	}
	smc.log.Debug("Set Game Port ID: ", id)
	return nil
}

// Set Game Seed
// The seed is a 2 byte number, 1-65535
func (smc *ModControl) SetGameSeed(seed uint16) error {
	_, err := smc.i2c.WriteBytes([]byte{0x18, byte(seed >> 8), byte(seed)})
	if err != nil {
		smc.log.Error("Failed to write Game Seed to bus: ", err)
		return err
	}
	smc.log.Debug("Set Game Seed: ", seed)
	time.Sleep(time.Millisecond * 20)
	return nil
}

// Get Module Type
func (smc *ModControl) GetModuleType() ([4]rune, error) {
	smc.log.Debug("Sending ModID Register Load")
	_, err := smc.i2c.WriteBytes([]byte{0x01})
	if err != nil {
		nothing := [4]rune{}
		smc.log.Error("Failed to write Get Module Type to bus: ", err)
		return nothing, err
	}
	buff := make([]byte, 4)
	smc.log.Debug("Reading Mod ID")
	numread, err := smc.i2c.ReadBytes(buff)
	if err != nil {
		nothing := [4]rune{}
		smc.log.Error("Failed to read Get Module Type from bus: ", err)
		return nothing, err
	}
	modtype := [4]rune{}
	for i := 0; i < numread && i < 4; i++ {
		modtype[i] = rune(buff[i])
	}
	smc.log.Debug("Get Module Type: ", modtype)
	return modtype, nil
}

/// Gets the solved status of the module
func (smc *ModControl) GetSolvedStatus() (int8, error) {
	_, err := smc.i2c.WriteBytes([]byte{0x02})
	if err != nil {
		smc.log.Error("Failed to write Get Solved Status to bus: ", err)
		return 0, err
	}
	buff := make([]byte, 1)
	_, err = smc.i2c.ReadBytes(buff)
	if err != nil {
		smc.log.Error("Failed to read Get Solved Status from bus: ", err)
		return 0, err
	}
	smc.log.Debug("Get Solved Status: ", int8(buff[0]))
	return int8(buff[0]), nil
}

// User Automation Functions
// Clear All Game Data from the specified module
func (smc *ModControl) ClearAllGameData() error {
	err := smc.StopGame()
	if err != nil {
		return err
	}
	err = smc.SetSolvedStatus(0)
	if err != nil {
		return err
	}
	err = smc.ClearGameSerialNumber()
	if err != nil {
		return err
	}
	err = smc.ClearGameLitIndicator()
	if err != nil {
		return err
	}
	err = smc.ClearGameNumBatteries()
	if err != nil {
		return err
	}
	err = smc.ClearGamePortIDS()
	if err != nil {
		return err
	}
	err = smc.ClearGameSeed()
	if err != nil {
		return err
	}

	return nil
}

// User Automation Functions
// Setup All Game Data from the specified module
func (smc *ModControl) SetupAllGameData(serialNumber [8]rune, litIndicators [][3]rune, numBatteries uint8, portID uint8, seed ...uint16) error {
	err := smc.ClearAllGameData()
	if err != nil {
		return err
	}
	err = smc.SetGameSerialNumber(serialNumber)
	if err != nil {
		return err
	}
	for i := 0; i < len(litIndicators); i++ {
		err = smc.SetGameLitIndicator(litIndicators[i])
		if err != nil {
			return err
		}
	}
	err = smc.SetGameNumBatteries(numBatteries)
	if err != nil {
		return err
	}
	err = smc.SetGamePortID(portID)
	if err != nil {
		return err
	}
	if len(seed) > 0 {
		err = smc.SetGameSeed(seed[0])
		if err != nil {
			return err
		}
	}
	err = smc.SetSolvedStatus(0)
	return err
}
