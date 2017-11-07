package keychain

import (
	"gitlab.com/distributed_lab/tokend/keychain/db2"
	"gitlab.com/distributed_lab/tokend/keychain/db2/core"
	"gitlab.com/distributed_lab/tokend/keychain/db2/keychain"
	"gitlab.com/distributed_lab/tokend/keychain/log"
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
