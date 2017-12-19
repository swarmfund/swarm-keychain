package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
	. "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/swarmfund/go/doorman"
	"gitlab.com/swarmfund/go/signcontrol"
)

type GetKeyRequest struct {
	Address  string
	Filename string
}

func NewGetKeyRequest(r *http.Request) (GetKeyRequest, error) {
	request := GetKeyRequest{
		Address:  chi.URLParam(r, "address"),
		Filename: chi.URLParam(r, "filename"),
	}
	return request, request.Validate()
}

func (r *GetKeyRequest) Validate() error {
	return ValidateStruct(r,
		Field(&r.Address, Required),
		Field(&r.Filename, Required),
	)
}

func GetKey(w http.ResponseWriter, r *http.Request) {
	request, err := NewGetKeyRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	// get key from database
	key, err := KeychainQ(r).Get(request.Address, request.Filename)
	if err != nil {
		Log(r).WithError(err).Error("failed to get key")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if key == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	// check allowed
	if err := Doorman(r,
		doorman.SignerOf(request.Address),
		doorman.SignerOf(CoreInfo(r).GetMasterAccountID()),
	); err != nil {
		RenderDoormanErr(w, err)
		return
	}

	// render response
	// TODO render JWK resource
	ape.Render(w, key)
}

func RenderDoormanErr(w http.ResponseWriter, err error) {
	switch err {
	case signcontrol.ErrNotAllowed, signcontrol.ErrNotSigned:
		ape.RenderErr(w, problems.NotAllowed())
	case nil:
		panic("expected not nil error")
	default:
		panic(errors.Wrap(err, "unexpected doorman error"))
	}
}
