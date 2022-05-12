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

type deleteSuite struct {
	suite.Suite
}

func TestDelete(t *testing.T) {
	suite.Run(t, new(deleteSuite))
}

func (suite *deleteSuite) BeforeTest(suiteName, testName string) {
	TestContext.BeforeTest()
}

func (suite *deleteSuite) TestDeleteWithValidSlugReturns204() {
	t := suite.T()

	testServer := TestContext.server
	testAPI := tdhttp.NewTestAPI(t, testServer)

	var slug string

	testAPI.PostJSON("/api/v1/shorturls", gin.H{"long_url": "https://www.google.com"}).
		CmpStatus(http.StatusCreated).
		CmpJSONBody(
			td.SuperJSONOf(
				`{"slug": "$slug"}`,
				td.Tag("slug", td.Catch(&slug, td.Ignore())),
				td.Tag("shortUrl", td.Ignore()),
			),
		)

	testAPI.Delete(fmt.Sprintf("/api/v1/shorturls/%s", slug), nil).
		CmpStatus(http.StatusNoContent)

	testAPI.Get(fmt.Sprintf("/%s", slug)).CmpStatus(http.StatusNotFound)
}

func (suite *deleteSuite) TestDeleteWithInvalidSlugReturns404() {
	t := suite.T()

	testServer := TestContext.server
	testAPI := tdhttp.NewTestAPI(t, testServer)

	testAPI.Delete("/api/v1/shorturls/invalid", nil).
		CmpStatus(http.StatusNotFound).
		CmpJSONBody(
			td.JSON(
				`{
				   "errors": [
					   {
						   "field": "Slug",
							 "reason": "not found"
						 }
					 ],
				 }`,
			),
		)
}
