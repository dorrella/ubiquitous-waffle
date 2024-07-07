package main

import (
	"github.com/dorrella/ubiquitous-waffle/service/types"
	"net/http"
)

// handler for liveness and readiness probes
type Liveness struct {
	App *types.App
}

// http up
func (l *Liveness) LivenessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// database ready
func (l *Liveness) ReadinessHandler(w http.ResponseWriter, r *http.Request) {
	err := l.App.Db.Ping()
	if err == nil {
		w.WriteHeader(http.StatusOK)
	} else {
		//maybe?
		w.WriteHeader(http.StatusPreconditionFailed)
	}

}
