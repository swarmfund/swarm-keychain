package keychain

import (
	"net/http"

	"bullioncoin.githost.io/development/keychain/context/requestid"
	"bullioncoin.githost.io/development/keychain/httpx"
	gctx "github.com/goji/context"
	"github.com/zenazn/goji/web"
	"golang.org/x/net/context"
)

func contextMiddleware(parent context.Context) func(c *web.C, next http.Handler) http.Handler {
	return func(c *web.C, next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := parent
			ctx = requestid.ContextFromC(ctx, c)
			ctx, cancel := httpx.RequestContext(ctx, w, r)

			gctx.Set(c, ctx)
			defer cancel()
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
