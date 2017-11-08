package core

import (
	sq "github.com/lann/squirrel"
)

type Statistics struct {
	Accountid      string `db:"account_id"`
	DailyOutcome   int64  `db:"daily_out"`
	WeeklyOutcome  int64  `db:"weekly_out"`
	MonthlyOutcome int64  `db:"monthly_out"`
	AnnualOutcome  int64  `db:"annual_out"`
}

// SignersByAddress loads all signer rows for `addy`
func (q *Q) StatisticsByAddress(dest interface{}, addy string) error {
	sql := selectStatistics.Where("account_id = ?", addy)
	return q.Get(dest, sql)
}

var selectStatistics = sq.Select(
	"st.account_id",
	"st.daily_out",
	"st.weekly_out",
	"st.monthly_out",
	"st.annual_out",
).From("statistics st").OrderBy("account_id DESC")
