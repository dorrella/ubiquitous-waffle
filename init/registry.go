// registry of data migrations
package main

import (
	mig1 "github.com/dorrella/ubiquitous-waffle/init/migrations/000001"
)

// registers data migrattion interfaces with schema version numbers
func RegisterMigrations(m *MigController) {
	m.Register(1, &mig1.CreateSomething{})
}
