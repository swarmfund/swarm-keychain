package config

import (
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure"
)

const (
	httpConfigKey = "http"
)

var (
	httpConfig *HTTP
)

type HTTP struct {
	Host string
	Port int
}

func (c *ViperConfig) HTTP() HTTP {
	if httpConfig == nil {
		httpConfig = &HTTP{
			Host: "localhost",
			Port: 8003,
		}
		config := c.GetStringMap(httpConfigKey)
		if err := figure.Out(httpConfig).From(config).Please(); err != nil {
			panic(errors.Wrap(err, "failed to figure out http"))
		}
		if httpConfig.Host == "" {
			panic(errors.New("host is required"))
		}
		if httpConfig.Port == 0 {
			panic(errors.New("port is required"))
		}
	}
	return *httpConfig
}
