package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ApiTestSuite struct {
	suite.Suite
	*sql.DB
	*httptest.Server
	Cleanup func()
}

func (suite *ApiTestSuite) SetupSuite() {
	ctx := context.Background()
	container, db, err := CreateTestContainer(ctx, "testdb")

	if err != nil {
		suite.T().Fatal(err)
	}

	config := ServerConfig{
		DB: db,
	}

	suite.Server = httptest.NewServer(SetupServer(&config))

	suite.Cleanup = func() {
		db.Close()
		container.Terminate(ctx)
		suite.Server.Close()
	}
}

func (suite *ApiTestSuite) TearDownSuite() {
	suite.Cleanup()
}

func (suite *ApiTestSuite) TestCreateWithNewShortUrlReturns201k() {
	t := suite.T()
	testServer := suite.Server

	postBody, err := json.Marshal(map[string]interface{}{
		"long_url": "https://www.cloudflare.com",
	})

	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.Post(
		fmt.Sprintf("%s/shorturls", testServer.URL),
		"application/json",
		bytes.NewBuffer(postBody),
	)

	defer resp.Body.Close()

	if err != nil {
		t.Fatal(err)
	}

	suite.Assert().Equal(201, resp.StatusCode)
}

func (suite *ApiTestSuite) TestCreateWithExistingLongUrlReturns409() {
	t := suite.T()
	testServer := suite.Server

	postBody, err := json.Marshal(map[string]interface{}{
		"long_url": "https://www.cloudflare.com",
	})

	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.Post(
		fmt.Sprintf("%s/shorturls", testServer.URL),
		"application/json",
		bytes.NewBuffer(postBody),
	)

	defer resp.Body.Close()

	if err != nil {
		t.Fatal(err)
	}

	suite.Assert().Equal(201, resp.StatusCode)
}

func TestMain(t *testing.T) {
	suite.Run(t, new(ApiTestSuite))
}
