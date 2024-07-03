package customer

import (
	"testing"
)

// basic get
func TestGetCustomer(t *testing.T) {
	tc := NewTestContext()
	tc.custDb.SeedCustomers(tc.ctx, 1)
	cust, err := tc.custDb.GetCustomer(tc.ctx, 1)
	if err != nil {
		t.Fatalf("unexpected error %s", err.Error())
	}
	if cust.Id != 1 {
		t.Errorf("did not get right customer")
	}
}

// get of user that does not exist returns nil
func TestGetNilCust(t *testing.T) {
	tc := NewTestContext()
	cust, err := tc.custDb.GetCustomer(tc.ctx, 1)
	if cust != nil {
		t.Error("returned non existent customer")
	}
	if err != nil {
		t.Errorf("returned unexpected error %s", err.Error())
	}
}

// get of deleted user returns nil
func TestGetDeletedCust(t *testing.T) {
	tc := NewTestContext()
	tc.custDb.SeedCustomers(tc.ctx, 1)
	cust, err := tc.custDb.GetCustomer(tc.ctx, 1)
	if err != nil {
		t.Fatalf("unexpected error %s", err.Error())
	}
	if cust.Id != 1 {
		t.Fatalf("did not get right customer")
	}
	_, err = tc.custDb.DeleteCustomer(tc.ctx, 1, 1)
	if err != nil {
		t.Fatalf("unexpected error %s", err.Error())
	}
	cust, err = tc.custDb.GetCustomer(tc.ctx, 1)
	if cust != nil {
		t.Error("returned deleted customer")
	}
	if err != nil {
		t.Errorf("returned unexpected error %s", err.Error())
	}
}
