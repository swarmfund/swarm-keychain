package data

import "gitlab.com/tokend/keychain/db2/keychain"

type KeychainQ interface {
	Create(key *keychain.Key) (bool, error)
	Get(address, filename string) (*keychain.Key, error)
}
