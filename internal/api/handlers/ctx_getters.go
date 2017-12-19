package handlers

import (
	"context"
	"net/http"

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

func CtxCoreInfo(conn data.CoreInfoI) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, CtxKeyCoreInfo, conn)
	}
}

func CoreInfo(r *http.Request) data.CoreInfoI {
	conn := r.Context().Value(CtxKeyCoreInfo)
	c := conn.(data.CoreInfoI)
	return c
}
