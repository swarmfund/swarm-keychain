package core

import (
	sq "github.com/lann/squirrel"
)

type TrustQ struct {
	Err    error
	parent *Q
	sql    sq.SelectBuilder
}

func (q *Q) Trusts() *TrustQ {
	return &TrustQ{
		parent: q,
		sql:    selectTrust,
	}
}

func (q *TrustQ) ForBalance(bid string) *TrustQ {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("t.balance_to_use = ?", bid)

	return q
}

func (q *TrustQ) ForAccount(aid string) *TrustQ {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("t.allowed_account = ?", aid)

	return q
}

// Select loads the results of the query specified by `q` into `dest`.
func (q *TrustQ) Select(dest interface{}) error {
	if q.Err != nil {
		return q.Err
	}

	q.Err = q.parent.Select(dest, q.sql)
	return q.Err
}

var selectTrust = sq.Select("t.allowed_account, t.balance_to_use").
	From("trusts t")
