package customer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dorrella/ubiquitous-waffle/service/types"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"net/http"
	"strconv"
)

// updates metics for deleted customers
func (ch *CustHandler) delCustMetric(ctx context.Context) {
	if !ch.App.Config.Telemetry.Enabled || !ch.App.Config.Telemetry.Metrics {
		//not enabled
		return
	}
	cust_attr := attribute.Int("customers.value", 1)
	ch.CustsDeleted.Add(ctx, 1, metric.WithAttributes(cust_attr))

}

// deletes customer by user id
func (ch *CustHandler) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()
	id_str := mux.Vars(r)["id"]
	customer_id, err := strconv.Atoi(id_str)
	if err != nil {
		ch.handleErr(ctx, w, r, err)
		return
	}
	cust, err := ch.CustDb.DeleteCustomer(ctx, int64(customer_id), modified_by)
	if err != nil {
		ch.handleErr(ctx, w, r, err)
		return
	} else if cust == nil {
		err = fmt.Errorf("customer %d not found", customer_id)
		ch.handleErr(ctx, w, r, err)
		return
	}

	resp, err := json.Marshal(types.CustResp{Customer: cust})
	if err != nil {
		ch.handleErr(ctx, w, r, err)
		return
	}

	msg := fmt.Sprintf("delted customer: %d", cust.Id)
	ch.App.Log.Info(ctx, msg)
	ch.delCustMetric(ctx)
	fmt.Fprintf(w, string(resp))
}
