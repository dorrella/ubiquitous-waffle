package customer

import (
	"net/http"
	"testing"
)

// validates get happy path
func TestGetCustomer(t *testing.T) {
	tc := newTestContext()
	tc.custDb.SeedCustomers(tc.ctx, 1)
	svr := tc.initServer()
	defer svr.Close()
	client := svr.Client()

	//get cust
	resp, err := client.Get(svr.URL + "/api/customer/1")
	if err != nil {
		t.Fatalf("unexepected error %s", err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("bad status code")
	}

	//parse response
	cust_resp := tc.parseCustomer(resp, t)

	//check for errors
	if cust_resp.Error != "" {
		t.Fatalf("unexptected error from server %s", cust_resp.Error)
	}

	//check fields
	if cust_resp.Customer.Id != 1 {
		t.Errorf("incorrect customer returned")
	}
	if cust_resp.Customer.NameFirst == "" {
		t.Errorf("incomplete customer returned")
	}
	if cust_resp.Customer.NameLast == "" {
		t.Errorf("incomplete customer returned")
	}
	if cust_resp.Customer.Email == "" {
		t.Errorf("incomplete customer returned")
	}
	if cust_resp.Customer.PhoneNumber == "" {
		t.Errorf("incomplete customer returned")
	}
}

// tests that it will return an error on an invalid customer
func TestGetErr(t *testing.T) {
	tc := newTestContext()
	svr := tc.initServer()
	defer svr.Close()
	client := svr.Client()

	//get of empty db
	resp, err := client.Get(svr.URL + "/api/customer/1")
	if err != nil {
		t.Fatalf("unexepected error %s", err.Error())
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("bad status code")
	}
	//parse response
	cust_resp := tc.parseCustomer(resp, t)

	//validate error
	if cust_resp.Error != "error: customer 1 not found" {
		t.Fatal("did not receive expected error from server")
	}
}
