package main

import (
	"log"

	"gitlab.com/tokend/keychain/db2/keychain/schema"
	"github.com/spf13/cobra"
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
	dbCmd.AddCommand(keychainCmd)
	keychainCmd.AddCommand(dbMigrateCmd)
}

var dbMigrateCmd = &cobra.Command{
	Use:   "migrate [up|down|redo] [COUNT]",
	Short: "migrate schema",
	Long:  "performs a schema migration command",
	Run: func(cmd *cobra.Command, args []string) {
		err := conf.Init()
		if err != nil {
			log.Fatal(err)
		}
		migrateDB(cmd, args, conf.KeychainDatabaseURL, schema.Migrate)
	},
}
