//database connection with pgx and sqlx

package database

import (
	"context"
	"fmt"
	otel "github.com/dorrella/ubiquitous-waffle/service/otel"
	"github.com/dorrella/ubiquitous-waffle/service/types"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"sync"
	"time"
)

var pgx_pool *pgxpool.Pool
var once sync.Once

// initializes pool. todo clean shutdown
func InitPool(ctx context.Context, app *types.App) error {
	//should only be called once, but just in case
	once.Do(func() {
		var postgres_url string = app.Config.GetPostgresUrl()

		pgx_config, err := pgxpool.ParseConfig(postgres_url)
		if err != nil {
			app.Log.Info(ctx, err.Error())
			return
		}
		//check for tracing
		if app.Config.Telemetry.Enabled && app.Config.Telemetry.Tracing {
			app.Log.Info(ctx, "enabling db traces")
			pgx_config.ConnConfig.Tracer = otel.GetDbTracer()
		}
		pool, err := pgxpool.NewWithConfig(context.Background(), pgx_config)
		if err != nil {
			app.Log.Info(ctx, err.Error())
			return
		}
		pgx_pool = pool
	})
	if pgx_pool == nil {
		return fmt.Errorf("failed to init db")
	}
	return nil
}

// gets a new sqlx interface from a pgxpool
func GetDB() *sqlx.DB {
	//sqlx requires a db/sql interface
	db := stdlib.OpenDBFromPool(pgx_pool)
	return sqlx.NewDb(db, "pgx")
}

// probably never called, but clean shutdown
func Close() {
	if pgx_pool != nil {
		pgx_pool.Close()
	}
}

// timestamp stripped of timezone
func TimeStamp() time.Time {
	return time.Now().UTC()
}
