package core

import (
	"database/sql"

	sq "github.com/lann/squirrel"
)

var selectFees = sq.Select("f.fee_type", "f.asset", "f.fixed", "f.percent", "f.lastmodified, f.period").
	From("fee_state f")

type FeeEntry struct {
	FeeType      int    `db:"fee_type"`
	Asset        string `db:"asset"`
	Fixed        int64  `db:"fixed"`
	Percent      int64  `db:"percent"`
	Period       int64  `db:"period"`
	LastModified int32  `db:"lastmodified"`
}

// Fees loads all row from `fee_state`
func (q *Q) Fees(dest interface{}) error {
	sql := selectFees
	return q.Get(dest, sql)
}

func (q *Q) FeeByTypeAsset(feeType int, asset string) (*FeeEntry, error) {
	var result FeeEntry
	query := selectFees.Limit(1).Where("f.fee_type = ? and f.asset = ?", feeType, asset)
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
type FeeEntryQ struct {
	Err    error
	parent *Q
	sql    sq.SelectBuilder
}

type FeeEntryQI interface {
	Select(dest interface{}) error
}

// FeeEntries provides a helper to filter the operations table with pre-defined
// filters.  See `FeeEntryQ` for the available filters.
func (q *Q) FeeEntries() FeeEntryQI {
	return &FeeEntryQ{
		parent: q,
		sql:    selectFees,
	}
}

// Select loads the results of the query specified by `q` into `dest`.
func (q *FeeEntryQ) Select(dest interface{}) error {
	if q.Err != nil {
		return q.Err
	}

	q.Err = q.parent.Select(dest, q.sql)
	return q.Err
}
