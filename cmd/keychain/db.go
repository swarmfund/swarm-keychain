package main

import (
	"github.com/spf13/cobra"
	"gitlab.com/tokend/api/log"
	"gitlab.com/tokend/keychain/config"
	"gitlab.com/tokend/keychain/db2/keychain/schema"
)

var dbCmd = &cobra.Command{
	Use:   "db [command]",
	Short: "commands to manage dbs",
}

var keychainCmd = &cobra.Command{
	Use:   "keychain [command]",
	Short: "commands to manage keychain's database",
}

func init() {

	keychainCmd.AddCommand(dbMigrateCmd)
}

var dbMigrateCmd = &cobra.Command{
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
