package customer

import (
	"net/http"
	"testing"
)

// test delete enpoint
func TestDelCustomer(t *testing.T) {
	tc := newTestContext()
	tc.custDb.SeedCustomers(tc.ctx, 1)
	svr := tc.initServer()
	defer svr.Close()
	client := svr.Client()

	//create request
	req, err := http.NewRequest(http.MethodDelete, svr.URL+"/api/customer/1", nil)
	if err != nil {
		t.Fatalf("unexpected error reading resp %s", err.Error())
	}

	//send request
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("unexepected error %s", err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("bad status code")
	}

	//parse response
	cust_resp := tc.parseCustomer(resp, t)

	//check error
	if cust_resp.Error != "" {
		t.Fatalf("unexptected error from server %s", cust_resp.Error)
	}

	//check fields
	if cust_resp.Customer.Id != 1 {
		t.Errorf("incorrect customer returned")
	}
}

// tests that it will return an error on an invalid customer
func TestDelErr(t *testing.T) {
	tc := newTestContext()
	svr := tc.initServer()
	defer svr.Close()
	client := svr.Client()

	//create request
	req, err := http.NewRequest(http.MethodDelete, svr.URL+"/api/customer/1", nil)
	if err != nil {
		t.Fatalf("unexpected error reading resp %s", err.Error())
	}

	//send request
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("unexepected error %s", err.Error())
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("bad status code")
	}

	//parse reponse
	cust_resp := tc.parseCustomer(resp, t)

	//validate error
	if cust_resp.Error != "error: invalid customer: customer not found" {
		t.Fatal("did not receive expected error from server")
	}
}
