package main

import (
	"github.com/dorrella/ubiquitous-waffle/service/types"
	"net/http"
)

type Liveness struct {
	App *types.App
}

func (l *Liveness) LivenessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (l *Liveness) ReadinessHandler(w http.ResponseWriter, r *http.Request) {
	err := l.App.Db.Ping()
	if err == nil {
		w.WriteHeader(http.StatusOK)
	} else {
		//maybe?
		w.WriteHeader(http.StatusPreconditionFailed)
	}

}
