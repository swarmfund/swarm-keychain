package keys

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/pkg/errors"
	"gitlab.com/tokend/keychain/db2/keychain"
)

func Generate(address, filename string) (keychain.Key, error) {
	rawKey := make([]byte, 32)
	_, err := rand.Read(rawKey)
	if err != nil {
		return keychain.Key{}, errors.Wrap(err, "faile to generate key")
	}

	encodedKey := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(rawKey)

	return keychain.Key{
		AccountID: address,
		Filename:  filename,
		Key:       encodedKey,
	}, nil
}
