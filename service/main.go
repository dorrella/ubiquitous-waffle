package main

import (
	"context"
	"errors"
	"github.com/dorrella/ubiquitous-waffle/service/api"
	conf "github.com/dorrella/ubiquitous-waffle/service/config"
	"github.com/dorrella/ubiquitous-waffle/service/database"
	log "github.com/dorrella/ubiquitous-waffle/service/logging"
	otel "github.com/dorrella/ubiquitous-waffle/service/otel"
	"github.com/dorrella/ubiquitous-waffle/service/types"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func setupRouter(app *types.App) *mux.Router {
	r := mux.NewRouter()
	otel.GetMuxMiddleware(r, app.Config)

	liveness := Liveness{app}
	r.HandleFunc("/live", liveness.LivenessHandler)
	r.HandleFunc("/ready", liveness.ReadinessHandler)
	api.InitRouter(r.PathPrefix("/api").Subrouter(), app)
	return r
}

func main() {
	// Handle SIGINT (CTRL+C) gracefully.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	//global interface
	app := &types.App{
		Log:       log.InitLogging(),
		Config:    conf.LoadConfig(),
		Telemetry: &types.Telemetry{},
	}
	app.Log.Info(ctx, "initialized logging")

	if app.Config.Telemetry.Enabled {
		// Set up OpenTelemetry.
		app.Log.Info(ctx, "telemetry enabled")
		otelShutdown, err := otel.SetupOTelSDK(ctx, app)
		if err != nil {
			app.Log.Info(ctx, err.Error())
			return
		}
		// Handle shutdown properly so nothing leaks.
		defer func() {
			err = errors.Join(err, otelShutdown(context.Background()))
		}()
	}

	err := database.InitPool(ctx, app)
	if err != nil {
		//already logged the real error
		return
	}
	app.Db = database.GetDB()

	// setup routes
	app.Log.Info(ctx, "initializing routing")
	r := setupRouter(app)

	srv := &http.Server{
		Handler:      r,
		Addr:         ":8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// start http server and wait for error
	srvErr := make(chan error, 1)
	go func() {
		app.Log.Info(ctx, "luanching http server")
		srvErr <- srv.ListenAndServe()
	}()

	// Wait for interruption.
	select {
	case err := <-srvErr:
		// Error when starting HTTP server.
		app.Log.Info(ctx, err.Error())
		return
	case <-ctx.Done():
		// Wait for first CTRL+C.
		// Stop receiving signal notifications as soon as possible.
		stop()
	}

	// When Shutdown is called, ListenAndServe immediately returns ErrServerClosed.
	err = srv.Shutdown(context.Background())
	app.Log.Info(ctx, err.Error())
}
