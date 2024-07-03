package customer

import (
	"testing"
)

// happey path
func TestFindEmailCustomer(t *testing.T) {
	tc := NewTestContext()
	tc.custDb.SeedCustomers(tc.ctx, 1)
	email := "DrReggiePierce_1@hotmail.com"

	cust, err := tc.custDb.FindByEmail(tc.ctx, email)
	if err != nil {
		t.Fatalf("unexpected error %s", err.Error())
	}
	if cust == nil {
		t.Fatal("could not find customer")
	}
	if cust.Email != email {
		t.Errorf("did not get right customer")
	}
}

// checks nil is returned when customer never existed
func TestFindEmailNilCust(t *testing.T) {
	tc := NewTestContext()
	email := "DrReggiePierce_1@hotmail.com"
	cust, err := tc.custDb.FindByEmail(tc.ctx, email)
	if cust != nil {
		t.Error("returned non existent customer")
	}
	if err != nil {
		t.Errorf("returned unexpected error %s", err.Error())
	}
}

// checks for nil when customer has been deleted
func TestFindEmailDeletedCust(t *testing.T) {
	tc := NewTestContext()
	tc.custDb.SeedCustomers(tc.ctx, 1)
	email := "DrReggiePierce_1@hotmail.com"
	//verify cust exists
	cust, err := tc.custDb.FindByEmail(tc.ctx, email)
	if err != nil {
		t.Fatalf("unexpected error %s", err.Error())
	}
	if cust.Email != email {
		t.Fatalf("did not get right customer")
	}
	_, err = tc.custDb.DeleteCustomer(tc.ctx, 1, 1)
	if err != nil {
		t.Fatalf("unexpected error %s", err.Error())
	}
	cust, err = tc.custDb.FindByEmail(tc.ctx, email)
	if cust != nil {
		t.Error("returned deleted customer")
	}
	if err != nil {
		t.Errorf("returned unexpected error %s", err.Error())
	}
}
