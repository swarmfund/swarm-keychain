package handlers

import (
	"net/http"

	"context"

	"gitlab.com/swarmfund/go/core"
	"gitlab.com/swarmfund/go/doorman"
	"gitlab.com/swarmfund/keychain/internal/api/data"
	"gitlab.com/swarmfund/keychain/log"
)

type CtxKey int

const (
	CtxKeyLog CtxKey = iota
	CtxKeyKeychainQ
	CtxKeyDoorman
	CtxKeyCoreInfo
)

func CtxLog(log *log.Entry) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, CtxKeyLog, log)
	}
}

func Log(r *http.Request) *log.Entry {
	return r.Context().
		Value(CtxKeyLog).(*log.Entry).
		WithField("path", r.URL.Path)

}

func CtxKeychainQ(q data.KeychainQ) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, CtxKeyKeychainQ, q)
	}
}

func KeychainQ(r *http.Request) data.KeychainQ {
	return r.Context().Value(CtxKeyKeychainQ).(data.KeychainQ)
}

func CtxDoorman(d doorman.Doorman) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, CtxKeyDoorman, d)
	}
}

func Doorman(r *http.Request, constraints ...doorman.SignerConstraint) error {
	d := r.Context().Value(CtxKeyDoorman).(doorman.Doorman)
	return d.Check(r, constraints...)
}

func CtxCoreInfo(ci *core.Info) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, CtxKeyCoreInfo, *ci)
	}
}
func CoreInfo(r *http.Request) core.Info {
	return r.Context().Value(CtxKeyCoreInfo).(core.Info)
}
