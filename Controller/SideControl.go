package controller

import (
	"github.com/d2r2/go-i2c"
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
}

func NewSideControl(address byte, bus int) (*SideControl, error) {
	i2c, err := i2c.NewI2C(address, bus)
	sc := &SideControl{

// Set Serial Number
func SetSerialNumber