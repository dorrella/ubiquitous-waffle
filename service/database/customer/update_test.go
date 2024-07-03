package customer

import (
	"errors"
	"github.com/dorrella/ubiquitous-waffle/service/types"
	"testing"
)

// basic update of existing customer
func TestUpdateCustomer(t *testing.T) {
	var updated_by int64 = 11
	tc := NewTestContext()
	tc.custDb.SeedCustomers(tc.ctx, 1)

	cust, err := tc.custDb.getAnyCustomer(tc.ctx, 1, true)
	if err != nil {
		t.Fatalf("could not find cust to update %s", err.Error())
	}

	cust.NamePrefix = ""
	cust.NameFirst = "Bob"
	cust.NameMiddle = ""
	cust.NameLast = "Dobolina"
	cust.NameSuffix = ""
	cust.Email = "bob@gmail.com"
	cust.PhoneNumber = "87654332"

	cust, err = tc.custDb.UpdateCustomer(tc.ctx, cust, updated_by)
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	if cust.NamePrefix != "" {
		t.Error(" not updated")
	}
	if cust.NameFirst != "Bob" {
		t.Error(" not updated")
	}
	if cust.NameMiddle != "" {
		t.Error(" not updated")
	}
	if cust.NameLast != "Dobolina" {
		t.Error(" not updated")
	}
	if cust.NameSuffix != "" {
		t.Error(" not updated")
	}
	if cust.Email != "bob@gmail.com" {
		t.Error(" not updated")
	}
	if cust.PhoneNumber != "87654332" {
		t.Error(" not updated")
	}
}

// checks for error when no customer
func TestUpdateNilCustomer(t *testing.T) {
	var updated_by int64 = 11
	tc := NewTestContext()
	_, err := tc.custDb.UpdateCustomer(tc.ctx, &types.Customer{}, updated_by)
	if !errors.Is(err, types.ErrCustomerValidation) {
		t.Errorf("unexpected error %s", err.Error())
	}
}

// tests that we get an error when updated a deleted customer
func TestUpdateDeletedCustomer(t *testing.T) {
	var updated_by int64 = 11
	tc := NewTestContext()
	tc.custDb.SeedCustomers(tc.ctx, 1)
	_, err := tc.custDb.DeleteCustomer(tc.ctx, 1, updated_by)
	if err != nil {
		t.Fatalf("test setup err %s", err.Error())
	}
	_, err = tc.custDb.UpdateCustomer(tc.ctx, &types.Customer{}, updated_by)
	if !errors.Is(err, types.ErrCustomerValidation) {
		t.Errorf("unexpected error %s", err.Error())
	}
}

// basic validation on update
func TestUpdateInvalidCustomer(t *testing.T) {
	var updated_by int64 = 11
	tc := NewTestContext()
	tc.custDb.SeedCustomers(tc.ctx, 1)
	//email := "DrReggiePierce_1@hotmail.com"

	cust, err := tc.custDb.getAnyCustomer(tc.ctx, 1, true)
	if err != nil {
		t.Fatalf("could not find cust to update %s", err.Error())
	}

	//no name
	cust.NamePrefix = ""
	cust.NameFirst = ""
	cust.NameMiddle = ""
	cust.NameLast = ""
	cust.NameSuffix = ""
	cust.Email = "bob@gmail.com"
	cust.PhoneNumber = "87654332"

	cust, err = tc.custDb.UpdateCustomer(tc.ctx, cust, updated_by)
	if !errors.Is(err, types.ErrCustomerValidation) {
		t.Fatalf("unexpected error: %s", err.Error())
	}
}

// basic email validation on update
func TestUpdateInvalidEmail(t *testing.T) {
	var updated_by int64 = 11
	tc := NewTestContext()
	tc.custDb.SeedCustomers(tc.ctx, 1)

	cust, err := tc.custDb.getAnyCustomer(tc.ctx, 1, true)
	if err != nil {
		t.Fatalf("could not find cust to update %s", err.Error())
	}

	//no name
	cust.NamePrefix = ""
	cust.NameFirst = "Bob"
	cust.NameMiddle = "B"
	cust.NameLast = "Jones"
	cust.NameSuffix = ""
	cust.Email = "bobby_b"
	cust.PhoneNumber = "87654332"

	cust, err = tc.custDb.UpdateCustomer(tc.ctx, cust, updated_by)
	if !errors.Is(err, types.ErrCustomerValidation) {
		t.Fatalf("unexpected error: %s", err.Error())
	}
}
