package core

type OrderBookEntry struct {
	BaseAmount  int64 `db:"base_amount"`
	QuoteAmount int64 `db:"quote_amount"`
	Price       int64 `db:"price"`
}
