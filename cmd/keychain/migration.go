package main

import (
	"bullioncoin.githost.io/development/keychain/db2"
	"database/sql"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strconv"
)

type Migrator func(*sql.DB, db2.MigrateDir, int) (int, error)

func migrateDB(cmd *cobra.Command, args []string, dbConnectionURL string, migrator Migrator) {

	// Allow invokations with 1 or 2 args.  All other args counts are erroneous.
	if len(args) < 1 || len(args) > 2 {
		cmd.Usage()
		os.Exit(1)
	}

	dir := db2.MigrateDir(args[0])
	count := 0

	// If a second arg is present, parse it to an int and use it as the count
	// argument to the migration call.
	if len(args) == 2 {
		var err error
		count, err = strconv.Atoi(args[1])
		if err != nil {
			log.Println(err)
			cmd.Usage()
			os.Exit(1)
		}
	}

	db, err := sql.Open("postgres", dbConnectionURL)
	if err != nil {
		log.Fatal(err)
	}

	count, err = migrator(db, dir, count)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Applied %d migrations", count)
}
