package otel

import (
	conf "github.com/dorrella/ubiquitous-waffle/service/config"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
)

// add tracing middleware, but only if configured
func GetMuxMiddleware(router *mux.Router, config *conf.Config) {
	if config.Telemetry.Enabled && config.Telemetry.Tracing {
		//otpional tracing middleware
		router.Use(otelmux.Middleware(config.Service.Name))
	}
}
