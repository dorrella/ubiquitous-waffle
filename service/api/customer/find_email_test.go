package customer

import (
	"net/http"
	"net/url"
	"testing"
)

// test finding a customer by their email
func TestFindEmailCustomer(t *testing.T) {
	tc := newTestContext()
	tc.custDb.SeedCustomers(tc.ctx, 1)
	email := "DrReggiePierce_1@hotmail.com"
	svr := tc.initServer()
	defer svr.Close()
	client := svr.Client()

	//setup query for email
	path, err := url.Parse(svr.URL + "/api/customer/by_email")
	if err != nil {
		t.Fatalf("unexepected error %s", err.Error())
	}
	values := url.Values{}
	values.Add("email", email)
	path.RawQuery = values.Encode()

	//run request
	resp, err := client.Get(path.String())
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
	if cust_resp.Customer.Email != email {
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
func TestFindEmailErr(t *testing.T) {
	tc := newTestContext()
	email := "DrReggiePierce_1@hotmail.com"
	svr := tc.initServer()
	defer svr.Close()
	client := svr.Client()

	//set request query
	path, err := url.Parse(svr.URL + "/api/customer/by_email")
	if err != nil {
		t.Fatalf("unexepected error %s", err.Error())
	}
	values := url.Values{}
	values.Add("email", email)
	path.RawQuery = values.Encode()

	//run request
	resp, err := client.Get(path.String())
	if err != nil {
		t.Fatalf("unexepected error %s", err.Error())
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("bad status code")
	}

	// parse response
	cust_resp := tc.parseCustomer(resp, t)

	//validate error
	if cust_resp.Error != "error: customer email DrReggiePierce_1@hotmail.com not found" {
		t.Fatal("did not receive expected error from server")
	}
}
