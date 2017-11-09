package handlers

import (
	"net/http"

	"gitlab.com/tokend/keychain/log"
)

const (
	LogCtxKey = "log"
)

func Log(r *http.Request) *log.Entry {
	return r.Context().Value(LogCtxKey).(*log.Entry)
}
