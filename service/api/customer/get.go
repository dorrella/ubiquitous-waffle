package customer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dorrella/ubiquitous-waffle/service/types"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// gets user by id
func (ch *CustHandler) GetCustomer(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()
	id_str := mux.Vars(r)["id"]
	customer_id, err := strconv.Atoi(id_str)
	if err != nil {
		ch.handleErr(ctx, w, r, err)
		return
	}
	cust, err := ch.CustDb.GetCustomer(ctx, int64(customer_id))
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
	msg := fmt.Sprintf("got customer: %d", cust.Id)
	ch.App.Log.Info(ctx, msg)
	fmt.Fprintf(w, string(resp))
}
