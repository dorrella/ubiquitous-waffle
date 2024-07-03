package customer

import (
	"bytes"
	"encoding/json"
	"github.com/dorrella/ubiquitous-waffle/service/types"
	"net/http"
	"testing"
)

// tests that we can update an existing customer
func TestUpdate(t *testing.T) {
	tc := newTestContext()
	tc.custDb.SeedCustomers(tc.ctx, 1)
	svr := tc.initServer()
	defer svr.Close()
	client := svr.Client()

	//create update and marshal
	cust, err := tc.custDb.GetCustomer(tc.ctx, 1)
	if err != nil {
		t.Fatalf("unexpected error getting customer %s", err.Error())
	}
	cust.NamePrefix = ""
	cust.NameFirst = "Max"
	cust.NameLast = "Powers"
	cust.NameSuffix = ""

	byte_slice, err := json.Marshal(cust)
	if err != nil {
		t.Fatalf("unexpected error reading resp %s", err.Error())
	}
	reader := bytes.NewReader(byte_slice)
	//create request
	req, err := http.NewRequest(http.MethodPut, svr.URL+"/api/customer/1", reader)
	if err != nil {
		t.Fatalf("unexpected error reading resp %s", err.Error())
	}

	//call update
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("unexepected error %s", err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("bad status code")
	}
	//parse response
	cust_resp := tc.parseCustomer(resp, t)

	//check for errors
	if cust_resp.Error != "" {
		t.Fatalf("received unexpected error from server %s", cust_resp.Error)
	}
	//validate results
	if cust_resp.Customer.Id != 1 {
		t.Errorf("wrong customer returned")
	}
	if cust_resp.Customer.NameFirst != "Max" {
		t.Errorf("update failed")
	}
	if cust_resp.Customer.NameLast != "Powers" {
		t.Errorf("update failed")
	}

}

// tests that update fails when no customer
func TestUpdateErr(t *testing.T) {
	tc := newTestContext()
	svr := tc.initServer()
	defer svr.Close()
	client := svr.Client()

	cust := &types.Customer{
		Id:        1,
		NameFirst: "Max",
		NameLast:  "Powers",
	}
	byte_slice, err := json.Marshal(cust)
	if err != nil {
		t.Fatalf("unexpected error reading resp %s", err.Error())
	}
	reader := bytes.NewReader(byte_slice)
	//create request
	req, err := http.NewRequest(http.MethodPut, svr.URL+"/api/customer/1", reader)
	if err != nil {
		t.Fatalf("unexpected error reading resp %s", err.Error())
	}

	// run request
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("unexepected error %s", err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("bad status code")
	}

	//parse response
	cust_resp := tc.parseCustomer(resp, t)

	//validate error
	if cust_resp.Error != "error: invalid customer: user not found" {
		t.Fatal("did not receive expected error from server")
	}
}
