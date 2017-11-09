package config

import (
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure"
)

const (
	keychainConfigKey = "keychain"
)

var (
	keychainConfig *Keychain
)

type Keychain struct {
	DatabaseURL        string
	SkipSignatureCheck bool
}

func (c *ViperConfig) Keychain() Keychain {
	if keychainConfig == nil {
		keychainConfig = &Keychain{}
		config := c.GetStringMap(keychainConfigKey)
		if err := figure.Out(keychainConfig).From(config).Please(); err != nil {
			panic(errors.Wrap(err, "failed to figure out keychain"))
		}
		if keychainConfig.DatabaseURL == "" {
			panic(errors.New("database_url is required"))
		}
	}
	return *keychainConfig
}
