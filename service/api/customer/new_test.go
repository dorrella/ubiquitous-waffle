package customer

import (
	"net/http"
	"net/url"
	"testing"
)

// tests happy path of customer creation
func TestNewCustomer(t *testing.T) {
	tc := newTestContext()
	svr := tc.initServer()
	defer svr.Close()
	client := svr.Client()

	//create form data
	values := url.Values{}
	values.Set("name_first", "Rusty")
	values.Set("name_middle", "M")
	values.Set("name_last", "Shackleford")
	values.Set("email", "rustys@yahoo.com")
	values.Set("phone_number", "4841234321")

	//post form
	resp, err := client.PostForm(svr.URL+"/api/customer", values)
	if err != nil {
		t.Fatalf("unexepected error %s", err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("bad status code")
	}

	//parse response
	cust_resp := tc.parseCustomer(resp, t)

	//check for error
	if cust_resp.Error != "" {
		t.Fatalf("unexptected error from server %s", cust_resp.Error)
	}

	//check fields
	if cust_resp.Customer.Id != 1 {
		t.Errorf("incorrect customer returned")
	}
	if cust_resp.Customer.NameFirst != "Rusty" {
		t.Errorf("incomplete customer returned")
	}
	if cust_resp.Customer.NameLast != "Shackleford" {
		t.Errorf("incomplete customer returned")
	}
	if cust_resp.Customer.Email != "rustys@yahoo.com" {
		t.Errorf("incomplete customer returned")
	}
	if cust_resp.Customer.PhoneNumber != "4841234321" {
		t.Errorf("incomplete customer returned")
	}
}

// test for an invalid customer.
func TestInvalidCust(t *testing.T) {
	tc := newTestContext()
	svr := tc.initServer()
	defer svr.Close()
	client := svr.Client()
	values := url.Values{}

	//no name set
	values.Set("email", "rustys@yahoo.com")
	values.Set("phone_number", "4841234321")

	//send request
	resp, err := client.PostForm(svr.URL+"/api/customer", values)
	if err != nil {
		t.Fatalf("unexepected error %s", err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("bad status code")
	}

	//parse response
	cust_resp := tc.parseCustomer(resp, t)

	//validate error
	if cust_resp.Error != "error: invalid customer: must have at least one valid name" {
		t.Errorf("unexptected error from server %s", cust_resp.Error)
	}
}
