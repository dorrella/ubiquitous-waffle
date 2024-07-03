package customer

import (
	"context"
	"fmt"
	"github.com/dorrella/ubiquitous-waffle/service/types"
)

const sqlListCust = `SELECT * from customers
WHERE id > $1 and deleted=false ORDER BY id LIMIT 25`

// Lists customers who have not been marked for deletion
//
// Returns lists of customers in groups of 25. and next starting id
// start with start_id=0 and go until emptry list or next_id == 0
func (cdb *CustDb) ListCustomers(ctx context.Context, start_id int64) (*[]types.Customer, int64, error) {
	var next_id int64 = 0
	results := []types.Customer{}
	err := cdb.Db.Select(&results, sqlListCust, start_id)
	if err != nil {
		wrapped := fmt.Errorf("%w: %w", types.ErrDatabaseErr, err)
		return nil, next_id, wrapped
	}
	if len(results) == 25 {
		//max len, assume there is 1 more.
		next_id = results[len(results)-1].Id
	}
	return &results, next_id, nil
}
