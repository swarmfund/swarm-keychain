package main

import (
	"github.com/spf13/cobra"
	"gitlab.com/tokend/keychain/config"
	"gitlab.com/tokend/keychain/db2/keychain/schema"
	"gitlab.com/tokend/keychain/log"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate [up|down|redo] [COUNT]",
	Short: "migrate schema",
	Long:  "performs a schema migration command",
	Run: func(cmd *cobra.Command, args []string) {
		c := config.NewViperConfig(configFile)
		if err := c.Init(); err != nil {
			log.WithField("service", "init").WithError(err).Fatal("failed to init config")
		}
		migrateDB(cmd, args, c.Keychain().DatabaseURL, schema.Migrate)
	},
}
