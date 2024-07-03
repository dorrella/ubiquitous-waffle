package customer

import (
	"context"
	"fmt"
	"github.com/dorrella/ubiquitous-waffle/service/database"
	"github.com/dorrella/ubiquitous-waffle/service/types"
)

const sqlUpdateCust = `UPDATE customers SET
      name_pref=:name_pref,
      name_first=:name_first,
      name_middle=:name_middle,
      name_last=:name_last,
      name_suffix=:name_suffix,
      email=:email,
      phone_number=:phone_number,
      updated_by=:updated_by,
      updated_at=:updated_at
  WHERE id=:id
`

// updated a given customer. id must be set for update
func (cdb *CustDb) UpdateCustomer(ctx context.Context, cust *types.Customer, updated_by int64) (*types.Customer, error) {
	//make sure customer has not been deleted
	current, err := cdb.getAnyCustomer(ctx, cust.Id, false)
	if err != nil {
		return nil, err
	}
	if current == nil {
		//user id not in db?
		return nil, fmt.Errorf("%w: user not found", types.ErrCustomerValidation)
	}

	err = cdb.validateCustomer(cust)
	if err != nil {
		return nil, err
	}
	//it is possible that the customer is trying to use another
	//customers email, but we will let postgress deal with that
	//since the email is marked unique
	cust.UpdatedBy = updated_by
	cust.UpdatedAt = database.TimeStamp()
	_, err = cdb.Db.NamedExecContext(ctx, sqlUpdateCust, cust)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", types.ErrDatabaseErr, err)
	}
	return cdb.GetCustomer(ctx, current.Id)
}
