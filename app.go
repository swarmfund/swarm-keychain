package keychain

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/pkg/errors"
	coreHelper "gitlab.com/tokend/go/core"
	"gitlab.com/tokend/go/doorman"
	"gitlab.com/tokend/keychain/config"
	"gitlab.com/tokend/keychain/db2"
	"gitlab.com/tokend/keychain/db2/core"
	"gitlab.com/tokend/keychain/db2/keychain"
	"gitlab.com/tokend/keychain/internal/api"
	"gitlab.com/tokend/keychain/log"
	"gitlab.com/tokend/keychain/render/sse"
	"golang.org/x/net/context"
)

// You can override this variable using: gb build -ldflags "-X main.version aabbccdd"
var version = ""

// App represents the root of the state of a horizon instance.
type App struct {
	CoreInfo *coreHelper.Info

	core      *coreHelper.Connector
	_config   config.ConfigI
	web       *Web
	coreQ     core.QInterface
	keychainQ *keychain.Q
	ctx       context.Context
	cancel    func()
	ticks     *time.Ticker

	horizonVersion string
}

// SetVersion records the provided version string in the package level `version`
// var, which will be used for the reported horizon version.
func SetVersion(v string) {
	version = v
}

// NewApp constructs an new App instance from the provided config.
func NewApp(config config.ConfigI) (*App, error) {
	app := &App{
		_config: config,
	}

	coreConnector, err := coreHelper.NewConnector(http.DefaultClient, config.Core().URL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get core connector")
	}
	app.core = coreConnector

	app.ticks = time.NewTicker(1 * time.Second)
	app.init()
	return app, nil
}

func (a *App) Config() config.ConfigI {
	return a._config
}

func (a *App) CoreAccountQ() *core.AccountQ {
	return core.NewAccountQ(a.CoreRepo(a.ctx))
}

func (a *App) Serve() {
	addr := fmt.Sprintf("%s:%d", a.Config().HTTP().Host, a.Config().HTTP().Port)

	router := api.Router(
		log.WithField("service", "api"),
		&doorman.Doorman{
			PassAllChecks: true,
			AccountQ: a.CoreAccountQ(),
		},
		a.KeychainQ(),
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
	a.ticks.Stop()

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

// UpdateStellarCoreInfo updates the value of coreVersion and networkPassphrase
// from the Stellar core API.
func (a *App) UpdateStellarCoreInfo() {
	info, err := a.core.GetCoreInfo()
	if err != nil {
		log.WithField("service", "core-info").WithError(err).Error("could not load stellar-core info")
		return
	}
	a.CoreInfo = info
}

// Tick triggers horizon to update all of it's background processes such as
// transaction submission, metrics, ingestion and reaping.
func (a *App) Tick() {
	var wg sync.WaitGroup
	log.Debug("ticking app")
	// update ledger state and stellar-core info in parallel
	wg.Add(1)
	go func() { a.UpdateStellarCoreInfo(); wg.Done() }()
	wg.Wait()

	sse.Tick()

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
		case <-a.ticks.C:
			a.Tick()
		case <-a.ctx.Done():
			log.Info("finished background ticker")
			return
		}
	}
}
