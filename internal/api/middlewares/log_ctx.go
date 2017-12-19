package middlewares

import (
	"context"
	"net/http"

	"gitlab.com/swarmfund/keychain/internal/api/handlers"
	"gitlab.com/swarmfund/keychain/log"
)

func LogCtx(log *log.Entry) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(
				r.Context(), handlers.LogCtxKey,
				log.WithField("path", r.URL.Path))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
