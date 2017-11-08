package db2

import (
	"fmt"
)

// TotalOrderID represents the ID portion of rows that are identified by the
// "TotalOrderID".  See total_order_id.go in the `db` package for details.
type TotalOrderID struct {
	ID int64 `db:"id"`
}

// PagingToken returns a cursor for this record
func (r *TotalOrderID) PagingToken() string {
	return fmt.Sprintf("%d", r.ID)
}
