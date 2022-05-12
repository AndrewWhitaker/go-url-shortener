package integration

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"os"
	"testing"
	"url-shortener/db"
	"url-shortener/server"
	"url-shortener/test/helpers"

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

	container, sqlDB, err := helpers.CreateTestContainer(ctx, "testdb")

	if err != nil {
		log.Fatal(err)
	}

	gormDB, err := db.ConnectDatabase(sqlDB)

	if err != nil {
		log.Fatal(err)
	}

	server := server.SetupServer(
		&server.ServerConfig{DB: gormDB},
	)

	testContext := &ApiTestContext{
		db:     sqlDB,
		server: server,
	}

	testContext.cleanup = func() {
		sqlDB.Close()
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
