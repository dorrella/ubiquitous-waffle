package customer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dorrella/ubiquitous-waffle/service/types"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"net/http"
)

// updates metics for new customers
func (ch *CustHandler) newCustMetric(ctx context.Context) {
	if !ch.App.Config.Telemetry.Enabled || !ch.App.Config.Telemetry.Metrics {
		//not enabled
		return
	}
	cust_attr := attribute.Int("customers.value", 1)
	ch.CustsCreated.Add(ctx, 1, metric.WithAttributes(cust_attr))

}

// new customer router. leaves validation up to db
func (ch *CustHandler) NewCustomer(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()
	cust := &types.Customer{
		NamePrefix:  r.FormValue("name_pref"),
		NameFirst:   r.FormValue("name_first"),
		NameMiddle:  r.FormValue("name_middle"),
		NameLast:    r.FormValue("name_last"),
		NameSuffix:  r.FormValue("name_suffix"),
		Email:       r.FormValue("email"),
		PhoneNumber: r.FormValue("phone_number"),
	}
	cust, err := ch.CustDb.NewCustomer(ctx, cust, modified_by)
	if err != nil {
		ch.handleErr(ctx, w, r, err)
		return
	}

	// user created
	resp, err := json.Marshal(types.CustResp{Customer: cust})
	if err != nil {
		ch.handleErr(ctx, w, r, err)
		return
	}

	msg := fmt.Sprintf("created customer: %d", cust.Id)
	ch.App.Log.Info(ctx, msg)
	ch.newCustMetric(ctx)
	fmt.Fprintf(w, string(resp))

}
