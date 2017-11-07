package config

import (
	"net/url"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Config is the configuration for horizon.  It get's populated by the
// app's main function and is provided to NewApp.
type Config struct {
	*Base
	KeychainDatabaseURL    string
	StellarCoreDatabaseURL string
	StellarCoreURL         string
	Port                   int
	LogLevel               logrus.Level
	LogToJSON              bool

	SlackWebhook *url.URL
	SlackChannel string

	//For developing without signatures
	SkipCheck bool
}

func (c *Config) DefineConfigStructure(cmd *cobra.Command) {
	c.Base = NewBase(nil, "")

	c.setDefault("port", 8000)
	c.setDefault("sign_checkskip", false)
	c.setDefault("log_level", "debug")

	c.bindEnv("port")
	c.bindEnv("keychain_database_url")
	c.bindEnv("stellar_core_database_url")
	c.bindEnv("stellar_core_url")
	c.bindEnv("sign_check_skip")
	c.bindEnv("log_level")
	c.bindEnv("log_to_json")

	c.bindEnv("slack_webhook")
	c.bindEnv("slack_channel")
}

func (c *Config) Init() error {
	c.Port = c.getInt("port")

	var err error
	c.KeychainDatabaseURL, err = c.getNonEmptyString("keychain_database_url")
	if err != nil {
		return err
	}

	c.StellarCoreDatabaseURL, err = c.getNonEmptyString("stellar_core_database_url")
	if err != nil {
		return err
	}

	c.StellarCoreURL, err = c.getNonEmptyString("stellar_core_url")
	if err != nil {
		return err
	}

	c.LogToJSON = c.getBool("log_to_json")
	c.LogLevel, err = logrus.ParseLevel(c.getString("log_level"))
	if err != nil {
		return err
	}

	c.SkipCheck = c.getBool("sign_check_skip")

	if c.getString("slack_webhook") != "" {
		c.SlackWebhook, err = c.getParsedURL("slack_webhook")
		if err != nil {
			return err
		}
		c.SlackChannel, err = c.getNonEmptyString("slack_channel")
		if err != nil {
			return err
		}
	}

	return nil
}
