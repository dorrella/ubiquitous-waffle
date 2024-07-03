package customer

import (
	//"github.com/dorrella/ubiquitous-waffle/service/types"
	"testing"
)

// checks paging of users
func TestListCustomer(t *testing.T) {
	tc := NewTestContext()
	tc.custDb.SeedCustomers(tc.ctx, 45)
	list, next_id, err := tc.custDb.ListCustomers(tc.ctx, 0)
	if err != nil {
		t.Fatalf("unexpected error %s", err.Error())
	}
	if len(*list) != 25 {
		t.Fatal("unexpected number of results")
	}
	if next_id != 25 {
		t.Fatalf("unexpected next id %d", next_id)
	}
	//get next batch
	list, next_id, err = tc.custDb.ListCustomers(tc.ctx, next_id)
	if err != nil {
		t.Fatalf("unexpected error %s", err.Error())
	}
	if len(*list) != 20 {
		t.Fatal("unexpected number of results")
	}
	//0 means probably over
	if next_id != 0 {
		t.Fatalf("unexpected next id %d", next_id)
	}
}

// checks that only non-deleted users are returned
func TestListDeleted(t *testing.T) {
	tc := NewTestContext()
	tc.custDb.SeedCustomers(tc.ctx, 45)
	for i := range 30 {
		_, err := tc.custDb.DeleteCustomer(tc.ctx, int64(i+1), 1)
		if err != nil {
			t.Fatalf("unexpected delete err: %s", err)
		}
	}
	list, next_id, err := tc.custDb.ListCustomers(tc.ctx, 0)
	if err != nil {
		t.Fatalf("unexpected error %s", err.Error())
	}
	if len(*list) != 15 {
		t.Fatal("unexpected number of results")
	}
	//0 means probably over
	if next_id != 0 {
		t.Fatalf("unexpected next id %d", next_id)
	}
}
