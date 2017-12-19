package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/swarmfund/go/doorman"
	"gitlab.com/swarmfund/keychain/internal/api/data"
	"gitlab.com/swarmfund/keychain/internal/api/handlers"
	. "gitlab.com/swarmfund/keychain/internal/api/middlewares"
	"gitlab.com/swarmfund/keychain/log"
)

func Router(entry *log.Entry, doorman *doorman.Doorman, keychainQ data.KeychainQ) chi.Router {
	r := chi.NewRouter()

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		ape.RenderErr(w, problems.NotFound())
	})

	r.Use(
		LogCtx(entry),
		Recover,
	)

	r.Route("/users/{address}/keys", func(r chi.Router) {
		r.Use(
			CheckAllowed("address", doorman.SignerOf, doorman.MasterSigner),
			KeychainQCtx(keychainQ),
		)

		r.Post("/", handlers.CreateKey)
		r.Get("/{filename}", handlers.GetKey)
	})

	return r
}
