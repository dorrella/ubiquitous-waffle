package customer

import (
	"errors"
	"github.com/dorrella/ubiquitous-waffle/service/types"
	"testing"
)

// checks that a new customer can be created and
// then validates the individual fields
func TestNewCustomer(t *testing.T) {
	var modified_by int64 = 10
	tc := NewTestContext()
	cust := &types.Customer{
		NamePrefix:  "Dr.",
		NameFirst:   "Reginald",
		NameMiddle:  "Franklin",
		NameLast:    "Pierce",
		NameSuffix:  "IV Phd. OBE",
		Email:       "coolhandluke@hotmail.com",
		PhoneNumber: "123456789",
	}
	cust, err := tc.custDb.NewCustomer(tc.ctx, cust, modified_by)
	// check fields
	if err != nil {
		t.Fatalf("NewCustomer returned error: %s", err.Error())
	}
	if cust.Id != 1 {
		t.Error("NewCustomer does not have id=1. errro with test setup?")
	}
	if cust.Deleted {
		t.Error("NewCustomer was marked for deletion")
	}
	if cust.CreatedAt.IsZero() {
		t.Error("NewCustomer created timestamp not set")
	}
	if cust.CreatedBy != modified_by {
		t.Error("NewCustomer created by not set")
	}
	if cust.UpdatedAt.IsZero() {
		t.Error("NewCustomer updated timestamp not set")
	}
	if cust.UpdatedBy != modified_by {
		t.Error("NewCustomer updated by not set")
	}
	if cust.NamePrefix != "Dr." {
		t.Error("NewCustomer name prefix not set")
	}
	if cust.NameFirst != "Reginald" {
		t.Error("NewCustomer first name not set")
	}
	if cust.NameMiddle != "Franklin" {
		t.Error("NewCustomer middle name not set")
	}
	if cust.NameLast != "Pierce" {
		t.Error("NewCustomer last name not set")
	}
	if cust.NameSuffix != "IV Phd. OBE" {
		t.Error("NewCustomer name suffix not set")
	}
	if cust.Email != "coolhandluke@hotmail.com" {
		t.Error("NewCustomer email not set")
	}
	if cust.PhoneNumber != "123456789" {
		t.Error("NewCustomer phone not set")
	}
}

// test email validation
func TestNewCustomerBadEmail(t *testing.T) {
	var modified_by int64 = 10
	tc := NewTestContext()
	cust := &types.Customer{
		NameFirst:   "Reginald",
		NameLast:    "Pierce",
		Email:       "junkmail",
		PhoneNumber: "123456789",
	}
	cust, err := tc.custDb.NewCustomer(tc.ctx, cust, modified_by)
	if !errors.Is(err, types.ErrEmailValidation) {
		t.Errorf("failed to make user: %s", err.Error())
	}
}

// checks emails must be unique
func TestNewCustomerDupeEmail(t *testing.T) {
	var modified_by int64 = 10
	tc := NewTestContext()
	cust1 := &types.Customer{
		NamePrefix:  "Dr.",
		NameFirst:   "Reginald",
		NameMiddle:  "Franklin",
		NameLast:    "Pierce",
		NameSuffix:  "IV Phd. OBE",
		Email:       "coolhandluke@hotmail.com",
		PhoneNumber: "123456789",
	}
	cust2 := &types.Customer{
		NameFirst:   "Bobby",
		NameLast:    "Hill",
		Email:       "coolhandluke@hotmail.com",
		PhoneNumber: "153549",
	}

	cust1, err := tc.custDb.NewCustomer(tc.ctx, cust1, modified_by)
	if err != nil {
		t.Fatalf("failed to make user: %s", err.Error())
	}
	cust2, err = tc.custDb.NewCustomer(tc.ctx, cust2, modified_by)
	if !errors.Is(err, types.ErrCustomerValidation) {
		t.Errorf("failed to make user: %s", err.Error())
	}
}
