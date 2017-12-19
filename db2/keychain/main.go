package keychain

import "gitlab.com/swarmfund/keychain/db2"

type Q struct {
	*db2.Repo
}

func (q *Q) GetRepo() *db2.Repo {
	return q.Repo
}

func (q *Q) Key() *KeyQ {
	return &KeyQ{}
}
