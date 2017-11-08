package core

import (
	sq "github.com/lann/squirrel"
)

type Asset struct {
	Code          string `db:"code"`
	CurrentPrice  int64  `db:"current_price"`
	PhysicalPrice int64  `db:"physical_price"`
	Policies      int32  `db:"policies"`
}

func (q *Q) Assets() ([]Asset, error) {
	sql := selectAsset
	var assets []Asset
	err := q.Select(&assets, sql)
	return assets, err
}

// TODO we can't be sure that joined quote asset is default quote asset in code
var selectAsset = sq.Select("a.code, ap.current_price, ap.physical_price, a.policies").From("asset a").Join("asset_pair ap ON ap.base = a.code")
