package main

import (
	"net/http"
)

func LivenessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func ReadinessHandler(w http.ResponseWriter, r *http.Request) {
	// todo ping db?
	w.WriteHeader(http.StatusOK)
}
