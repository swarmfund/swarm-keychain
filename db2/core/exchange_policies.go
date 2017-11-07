package core

import (
	sq "github.com/lann/squirrel"
)

type ExchangePolicies struct {
	AccountID string `db:"account_id"`
	Asset     string `db:"asset"`
	Policies  int32  `db:"policies"`
}

var selectExchangePolicies = sq.Select("ep.account_id, ep.asset, ep.policies").
	From("exchange_policies ep")

func (q *Q) PoliciesByExchangeID(dest interface{}, addy string) error {
	sql := selectExchangePolicies.Where("account_id = ?", addy)
	return q.Select(dest, sql)
}
