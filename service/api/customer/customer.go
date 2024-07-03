package customer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dorrella/ubiquitous-waffle/service/database/customer"
	"github.com/dorrella/ubiquitous-waffle/service/types"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel/metric"
	"net/http"
)

// todo pull user id from token of whomever called this api
// user used for auditing cusomer changes
const modified_by = 1

type CustHandler struct {
	App          *types.App
	CustDb       *customer.CustDb
	CustsCreated metric.Int64Counter
	CustsDeleted metric.Int64Counter
	CustsUpdated metric.Int64Counter
}

// hanldes error for customer responses
func (ch *CustHandler) handleErr(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(http.StatusBadRequest)
	msg := fmt.Sprintf("error: %s", err.Error())
	ch.App.Log.Info(ctx, msg)
	resp, err := json.Marshal(types.CustResp{Error: msg})
	if err != nil {
		msg := fmt.Sprintf("error marshaling json: %s", err.Error())
		ch.App.Log.Info(ctx, msg)
		return
	}
	fmt.Fprintf(w, string(resp))
}

// initializes route handlers
func InitRouter(r *mux.Router, app *types.App) {
	ch := &CustHandler{
		App:    app,
		CustDb: &customer.CustDb{app.Db},
	}
	r.HandleFunc("", ch.NewCustomer).Methods("POST")
	r.HandleFunc("/{id:[0-9]+}", ch.GetCustomer).Methods("GET")
	r.HandleFunc("/{id:[0-9]+}", ch.UpdateCustomer).Methods("PUT")
	r.HandleFunc("/{id:[0-9]+}", ch.DeleteCustomer).Methods("DELETE")
	r.HandleFunc("/list", ch.ListCustomers).Methods("GET")
	r.HandleFunc("/by_email", ch.GetFromEmail).Methods("GET")

	//optional metrics
	if app.Config.Telemetry.Enabled && app.Config.Telemetry.Metrics {
		//setup metrics panic on err
		var err error
		meter := app.Telemetry.GetMeter(app.Config.Service.Name)
		ch.CustsCreated, err = meter.Int64Counter("customers_created",
			metric.WithDescription("customers created by this instance"),
			metric.WithUnit("{customers}"))
		if err != nil {
			panic(err)
		}
		//customers deleted
		ch.CustsDeleted, err = meter.Int64Counter("customers_deleted",
			metric.WithDescription("customers deleted by this instance"),
			metric.WithUnit("{customers}"))
		if err != nil {
			panic(err)
		}
		//customers updated
		ch.CustsUpdated, err = meter.Int64Counter("customers_updated",
			metric.WithDescription("customers updated by this instance"),
			metric.WithUnit("{customers}"))
		if err != nil {
			panic(err)
		}
	}
}
