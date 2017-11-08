package keychain

import (
	"gitlab.com/tokend/keychain/db2/keychain"
	"gitlab.com/tokend/keychain/render/hal"
	"gitlab.com/tokend/keychain/render/problem"
	"gitlab.com/tokend/keychain/resource"
)

type GetKeyAction struct {
	Action

	AccountID string
	Filename  string

	Record   *keychain.Key
	Resource resource.JWK
}

func (action *GetKeyAction) JSON() {
	action.Do(
		action.loadParams,
		action.checkAllowed,
		action.loadRecord,
		action.loadResource,
		func() {
			hal.Render(action.W, action.Resource)
		},
	)
}

func (action *GetKeyAction) loadParams() {
	action.AccountID = action.GetNonEmptyString("id")
	action.Filename = action.GetNonEmptyString("fn")
}

func (action *GetKeyAction) checkAllowed() {
	action.IsAllowed(action.AccountID)
}

func (action *GetKeyAction) loadRecord() {
	key, err := action.KeychainQ().Key().Get(action.AccountID, action.Filename)
	if err != nil {
		action.Log.WithError(err).Error("failed to get key")
		action.Err = &problem.ServerError
		return
	}

	if key == nil {
		action.Err = &problem.NotFound
		return
	}

	action.Record = key
}

func (action *GetKeyAction) loadResource() {
	action.Resource.Populate(action.Record)
}
