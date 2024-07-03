package customer

import (
	"context"
	"fmt"
	"github.com/dorrella/ubiquitous-waffle/service/types"
)

const sqlFindEmailAny = "SELECT * from customers where email=$1"
const sqlFindEmail = sqlFindEmailAny + " AND deleted=false"

// used to check for deleted customers. only used internally
func (cdb *CustDb) findAnyEmail(ctx context.Context, email string, any bool) (*types.Customer, error) {
	results := []types.Customer{}
	stmt := sqlFindEmail
	if any {
		stmt = sqlFindEmailAny
	}
	err := cdb.Db.Select(&results, stmt, email)
	if err != nil {
		wrapped := fmt.Errorf("%w: %w", types.ErrDatabaseErr, err)
		return nil, wrapped
	}
	if len(results) == 1 {
		//found customer
		return &results[0], nil
	} else if len(results) > 1 {
		//should never happen. email is supposed to be unique
		return nil, fmt.Errorf("%w: too many results for email", types.ErrUnexpectedResult)
	}
	// none found
	return nil, nil
}

// finds customer from email as long as the
// customer has not been marked for deletion
func (cdb *CustDb) FindByEmail(ctx context.Context, email string) (*types.Customer, error) {
	err := cdb.validateEmail(email)
	if err != nil {
		return nil, err
	}

	return cdb.findAnyEmail(ctx, email, false)
}
