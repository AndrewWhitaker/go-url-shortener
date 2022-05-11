package integrationtests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/maxatome/go-testdeep/helpers/tdhttp"
	"github.com/maxatome/go-testdeep/td"
	"github.com/stretchr/testify/suite"
)

type clicksSuite struct {
	suite.Suite
}

func TestClicks(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	suite.Run(t, new(clicksSuite))
}

func (suite *clicksSuite) BeforeTest(suiteName, testName string) {
	TestContext.BeforeTest()
}

func (suite *clicksSuite) TestGetClicksWithValidSlugReturns200() {
	t := suite.T()

	testServer := TestContext.server
	testAPI := tdhttp.NewTestAPI(t, testServer)

	var slug string

	testAPI.PostJSON("/api/v1/shorturls", gin.H{"long_url": "https://www.cloudflare.com"}).
		CmpStatus(http.StatusCreated).
		CmpJSONBody(
			td.JSON(
				`{
				   "short_url": "$shortUrl",
					 "slug": "$slug"
				 }`,
				td.Tag("slug", td.Catch(&slug, td.Ignore())),
				td.Tag("shortUrl", td.Ignore()),
			),
		)

	for i := 0; i < 5; i++ {
		testAPI.Get(fmt.Sprintf("/%s", slug)).
			CmpStatus(http.StatusMovedPermanently)
	}

	testAPI.Get(fmt.Sprintf("/api/v1/shorturls/%s/clicks?time_period=ALL_TIME", slug)).
		CmpStatus(http.StatusOK).
		CmpJSONBody(
			td.JSON(
				`{
				   "count": 5,
					 "time_period": "ALL_TIME"
				 }`,
			),
		)
}

func (suite *clicksSuite) TestGetClicksWithInvalidSlugReturns404() {
	t := suite.T()

	testServer := TestContext.server
	testAPI := tdhttp.NewTestAPI(t, testServer)

	testAPI.Get("/api/v1/shorturls/invalid/clicks?time_period=ALL_TIME").
		CmpStatus(http.StatusNotFound).
		CmpJSONBody(
			td.JSON(
				`{
				   "errors": [
					   {
						   "field": "Slug",
							 "reason": "not found"
						 }
					 ]
				 }`,
			),
		)
}
