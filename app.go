package keychain

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"gitlab.com/tokend/go/build"
	"gitlab.com/tokend/keychain/config"
	coreHelper "gitlab.com/tokend/keychain/core"
	"gitlab.com/tokend/keychain/db2"
	"gitlab.com/tokend/keychain/db2/core"
	"gitlab.com/tokend/keychain/db2/keychain"
	"gitlab.com/tokend/keychain/log"
	"gitlab.com/tokend/keychain/render/sse"
	"golang.org/x/net/context"
	"golang.org/x/net/http2"
	"gopkg.in/tylerb/graceful.v1"
)

// You can override this variable using: gb build -ldflags "-X main.version aabbccdd"
var version = ""

// App represents the root of the state of a horizon instance.
type App struct {
	config         config.Config
	web            *Web
	coreQ          core.QInterface
	keychainQ      *keychain.Q
	ctx            context.Context
	cancel         func()
	ticks          *time.Ticker
	CoreInfo       coreHelper.Info
	horizonVersion string
}

// SetVersion records the provided version string in the package level `version`
// var, which will be used for the reported horizon version.
func SetVersion(v string) {
	version = v
}

// NewApp constructs an new App instance from the provided config.
func NewApp(config config.Config) (*App, error) {

	result := &App{config: config}
	result.horizonVersion = version
	result.CoreInfo.NetworkPassphrase = build.DefaultNetwork.Passphrase
	result.ticks = time.NewTicker(1 * time.Second)
	result.init()
	return result, nil
}

// Serve starts the horizon web server, binding it to a socket, setting up
// the shutdown signals.
func (a *App) Serve() {

	a.web.router.Compile()
	http.Handle("/", a.web.router)

	addr := fmt.Sprintf(":%d", a.config.Port)

	srv := &graceful.Server{
		Timeout: 10 * time.Second,

		Server: &http.Server{
			Addr:    addr,
			Handler: http.DefaultServeMux,
		},

		ShutdownInitiated: func() {
			log.Info("received signal, gracefully stopping")
			a.Close()
		},
	}

	http2.ConfigureServer(srv.Server, nil)

	log.Infof("Starting horizon on %s", addr)

	go a.run()

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
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

func (a *App) KeychainQ() *keychain.Q {
	return a.keychainQ
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
	if a.config.StellarCoreURL == "" {
		return
	}

	var err error
	a.CoreInfo, err = coreHelper.GetStellarCoreInfo(a.config.StellarCoreURL)
	if err != nil {
		log.WithField("service", "core-info").WithError(err).Error("could not load stellar-core info")
		return
	}
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
