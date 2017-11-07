package core

import (
	"fmt"
)

type CoinsEmissionRequest struct {
	Issuer    string `db:"issuer"`
	RequestID string `db:"request_id"`
	Amount    int64  `db:"amount"`
	Approved  bool   `db:"is_approved"`
}

// PagingToken returns a suitable paging token for the Offer
func (r *CoinsEmissionRequest) PagingToken() string {
	return fmt.Sprintf("%d", r.RequestID)
}
