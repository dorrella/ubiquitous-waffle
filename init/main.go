package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	conf "github.com/dorrella/ubiquitous-waffle/service/config"
	"github.com/golang-migrate/migrate/v4"
	pgxmig "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"time"
)

// checks for testing config
func initConfig() *conf.Config {
	testing := flag.Bool("testing", false, "enable test config for local db schema init")
	flag.Parse()
	if *testing {
		return conf.TestConfig()
	}
	return conf.LoadConfig()

}

// setup a pgxpool
func getDbPool(config *conf.Config) *pgxpool.Pool {
	var postgres_url string = config.GetPostgresUrl()
	pool, err := pgxpool.New(context.Background(), postgres_url)
	if err != nil {
		panic(err)
	}
	return pool
}

func waitDb(db *sql.DB) {
	var current = time.Now().UTC()
	var max = current.Add(3 * time.Minute)
	var err error
	fmt.Println("pinging database")
	for current.Before(max) {
		err = db.Ping()
		if err == nil {
			return
		}
		time.Sleep(10 * time.Second)
		current = time.Now()
	}
	panic(err)
}

// setup go gration object
func getGoMig(db *sql.DB) *migrate.Migrate {
	driver, err := pgxmig.WithInstance(db, &pgxmig.Config{})
	if err != nil {
		panic(err)
	}
	migrate, err := migrate.NewWithDatabaseInstance(
		"file://migration_files", "postgres", driver)
	if err != nil {
		panic(err)
	}
	return migrate
}

func main() {
	config := initConfig()
	pool := getDbPool(config)
	// go-migrate and sqlx both require the std interface
	db := stdlib.OpenDBFromPool(pool)
	waitDb(db)
	migrate := getGoMig(db)

	version, dirty, err := migrate.Version()
	fmt.Println(fmt.Sprintf("starting version %v. dirty: %v", version, dirty))
	if err != nil {
		if version == 0 && err.Error() == "no migration" {
			fmt.Println("first migration")
		} else {
			panic(err)
		}
	}
	controller := NewMigController(migrate, db)
	controller.Update()
	migrate.Close()
	db.Close()
	pool.Close()
}
