package data

import "gitlab.com/swarmfund/keychain/db2/keychain"

type KeychainQ interface {
	Create(key *keychain.Key) (bool, error)
	Get(address, filename string) (*keychain.Key, error)
}

type CoreInfoI interface {
	GetMasterAccountID() string
}
