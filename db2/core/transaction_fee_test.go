package core

import (
	"testing"

	"gitlab.com/distributed_lab/tokend/keychain/test"
)

func TestTransactionFeesByLedger(t *testing.T) {
	tt := test.Start(t).Scenario("base")
	defer tt.Finish()
	q := &Q{tt.CoreRepo()}

	var fees []TransactionFee
	err := q.TransactionFeesByLedger(&fees, 2)

	if tt.Assert.NoError(err) {
		tt.Assert.Len(fees, 3)
	}
}
