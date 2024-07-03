package customer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dorrella/ubiquitous-waffle/service/types"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"io"
	"net/http"
)

// updates metics for updated customers
func (ch *CustHandler) updateCustMetric(ctx context.Context) {
	if !ch.App.Config.Telemetry.Enabled || !ch.App.Config.Telemetry.Metrics {
		//not enabled
		return
	}
	cust_attr := attribute.Int("customers.value", 1)
	ch.CustsUpdated.Add(ctx, 1, metric.WithAttributes(cust_attr))

}

// updates a customer using a json object
func (ch *CustHandler) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()
	defer r.Body.Close()
	body_bytes, err := io.ReadAll(r.Body)
	if err != nil {
		ch.handleErr(ctx, w, r, err)
		return
	}
	//should probably wrap the req in a json object
	cust := &types.Customer{}
	err = json.Unmarshal(body_bytes, cust)
	if err != nil {
		ch.handleErr(ctx, w, r, err)
		return
	}
	cust, err = ch.CustDb.UpdateCustomer(ctx, cust, modified_by)
	if err != nil {
		ch.handleErr(ctx, w, r, err)
		return
	}
	resp, err := json.Marshal(types.CustResp{Customer: cust})
	if err != nil {
		ch.handleErr(ctx, w, r, err)
		return
	}
	msg := fmt.Sprintf("updated customer: %d", cust.Id)
	ch.App.Log.Info(ctx, msg)
	ch.updateCustMetric(ctx)
	fmt.Fprintf(w, string(resp))
}
