package config

import (
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure"
)

const (
	coreConfigKey = "core"
)

var (
	coreConfig *Core
)

type Core struct {
	DatabaseURL string
	URL         string
}

func (c *ViperConfig) Core() Core {
	if coreConfig == nil {
		coreConfig = &Core{}
		config := c.GetStringMap(coreConfigKey)
		if err := figure.Out(coreConfig).From(config).Please(); err != nil {
			panic(errors.Wrap(err, "failed to figure out keychain"))
		}
		if coreConfig.DatabaseURL == "" {
			panic(errors.New("database_url is required"))
		}
		if coreConfig.URL == "" {
			panic(errors.New("url is required"))
		}
	}
	return *coreConfig
}
