package keychain

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"gitlab.com/swarmfund/go/doorman"
	"gitlab.com/swarmfund/keychain/config"
	"gitlab.com/swarmfund/keychain/coreinfo"
	"gitlab.com/swarmfund/keychain/db2"
	"gitlab.com/swarmfund/keychain/db2/core"
	"gitlab.com/swarmfund/keychain/db2/keychain"
	"gitlab.com/swarmfund/keychain/internal/api"
	"gitlab.com/swarmfund/keychain/log"
	"golang.org/x/net/context"
)

// You can override this variable using: gb build -ldflags "-X main.version aabbccdd"
var version = ""

// App represents the root of the state of a horizon instance.
type App struct {
	_config   config.Config
	coreQ     core.QInterface
	keychainQ *keychain.Q
	ctx       context.Context
	cancel    func()

	horizonVersion string
}

// SetVersion records the provided version string in the package level `version`
// var, which will be used for the reported horizon version.
func SetVersion(v string) {
	version = v
}

// NewApp constructs an new App instance from the provided config.
func NewApp(config config.Config) (*App, error) {
	app := &App{
		_config: config,
	}

	app.init()
	return app, nil
}

func (a *App) Config() config.Config {
	return a._config
}

func (a *App) CoreAccountQ() *core.AccountQ {
	return core.NewAccountQ(a.CoreRepo(a.ctx))
}

func (a *App) Serve() {
	addr := fmt.Sprintf("%s:%d", a.Config().HTTP().Host, a.Config().HTTP().Port)
	coreConnector, err := coreinfo.NewConnector(a._config.Core().URL)

	if err != nil {
		panic(errors.Wrap(err, "failed to get core connector"))
	}

	router := api.Router(
		log.WithField("service", "api"),
		doorman.New(false, a.CoreAccountQ()),
		a.KeychainQ(),
		coreConnector,
	)

	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	log.WithFields(log.F{
		"addr":    addr,
		"service": "api",
	}).Info("listening")

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}

// Close cancels the app and forces the closure of db connections
func (a *App) Close() {
	a.cancel()

	a.keychainQ.GetRepo().DB.Close()
	a.coreQ.GetRepo().DB.Close()

}

// CoreQ returns a helper object for performing sql queries against the
// stellar core database.
func (a *App) CoreQ() core.QInterface {
	return a.coreQ
}

func (a *App) KeychainQ() *keychain.KeyQ {
	return &keychain.KeyQ{
		Repo: a.KeychainRepo(a.ctx),
	}
}

// CoreRepo returns a new repo that loads data from the stellar core
// database. The returned repo is bound to `ctx`.
func (a *App) CoreRepo(ctx context.Context) *db2.Repo {
	return &db2.Repo{DB: a.coreQ.GetRepo().DB, Ctx: ctx}
}

func (a *App) KeychainRepo(ctx context.Context) *db2.Repo {
	return &db2.Repo{DB: a.keychainQ.GetRepo().DB, Ctx: ctx}
}

// Tick triggers horizon to update all of it's background processes such as
// transaction submission, metrics, ingestion and reaping.
func (a *App) Tick() {
	log.Debug("ticking app")

	log.Debug("finished ticking app")
}

// Init initializes app, using the config to populate db connections and
// whatnot.
func (a *App) init() {
	appInit.Run(a)
}

// run is the function that runs in the background that triggers Tick each
// second
func (a *App) run() {
	for {
		select {
		case <-a.ctx.Done():
			log.Info("finished background ticker")
			return
		}
	}
}
