package core

import (
	"gitlab.com/swarmfund/go/xdr"
	sq "github.com/lann/squirrel"
)

type Balance struct {
	BalanceID                string `db:"balance_id"`
	AccountID                string `db:"account_id"`
	ExchangeID               string `db:"exchange_id"`
	Asset                    string `db:"asset"`
	Amount                   int64  `db:"amount"`
	Locked                   int64  `db:"locked"`
	ExchangeName             string `db:"exchange_name"`
	RequireReview            bool   `db:"require_review"`
	StorageFee               int64  `db:"storage_fee"`
	FeesPaid                 int64  `db:"fees_paid"`
	StorageFeeLastCharged    uint64 `db:"storage_fee_last_charged"`
	StorageFeeLastCalculated uint64 `db:"storage_fee_last_calc"`
	FeePercent               int64  `db:"percent"`
	FeePeriod                int64  `db:"period"`
}

func (q *Q) BalancesByAddress(dest interface{}, addy string) error {
	sql := selectBalance.Where("ba.account_id = ?", addy)
	return q.Select(dest, sql)
}

func (q *Q) BalanceByID(dest interface{}, bid string) error {
	sql := selectBalance.Where("ba.balance_id = ?", bid)
	return q.Get(dest, sql)
}

var selectBalance = sq.Select(
	"ba.balance_id",
	"ba.account_id",
	"ba.exchange_id",
	"ba.asset",
	"ba.amount",
	"ba.locked",
	"ba.storage_fee",
	"ba.fees_paid",
	"ba.storage_fee_last_charged",
	"ba.storage_fee_last_calc",
	"ex.name as exchange_name",
	"ex.require_review",
	"fee.percent",
	"fee.period",
).From("balance ba").
	Join("exchanges ex ON ex.account_id=ba.exchange_id").
	Join("fee_state fee ON fee.asset=ba.asset").
	Where("fee.fee_type=?", xdr.FeeTypeStorageFee).
	OrderBy("ba.balance_id")

var selectCoinsInCirculationAmounts = sq.Select(
	"b.asset as asset, sum(b.amount + b.locked) as amount").
	From("balance b").
	GroupBy("b.asset")

func (q *Q) CoinsInCirculation(dest interface{}, masterAccountID string) error {
	sql := selectCoinsInCirculationAmounts.Where("account_id != ?", masterAccountID)

	return q.Select(dest, sql)
}

func (q *Q) CoinsInCirculationForAsset(dest interface{}, masterAccountID, asset string) error {
	sql := selectCoinsInCirculationAmounts.Where("account_id = ?", masterAccountID).
		Where("asset = ?", asset)

	return q.Get(dest, sql)
}

var selectBalanceAmounts = sq.Select(
	"b.asset as asset, b.amount as amount").
	From("balance b")

func (q *Q) AvailableEmissions(dest interface{}, masterAccountID string) error {
	sql := selectBalanceAmounts.Where("account_id = ?", masterAccountID)

	return q.Select(dest, sql)
}
