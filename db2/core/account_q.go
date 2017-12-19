package core

import (
	"gitlab.com/swarmfund/go/resources"
	"gitlab.com/swarmfund/keychain/db2"
)

type AccountQ struct {
	repo *db2.Repo
}

func NewAccountQ(repo *db2.Repo) *AccountQ {
	return &AccountQ{
		repo: repo,
	}
}

func (q *AccountQ) Signers(address string) ([]resources.Signer, error) {
	return nil, nil
}
