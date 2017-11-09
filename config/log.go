package config

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gitlab.com/distributed_lab/figure"
)

const (
	logConfigKey = "log"
)

var (
	logConfig    *Log
	logLevelHook = figure.Hooks{
		"logrus.Level": func(value interface{}) (reflect.Value, error) {
			switch v := value.(type) {
			case string:
				lvl, err := logrus.ParseLevel(v)
				if err != nil {
					return reflect.Value{}, errors.Wrap(err, "failed to parse log level")
				}
				return reflect.ValueOf(lvl), nil
			case nil:
				return reflect.ValueOf(nil), nil
			default:
				return reflect.Value{}, fmt.Errorf("unsupported conversion from %T", value)
			}
		},
	}
)

type Log struct {
	// Level log level, default is `warn`
	Level logrus.Level
}

func (c *ViperConfig) Log() Log {
	if logConfig == nil {
		logConfig = &Log{
			Level: logrus.WarnLevel,
		}
		config := c.GetStringMap(logConfigKey)
		if err := figure.Out(logConfig).With(logLevelHook).From(config).Please(); err != nil {
			panic(errors.Wrap(err, "failed to figure out log"))
		}
	}
	return *logConfig
}
