package middlewares

import (
	"net/http"

	"gitlab.com/tokend/keychain/internal/api/handlers"
)

func CheckAllowed(key string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			handlers.Log(r).Warn("signature check not implemented")
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
