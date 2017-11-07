package actions

import (
	"net/http"

	"encoding/base64"

	gctx "github.com/goji/context"

	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"

	"bullioncoin.githost.io/development/keychain/render"
	"bullioncoin.githost.io/development/keychain/render/problem"
	"bullioncoin.githost.io/development/keychain/render/sse"
	"github.com/zenazn/goji/web"
	"golang.org/x/net/context"
)

// Base is a helper struct you can use as part of a custom action via
// composition.
//
// TODO: example usage
type Base struct {
	Ctx     context.Context
	GojiCtx web.C
	W       http.ResponseWriter
	R       *http.Request
	Err     error

	Signer string

	isSetup bool
	IsAdmin bool

	isJson   bool
	jsonBody map[string]string

	SkipCheck bool //flag for developing without signing, false - for checking signatures
}

// Prepare established the common attributes that get used in nearly every
// action.  "Child" actions may override this method to extend action, but it
// is advised you also call this implementation to maintain behavior.
func (base *Base) Prepare(c web.C, w http.ResponseWriter, r *http.Request) {
	base.Ctx = gctx.FromC(c)
	base.GojiCtx = c
	base.W = w
	base.R = r
}

func (base *Base) JsonValue(name string) string {
	if !base.isJson {
		return ""
	}

	if base.jsonBody == nil {
		decoder := json.NewDecoder(base.R.Body)
		err := decoder.Decode(&base.jsonBody)
		if err != nil {
			base.jsonBody = map[string]string{}
		}
	}

	value, ok := base.jsonBody[name]
	if !ok {
		return ""
	}
	return value
}

// Execute trigger content negottion and the actual execution of one of the
// action's handlers.
func (base *Base) Execute(action interface{}) {
	contentType := render.Negotiate(base.Ctx, base.R)

	switch contentType {
	case render.MimeHal, render.MimeJSON:
		action, ok := action.(JSON)

		if !ok {
			goto NotAcceptable
		}

		action.JSON()

		if base.Err != nil {
			problem.Render(base.Ctx, base.W, base.Err)
			return
		}

	case render.MimeEventStream:
		action, ok := action.(SSE)
		if !ok {
			goto NotAcceptable
		}

		stream := sse.NewStream(base.Ctx, base.W, base.R)

		for {
			action.SSE(stream)

			if base.Err != nil {
				// in the case that we haven't yet sent an event, is also means we
				// havent sent the preamble, meaning we should simply return the normal
				// error.
				if stream.SentCount() == 0 {
					problem.Render(base.Ctx, base.W, base.Err)
					return
				}

				stream.Err(base.Err)
			}

			select {
			case <-base.Ctx.Done():
				return
			case <-sse.Pumped():
				//no-op, continue onto the next iteration
			}
		}
	case render.MimeRaw:
		action, ok := action.(Raw)

		if !ok {
			goto NotAcceptable
		}

		action.Raw()

		if base.Err != nil {
			problem.Render(base.Ctx, base.W, base.Err)
			return
		}
	default:
		goto NotAcceptable
	}
	return

NotAcceptable:
	problem.Render(base.Ctx, base.W, problem.NotAcceptable)
	return
}

// Do executes the provided func iff there is no current error for the action.
// Provides a nicer way to invoke a set of steps that each may set `action.Err`
// during execution
func (base *Base) Do(fns ...func()) {
	for _, fn := range fns {
		if base.Err != nil {
			return
		}

		fn()
	}
}

// Setup runs the provided funcs if and only if no call to Setup() has been
// made previously on this action.
func (base *Base) Setup(fns ...func()) {
	if base.isSetup {
		return
	}
	base.Do(fns...)
	base.isSetup = true
}

func (base *Base) GetByteArray(name string, length int) string {
	rawValue := base.GetNonEmptyString(name)
	if base.Err != nil {
		return ""
	}

	value, err := base64.StdEncoding.DecodeString(rawValue)

	if err != nil {
		base.Err = err
		return ""
	}

	if len(value) != length {
		base.SetInvalidField(name, errors.New(" is not "+string(length)+"byte length"))
		return ""
	}
	return base64.StdEncoding.EncodeToString(value)
}

func (base *Base) ValidateHash(orig, hash string) bool {
	rawOrig := []byte(orig)
	hasher := sha1.New()
	hasher.Write(rawOrig)
	hashed := hex.EncodeToString(hasher.Sum(nil))
	return hashed == hash
}
