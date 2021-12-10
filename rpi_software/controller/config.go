package controller

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	// Default Config
	defaultMultiCastIP        = "228.74.66.248"
	defaultMultiCastPort      = "26960"
	defaultI2cBusNumber       = 1
	defaultbuzzerPinNum       = 18
	defaultstrike1PinNum      = 22
	defaultstrike2PinNum      = 23
	defaultmodInterruptPinNum = 17
	defaultMFBStartPinNum     = 27
	defaultSevenSegAddress    = 0x70
	defaultUseMultiCast       = false
)

type Config struct {
	configReader *viper.Viper
	log          *zap.SugaredLogger
	// ConfigValues
	Network struct {
		UseMulticast  bool
		MultiCastIP   string
		MultiCastPort int
	}
	Shield struct {
		I2cBusNumber       uint8
		BuzzerPinNum       uint8
		Strike1PinNum      uint8
		Strike2PinNum      uint8
		ModInterruptPinNum uint8
		MfbStartPinNum     uint8
		SevenSegAddress    uint8
	}
}

func NewConfig(logger *zap.SugaredLogger) *Config {
	logger.Named("Config")
	c := &Config{
		configReader: viper.New(),
		log:          logger,
	}
	c.configReader.SetConfigName("config")
	c.configReader.SetConfigType("yaml")
	c.configReader.AddConfigPath(".")
	c.configReader.AddConfigPath("/etc/mktne/")

	c.configReader.SetDefault("Network", map[string]string{
		"UseMultiCast":  "false",
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

	return c
}
func (c *Config) Load() {
	err := c.configReader.ReadInConfig()
	if err != nil {
		c.log.Error("Error loading config", err)
	}

	c.populateFromVipers()
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
