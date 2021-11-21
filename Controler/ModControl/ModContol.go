package ModControl

import (
	"math"

	"github.com/d2r2/go-i2c"
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
}

// Returns a new module controller
func NewModControl(address byte, bus int) (*ModControl, error) {
	i2c, err := i2c.NewI2C(address, bus)
	mc := &ModControl{
		i2c: i2c,
	}
	return mc
}

// Safely closes the i2c connection
func (self *ModControl) Close() {
	self.i2c.Close()
}

// Sends 1 byte of data to the module
// If the byte is not writen then the module is not installed
func (self *ModControl) TestIfPresent() bool {
	bytesWritten, err := self.i2c.WriteBytes([]byte{0x00})
	if err != nil {
		return false
	}
	if bytesWritten < 1 {
		return false
	}
	return true
}

// Game Status: Stop all gameplay functions
func (self *ModControl) StopGame() error {
	_, err := self.i2c.WriteBytes([]byte{0x40})
	if err != nil {
		return err
	}
	return nil
}

// Game Status: Start all gameplay functions
/// Tells the module to start the game and start its internal timer
/// if no game seed is set a random seed will be generated
func (self *ModControl) StartGame() error {
	_, err := self.i2c.WriteBytes([]byte{0x30})
	if err != nil {
		return err
	}
	return nil
}

// Clears out the gameplay serial number
func (self *ModControl) ClearGameSerialNumber() error {
	_, err := self.i2c.WriteBytes([]byte{0x24})
	if err != nil {
		return err
	}
	return nil
}

// Clears out all gameplay lit indicators
func (self *ModControl) ClearGameLitIndicator() error {
	_, err := self.i2c.WriteBytes([]byte{0x25})
	if err != nil {
		return err
	}
	return nil
}

// Clears out the number of batteries
func (self *ModControl) ClearGameNumBatteries() error {
	_, err := self.i2c.WriteBytes([]byte{0x26})
	if err != nil {
		return err
	}
	return nil
}

// Clears out all ports from the module
func (self *ModControl) ClearGamePortIDS() error {
	_, err := self.i2c.WriteBytes([]byte{0x27})
	if err != nil {
		return err
	}
	return nil
}

// Clears out the game seed from the module
func (self *ModControl) ClearGameSeed() error {
	_, err := self.i2c.WriteBytes([]byte{0x28})
	if err != nil {
		return err
	}
	return nil
}

// Sets the solved status of the module
// 0 = Unsolved
// 1 = Solved
// anything negative is the number of strikes up to -128
func (self *ModControl) SetSolvedStatus(status int8) error {
	_, err := self.i2c.WriteBytes([]byte{0x11, byte(status)})
	if err != nil {
		return err
	}
	return nil
}

// Set: Sync Game Time
// Sets the game time to the current time
func (self *ModControl) SyncGameTime(value uint32) error {
	_, err := self.i2c.WriteBytes([]byte{0x20, byte(value >> 24), byte(value >> 16), byte(value >> 8), byte(value)})
	if err != nil {
		return err
	}
	return nil
}

// Sets how fast the module should accelerate per strike
// should be 0 <= x < 1
// defaults to 0.25
func (self *ModControl) SetStrikeReductionRate(rate float32) error {
	n := math.Float32bits(rate)
	_, err := self.i2c.WriteBytes([]byte{0x12, byte(n >> 24), byte(n >> 16), byte(n >> 8), byte(n)})
	if err != nil {
		return err
	}
	return nil
}

// Set Game Serial Number
func (self *ModControl) SetGameSerialNumber(serialnumber [8]rune) error {
	buff := []byte{0x13}
	for i := 0; i < 8; i++ {
		buff = append(buff, byte(serialnumber[i]))
	}
	_, err := self.i2c.WriteBytes(buff)
	if err != nil {
		return err
	}
	return nil
}

// Sets a lit indicator, only provide indicators that are lit
// Each time a lit indicator is sent, the module will append it to the list of indicators
// There is a max of 6 indicators that can be lit at a time
// once all are sent the last indicator will be overwritten if more data is sent
// use clearGameLitIndicators to clear the list
// the indicator label is exactly 3 characters long
func (self *ModControl) SetGameLitIndicator(indlabel [3]rune) error {
	buff := []byte{0x14}
	for i := 0; i < 3; i++ {
		buff = append(buff, byte(indlabel[i]))
	}
	_, err := self.i2c.WriteBytes(buff)
	if err != nil {
		return err
	}
	return nil
}

// Set Game Num Batteries
// 0 - 255
func (self *ModControl) SetGameNumBatteries(num uint8) error {
	_, err := self.i2c.WriteBytes([]byte{0x15, byte(num)})
	if err != nil {
		return err
	}
	return nil
}

// Set Game Port IDS
// 0x1: DVI-D, 0x2: Parallel, 0x3: PS2, 0x4: RJ-45, 0x5: Serial, 0x6: Stereo RCA
/// There is a max of 6 ports that can be set at a time
// Once all ports are set the last port will be overwritten if more ports are set
// use clearGamePortIDS to clear the list
// Only send that specific port ID once, thats all that matters for the game logic
func (self *ModControl) SetGamePortID(id byte) error {
	_, err := self.i2c.WriteBytes([]byte{0x16, id})
	if err != nil {
		return err
	}
	return nil
}

// Set Game Seed
// The seed is a 2 byte number, 1-65535
func (self *ModControl) SetGameSeed(seed uint16) error {
	_, err := self.i2c.WriteBytes([]byte{0x17, byte(seed >> 8), byte(seed)})
	if err != nil {
		return err
	}
	return nil
}

// Get Module Type
func (self *ModControl) GetModuleType() ([4]rune, error) {
	_, err := self.i2c.WriteBytes([]byte{0x00})
	if err != nil {
		nothing := [4]rune{}
		return nothing, err
	}
	buff := make([]byte, 4)
	numread, err := self.i2c.ReadBytes(buff)
	if err != nil {
		nothing := [4]rune{}
		return nothing, err
	}
	modtype := [4]rune{}
	for i := 0; i < numread && i < 4; i++ {
		modtype[i] = rune(buff[i])
	}
	return modtype, nil
}

/// Gets the solved status of the module
func (self *ModControl) GetSolvedStatus() (int8, error) {
	_, err := self.i2c.WriteBytes([]byte{0x01})
	if err != nil {
		return 0, err
	}
	buff := make([]byte, 1)
	_, err = self.i2c.ReadBytes(buff)
	if err != nil {
		return 0, err
	}
	return int8(buff[0]), nil
}

// User Automation Functions
// Clear All Game Data from the specified module
func (self *ModControl) ClearAllGameData() error {
	err := self.ClearGameSerialNumber()
	if err != nil {
		return err
	}
	err = self.ClearGameLitIndicator()
	if err != nil {
		return err
	}
	err = self.ClearGameNumBatteries()
	if err != nil {
		return err
	}
	err = self.ClearGamePortIDS()
	if err != nil {
		return err
	}
	err = self.ClearGameSeed()
	if err != nil {
		return err
	}
	err = self.SetSolvedStatus(0)
	if err != nil {
		return err
	}
	err = self.StopGame()
	if err != nil {
		return err
	}
	return nil
}

// User Automation Functions
// Setup All Game Data from the specified module
func (self *ModControl) SetupAllGameData(serialNumber [8]rune, litIndicators [][3]rune, numBatteries uint8, portIDs []uint8, uint16_t seed = nil) error{
    err := ClearGameFromMod()
	if err := nil {
		return err
	}
    err = SetGameSerialNumber(serialNumber)
	if err := nil {
		return err
	}
    for (int i = 0; i < length(indlabel); i++){
        err = SetGameLitIndicator(indlabel[i])
		if err := nil {
			return err
		}
    }
	err = setGameNumBatteries(numBatteries)
	if err := nil {
		return err
	}
    for i := 0; i < length(portIDs); i++){
        err = setGamePortID(portIDs[i])
		if err := nil {
			return err
		}
    }
    if (seed != nil) {
        err = setGameSeed(seed)
		if err := nil {
			return err
		}
    }
    err = setSolvedStatus(0);
    return 1;
}