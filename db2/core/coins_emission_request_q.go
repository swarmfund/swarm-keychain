package core

import (
	"gitlab.com/swarmfund/keychain/db2"
	sq "github.com/lann/squirrel"
)

// CoinsEmissionRequestQ is a helper struct to aid in configuring queries that loads
// slices of CoinsEmissionRequest structs.
type CoinsEmissionRequestQ struct {
	Err    error
	parent *Q
	sql    sq.SelectBuilder
}

// CoinsEmissionRequests provides a helper to filter the operations table with pre-defined
// filters.  See `CoinsEmissionRequestQ` for the available filters.
func (q *Q) CoinsEmissionRequests() *CoinsEmissionRequestQ {
	return &CoinsEmissionRequestQ{
		parent: q,
		sql:    selectCoinsEmissionRequest,
	}
}

func (q *Q) EmissionRequestByExchangeAndRef(exchange, reference string) (exists *bool, err error) {
	sql := sq.Expr("SELECT EXISTS(SELECT * FROM coins_emission_request cemr WHERE "+
		" cemr.issuer = ? AND cemr.reference = ?)", exchange, reference)
	err = q.Get(&exists, sql)
	if err != nil {
		return nil, err
	}
	return exists, nil
}

// OnlyApproved filters the query being built to only include requests that
// are approved.
func (q *CoinsEmissionRequestQ) OnlyApproved() *CoinsEmissionRequestQ {
	q.sql = q.sql.Where("cemr.is_approved")
	return q
}

// ForAccount filters the operations collection to a specific account
func (q *CoinsEmissionRequestQ) ForAccount(aid string) *CoinsEmissionRequestQ {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("cemr.issuer = ?", aid)

	return q
}

func (q *CoinsEmissionRequestQ) ByID(id uint64) (*CoinsEmissionRequest, error) {
	if q.Err != nil {
		return nil, q.Err
	}

	q.sql = q.sql.Where("cemr.request_id = ?", id)
	var result CoinsEmissionRequest
	err := q.parent.Get(&result, q.sql)
	if q.parent.NoRows(err) {
		return nil, nil
	}

	return &result, err
}

// Page specifies the paging constraints for the query being built by `q`.
func (q *CoinsEmissionRequestQ) Page(page db2.PageQuery) *CoinsEmissionRequestQ {
	if q.Err != nil {
		return q
	}

	q.sql, q.Err = page.ApplyTo(q.sql, "cemr.request_id")
	return q
}

// Select loads the results of the query specified by `q` into `dest`.
func (q *CoinsEmissionRequestQ) Select(dest interface{}) error {
	if q.Err != nil {
		return q.Err
	}

	q.Err = q.parent.Select(dest, q.sql)
	return q.Err
}

var selectCoinsEmissionRequest = sq.Select("cemr.issuer, cemr.request_id, cemr.amount, cemr.is_approved").From("coins_emission_request cemr")
