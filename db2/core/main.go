// Package core contains database record definitions useable for
// reading rows from a Stellar Core db
package core

import (
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/keychain/db2"
	"github.com/jmoiron/sqlx"
	sq "github.com/lann/squirrel"
)

type ExchangeData struct {
	Name          *string
	RequireReview *bool `db:"require_review"`
}

// Account is a row of data from the `accounts` table
type Account struct {
	Accountid    string
	Balance      xdr.Int64
	Thresholds   xdr.Thresholds
	AccountType  int32 `db:"account_type"`
	BlockReasons int32 `db:"block_reasons"`
	ExchangeData
}

// LedgerHeader is row of data from the `ledgerheaders` table
type LedgerHeader struct {
	LedgerHash     string           `db:"ledgerhash"`
	PrevHash       string           `db:"prevhash"`
	BucketListHash string           `db:"bucketlisthash"`
	CloseTime      int64            `db:"closetime"`
	Sequence       uint32           `db:"ledgerseq"`
	Data           xdr.LedgerHeader `db:"data"`
}

// Q is a helper struct on which to hang common queries against a stellar
// core database.
type Q struct {
	*db2.Repo

	err error
	sql sq.SelectBuilder
}

func NewQ(repo *db2.Repo) *Q {
	return &Q{
		Repo: repo,
	}
}

func (q *Q) GetRepo() *db2.Repo {
	return q.Repo
}

// Q interface helper for testing purposes mainly

type QInterface interface {
	GetRepo() *db2.Repo
	// DEPRECATED use `ByAddress` with explicit return value
	AccountByAddress(dest interface{}, addy string) error
	ExchangeName(addy string) (*string, error)
	ByAddress(address string) (*Account, error)
	SequencesForAddresses(dest interface{}, addys []string) error
	SequenceProvider() *SequenceProvider
	LedgerHeaderBySequence(dest interface{}, seq int32) error
	ElderLedger(dest *int32) error
	LatestLedger(dest interface{}) error
	SignersByAddress(dest interface{}, addy string) error
	PoliciesByExchangeID(dest interface{}, addy string) error
	LimitsByAddress(dest interface{}, addy string) error
	StatisticsByAddress(dest interface{}, addy string) error
	BalancesByAddress(dest interface{}, addy string) error
	BalanceByID(dest interface{}, bid string) error
	TransactionByHash(dest interface{}, hash string) error
	TransactionsByLedger(dest interface{}, seq int32) error
	TransactionFeesByLedger(dest interface{}, seq int32) error
	FeeEntries() FeeEntryQI
	AccountTypeLimits() AccountTypeLimitsQI
	Query(query sq.Sqlizer) (*sqlx.Rows, error)
	NoRows(err error) bool
	// Returns nil, if not found
	FeeByTypeAsset(feeType int, asset string) (*FeeEntry, error)
	LimitsByAccountType(accountType int) (*AccountTypeLimits, error)
	Assets() ([]Asset, error)

	AvailableEmissions(dest interface{}, masterAccountID string) error

	EmissionRequestByExchangeAndRef(exchange, reference string) (*bool, error)

	CoinsInCirculation(dest interface{}, masterAccountID string) error
	CoinsInCirculationForAsset(dest interface{}, masterAccountID, asset string) error

	// should probably be separate accounts repo
	Accounts() QInterface
	ForTypes(types []xdr.AccountType) QInterface
	Filter(dest interface{}) error

	CoinsEmissions() *CoinsEmissionQ
	Trusts() *TrustQ

	Offers() *OfferQ
	OrderBook() *OrderBookQ

	AssetPair(base, quote string) (*AssetPair, error)
	AssetPairs() ([]AssetPair, error)
}

// PriceLevel represents an aggregation of offers to trade at a certain
// price.
type PriceLevel struct {
	Pricen int32   `db:"pricen"`
	Priced int32   `db:"priced"`
	Pricef float64 `db:"pricef"`
	Amount int64   `db:"amount"`
}

// SequenceProvider implements `txsub.SequenceProvider`
type SequenceProvider struct {
	Q *Q
}

// Transaction is row of data from the `txhistory` table from stellar-core
type Transaction struct {
	TransactionHash string                    `db:"txid"`
	LedgerSequence  int32                     `db:"ledgerseq"`
	Index           int32                     `db:"txindex"`
	Envelope        xdr.TransactionEnvelope   `db:"txbody"`
	Result          xdr.TransactionResultPair `db:"txresult"`
	ResultMeta      xdr.TransactionMeta       `db:"txmeta"`
}

// TransactionFee is row of data from the `txfeehistory` table from stellar-core
type TransactionFee struct {
	TransactionHash string                 `db:"txid"`
	LedgerSequence  int32                  `db:"ledgerseq"`
	Index           int32                  `db:"txindex"`
	Changes         xdr.LedgerEntryChanges `db:"txchanges"`
}

// ElderLedger represents the oldest "ingestable" ledger known to the
// stellar-core database this ingestion system is communicating with.  Horizon,
// which wants to operate on a contiguous range of ledger data (i.e. free from
// gaps) uses the elder ledger to start importing in the case of an empty
// database.
//
// Due to the design of stellar-core, ledger 1 will _always_ be in the database,
// even when configured to catchup minimally, so we cannot simply take
// MIN(ledgerseq). Instead, we can find whether or not 1 is the elder ledger by
// checking for the presence of ledger 2.
func (q *Q) ElderLedger(dest *int32) error {
	var found bool
	err := q.GetRaw(&found, `
		SELECT CASE WHEN EXISTS (
		    SELECT *
		    FROM ledgerheaders
		    WHERE ledgerseq = 2
		)
		THEN CAST(1 AS BIT)
		ELSE CAST(0 AS BIT) END
	`)

	if err != nil {
		return err
	}

	// if ledger 2 is present, use it 1 as the elder ledger (since 1 is guaranteed
	// to be present)
	if found {
		*dest = 1
		return nil
	}

	err = q.GetRaw(dest, `
		SELECT COALESCE(MIN(ledgerseq), 0)
		FROM ledgerheaders
		WHERE ledgerseq > 2
	`)

	return err
}

// LatestLedger loads the latest known ledger
func (q *Q) LatestLedger(dest interface{}) error {
	return q.GetRaw(dest, `SELECT COALESCE(MAX(ledgerseq), 0) FROM ledgerheaders`)
}
