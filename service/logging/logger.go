package logging

import (
	"context"
	conf "github.com/dorrella/ubiquitous-waffle/service/config"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	slog "log/slog"
)

type Logger interface {
	Info(ctx context.Context, msg string, args ...interface{})
	InitOtelLogging(config *conf.Config)
}

type logger struct {
	otel *slog.Logger
}

type testLogger struct {
}

func InitLogging() Logger {
	return &logger{}
}

func (l *logger) InitOtelLogging(config *conf.Config) {
	//should already have been called, but just in case
	InitLogging()
	l.otel = otelslog.NewLogger(config.Service.Name)
}

// Todo other log levels
func (l *logger) Info(ctx context.Context, msg string, args ...interface{}) {
	slog.InfoContext(ctx, msg, args...)

	if l.otel != nil {
		l.otel.InfoContext(ctx, msg, args...)
	}
}

func InitTestLogging() Logger {
	return &testLogger{}
}

func (l *testLogger) InitOtelLogging(config *conf.Config) {}

// Todo other log levels
func (l *testLogger) Info(ctx context.Context, msg string, args ...interface{}) {}
