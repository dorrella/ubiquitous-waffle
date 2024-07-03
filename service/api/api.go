package api

import (
	"github.com/dorrella/ubiquitous-waffle/service/api/customer"
	"github.com/dorrella/ubiquitous-waffle/service/types"
	"github.com/gorilla/mux"
)

// router of routers
func InitRouter(r *mux.Router, app *types.App) {
	customer.InitRouter(r.PathPrefix("/customer").Subrouter(), app)
}
