package customer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dorrella/ubiquitous-waffle/service/types"
	"net/http"
	"strconv"
)

// wraps list error and sends CustList response
func (ch *CustHandler) handleListErr(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(http.StatusBadRequest)
	msg := fmt.Sprintf("error: %s", err.Error())
	ch.App.Log.Info(ctx, msg)
	resp, err := json.Marshal(types.CustList{Error: msg})
	if err != nil {
		msg := fmt.Sprintf("error marshaling json: %s", err.Error())
		ch.App.Log.Info(ctx, msg)
		return
	}
	fmt.Fprintf(w, string(resp))
}

// lists customers in batches of 25. get pages with next query
func (ch *CustHandler) ListCustomers(w http.ResponseWriter, r *http.Request) {
	var next_id int
	var err error
	ctx := context.TODO()
	next_str := r.URL.Query().Get("next")

	if next_str == "" {
		next_id = 0
	} else {
		next_id, err = strconv.Atoi(next_str)
		if err != nil {
			ch.handleListErr(ctx, w, r, err)
			return
		}
	}

	cust_list, new_id, err := ch.CustDb.ListCustomers(ctx, int64(next_id))
	resp, err := json.Marshal(types.CustList{Customers: cust_list, Next: new_id})
	if err != nil {
		ch.handleListErr(ctx, w, r, err)
		return
	}
	fmt.Fprintf(w, string(resp))
}
