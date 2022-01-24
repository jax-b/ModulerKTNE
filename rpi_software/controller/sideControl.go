package controller

import (
	"fmt"

	"github.com/d2r2/go-i2c"
	"go.uber.org/zap"
)

const (
	// Side Panel addresses
	SIDE_PANEL byte = 0x50
)

type SideControl struct {
	i2c *i2c.I2C
	log *zap.SugaredLogger
}

func NewSideControl(logger *zap.SugaredLogger, bus int) *SideControl {
	logger = logger.Named("SideControl-" + fmt.Sprintf("%02X", SIDE_PANEL))
	i2c, err := i2c.NewI2C(SIDE_PANEL, bus)
	if err != nil {
		logger.Error("Error opening i2c bus", err)
	}
	sc := &SideControl{
		i2c: i2c,
		log: logger,
	}
	return sc
}

func (ssc *SideControl) Close() {
	ssc.ClearAll()
	ssc.i2c.Close()
}

// Sends 1 byte of data to the module
// If the byte is not writen then the module is not installed
func (ssc *SideControl) TestIfPresent() bool {
	bytesWritten, err := ssc.i2c.WriteBytes([]byte{0x0})
	if err != nil {
		ssc.log.Error("Error writing to i2c bus", err)
	}
	if bytesWritten < 1 {
		return false
	}
	return true
}

// Set Serial Number
func (ssc *SideControl) SetSerialNumber(serialnumber string) error {
	buff := []byte{0x10}
	for i := 0; i < 8; i++ {
		buff = append(buff, byte(serialnumber[i]))
	}
	_, err := ssc.i2c.WriteBytes(buff)
	if err != nil {
		return err
	}
	return nil
}

// Set Lit Indicator
// Max 2 per side
// Once one it set caling this function again will set the other indicator for that panel
// If both are set the last one will be replaced with the new value
func (ssc *SideControl) SetIndicator(lit bool, indlabel [3]rune) error {
	var bitlit byte
	if lit {
		bitlit = 1
	}
	buff := []byte{0x11, bitlit}
	for i := 0; i < 3; i++ {
		buff = append(buff, byte(indlabel[i]))
	}
	_, err := ssc.i2c.WriteBytes(buff)
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
func (ssc *SideControl) SetSideArt(artcode byte) error {
	_, err := ssc.i2c.WriteBytes([]byte{0x12, artcode})
	if err != nil {
		return err
	}
	return nil
}

// Clears the serial number
// only works if the address of this instance is the RIGHT_PANEL address
func (ssc *SideControl) ClearSerialNumber() error {
	_, err := ssc.i2c.WriteBytes([]byte{0x20})
	if err != nil {
		return err
	}
	return nil
}

// Clears all set indicators from the panel
func (ssc *SideControl) ClearAllIndicator() error {
	_, err := ssc.i2c.WriteBytes([]byte{0x21})
	if err != nil {
		return err
	}
	return nil
}

// Clears all SideArt from the panel
func (ssc *SideControl) ClearAllSideArt() error {
	_, err := ssc.i2c.WriteBytes([]byte{0x22})
	if err != nil {
		return err
	}
	return nil
}

// Clears Everything from the panel
func (ssc *SideControl) ClearAll() error {
	err := ssc.ClearSerialNumber()
	if err != nil {
		return err
	}
	err = ssc.ClearAllIndicator()
	if err != nil {
		return err
	}
	err = ssc.ClearAllSideArt()
	if err != nil {
		return err
	}
	return nil
}
