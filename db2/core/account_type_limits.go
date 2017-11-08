package core

import (
	"database/sql"

	sq "github.com/lann/squirrel"
)

var selectAccountTypeLimits = sq.Select("atl.account_type, atl.daily_out, atl.weekly_out," +
	"atl.monthly_out, atl.annual_out").
	From("account_type_limits atl")

type AccountTypeLimits struct {
	AccountType int `db:"account_type"`
	Limits
}

func (q *Q) DefaultLimits(dest interface{}) error {
	sql := selectAccountTypeLimits
	return q.Get(dest, sql)
}

func (q *Q) LimitsByAccountType(accountType int) (*AccountTypeLimits, error) {
	var result AccountTypeLimits
	query := selectAccountTypeLimits.Limit(1).Where("atl.account_type = ?", accountType)
	err := q.Get(&result, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &result, nil
}

// FeeEntryQ is a helper struct to aid in configuring queries that loads
// slices of FeeEntry structs.
type AccountTypeLimitsQ struct {
	Err    error
	parent *Q
	sql    sq.SelectBuilder
}

type AccountTypeLimitsQI interface {
	Select(dest interface{}) error
}

func (q *Q) AccountTypeLimits() AccountTypeLimitsQI {
	return &AccountTypeLimitsQ{
		parent: q,
		sql:    selectAccountTypeLimits,
	}
}

// Select loads the results of the query specified by `q` into `dest`.
func (q *AccountTypeLimitsQ) Select(dest interface{}) error {
	if q.Err != nil {
		return q.Err
	}

	q.Err = q.parent.Select(dest, q.sql)
	return q.Err
}
