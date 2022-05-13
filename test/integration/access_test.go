package integration

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/maxatome/go-testdeep/helpers/tdhttp"
	"github.com/maxatome/go-testdeep/td"
	"github.com/stretchr/testify/suite"
)

type accessSuite struct {
	suite.Suite
}

func TestAccess(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	suite.Run(t, new(accessSuite))
}

func (suite *accessSuite) BeforeTest(suiteName, testName string) {
	TestContext.BeforeTest()
}

func (suite *accessSuite) TestAccessWithValidSlugReturns301() {
	t := suite.T()

	testServer := TestContext.server
	testAPI := tdhttp.NewTestAPI(t, testServer)

	var slug string

	testAPI.PostJSON("/api/v1/shorturls", gin.H{"long_url": "https://www.cloudflare.com"}).
		CmpStatus(http.StatusCreated).
		CmpJSONBody(
			td.SuperJSONOf(`{"slug": "$slug"}`, td.Tag("slug", td.Catch(&slug, td.Ignore()))),
		)

	testAPI.Get(fmt.Sprintf("/%s", slug)).
		CmpStatus(http.StatusMovedPermanently).
		CmpHeader(http.Header{
			"Location":      []string{"https://www.cloudflare.com"},
			"Cache-Control": []string{"no-cache"},
		})
}

func (suite *accessSuite) TestAccessWithInvalidSlugReturns404() {
	t := suite.T()

	testServer := TestContext.server
	testAPI := tdhttp.NewTestAPI(t, testServer)

	testAPI.Get("/invalid").CmpStatus(http.StatusNotFound)
}
