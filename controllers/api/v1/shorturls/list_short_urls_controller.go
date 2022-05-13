package shorturls

import (
	"net/http"
	"url-shortener/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ListShortUrlsController struct {
	DB *gorm.DB
}

// ListShortUrls  godoc
// @Summary      List all short URLs
// @Description  List all short URLs
// @Tags         shorturls
// @Accept       json
// @Produce      json
// @Success      200  {array}  models.ShortUrlReadFields
// @Failure      500
// @Router       /shorturls [get]
func (controller *ListShortUrlsController) HandleRequest(c *gin.Context) {
	var allShortUrls []models.ShortUrl

	listResult := controller.DB.
		Order("created_at ASC").
		Find(&allShortUrls)

	var jsonResults []shortUrlResponseHelper

	if listResult.Error == nil {
		for _, shortUrl := range allShortUrls {
			jsonResults = append(jsonResults, shortUrlResponseHelper{
				Host:     c.Request.Host,
				ShortUrl: shortUrl,
			})
		}

		c.JSON(http.StatusOK, jsonResults)
		return
	}

	c.Writer.WriteHeader(http.StatusInternalServerError)
}

func (controller *ListShortUrlsController) Register(r *gin.Engine) {
	r.GET("/api/v1/shorturls", controller.HandleRequest)
}
