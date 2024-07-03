package customer

import (
	"context"
	"fmt"
	"github.com/dorrella/ubiquitous-waffle/service/database"
	"github.com/dorrella/ubiquitous-waffle/service/types"
)

const sqlDeleteCust = `UPDATE customers SET
      deleted=:deleted,
      updated_by=:updated_by,
      updated_at=:updated_at
  WHERE id=:id
`

// Deletes a customer, if able.
func (cdb *CustDb) DeleteCustomer(ctx context.Context, id int64, deleted_by int64) (*types.Customer, error) {
	cust, err := cdb.getAnyCustomer(ctx, id, true)
	if err != nil {
		return nil, err
	}
	if cust == nil || cust.Deleted {
		err = fmt.Errorf("%w: customer not found", types.ErrCustomerValidation)
		return nil, err
	}

	//actually delete customer
	cust.Deleted = true
	cust.UpdatedBy = deleted_by
	cust.UpdatedAt = database.TimeStamp()
	_, err = cdb.Db.NamedExecContext(ctx, sqlDeleteCust, &cust)
	if err != nil {
		wrapped := fmt.Errorf("%w: %w", types.ErrDatabaseErr, err)
		return nil, wrapped
	}
	return cdb.getAnyCustomer(ctx, id, true)
}
