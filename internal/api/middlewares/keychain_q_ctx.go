package middlewares

import (
	"context"
	"net/http"

	"gitlab.com/swarmfund/keychain/internal/api/data"
	"gitlab.com/swarmfund/keychain/internal/api/handlers"
)

func KeychainQCtx(q data.KeychainQ) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), handlers.KeychainQCtxKey, q)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
