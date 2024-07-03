package main

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

// interface for running migrations. all data migrations
// must implement this and use it as an entrypoint
type Migration interface {
	Run(db *sqlx.DB)
}

// migration context
type MigController struct {
	mig      *migrate.Migrate
	db       *sql.DB
	registry map[uint]Migration
}

// creates a migration controller and then registers
// the data migrations from the registry
func NewMigController(mig *migrate.Migrate, db *sql.DB) *MigController {
	m := &MigController{mig: mig,
		db:       db,
		registry: make(map[uint]Migration)}
	RegisterMigrations(m)
	return m
}

// register a data migration to a particular revision
func (m *MigController) Register(version uint, mig Migration) {
	m.registry[version] = mig
}

// bring db to latest schema migration, while running
// any registered data migrations
func (m *MigController) Update() {
	//sqlx interface to ser/der into classes
	db := sqlx.NewDb(m.db, "pgx")
	defer db.Close()
	for true {
		//step forward a single migration until we hit
		//a file not found error
		err := m.mig.Steps(1)
		if err != nil {
			if err.Error() == "file does not exist" {
				fmt.Println("finished schema migrations")
				return
			} else {
				//mystery error
				panic(err)
			}
		}

		//get version and check for a data migration
		version, dirty, err := m.mig.Version()
		fmt.Println(fmt.Sprintf("at version %v. dirty: %v", version, dirty))
		if err != nil {
			panic(err)
		}
		mig, ok := m.registry[version]
		if ok {
			fmt.Println("running migration")
			mig.Run(db)
		}
	}
}
