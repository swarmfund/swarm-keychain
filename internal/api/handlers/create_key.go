package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/keychain/internal/keys"
)

type CreateKeyRequest struct {
	Address  string
	Filename string
}

func NewCreateKeyRequest(r *http.Request) *CreateKeyRequest {
	return &CreateKeyRequest{
		Address:  chi.URLParam(r, "address"),
		Filename: chi.URLParam(r, "filename"),
	}
}

func (r *CreateKeyRequest) Validate() error {
	return nil
}

func CreateKey(w http.ResponseWriter, r *http.Request) {
	request := NewCreateKeyRequest(r)
	if err := request.Validate(); err != nil {
		// TODO bad request
		panic(400)
	}

	// generate encryption key
	key, err := keys.Generate(request.Address, request.Filename)
	if err != nil {
		Log(r).WithError(err).Error("failed to generate key")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	// save key
	ok, err := KeychainQ(r).Create(&key)
	if err != nil {
		Log(r).WithError(err).Error("failed to save key")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if !ok {
		ape.RenderErr(w, problems.Conflict())
		return
	}

	// render response
	ape.Render(w, key)
}
