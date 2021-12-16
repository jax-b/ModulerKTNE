package netdisplay

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	// Default Config
	defaultMultiCastIP   = "228.246.151.248"
	defaultMultiCastPort = "26960"
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
		"UseMultiCast":  "true",
		"MultiCastIP":   defaultMultiCastIP,
		"MultiCastPort": defaultMultiCastPort,
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
}
