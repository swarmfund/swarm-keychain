package core

import sq "github.com/lann/squirrel"

type AssetPair struct {
	BaseAsset               string `db:"base"`
	QuoteAsset              string `db:"quote"`
	CurrentPrice            int64  `db:"current_price"`
	PhysicalPrice           int64  `db:"physical_price"`
	PhysicalPriceCorrection int64  `db:"physical_price_correction"`
	MaxPriceStep            int64  `db:"max_price_step"`
	Policies                int32  `db:"policies"`
}

func (q *Q) AssetPairs() ([]AssetPair, error) {
	sql := selectAssetPair
	var assetPairs []AssetPair
	err := q.Select(&assetPairs, sql)
	return assetPairs, err
}

// returns nil, if not found
func (q *Q) AssetPair(base, quote string) (*AssetPair, error) {
	sql := selectAssetPair.Where("base = ? AND quote = ?", base, quote)
	var result AssetPair
	err := q.Get(&result, sql)
	if q.Repo.NoRows(err) {
		return nil, nil
	}

	return &result, err
}

var selectAssetPair = sq.Select("a.base, a.quote, a.current_price, a.physical_price, a.physical_price_correction, a.max_price_step, a.policies").From("asset_pair a")
