package customer

import (
	"context"
	"fmt"
	"github.com/dorrella/ubiquitous-waffle/service/types"
)

const sqlGetAny = "SELECT * from customers where id=$1"
const sqlGet = sqlGetAny + " AND deleted=false"

// used to check for deleted customers. only used internally.
func (cdb *CustDb) getAnyCustomer(ctx context.Context, id int64, any bool) (*types.Customer, error) {
	results := []types.Customer{}
	stmt := sqlGet
	if any {
		stmt = sqlGetAny
	}
	err := cdb.Db.Select(&results, stmt, id)
	if err != nil {
		wrapped := fmt.Errorf("%w: %w", types.ErrDatabaseErr, err)
		return nil, wrapped
	}
	if len(results) == 1 {
		//found customer
		return &results[0], nil
	} else if len(results) > 1 {
		//should never happen. email is supposed to be unique
		return nil, fmt.Errorf("%w: too many results for id", types.ErrUnexpectedResult)
	}
	//not found
	return nil, nil
}

// get a customer by id. if none, returns nil
func (cdb *CustDb) GetCustomer(ctx context.Context, id int64) (*types.Customer, error) {
	return cdb.getAnyCustomer(ctx, id, false)
}
