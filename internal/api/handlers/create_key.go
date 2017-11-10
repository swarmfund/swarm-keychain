package handlers

import (
	"net/http"

	"encoding/json"

	"github.com/go-chi/chi"
	. "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/keychain/internal/keys"
)

type CreateKeyRequest struct {
	Address  string `json:"-"`
	Filename string `json:"filename"`
}

func NewCreateKeyRequest(r *http.Request) (CreateKeyRequest, error) {
	request := CreateKeyRequest{
		Address: chi.URLParam(r, "address"),
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}
	return request, request.Validate()
}

func (r *CreateKeyRequest) Validate() error {
	return ValidateStruct(r,
		Field(&r.Address, Required),
		Field(&r.Filename, Required))
}

func CreateKey(w http.ResponseWriter, r *http.Request) {
	request, err := NewCreateKeyRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
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
	ape.Render(w, &key)
}
