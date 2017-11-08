package core

import (
	sq "github.com/lann/squirrel"
)

type Limits struct {
	DailyOut   int64 `db:"daily_out"`
	WeeklyOut  int64 `db:"weekly_out"`
	MonthlyOut int64 `db:"monthly_out"`
	AnnualOut  int64 `db:"annual_out"`
}

type AccountLimits struct {
	Accountid string `db:"accountid"`
	Limits
}

// SignersByAddress loads all signer rows for `addy`
func (q *Q) LimitsByAddress(dest interface{}, addy string) error {
	sql := selectLimit.Where("accountid = ?", addy)
	return q.Get(dest, sql)
}

var selectLimit = sq.Select(
	"li.daily_out",
	"li.weekly_out",
	"li.monthly_out",
	"li.annual_out",
).From("account_limits li").OrderBy("accountid DESC")
