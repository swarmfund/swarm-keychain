package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	coreHelper "gitlab.com/swarmfund/go/core"
	"gitlab.com/swarmfund/go/doorman"
	"gitlab.com/swarmfund/keychain/internal/api/data"
	"gitlab.com/swarmfund/keychain/internal/api/handlers"
	. "gitlab.com/swarmfund/keychain/internal/api/middlewares"
	"gitlab.com/swarmfund/keychain/log"
)

func Router(entry *log.Entry, doorman doorman.Doorman, keychainQ data.KeychainQ, coreInfoGetter func() *coreHelper.Info) chi.Router {
	r := chi.NewRouter()

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		ape.RenderErr(w, problems.NotFound())
	})

	r.Use(
		Recover,
	)

	r.Route("/users/{address}/keys", func(r chi.Router) {
		r.Use(
			Ctx(
				handlers.CtxLog(entry),
				handlers.CtxKeychainQ(keychainQ),
				handlers.CtxDoorman(doorman),
				handlers.CtxCoreInfo(coreInfoGetter()),
			),
		)

		r.Post("/", handlers.CreateKey)
		r.Get("/{filename}", handlers.GetKey)
	})

	return r
}
