package api

import (
	"github.com/go-chi/chi"
	"gitlab.com/tokend/keychain/log"
	"gitlab.com/tokend/keychain/internal/api/handlers"
	. "gitlab.com/tokend/keychain/internal/api/middlewares"
)

func Router(entry *log.Entry) chi.Router {
	r := chi.NewRouter()

	r.Use(LogCtx(entry))

	r.Route("/users/{address}/keys", func(r chi.Router) {
		r.Use(CheckAllowed("address"))
		r.Post("/", handlers.CreateKey)
		r.Get("/{key}", handlers.GetKey)
	})

	return r
}
