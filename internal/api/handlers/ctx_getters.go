package handlers

import (
	"net/http"

	"gitlab.com/tokend/keychain/internal/api/data"
	"gitlab.com/tokend/keychain/log"
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
