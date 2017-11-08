package keychain

import (
	"net/http"
	"time"

	"golang.org/x/net/context"

	"gitlab.com/tokend/keychain/log"
	gctx "github.com/goji/context"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
	"github.com/zenazn/goji/web/mutil"
)

// LoggerMiddleware is the middleware that logs http requests and resposnes
// to the logging subsytem of horizon.
func LoggerMiddleware(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := gctx.FromC(*c)
		mw := mutil.WrapWriter(w)

		logger := log.WithField("req", middleware.GetReqID(*c))

		ctx = log.Set(ctx, logger)
		gctx.Set(c, ctx)

		logStartOfRequest(ctx, r)

		then := time.Now()
		h.ServeHTTP(mw, r)
		duration := time.Now().Sub(then)

		logEndOfRequest(ctx, duration, mw)
	}

	return http.HandlerFunc(fn)
}

func logStartOfRequest(ctx context.Context, r *http.Request) {
	log.Ctx(ctx).WithFields(log.F{
		"path":   r.URL.String(),
		"method": r.Method,
		"ip":     r.RemoteAddr,
		"host":   r.Host,
	}).Info("Starting request")
}

func logEndOfRequest(ctx context.Context, duration time.Duration, mw mutil.WriterProxy) {
	log.Ctx(ctx).WithFields(log.F{
		"status":   mw.Status(),
		"bytes":    mw.BytesWritten(),
		"duration": duration,
	}).Info("Finished request")
}
