package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

	suite.DB = db

	suite.Server = httptest.NewServer(
		SetupServer(&ServerConfig{
			DB: db,
		}),
	)

	suite.Cleanup = func() {
		db.Close()
		container.Terminate(ctx)
		suite.Server.Close()
	}
}

func (suite *ApiTestSuite) TearDownSuite() {
	suite.Cleanup()
}

func (suite *ApiTestSuite) BeforeTest(suiteName, testName string) {
	db := suite.DB
	t := suite.T()

	_, err := db.Exec("TRUNCATE TABLE short_urls")

	if err != nil {
		t.Fatal("Failed to truncate short_urls table:", err)
	}

	_, err = db.Exec("ALTER SEQUENCE short_urls_id_seq RESTART")

	if err != nil {
		t.Fatal("Failed to reset short_urls_id_seq", err)
	}
}

func (suite *ApiTestSuite) TestCreateWithNewShortUrlReturns201WithShortUrl() {
	t := suite.T()
	testServer := suite.Server
	assert := suite.Assert()

	// Create first shortened url:
	result, err := createShortUrl("https://www.cloudflare.com", testServer.URL)

	if err != nil {
		t.Fatal(err)
	}

	expectedUrl := fmt.Sprintf("%s/1", testServer.URL)

	assert.Equal(201, result.Response.StatusCode)
	assert.Equal(expectedUrl, result.Data["short_url"])
}

func (suite *ApiTestSuite) TestCreateWithExistingLongUrlReturns200WithExistingShortUrl() {
	t := suite.T()
	testServer := suite.Server
	assert := suite.Assert()

	// Create first shortened url:
	result, err := createShortUrl("https://www.cloudflare.com", testServer.URL)

	if err != nil {
		t.Fatal(err)
	}

	resp := result.Response
	data := result.Data

	assert.Equal(201, resp.StatusCode, "Failed to create first URL")

	shortUrl := data["short_url"]

	// Create second shortened url, should get a 200 with existing shortened url:
	result, err = createShortUrl("https://www.cloudflare.com", testServer.URL)

	if err != nil {
		t.Fatal(err)
	}

	resp = result.Response
	data = result.Data

	assert.Equal(200, resp.StatusCode)
	assert.Equal(shortUrl, data["short_url"])
}

func (suite *ApiTestSuite) TestCreateWithInvalidJsonReturns400WithErrorMessage() {
	t := suite.T()
	testServer := suite.Server
	assert := suite.Assert()

	result, err := makeCreateRequest(map[string]interface{}{
		"foo": "bar",
	}, testServer.URL)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(400, result.Response.StatusCode)

	errors := result.Data["errors"]

	assert.NotNil(errors)
}

func (suite *ApiTestSuite) TestCreateWithSlugReturns201CreatedWithShortUrl() {
	t := suite.T()
	testServer := suite.Server
	assert := suite.Assert()

	result, err := createShortUrlWithSlug("https://www.cloudflare.com", "cf", testServer.URL)

	if err != nil {
		t.Fatal(err)
	}

	expectedUrl := fmt.Sprintf("%s/cf", testServer.URL)

	resp := result.Response
	data := result.Data

	assert.Equal(201, resp.StatusCode)
	assert.Equal(expectedUrl, data["short_url"])
}

type CreateShortUrlResult struct {
	Data     map[string]interface{}
	Response *http.Response
}

func createShortUrl(longUrl, url string) (*CreateShortUrlResult, error) {
	return makeCreateRequest(map[string]interface{}{
		"long_url": longUrl,
	}, url)
}

func createShortUrlWithSlug(longUrl, slug, url string) (*CreateShortUrlResult, error) {
	return makeCreateRequest(map[string]interface{}{
		"long_url": longUrl,
		"slug":     slug,
	}, url)
}

func makeCreateRequest(data map[string]interface{}, url string) (*CreateShortUrlResult, error) {
	postBody, err := json.Marshal(data)

	if err != nil {
		return nil, err
	}

	resp, err := http.Post(
		fmt.Sprintf("%s/shorturls", url),
		"application/json",
		bytes.NewBuffer(postBody),
	)
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)

	fmt.Printf("RESPONSE: %s\n", string(bytes))

	if err != nil {
		return nil, err
	}

	var responseBody map[string]interface{}

	err = json.Unmarshal([]byte(bytes), &responseBody)

	if err != nil {
		return nil, err
	}

	return &CreateShortUrlResult{
		Data:     responseBody,
		Response: resp,
	}, nil
}

func TestMain(t *testing.T) {
	suite.Run(t, new(ApiTestSuite))
}
