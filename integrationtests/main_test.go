package integrationtests

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"os"
	"testing"
	"url-shortener/server"
	"url-shortener/testhelpers"

	"github.com/gin-gonic/gin"
)

var TestContext *ApiTestContext

func TestMain(m *testing.M) {
	flag.Parse()

	if testing.Short() {
		return
	}

	TestContext = BuildTestContext()
	defer TestContext.cleanup()

	os.Exit(m.Run())
}

func BuildTestContext() *ApiTestContext {
	ctx := context.Background()

	container, db, err := testhelpers.CreateTestContainer(ctx, "testdb")

	if err != nil {
		log.Fatal(err)
	}

	server := server.SetupServer(
		&server.ServerConfig{DB: db},
	)

	testContext := &ApiTestContext{
		db:     db,
		server: server,
	}

	testContext.cleanup = func() {
		db.Close()
		container.Terminate(ctx)
	}

	return testContext
}

type ApiTestContext struct {
	db      *sql.DB
	server  *gin.Engine
	cleanup func()
}

func (ctx *ApiTestContext) BeforeTest() {
	db := ctx.db

	_, err := db.Exec("TRUNCATE TABLE short_urls CASCADE")

	if err != nil {
		log.Fatal("Failed to truncate short_urls table:", err)
	}

	_, err = db.Exec("ALTER SEQUENCE short_urls_id_seq RESTART")

	if err != nil {
		log.Fatal("Failed to reset short_urls_id_seq", err)
	}
}
