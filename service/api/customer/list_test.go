package customer

import (
	"net/http"
	"net/url"
	"strconv"
	"testing"
)

// Tests list api route
func TestListCustomers(t *testing.T) {
	tc := newTestContext()
	//two batches of 25 and 2-
	tc.custDb.SeedCustomers(tc.ctx, 45)
	svr := tc.initServer()
	defer svr.Close()
	client := svr.Client()
	//calling client without query for first page
	resp, err := client.Get(svr.URL + "/api/customer/list")
	if err != nil {
		t.Fatalf("unexepected error %s", err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("bad status code")
	}

	//parse response
	cust_resp := tc.parseCustomerList(resp, t)

	//check for server errors
	if cust_resp.Error != "" {
		t.Fatalf("unexptected error from server %s", cust_resp.Error)
	}
	// check expected response
	if len(*cust_resp.Customers) != 25 {
		t.Error("unexpected number of customers returned")
	}
	if cust_resp.Next != 25 {
		t.Error("invalid next id")
	}

	//try next page
	path, err := url.Parse(svr.URL + "/api/customer/list")
	if err != nil {
		t.Fatalf("unexepected error %s", err.Error())
	}
	values := url.Values{}
	//Itoa does not like int64, but is a wrapper around FormatInt
	values.Add("next", strconv.FormatInt(cust_resp.Next, 10))
	path.RawQuery = values.Encode()
	resp, err = client.Get(path.String())
	if err != nil {
		t.Fatalf("unexepected error %s", err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("bad status code")
	}

	//parse rew resp
	cust_resp = tc.parseCustomerList(resp, t)

	// check for errors
	if cust_resp.Error != "" {
		t.Fatalf("unexptected error from server %s", cust_resp.Error)
	}

	//validate response
	if len(*cust_resp.Customers) != 20 {
		t.Error("unexpected number of customers returned")
	}
	if cust_resp.Next != 0 {
		t.Error("invalid next id")
	}

}
