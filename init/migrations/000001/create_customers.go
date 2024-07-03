// sample data migration
package mig000001

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type CreateSomething struct{}

func (c *CreateSomething) Run(db *sqlx.DB) {
	fmt.Println("example")
}
