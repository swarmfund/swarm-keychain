package keychain

import (
	"bullioncoin.githost.io/development/keychain/db2"
	"bullioncoin.githost.io/development/keychain/db2/core"
	"bullioncoin.githost.io/development/keychain/db2/keychain"
	"bullioncoin.githost.io/development/keychain/log"
)

func initKeychainDb(app *App) {
	repo, err := db2.Open(app.config.KeychainDatabaseURL)

	if err != nil {
		log.Panic(err)
	}
	repo.DB.SetMaxIdleConns(4)
	repo.DB.SetMaxOpenConns(12)

	app.keychainQ = &keychain.Q{Repo: repo}
}

func initCoreDb(app *App) {
	repo, err := db2.Open(app.config.StellarCoreDatabaseURL)

	if err != nil {
		log.Panic(err)
	}

	repo.DB.SetMaxIdleConns(4)
	repo.DB.SetMaxOpenConns(12)
	app.coreQ = core.NewQ(repo)
}

func init() {
	appInit.Add("keychain-db", initKeychainDb, "app-context", "log")
	appInit.Add("core-db", initCoreDb, "app-context", "log")
}
