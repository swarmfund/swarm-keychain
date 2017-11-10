package keychain

import (
	"strings"

	"database/sql"

	"github.com/lann/squirrel"
	"gitlab.com/tokend/keychain/db2"
)

type Key struct {
	ID        int64  `db:"id" jsonapi:"primary,key"`
	AccountID string `db:"account_id" jsonapi:"attr,account_id"`
	Filename  string `db:"filename" jsonapi:"attr,filename"`
	Key       string `db:"key" jsonapi:"attr,key"`
}

type KeyQ struct {
	*db2.Repo
}

func (q *KeyQ) Create(key *Key) (bool, error) {
	stmt := squirrel.Insert("keys").SetMap(map[string]interface{}{
		"account_id": key.AccountID,
		"filename":   key.Filename,
		"key":        key.Key,
	}).Suffix("returning id")

	err := q.Repo.Get(&(key.ID), stmt)
	if err != nil {
		if strings.Contains(err.Error(), "unique_account_id_filename") {
			return false, nil
		}
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (q *KeyQ) Get(accountID, filename string) (*Key, error) {
	var result Key
	stmt := squirrel.Select("*").From("keys").
		Where("account_id = ? and filename = ?", accountID, filename)

	err := q.Repo.Get(&result, stmt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &result, err
}
