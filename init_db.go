package keychain

import (
	"gitlab.com/swarmfund/keychain/db2"
	"gitlab.com/swarmfund/keychain/db2/core"
	"gitlab.com/swarmfund/keychain/db2/keychain"
	"gitlab.com/swarmfund/keychain/log"
)

func initKeychainDb(app *App) {
	repo, err := db2.Open(app.Config().Keychain().DatabaseURL)

	if err != nil {
		log.Panic(err)
	}
	repo.DB.SetMaxIdleConns(4)
	repo.DB.SetMaxOpenConns(12)

	app.keychainQ = &keychain.Q{Repo: repo}
}

func initCoreDb(app *App) {
	repo, err := db2.Open(app.Config().Core().DatabaseURL)

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
