package shorturls

import (
	"errors"
	"net/http"
	"url-shortener/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GetShortUrlController struct {
	DB *gorm.DB
}

// GetShortUrl godoc
// @Summary      Get information about an existing short URL
// @Description  Get information about an existing short URL
// @Tags         shorturls
// @Accept       json
// @Produce      json
// @Param        slug  path      string  true  "slug of short URL to get information about"
// @Success      200   {object}  models.ShortUrlReadFields
// @Failure      404   {object}  e.ErrorResponse
// @Failure      500
// @Router       /shorturls/{slug} [get]
func (controller *GetShortUrlController) HandleRequest(c *gin.Context) {
	slug := c.Param("slug")

	var shortUrl models.ShortUrl

	whereClause := models.ShortUrl{}
	whereClause.Slug = slug

	err := controller.DB.
		Where(&whereClause).
		First(&shortUrl).Error

	if err == nil {
		c.JSON(http.StatusOK, shortUrlResponseHelper{
			Host:     c.Request.Host,
			ShortUrl: shortUrl,
		})

		return
	}

	status := http.StatusInternalServerError

	if errors.Is(err, gorm.ErrRecordNotFound) {
		status = http.StatusNotFound
	}

	c.Writer.WriteHeader(status)
}

func (controller *GetShortUrlController) Register(r *gin.Engine) {
	r.GET("/api/v1/shorturls/:slug", controller.HandleRequest)
}
