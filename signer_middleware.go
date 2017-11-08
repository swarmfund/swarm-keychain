package keychain

import (
	"net/http"

	"gitlab.com/tokend/go/signcontrol"
	"gitlab.com/tokend/keychain/render/problem"
	"github.com/zenazn/goji/web"
)

type SignatureValidator struct {
	SkipCheck bool
}

func (v *SignatureValidator) Check(r *http.Request) error {
	if v.SkipCheck {
		return nil
	}

	// dropping header value just in case
	//r.Header.Set(IsSignedHeader, "")

	_, err := signcontrol.CheckSignature(r)
	if err != nil {
		return err
	}

	//r.Header.Set(IsSignedHeader, IsSignedValue)

	return nil
}

func (v *SignatureValidator) Middleware(c *web.C, next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		err := v.Check(r)
		if err != nil {
			switch err {
			case signcontrol.ErrNotSigned:
				// passing not signed requests through w/o setting any headers
				next.ServeHTTP(w, r)
			default:
				problem.Render(r.Context(), w, &problem.BadRequest)
			}
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
