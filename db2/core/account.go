package core

import (
	"gitlab.com/swarmfund/go/xdr"
	sq "github.com/lann/squirrel"
)

// AccountByAddress loads a row from `accounts`, by address
func (q *Q) AccountByAddress(dest interface{}, addy string) error {
	sql := selectAccount.Limit(1).Where("accountid = ?", addy)

	return q.Get(dest, sql)
}

func (q *Q) ByAddress(address string) (*Account, error) {
	result := new(Account)
	sql := selectAccount.Limit(1).Where("accountid = ?", address)
	err := q.Get(result, sql)
	if err != nil {
		if q.NoRows(err) {
			return nil, nil
		}
	}
	return result, err
}

func (q *Q) Accounts() QInterface {
	if q.err != nil {
		return q
	}

	q.sql = selectAccount
	return q
}

func (q *Q) ForTypes(types []xdr.AccountType) QInterface {
	if q.err != nil {
		return q
	}
	q.sql = q.sql.Where(sq.Eq{"account_type": types})
	return q
}

func (q *Q) Filter(dest interface{}) error {
	if q.err != nil {
		return q.err
	}
	return q.Repo.Select(dest, q.sql)
}

// SequencesForAddresses loads the current sequence number for every accountid
// specified in `addys`
func (q *Q) SequencesForAddresses(dest interface{}, addys []string) error {
	sql := sq.
		Select("accountid as address").
		From("accounts").
		Where(sq.Eq{"accountid": addys})

	return q.Select(dest, sql)
}

// SequenceProvider returns a new sequence provider.
func (q *Q) SequenceProvider() *SequenceProvider {
	return &SequenceProvider{Q: q}
}

// Get implements `txsub.SequenceProvider`
func (sp *SequenceProvider) Get(addys []string) (map[string]uint64, error) {
	rows := []struct {
		Address  string
		Sequence uint64
	}{}

	err := sp.Q.SequencesForAddresses(&rows, addys)
	if err != nil {
		return nil, err
	}

	results := make(map[string]uint64)
	for _, r := range rows {
		results[r.Address] = r.Sequence
	}
	return results, nil
}

var selectAccount = sq.Select(
	"a.accountid",
	"a.thresholds",
	"a.account_type",
	"a.block_reasons",
	"ex.name",
	"ex.require_review",
).From("accounts a").
	LeftJoin("exchanges ex ON a.accountid=ex.account_id")
