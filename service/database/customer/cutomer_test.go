package customer

import (
	"context"
	conf "github.com/dorrella/ubiquitous-waffle/service/config"
	"github.com/dorrella/ubiquitous-waffle/service/database"
	"github.com/dorrella/ubiquitous-waffle/service/logging"
	"github.com/dorrella/ubiquitous-waffle/service/types"
)

// wrapper around testing objects
type TestContext struct {
	app    *types.App
	custDb *TestCustDb
	ctx    context.Context
}

// creates a testing config and resets
// the customer table
func NewTestContext() *TestContext {
	ctx := context.Background()
	app := &types.App{
		Config: conf.TestConfig(),
		Log:    logging.InitTestLogging(),
	}
	database.InitPool(ctx, app)
	app.Db = database.GetDB()
	tc := &TestContext{
		app:    app,
		custDb: &TestCustDb{CustDb{app.Db}},
		ctx:    ctx,
	}
	tc.custDb.Reset()
	return tc
}
