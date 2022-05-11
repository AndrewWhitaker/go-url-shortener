package integrationtests

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/maxatome/go-testdeep/helpers/tdhttp"
	"github.com/maxatome/go-testdeep/td"
	"github.com/stretchr/testify/suite"
)

type createSuite struct {
	suite.Suite
}

func TestCreate(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

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
					 "slug": "$slug"
				 }`,
				td.Tag("shortUrl", td.Re("http:\\/\\/example\\.com\\/([A-Za-z0-9]{8})")),
				td.Tag("slug", td.Re("[A-Za-z0-9]{8}")),
			),
		)
}

func (suite *createSuite) TestCreateWithExistingLongUrlReturns200() {
	t := suite.T()
	testServer := TestContext.server

	testAPI := tdhttp.NewTestAPI(t, testServer)

	var slug string
	var shortUrl string

	// Initial POST should return a 201 CREATED
	testAPI.PostJSON("/api/v1/shorturls", gin.H{"long_url": "https://www.cloudflare.com"}).
		CmpStatus(http.StatusCreated).
		CmpJSONBody(
			td.JSON(
				`{
				   "short_url": "$shortUrl",
					 "slug": "$slug"
				 }`,
				td.Tag("slug", td.Catch(&slug, td.Re("[A-Za-z0-9]{8}"))),
				td.Tag("shortUrl", td.Catch(&shortUrl, td.Ignore())),
			),
		)

	// Second POST should return a 200 OK with information about existing long URL
	testAPI.PostJSON("/api/v1/shorturls", gin.H{"long_url": "https://www.cloudflare.com"}).
		CmpStatus(http.StatusOK).
		CmpJSONBody(
			td.JSON(
				`{
				   "short_url": "$shortUrl",
					 "slug": "$slug"
				 }`,
				td.Tag("shortUrl", shortUrl),
				td.Tag("slug", slug),
			),
		)
}

func (suite *createSuite) TestCreateWithExistingSlugReturns409() {
	t := suite.T()
	testServer := TestContext.server

	testAPI := tdhttp.NewTestAPI(t, testServer)

	testAPI.PostJSON(
		"/api/v1/shorturls",
		gin.H{
			"long_url": "https://www.cloudflare.com",
			"slug":     "cf",
		},
	).
		CmpStatus(http.StatusCreated).
		CmpJSONBody(
			td.JSON(
				`{
				   "short_url": "$shortUrl",
			     "slug": "$slug"
			   }`,
				td.Tag("slug", "cf"),
				td.Tag("shortUrl", td.Ignore()),
			),
		)

	testAPI.PostJSON(
		"/api/v1/shorturls",
		gin.H{
			"long_url": "https://www.stackoverflow.com",
			"slug":     "cf",
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
