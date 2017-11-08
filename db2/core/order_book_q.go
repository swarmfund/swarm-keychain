package core

import (
	sq "github.com/lann/squirrel"
)

type OrderBookQ struct {
	Err    error
	parent *Q
	sql    sq.SelectBuilder
}

func (q *Q) OrderBook() *OrderBookQ {
	return &OrderBookQ{
		parent: q,
		sql: sq.Select(
			"SUM(o.base_amount) as base_amount",
			"SUM(o.quote_amount) as quote_amount",
			"o.price as price",
		).From("offer o").GroupBy("o.price"),
	}
}

func (q *OrderBookQ) ForAssets(base, quote string) *OrderBookQ {
	if q.Err != nil {
		return q
	}

	q.sql = q.sql.Where("base_asset_code = ? AND quote_asset_code = ?", base, quote)
	return q
}

func (q *OrderBookQ) Direction(isBuy bool) *OrderBookQ {
	if q.Err != nil {
		return q
	}

	orderDirection := "ASC"
	if isBuy {
		orderDirection = "DESC"
	}

	q.sql = q.sql.Where("is_buy = ?", isBuy).OrderBy("price " + orderDirection)
	return q
}

// Select loads the results of the query specified by `q` into `dest`.
func (q *OrderBookQ) Select(dest interface{}) error {
	if q.Err != nil {
		return q.Err
	}

	q.Err = q.parent.Select(dest, q.sql)
	return q.Err
}
