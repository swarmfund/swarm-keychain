package core

import (
	"database/sql"

	sq "github.com/lann/squirrel"
)

type Exchange struct {
	AccountID     string `db:"account_id"`
	Name          string `db:"name"`
	RequireReview bool   `db:"require_review"`
}

func (q *Q) ExchangeName(address string) (*string, error) {
	var result string
	stmt := sq.Select("name").
		From("exchanges ex").
		Limit(1).
		Where("account_id = ?", address)
	err := q.Get(&result, stmt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &result, err
}
