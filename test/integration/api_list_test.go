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

type listSuite struct {
	suite.Suite
}

func TestList(t *testing.T) {
	suite.Run(t, new(listSuite))
}

func (suite *listSuite) BeforeTest(suiteName, testName string) {
	TestContext.BeforeTest()
}

func (suite *listSuite) TestListReturns200() {
	t := suite.T()
	testServer := TestContext.server

	testAPI := tdhttp.NewTestAPI(t, testServer)

	var slug string
	var shortUrl string
	var longUrl string
	var createdAt time.Time

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

	testAPI.Get("/api/v1/shorturls").
		CmpStatus(http.StatusOK).
		CmpJSONBody(
			td.JSON(
				`[
				   {
						 "short_url": "$shortUrl",
						 "slug": "$slug",
						 "long_url": "$longUrl",
						 "expires_on": "$expiresOn",
						 "created_at": "$createdAt"
					 }
				 ]`,
				td.Tag("shortUrl", shortUrl),
				td.Tag("slug", slug),
				td.Tag("longUrl", longUrl),
				td.Tag("expiresOn", td.Nil()),
				td.Tag("createdAt", td.Smuggle(parseDateTime, createdAt)),
			),
		)
}
