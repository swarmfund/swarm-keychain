package middlewares

import (
	"net/http"

	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/go/doorman"
	"gitlab.com/tokend/go/signcontrol"
)

func CheckAllowed(key string, checks ...func(string) doorman.SignerConstraint) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			value := chi.URLParam(r, key)
			constraints := make([]doorman.SignerConstraint, 0, len(checks))
			for _, check := range checks {
				constraints = append(constraints, check(value))
			}
			switch err := doorman.Check(r, constraints...); err {
			case signcontrol.ErrNotAllowed, signcontrol.ErrNotSigned:
				ape.RenderErr(w, problems.NotAllowed())
			case nil:
				next.ServeHTTP(w, r)
			default:
				ape.RenderErr(w, problems.InternalError())
			}
		}
		return http.HandlerFunc(fn)
	}
}
