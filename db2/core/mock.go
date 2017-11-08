package core

import (
	"gitlab.com/tokend/keychain/db2"
	"github.com/jmoiron/sqlx"
	sq "github.com/lann/squirrel"
	"github.com/stretchr/testify/mock"
)

type CoreQMock struct {
	mock.Mock
}

func (q *CoreQMock) GetRepo() *db2.Repo {
	return nil
}
func (q *CoreQMock) AccountByAddress(dest interface{}, addy string) error {
	args := q.Called(dest, addy)
	return args.Error(0)
}
func (q *CoreQMock) ExchangeName(addy string) (*string, error) {
	return nil, nil
}
func (q *CoreQMock) SequencesForAddresses(dest interface{}, addys []string) error {
	return nil
}

func (q *CoreQMock) SequenceProvider() *SequenceProvider {
	return nil
}
func (q *CoreQMock) LedgerHeaderBySequence(dest interface{}, seq int32) error {
	return nil
}
func (q *CoreQMock) ElderLedger(dest *int32) error {
	return nil
}
func (q *CoreQMock) LatestLedger(dest interface{}) error {
	return nil
}
func (q *CoreQMock) SignersByAddress(dest interface{}, addy string) error {
	args := q.Called(dest, addy)
	return args.Error(0)
}
func (q *CoreQMock) PoliciesByExchangeID(dest interface{}, addy string) error {
	return nil
}
func (q *CoreQMock) TransactionByHash(dest interface{}, hash string) error {
	return nil
}
func (q *CoreQMock) TransactionsByLedger(dest interface{}, seq int32) error {
	return nil
}
func (q *CoreQMock) TransactionFeesByLedger(dest interface{}, seq int32) error {
	return nil
}
func (q *CoreQMock) FeeEntries() FeeEntryQI {
	return nil
}
func (q *CoreQMock) Query(query sq.Sqlizer) (*sqlx.Rows, error) {
	return nil, nil
}
func (q *CoreQMock) NoRows(err error) bool {
	return false
}
func (q *CoreQMock) FeeByOperationType(dest interface{}, operationType int) error {
	return nil
}
