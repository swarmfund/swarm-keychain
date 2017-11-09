package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

func GetKey(w http.ResponseWriter, r *http.Request) {
	address := chi.URLParam(r, "address")
	key := chi.URLParam(r, "key")
	w.Write([]byte(fmt.Sprintf("%s -> %s", address, key)))
}
