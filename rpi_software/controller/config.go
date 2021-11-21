package controller

import (
	"github.com/spf13/viper"
)

const (
	// Default Config
	defaultMultiCastIP        = "228.246.151.248"
	defaultMultiCastPort      = "26960"
	defaultI2cBusNumber       = 1
	defaultbuzzerPinNum       = 18
	defaultstrike1PinNum      = 22
	defaultstrike2PinNum      = 23
	defaultmodInterruptPinNum = 17
	defaultMFBStartPinNum     = 27
	defaultSevenSegAddress    = 0x70
)

type Config struct {
	configReader *viper.Viper
	// ConfigValues
	Network struct {
		MultiCastIP   string `yaml:"multicast_ip"`
		MultiCastPort int    `yaml:"multicast_port"`
	} `yaml:"Network"`
	Shield struct {
		I2cBusNumber       uint8 `yaml:"i2c_bus_number"`
		BuzzerPinNum       uint8 `yaml:"buzzer_pin"`
		Strike1PinNum      uint8 `yaml:"strike_1_pin"`
		Strike2PinNum      uint8 `yaml:"strike_2_pin"`
		ModInterruptPinNum uint8 `yaml:"module_interupt_pin"`
		MfbStartPinNum     uint8 `yaml:"mfb_pin"`
		SevenSegAddress    uint8 `yaml:"seven_segment_address"`
	} `yaml:Shield`
}

func NewConfig() (*Config, error) {
	c := &Config{
		configReader: viper.New(),
	}
	c.configReader.SetConfigName("config")
	c.configReader.SetConfigType("yaml")
	c.configReader.AddConfigPath(".")
	c.configReader.AddConfigPath("/etc/mktne/")

	c.configReader.SetDefault("Network", map[string]string{
		"MultiCastIP":   defaultMultiCastIP,
		"MultiCastPort": defaultMultiCastPort,
	})
	c.configReader.SetDefault("Shield", map[string]uint8{
		"I2cBusNumber":       defaultI2cBusNumber,
		"BuzzerPinNum":       defaultbuzzerPinNum,
		"Strike1PinNum":      defaultstrike1PinNum,
		"Strike2PinNum":      defaultstrike2PinNum,
		"ModInterruptPinNum": defaultmodInterruptPinNum,
		"MfbStartPinNum":     defaultMFBStartPinNum,
		"SevenSegAddress":    defaultSevenSegAddress,
	})

	return c, nil
}
func (c *Config) Load() error {
	err := c.configReader.ReadInConfig()
	if err != nil {
		return err
	}

	c.populateFromVipers()

	return nil
}
func (c *Config) populateFromVipers() {
	c.Network.MultiCastIP = c.configReader.GetString("Network.MultiCastIP")
	c.Network.MultiCastPort = c.configReader.GetInt("Network.MultiCastPort")
	c.Shield.I2cBusNumber = uint8(c.configReader.GetInt("Shield.I2cBusNumber"))
	c.Shield.BuzzerPinNum = uint8(c.configReader.GetInt("Shield.BuzzerPinNum"))
	c.Shield.Strike1PinNum = uint8(c.configReader.GetInt("Shield.Strike1PinNum"))
	c.Shield.Strike2PinNum = uint8(c.configReader.GetInt("Shield.Strike2PinNum"))
	c.Shield.ModInterruptPinNum = uint8(c.configReader.GetInt("Shield.ModInterruptPinNum"))
	c.Shield.MfbStartPinNum = uint8(c.configReader.GetInt("Shield.MfbStartPinNum"))
	c.Shield.SevenSegAddress = uint8(c.configReader.GetInt("Shield.SevenSegAddress"))
}
