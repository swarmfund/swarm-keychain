package middlewares

import (
	"net/http"
	"runtime/debug"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/keychain/errors"
	"gitlab.com/tokend/keychain/internal/api/handlers"
)

func Recover(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				handlers.Log(r).
					WithField("stack", string(debug.Stack())).
					WithError(errors.FromPanic(rvr)).Error("handler panicked")
				ape.RenderErr(w, problems.InternalError())
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
