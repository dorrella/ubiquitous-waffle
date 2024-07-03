package otel

import (
	"github.com/exaring/otelpgx"
)

func GetDbTracer() *otelpgx.Tracer {
	return otelpgx.NewTracer()
}
