package customer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dorrella/ubiquitous-waffle/service/types"
	"net/http"
)

// Gets user from unique email
func (ch *CustHandler) GetFromEmail(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()
	email := r.URL.Query().Get("email")
	if email == "" {
		err := fmt.Errorf("no email supplied")
		ch.handleErr(ctx, w, r, err)
		return
	}
	cust, err := ch.CustDb.FindByEmail(ctx, email)
	if err != nil {
		ch.handleErr(ctx, w, r, err)
		return
	} else if cust == nil {
		err = fmt.Errorf("customer email %s not found", email)
		ch.handleErr(ctx, w, r, err)
		return
	}

	resp, err := json.Marshal(types.CustResp{Customer: cust})
	if err != nil {
		ch.handleErr(ctx, w, r, err)
		return
	}
	msg := fmt.Sprintf("got customer by email: %d", cust.Id)
	ch.App.Log.Info(ctx, msg)
	fmt.Fprintf(w, string(resp))
}
