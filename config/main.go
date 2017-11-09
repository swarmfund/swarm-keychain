package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type ConfigI interface {
	Init() error
	Keychain() Keychain
	HTTP() HTTP
	Core() Core
	Log() Log
}

type ViperConfig struct {
	*viper.Viper
}

func NewViperConfig(fn string) ConfigI {
	config := ViperConfig{
		viper.GetViper(),
	}
	config.SetConfigFile(fn)
	return &config
}

func (c *ViperConfig) Init() error {
	if err := viper.ReadInConfig(); err != nil {
		return errors.Wrap(err, "failed to read config file")
	}
	return nil
}
