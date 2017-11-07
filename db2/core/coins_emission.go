package core

type CoinsEmission struct {
	SerialNumber string `db:"serial_number"`
	Amount       int64  `db:"amount"`
	LastModified int64  `db:"lastmodified"`
}

type AssetAmount struct {
	Asset  string `db:"asset"`
	Amount int64  `db:"amount"`
}
