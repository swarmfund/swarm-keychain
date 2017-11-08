package core

import (
	"gitlab.com/tokend/keychain/db2"
	sq "github.com/lann/squirrel"
)

type OfferQ struct {
	Err    error
	parent *Q
	sql    sq.SelectBuilder
}

func (q *Q) Offers() *OfferQ {
	return &OfferQ{
		parent: q,
		sql:    selectOffer,
	}
}

func (q *OfferQ) ForAccount(accountID string) *OfferQ {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("owner_id = ?", accountID)
	return q
}

func (q *OfferQ) ForAssets(base, quote string) *OfferQ {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("base_asset_code = ? AND quote_asset_code = ?", base, quote)
	return q
}

func (q *OfferQ) IsBuy(isBuy bool) *OfferQ {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("is_buy = ?", isBuy)
	return q
}

// Page specifies the paging constraints for the query being built by `q`.
func (q *OfferQ) Page(page db2.PageQuery) *OfferQ {
	if q.Err != nil {
		return q
	}

	q.sql, q.Err = page.ApplyTo(q.sql, "o.offer_id")
	return q
}

// Select loads the results of the query specified by `q` into `dest`.
func (q *OfferQ) Select(dest interface{}) error {
	if q.Err != nil {
		return q.Err
	}

	q.Err = q.parent.Select(dest, q.sql)
	return q.Err
}

var selectOffer = sq.Select(
	"o.owner_id",
	"o.offer_id",
	"o.base_asset_code",
	"o.quote_asset_code",
	"o.is_buy",
	"o.base_amount",
	"o.quote_amount",
	"o.price",
	"o.base_balance_id",
	"o.quote_balance_id",
	"o.created_at",
).From("offer o")
