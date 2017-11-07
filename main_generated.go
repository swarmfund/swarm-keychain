package keychain

import (
	"net/http"

	"github.com/zenazn/goji/web"
)

// ServeHTTPC is a method for web.Handler
func (action CreateKeyAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "CreateKeyAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action GetKeyAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "GetKeyAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action NotFoundAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "NotFoundAction")
	ap.Execute(&action)
}

// ServeHTTPC is a method for web.Handler
func (action NotImplementedAction) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(c, w, r)
	action.Log = action.Log.WithField("action", "NotImplementedAction")
	ap.Execute(&action)
}
