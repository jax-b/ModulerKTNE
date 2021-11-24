package controller

import (
	"errors"
	"fmt"

	"github.com/d2r2/go-i2c"
	"go.uber.org/zap"
)

const (
	// Side Panel addresses
	TOP_PANEL    byte = 0x50
	RIGHT_PANEL  byte = 0x51
	BOTTOM_PANEL byte = 0x52
	LEFT_PANEL   byte = 0x53
)

type SideControl struct {
	i2c *i2c.I2C
	log *zap.SugaredLogger
}

func NewSideControl(logger *zap.SugaredLogger, address byte, bus int) *SideControl {
	logger = logger.Named("SideControl-" + fmt.Sprintf("%02X", address))
	i2c, err := i2c.NewI2C(address, bus)
	if err != nil {
		logger.Error("Error opening i2c bus", err)
	}
	sc := &SideControl{
		i2c: i2c,
		log: logger,
	}
	return sc
}

func (self *SideControl) Close() {
	self.ClearAll()
	self.i2c.Close()
}

// Sends 1 byte of data to the module
// If the byte is not writen then the module is not installed
func (self *SideControl) TestIfPresent() bool {
	bytesWritten, err := self.i2c.WriteBytes([]byte{0x00})
	if err != nil {
		self.log.Error("Error writing to i2c bus", err)
	}
	if bytesWritten < 1 {
		return false
	}
	return true
}

// Set Serial Number
func (self *SideControl) SetSerialNumber(serialnumber [8]rune) error {
	if self.i2c.GetAddr() != RIGHT_PANEL {
		return errors.New("Serial Number is not supported on this side")
	}
	buff := []byte{0x10}
	for i := 0; i < 8; i++ {
		buff = append(buff, byte(serialnumber[i]))
	}
	_, err := self.i2c.WriteBytes(buff)
	if err != nil {
		return err
	}
	return nil
}

// Set Lit Indicator
// Max 2 per side
// Once one it set caling this function again will set the other indicator for that panel
// If both are set the last one will be replaced with the new value
func (self *SideControl) SetIndicator(lit bool, indlabel [3]rune) error {
	var bitlit byte
	if lit {
		bitlit = 1
	}
	buff := []byte{0x11, bitlit}
	for i := 0; i < 3; i++ {
		buff = append(buff, byte(indlabel[i]))
	}
	_, err := self.i2c.WriteBytes(buff)
	if err != nil {
		return err
	}
	return nil
}

// Set Side Art
// Will start setting art on that panel then once the panel is full it will overwrite the last panel
// first bit equals art type Battery or Port
// 0 = Battery
// 1 = Port
// for Battery the last 2 bits sets the battery type
// 0 = Battery
// 0 = not used
// 0 = not used
// 0 = not used
// 0 = not used
// 0 = not used
// 1 = 2xAA
// 1 = D
// only one of the two last bits can be set
// for Port the last six bits computes what ports are shown
// 1 = Port
// 0 = not used
// 1 = DVI
// 1 = Parallel
// 1 = PS/2
// 1 = RJ45
// 1 = Serial
// 1 = SteroRCA
// 0x8b would be a parallel port with PS/2 and RJ45
func (self *SideControl) SetSideArt(artcode byte) error {
	_, err := self.i2c.WriteBytes([]byte{0x12, artcode})
	if err != nil {
		return err
	}
	return nil
}

// Clears the serial number
// only works if the address of this instance is the RIGHT_PANEL address
func (self *SideControl) ClearSerialNumber() error {
	if self.i2c.GetAddr() != RIGHT_PANEL {
		return errors.New("Serial Number is not supported on this side")
	}
	_, err := self.i2c.WriteBytes([]byte{0x20})
	if err != nil {
		return err
	}
	return nil
}

// Clears all set indicators from the panel
func (self *SideControl) ClearAllIndicator() error {
	_, err := self.i2c.WriteBytes([]byte{0x21})
	if err != nil {
		return err
	}
	return nil
}

// Clears all SideArt from the panel
func (self *SideControl) ClearAllSideArt() error {
	_, err := self.i2c.WriteBytes([]byte{0x22})
	if err != nil {
		return err
	}
	return nil
}

// Clears Everything from the panel
func (self *SideControl) ClearAll() error {
	if self.i2c.GetAddr() == RIGHT_PANEL {
		err := self.ClearSerialNumber()
		if err != nil {
			return err
		}
	}
	err := self.ClearAllIndicator()
	if err != nil {
		return err
	}
	err = self.ClearAllSideArt()
	if err != nil {
		return err
	}
	return nil
}
