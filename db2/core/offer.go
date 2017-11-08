package core

type Offer struct {
	OwnerID        string `db:"owner_id"`
	OfferID        uint64 `db:"offer_id"`
	BaseAssetCode  string `db:"base_asset_code"`
	QuoteAssetCode string `db:"quote_asset_code"`
	IsBuy          bool   `db:"is_buy"`
	BaseBalanceID  string `db:"base_balance_id"`
	QuoteBalanceID string `db:"quote_balance_id"`
	CreatedAt      int64  `db:"created_at"`
	OrderBookEntry
}
