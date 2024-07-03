package customer

import (
	"context"
	"encoding/json"
	conf "github.com/dorrella/ubiquitous-waffle/service/config"
	"github.com/dorrella/ubiquitous-waffle/service/database"
	custdb "github.com/dorrella/ubiquitous-waffle/service/database/customer"
	"github.com/dorrella/ubiquitous-waffle/service/logging"
	"github.com/dorrella/ubiquitous-waffle/service/types"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testContext struct {
	app    *types.App
	ctx    context.Context
	custDb *custdb.TestCustDb
}

// creates a http server and mimics routing for /api/customers
// should we do the whole /api route?
func (tc *testContext) initServer() *httptest.Server {
	r := mux.NewRouter()
	s := r.PathPrefix("/api/customer").Subrouter()
	InitRouter(s, tc.app)
	return httptest.NewServer(r)
}

// unmarshals CustomerResponse fails test on error
func (tc *testContext) parseCustomer(resp *http.Response, t *testing.T) *types.CustResp {
	defer resp.Body.Close()
	body_bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("unexpected error reading resp %s", err.Error())
	}
	cust_resp := types.CustResp{}
	err = json.Unmarshal(body_bytes, &cust_resp)
	if err != nil {
		t.Fatalf("unexpected error unmarshaling resp %s", err.Error())
	}
	return &cust_resp
}

// unmarshals CustomerResponse fails test on error
func (tc *testContext) parseCustomerList(resp *http.Response, t *testing.T) *types.CustList {
	defer resp.Body.Close()
	body_bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("unexpected error reading resp %s", err.Error())
	}
	cust_resp := types.CustList{}
	err = json.Unmarshal(body_bytes, &cust_resp)
	if err != nil {
		t.Fatalf("unexpected error unmarshaling resp %s", err.Error())
	}
	return &cust_resp
}

// creates a new testing context and clears the customerdb
func newTestContext() *testContext {
	app := &types.App{
		Config: conf.TestConfig(),
		Log:    logging.InitTestLogging(),
	}
	ctx := context.Background()
	database.InitPool(ctx, app)
	app.Db = database.GetDB()
	tc := &testContext{
		app:    app,
		ctx:    ctx,
		custDb: &custdb.TestCustDb{custdb.CustDb{app.Db}},
	}
	tc.custDb.Reset()
	return tc
}
