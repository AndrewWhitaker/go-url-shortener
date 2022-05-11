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
	"net/url"
	"testing"
	"url-shortener/server"
	"url-shortener/testhelpers"

	"github.com/stretchr/testify/suite"
)

type ErrorResponse struct {
	Errors []struct {
		Field  string `json:"field"`
		Reason string `json:"reason"`
	} `json:"errors"`
}

type apiTestSuite struct {
	suite.Suite
	db      *sql.DB
	server  *httptest.Server
	cleanup func()
}

func (suite *apiTestSuite) SetupSuite() {
	ctx := context.Background()
	container, db, err := testhelpers.CreateTestContainer(ctx, "testdb")

	if err != nil {
		suite.T().Fatal(err)
	}

	suite.db = db

	suite.server = httptest.NewServer(
		server.SetupServer(&server.ServerConfig{
			DB: db,
		}),
	)

	suite.cleanup = func() {
		db.Close()
		container.Terminate(ctx)
		suite.server.Close()
	}
}

func (suite *apiTestSuite) TearDownSuite() {
	suite.cleanup()
}

func (suite *apiTestSuite) BeforeTest(suiteName, testName string) {
	db := suite.db
	t := suite.T()

	_, err := db.Exec("TRUNCATE TABLE short_urls CASCADE")

	if err != nil {
		t.Fatal("Failed to truncate short_urls table:", err)
	}

	_, err = db.Exec("ALTER SEQUENCE short_urls_id_seq RESTART")

	if err != nil {
		t.Fatal("Failed to reset short_urls_id_seq", err)
	}
}

func (suite *apiTestSuite) TestGetClicks() {
	t := suite.T()
	testServer := suite.server
	assert := suite.Assert()

	result, err := createShortUrl("https://www.google.com", testServer.URL)

	if err != nil {
		t.Fatal(err)
	}

	shortUrl := result.data["short_url"]

	// Access the URL a few times to generate statistics
	for i := 0; i < 5; i++ {
		_, err = http.Get(fmt.Sprintf("%s", shortUrl))
		if err != nil {
			t.Fatal(err)
		}
	}

	slug := fmt.Sprintf("%v", result.data["slug"])

	type test struct {
		expectedStatus int
		expectedCount  int
		slug           string
	}

	tests := []test{
		{expectedStatus: 200, expectedCount: 5, slug: slug},
		{expectedStatus: 404, expectedCount: 0, slug: "invalid"},
	}

	for _, tc := range tests {
		clicksUrl, err := url.Parse(testServer.URL)

		if err != nil {
			t.Fatal(err)
		}

		clicksUrl.Path = fmt.Sprintf("api/v1/shorturls/%s/clicks", tc.slug)

		query := clicksUrl.Query()
		query.Add("time_period", "ALL_TIME")

		clicksUrl.RawQuery = query.Encode()

		resp, err := http.Get(clicksUrl.String())
		defer resp.Body.Close()

		bytes, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			t.Fatal(err)
		}

		var responseBody map[string]interface{}
		err = json.Unmarshal([]byte(bytes), &responseBody)

		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(tc.expectedStatus, resp.StatusCode)
		//assert.Equal(float64(tc.expectedCount), responseBody["count"])
		//assert.Equal("ALL_TIME", responseBody["time_period"])
	}
}

func (suite *apiTestSuite) TestGetClicksWithValidSlugAndTimePeriodReturns200AndClicks() {
	t := suite.T()
	testServer := suite.server
	assert := suite.Assert()

	result, err := createShortUrl("https://www.google.com", testServer.URL)

	if err != nil {
		t.Fatal(err)
	}

	shortUrl := result.data["short_url"]

	// Access url a few times
	for i := 0; i < 5; i++ {
		_, err = http.Get(fmt.Sprintf("%s", shortUrl))
		if err != nil {
			t.Fatal(err)
		}
	}

	slug := result.data["slug"]

	resp, err := http.Get(fmt.Sprintf("%s/api/v1/shorturls/%s/clicks?time_period=ALL_TIME", testServer.URL, slug))

	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("RESPONSE=%s\n", string(bytes))

	var responseBody map[string]interface{}
	err = json.Unmarshal([]byte(bytes), &responseBody)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(200, resp.StatusCode)
	assert.Equal(float64(5), responseBody["count"])
	assert.Equal("ALL_TIME", responseBody["time_period"])
}

type createShortUrlResult struct {
	data     map[string]interface{}
	response *http.Response
}

func createShortUrl(longUrl, url string) (*createShortUrlResult, error) {
	return makeCreateRequest(map[string]interface{}{
		"long_url": longUrl,
	}, url)
}

func createShortUrlWithSlug(longUrl, slug, url string) (*createShortUrlResult, error) {
	return makeCreateRequest(map[string]interface{}{
		"long_url": longUrl,
		"slug":     slug,
	}, url)
}

func makeCreateRequest(data map[string]interface{}, host string) (*createShortUrlResult, error) {
	postBody, err := json.Marshal(data)

	if err != nil {
		return nil, err
	}

	u, err := url.Parse(host)

	if err != nil {
		return nil, err
	}

	u.Path = "/api/v1/shorturls"

	resp, err := http.Post(u.String(), "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		return nil, err
	}

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

	return &createShortUrlResult{
		data:     responseBody,
		response: resp,
	}, nil
}

func TestMain(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	suite.Run(t, new(apiTestSuite))
}
