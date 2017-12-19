package handlers

import (
	"net/http"

	"gitlab.com/swarmfund/keychain/internal/api/data"
	"gitlab.com/swarmfund/keychain/log"
)

type CtxKey int

const (
	LogCtxKey CtxKey = iota
	KeychainQCtxKey
)

func Log(r *http.Request) *log.Entry {
	return r.Context().Value(LogCtxKey).(*log.Entry)
}

func KeychainQ(r *http.Request) data.KeychainQ {
	return r.Context().Value(KeychainQCtxKey).(data.KeychainQ)
}
