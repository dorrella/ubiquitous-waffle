package otel

import (
	"context"
	"errors"
	"fmt"
	"github.com/dorrella/ubiquitous-waffle/service/types"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type otelContext struct {
	ctx           context.Context
	app           *types.App
	res           *resource.Resource
	conn          *grpc.ClientConn
	err           error
	shutdown      func(ctx context.Context) error
	shutdownFuncs []func(context.Context) error
}

// wraps all errors for the shutdown path
func (octx *otelContext) handleErr(inErr error) {
	octx.err = errors.Join(inErr, octx.shutdown(octx.ctx))
}

// setupOTelSDK bootstraps the OpenTelemetry pipeline.
// If it does not return an error, make sure to call shutdown for proper cleanup.
func SetupOTelSDK(ctx context.Context, app *types.App) (shutdown func(context.Context) error, err error) {
	var octx *otelContext = &otelContext{
		app: app,
		ctx: ctx,
	}

	app.Log.Info(ctx, "initializing open telemetry")
	octx.conn, octx.err = initGrpc(octx)
	if octx.err != nil {
		return func(ctx context.Context) error { return nil }, err
	}
	app.Log.Info(ctx, "connect to grpc collector initializing open telemetry")
	octx.res, octx.err = resource.New(ctx,
		resource.WithAttributes(
			// The service name used to display traces in backends
			semconv.ServiceNameKey.String(octx.app.Config.Service.Name),
		),
	)
	if err != nil {
		return func(ctx context.Context) error { return nil }, err
	}

	// shutdown calls cleanup functions registered via shutdownFuncs.
	// The errors from the calls are joined.
	// Each registered cleanup will be invoked once.
	octx.shutdown = func(ctx context.Context) error {
		var err error
		for _, fn := range octx.shutdownFuncs {
			octx.err = errors.Join(err, fn(ctx))
		}
		octx.shutdownFuncs = nil
		return err
	}

	// Set up propagator.
	prop := newPropagator()
	otel.SetTextMapPropagator(prop)
	app.Log.Info(ctx, "initializd propagator")
	initTelemetry(octx)
	return octx.shutdown, octx.err
}

// init each component
func initTelemetry(octx *otelContext) {
	// Set up logger provider.
	if octx.app.Config.Telemetry.Logs {
		octx.app.Log.Info(octx.ctx, "initializing open telemetry logger")
		loggerProvider, err := newLoggerProvider(octx)
		if err != nil {
			octx.handleErr(err)
			return
		}
		octx.shutdownFuncs = append(octx.shutdownFuncs, loggerProvider.Shutdown)
		global.SetLoggerProvider(loggerProvider)
		octx.app.Log.InitOtelLogging(octx.app.Config)
	}

	// Set up trace provider.
	if octx.app.Config.Telemetry.Tracing {
		octx.app.Log.Info(octx.ctx, "initializing open telemetry tracer")
		tracerProvider, err := newTracerProvider(octx)
		if err != nil {
			octx.handleErr(err)
			return
		}
		octx.shutdownFuncs = append(octx.shutdownFuncs, tracerProvider.Shutdown)
		otel.SetTracerProvider(tracerProvider)
	}

	// Set up meter provider.
	if octx.app.Config.Telemetry.Metrics {
		octx.app.Log.Info(octx.ctx, "initializing open telemetry meter")
		meterProvider, err := newMeterProvider(octx)
		if err != nil {
			octx.handleErr(err)
			return
		}
		octx.shutdownFuncs = append(octx.shutdownFuncs, meterProvider.Shutdown)
		otel.SetMeterProvider(meterProvider)
	}
}

func initGrpc(octx *otelContext) (*grpc.ClientConn, error) {
	// It connects the OpenTelemetry Collector through local gRPC connection.
	url := fmt.Sprintf("%s:%d", octx.app.Config.Telemetry.Collector.Url, octx.app.Config.Telemetry.Collector.Port)
	msg := fmt.Sprintf("using collector url: %s", url)
	octx.app.Log.Info(octx.ctx, msg)
	conn, err := grpc.NewClient(url,
		// Note the use of insecure transport here. TLS is recommended in production.
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	return conn, err
}

// setupup propogator for context
func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

// Initializes an OTLP exporter, and configures the corresponding trace provider.
func newTracerProvider(octx *otelContext) (*sdktrace.TracerProvider, error) {
	traceExporter, err := otlptracegrpc.New(octx.ctx, otlptracegrpc.WithGRPCConn(octx.conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(octx.res),
		sdktrace.WithSpanProcessor(bsp),
	)

	otel.SetTextMapPropagator(propagation.TraceContext{})
	return tracerProvider, nil
}

// Initializes Meter progiver over GRPC
func newMeterProvider(octx *otelContext) (*sdkmetric.MeterProvider, error) {
	metricExporter, err := otlpmetricgrpc.New(octx.ctx, otlpmetricgrpc.WithGRPCConn(octx.conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics exporter: %w", err)
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter)),
		sdkmetric.WithResource(octx.res),
	)
	return meterProvider, nil
}

// Initializes Logger progiver over GRPC
func newLoggerProvider(octx *otelContext) (*sdklog.LoggerProvider, error) {
	logExporter, err := otlploggrpc.New(octx.ctx, otlploggrpc.WithGRPCConn(octx.conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics exporter: %w", err)
	}

	logProvider := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(sdklog.NewBatchProcessor(logExporter)),
		sdklog.WithResource(octx.res),
	)
	return logProvider, nil
}
