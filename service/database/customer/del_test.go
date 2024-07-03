package customer

import (
	"errors"
	"github.com/dorrella/ubiquitous-waffle/service/types"
	"testing"
)

func TestDelCustomer(t *testing.T) {
	var deleted_by int64 = 10
	tc := NewTestContext()
	tc.custDb.SeedCustomers(tc.ctx, 1)
	cust, err := tc.custDb.getAnyCustomer(tc.ctx, 1, true)
	if err != nil {
		t.Errorf("Customer not seeded: %s", err.Error())
	}
	if cust.Deleted == true {
		t.Error("cust deleted")
	}
	cust, err = tc.custDb.DeleteCustomer(tc.ctx, 1, deleted_by)
	if err != nil {
		t.Error("unexpected err during del")
	}
	if cust.Deleted == false {
		t.Error("delete failed")
	}
	if cust.UpdatedBy != deleted_by {
		t.Error("cust was not updated")
	}
	//double check db
	cust, err = tc.custDb.getAnyCustomer(tc.ctx, 1, true)
	if err != nil {
		t.Errorf("get failed: %s", err.Error())
	}
	if cust.Deleted == false {
		t.Error("cust not really deleted")
	}
}

func TestEmptyDelete(t *testing.T) {
	var deleted_by int64 = 10
	tc := NewTestContext()
	_, err := tc.custDb.DeleteCustomer(tc.ctx, 1, deleted_by)
	if !errors.Is(err, types.ErrCustomerValidation) {
		t.Error("unexpected err during del")
	}
}

func TestDoubleDelete(t *testing.T) {
	var deleted_by int64 = 10
	tc := NewTestContext()
	tc.custDb.SeedCustomers(tc.ctx, 1)

	_, err := tc.custDb.DeleteCustomer(tc.ctx, 1, deleted_by)
	if err != nil {
		t.Errorf("unexepected err %s", err.Error())
	}
	_, err = tc.custDb.DeleteCustomer(tc.ctx, 1, deleted_by)
	if !errors.Is(err, types.ErrCustomerValidation) {
		t.Errorf("unexepected err %s", err.Error())
	}

}
