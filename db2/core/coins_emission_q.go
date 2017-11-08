package core

import (
	"database/sql"
	sq "github.com/lann/squirrel"
)

var selectCoinsEmission = sq.Select("cem.serial_number, cem.amount, cem.lastmodified").From("coins_emission cem")

// CoinsEmissionQ is a helper struct to aid in configuring queries that loads
// slices of CoinsEmission structs.
type CoinsEmissionQ struct {
	Err    error
	parent *Q
	sql    sq.SelectBuilder
}

type CoinsEmissionQI interface {
	BySerialNumber(serialNumber string) (*CoinsEmission, error)
}

// CoinsEmissions provides a helper to filter the operations table with pre-defined
// filters.  See `CoinsEmissionQ` for the available filters.
func (q *Q) CoinsEmissions() *CoinsEmissionQ {
	return &CoinsEmissionQ{
		parent: q,
		sql:    selectCoinsEmission,
	}
}

func (q *CoinsEmissionQ) BySerialNumber(serialNumber string) (*CoinsEmission, error) {
	if q.Err != nil {
		return nil, q.Err
	}

	q.sql = q.sql.Where("cem.serial_number = ?", serialNumber)
	result := new(CoinsEmission)

	q.Err = q.parent.Get(result, q.sql)
	if q.Err == sql.ErrNoRows {
		return nil, nil
	}

	return result, q.Err
}
