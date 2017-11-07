package keychain

import (
	"crypto/rand"

	"encoding/base64"

	"strings"

	"bullioncoin.githost.io/development/keychain/db2/keychain"
	"bullioncoin.githost.io/development/keychain/render/hal"
	"bullioncoin.githost.io/development/keychain/render/problem"
	"bullioncoin.githost.io/development/keychain/resource"
)

type CreateKeyAction struct {
	Action

	AccountID string
	Filename  string

	Resource resource.JWK
}

func (action *CreateKeyAction) JSON() {
	action.Do(
		action.ValidateBodyType,
		action.loadParams,
		action.checkAllowed,
		action.performRequest,
		func() {
			hal.Render(action.W, action.Resource)
		},
	)
}

func (action *CreateKeyAction) loadParams() {
	action.AccountID = action.GetNonEmptyString("id")
	action.Filename = action.GetNonEmptyString("fn")
}

func (action *CreateKeyAction) checkAllowed() {
	action.IsAllowed(action.AccountID)
}

func (action *CreateKeyAction) performRequest() {
	// generate key
	rawKey := make([]byte, 32)
	_, err := rand.Read(rawKey)
	if err != nil {
		action.Log.WithError(err).Error("failed to generate key")
		action.Err = &problem.ServerError
		return
	}

	encodedKey := strings.Replace(strings.Replace(base64.StdEncoding.WithPadding(base64.NoPadding).EncodeToString(rawKey), "+", "-", -1), "/", "_", -1)

	key := keychain.Key{
		AccountID: action.AccountID,
		Filename:  action.Filename,
		Key:       encodedKey,
	}

	// save key

	ok, err := action.KeychainQ().Key().Create(&key)
	if err != nil {
		action.Log.WithError(err).Error("failed to save key")
		action.Err = &problem.ServerError
		return
	}

	if !ok {
		action.Err = &problem.Conflict
		return
	}

	// populate resource

	action.Resource.Populate(&key)
}
