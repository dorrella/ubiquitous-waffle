package customer

import (
	"context"
	"fmt"
	"github.com/dorrella/ubiquitous-waffle/service/database"
	"github.com/dorrella/ubiquitous-waffle/service/types"
	"net/mail"
)

const sqlNewCust = `INSERT INTO customers
    (
      name_pref,
      name_first,
      name_middle,
      name_last,
      name_suffix,
      email,
      phone_number,
      deleted,
      created_by,
      created_at,
      updated_by,
      updated_at
    )
  VALUES
    (
      :name_pref,
      :name_first,
      :name_middle,
      :name_last,
      :name_suffix,
      :email,
      :phone_number,
      :deleted,
      :created_by,
      :created_at,
      :updated_by,
      :updated_at
    )
`

const sqlReactivateCust = `UPDATE customers SET
      name_pref=:name_pref,
      name_first=:name_first,
      name_middle=:name_middle,
      name_last=:name_last,
      name_suffix=:name_suffix,
      phone_number=:phone_number,
      deleted=:deleted,
      created_by=:created_by,
      created_at=:created_at,
      updated_by=:updated_by,
      updated_at=:updated_at
   WHERE email=:email

`

// validates the email address. chechs it exists and
// passes net/mail parsing
func (cdb *CustDb) validateEmail(email string) error {
	if len(email) < 1 {
		return fmt.Errorf("%w: email is required", types.ErrEmailValidation)
	}

	//check that the email can be parsed
	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("%w: cannot parse: %w", types.ErrEmailValidation, err)
	}
	return nil
}

// validate the fields for a new customer
func (cdb *CustDb) validateCustomer(cust *types.Customer) error {
	var valid bool = false
	// maybe just first and last name > 0?
	for _, token := range []string{cust.NamePrefix, cust.NameFirst, cust.NameMiddle, cust.NameLast, cust.NameSuffix} {
		if len(token) > 0 {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("%w: must have at least one valid name", types.ErrCustomerValidation)
	}

	err := cdb.validateEmail(cust.Email)
	if err != nil {
		return fmt.Errorf("%w: %w", types.ErrCustomerValidation, err)
	}
	//is phone required?
	if len(cust.PhoneNumber) < 1 {
		return fmt.Errorf("%w: phone number is required", types.ErrCustomerValidation)
	}
	return nil
}

// creates a new user by updating an old customer. works
// on the assumption that emails are unique
func (cdb *CustDb) reactivateCustomer(ctx context.Context, old *types.Customer, new *types.Customer, created_by int64) (*types.Customer, error) {
	if old.Deleted == false {
		return nil, fmt.Errorf("%w: cannot reactivate active customer", types.ErrCustomerValidation)
	}
	timestamp := database.TimeStamp()
	new.Id = old.Id
	new.Deleted = false
	new.CreatedBy = created_by
	new.UpdatedBy = created_by
	new.CreatedAt = timestamp
	new.UpdatedAt = timestamp
	_, err := cdb.Db.NamedExec(sqlReactivateCust, new)
	if err != nil {
		wrapped := fmt.Errorf("%w: %w", types.ErrDatabaseErr, err)
		return nil, wrapped
	}
	// driver doesnt support res.LastIntertedId()
	// so reload id from db
	new, err = cdb.FindByEmail(ctx, new.Email)
	if err != nil {
		return nil, err
	}
	return new, nil

}

// creates a new customer and inserts it into the customer db
//
// the customer must have at least 1 valid name
// the customer must have a phone numver
// the customer must have a unique and valid email
func (cdb *CustDb) NewCustomer(ctx context.Context, cust *types.Customer, created_by int64) (*types.Customer, error) {
	//validate fields
	err := cdb.validateCustomer(cust)
	if err != nil {
		return nil, err
	}

	//check that it is unique.
	existing, err := cdb.findAnyEmail(ctx, cust.Email, true)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		if existing.Deleted {
			return cdb.reactivateCustomer(ctx, existing, cust, created_by)
		}
		return nil, fmt.Errorf("%w: there is an existing user with that email", types.ErrCustomerValidation)
	}

	//insert new user
	timestamp := database.TimeStamp()
	cust.Deleted = false
	cust.CreatedBy = created_by
	cust.UpdatedBy = created_by
	cust.CreatedAt = timestamp
	cust.UpdatedAt = timestamp
	_, err = cdb.Db.NamedExec(sqlNewCust, cust)
	if err != nil {
		return nil, err
	}
	// driver doesnt support res.LastIntertedId()
	// so reload id from db
	cust, err = cdb.FindByEmail(ctx, cust.Email)
	if err != nil {
		wrapped := fmt.Errorf("%w: %w", types.ErrDatabaseErr, err)
		return nil, wrapped
	}
	return cust, nil

}
