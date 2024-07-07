// customer database contoller
package customer

import (
	"context"
	"github.com/dorrella/ubiquitous-waffle/service/types"
	"github.com/jmoiron/sqlx"
	"strconv"
)

// todo use this instead for testing dep injections
type CustInterface interface {
	GetCustomer(ctx context.Context, id int64) (*types.Customer, error)
	FindByEmail(ctx context.Context, email string) (*types.Customer, error)
	DeleteCustomer(ctx context.Context, id int64, deleted_by int64) (*types.Customer, error)
	NewCustomer(ctx context.Context, cust *types.Customer, created_by int64) (*types.Customer, error)
	ListCustomers(ctx context.Context, start_id int64) (*[]types.Customer, int64, error)
}

// wrapper for customer db interactions
type CustDb struct {
	Db *sqlx.DB
}

type TestCustDb struct {
	CustDb
}

// wipes db for tests
func (tc *TestCustDb) Reset() {
	_, err := tc.Db.Exec("TRUNCATE TABLE customers RESTART IDENTITY")
	if err != nil {
		panic(err)
	}
}

// makes a number of unique customers
func (tc *TestCustDb) SeedCustomers(ctx context.Context, num int) {
	for idx := range num {
		//ofset by 1 since primary keys start at 1
		num_str := strconv.Itoa(idx + 1)
		email := "DrReggiePierce_" + num_str + "@hotmail.com"
		cust := &types.Customer{
			NamePrefix:  "Dr.",
			NameFirst:   "Reginald",
			NameMiddle:  "Franklin",
			NameLast:    "Pierce",
			NameSuffix:  "IV Phd. OBE",
			Email:       email,
			PhoneNumber: "123" + num_str,
		}
		_, err := tc.NewCustomer(ctx, cust, 1)
		if err != nil {
			panic(err)
		}
	}

}
