// common types for database and api
package types

import (
	conf "github.com/dorrella/ubiquitous-waffle/service/config"
	"github.com/dorrella/ubiquitous-waffle/service/logging"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

// common apt context
type App struct {
	Config    *conf.Config
	Db        *sqlx.DB
	Telemetry *Telemetry
	Log       logging.Logger
}

// wrapper around otel Tracer and Meter
// here to prevent circual import with otel lib
//
// todo use interface instead
type Telemetry struct {
}

func (t *Telemetry) GetTracer(name string) trace.Tracer {
	return otel.Tracer(name)
}

func (t *Telemetry) GetMeter(name string) metric.Meter {
	return otel.Meter(name)
}
