package core

type Trust struct {
	AllowedAccount string `db:"allowed_account"`
	BalanceToUse   string `db:"balance_to_use"`
}
