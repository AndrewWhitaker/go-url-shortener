package integration

import (
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/maxatome/go-testdeep/helpers/tdhttp"
	"github.com/maxatome/go-testdeep/td"
	"github.com/stretchr/testify/suite"
)

type createSuite struct {
	suite.Suite
}

func TestCreate(t *testing.T) {
	suite.Run(t, new(createSuite))
}

func (suite *createSuite) BeforeTest(suiteName, testName string) {
	TestContext.BeforeTest()
}

func (suite *createSuite) TestCreateWithNewShortUrlReturns201() {
	t := suite.T()
	testServer := TestContext.server

	testAPI := tdhttp.NewTestAPI(t, testServer)

	testAPI.PostJSON("/api/v1/shorturls", gin.H{"long_url": "https://www.cloudflare.com"}).
		CmpStatus(http.StatusCreated).
		CmpJSONBody(
			td.JSON(
				`{
				   "short_url": "$shortUrl",
					 "slug": "$slug",
					 "long_url": "$longUrl",
					 "expires_on": "$expiresOn",
					 "created_at": "$createdAt"
				 }`,
				td.Tag("shortUrl", td.Re("http:\\/\\/example\\.com\\/([A-Za-z0-9]{8})")),
				td.Tag("slug", td.Re("[A-Za-z0-9]{8}")),
				td.Tag("longUrl", "https://www.cloudflare.com"),
				td.Tag("expiresOn", td.Nil()),
				td.Tag("createdAt", td.Smuggle(parseDateTime, td.Between(testAPI.SentAt(), time.Now()))),
			),
		)
}

func (suite *createSuite) TestCreateWithExpirationDateReturns201() {
	t := suite.T()

	testServer := TestContext.server
	testAPI := tdhttp.NewTestAPI(t, testServer)

	expirationDateTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	testAPI.PostJSON(
		"/api/v1/shorturls",
		gin.H{
			"long_url":   "https://www.cloudflare.com",
			"expires_on": expirationDateTime.Format(time.RFC3339),
		},
	).
		CmpStatus(http.StatusCreated).
		CmpJSONBody(
			td.JSON(
				`{
				   "short_url": "$shortUrl",
					 "slug": "$slug",
					 "long_url": "$longUrl",
					 "expires_on": "$expiresOn",
					 "created_at": "$createdAt"
				 }`,
				td.Tag("shortUrl", td.Re("http:\\/\\/example\\.com\\/([A-Za-z0-9]{8})")),
				td.Tag("slug", td.Re("[A-Za-z0-9]{8}")),
				td.Tag("longUrl", "https://www.cloudflare.com"),
				td.Tag("expiresOn", td.Smuggle(parseDateTime, expirationDateTime)),
				td.Tag("createdAt", td.Smuggle(parseDateTime, td.Between(testAPI.SentAt(), time.Now()))),
			),
		)
}

func (suite *createSuite) TestCreateWithExistingLongUrlReturns200() {
	t := suite.T()
	testServer := TestContext.server

	testAPI := tdhttp.NewTestAPI(t, testServer)

	var slug string
	var shortUrl string
	var longUrl string
	var createdAt time.Time

	// Initial POST should return a 201 CREATED
	testAPI.PostJSON("/api/v1/shorturls", gin.H{"long_url": "https://www.cloudflare.com"}).
		CmpStatus(http.StatusCreated).
		CmpJSONBody(
			td.JSON(
				`{
				   "short_url": "$shortUrl",
					 "slug": "$slug",
					 "long_url": "$longUrl",
					 "expires_on": "$expiresOn",
					 "created_at": "$createdAt"
				 }`,
				td.Tag("slug", td.Catch(&slug, td.Re("[A-Za-z0-9]{8}"))),
				td.Tag("shortUrl", td.Catch(&shortUrl, td.Ignore())),
				td.Tag("longUrl", td.Catch(&longUrl, "https://www.cloudflare.com")),
				td.Tag("expiresOn", td.Nil()),
				td.Tag("createdAt", td.Smuggle(parseDateTime, td.Catch(&createdAt, td.Ignore()))),
			),
		)

	// Second POST should return a 200 OK with information about existing long URL
	testAPI.PostJSON("/api/v1/shorturls", gin.H{"long_url": "https://www.cloudflare.com"}).
		CmpStatus(http.StatusOK).
		CmpJSONBody(
			td.JSON(
				`{
				   "short_url": "$shortUrl",
					 "slug": "$slug",
					 "long_url": "$longUrl",
					 "expires_on": "$expiresOn",
					 "created_at": "$createdAt"
				 }`,
				td.Tag("shortUrl", shortUrl),
				td.Tag("slug", slug),
				td.Tag("longUrl", longUrl),
				td.Tag("expiresOn", td.Nil()),
				td.Tag("createdAt", td.Smuggle(parseDateTime, createdAt)),
			),
		)
}

func (suite *createSuite) TestCreateWithExistingSlugReturns409() {
	t := suite.T()

	testServer := TestContext.server
	testAPI := tdhttp.NewTestAPI(t, testServer)

	slug := "cf"

	testAPI.PostJSON(
		"/api/v1/shorturls",
		gin.H{
			"slug":     slug,
			"long_url": "https://www.cloudflare.com",
		},
	).
		CmpStatus(http.StatusCreated).
		CmpJSONBody(
			td.SuperJSONOf(`{"slug": "$slug"}`, td.Tag("slug", slug)),
		)

	testAPI.PostJSON(
		"/api/v1/shorturls",
		gin.H{
			"long_url": "https://www.stackoverflow.com",
			"slug":     slug,
		},
	).
		CmpStatus(http.StatusConflict).
		CmpJSONBody(
			td.JSON(
				`{
				   "errors": [
					   {
						   "field": "Slug",
							 "reason": "must be unique"
						 }
					 ],
				 }`,
			),
		)
}

func (suite *createSuite) TestCreateWithInvalidJSONReturns400() {
	t := suite.T()
	testServer := TestContext.server

	testAPI := tdhttp.NewTestAPI(t, testServer)

	testAPI.PostJSON("/api/v1/shorturls", gin.H{"foo": "bar"}).
		CmpStatus(http.StatusBadRequest).
		CmpJSONBody(
			td.JSON(
				`{
				   "errors": [
					   {
						   "field": "LongUrl",
							 "reason": "required",
						 }
					 ]
			   }`,
				td.Tag("slug", "cf"),
				td.Tag("shortUrl", td.Ignore()),
			),
		)
}
