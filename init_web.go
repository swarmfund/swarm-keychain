package keychain

import (
	"database/sql"

	"github.com/rcrowley/go-metrics"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
	"gitlab.com/distributed_lab/tokend/keychain/render/problem"
)

// Web contains the http server related fields for horizon: the router,
// rate limiter, etc.
type Web struct {
	router *web.Mux

	requestTimer metrics.Timer
	failureMeter metrics.Meter
	successMeter metrics.Meter
}

// initWeb installed a new Web instance onto the provided app object.
func initWeb(app *App) {
	mux := web.New()
	app.web = &Web{
		router: mux,
	}

	// register problems
	problem.RegisterError(sql.ErrNoRows, problem.NotFound)
}

// initWebMiddleware installs the middleware stack used for horizon onto the
// provided app.
func initWebMiddleware(app *App) {
	r := app.web.router

	r.Use(stripTrailingSlashMiddleware())
	r.Use(middleware.EnvInit)
	r.Use(app.Middleware)
	r.Use(middleware.RequestID)
	r.Use(contextMiddleware(app.ctx))
	r.Use(LoggerMiddleware)
	r.Use(RecoverMiddleware)
	r.Use(middleware.AutomaticOptions)

	signatureValidator := &SignatureValidator{app.config.SkipCheck}

	r.Use(signatureValidator.Middleware)
}

// initWebActions installs the routing configuration of horizon onto the
// provided app.  All route registration should be implemented here.
func initWebActions(app *App) {
	r := app.web.router

	r.Post("/users/:id/keys/:fn", &CreateKeyAction{})
	r.Get("/users/:id/keys/:fn", &GetKeyAction{})
}

func init() {
	appInit.Add(
		"web.init",
		initWeb,
		"app-context", "stellarCoreInfo",
	)

	appInit.Add(
		"web.middleware",
		initWebMiddleware,
		"web.init",
	)

	appInit.Add(
		"web.actions",
		initWebActions,
		"web.init",
	)
}
