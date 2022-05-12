package controllers

import (
	"errors"
	"net/http"
	"url-shortener/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AccessShortUrlController struct {
	DB *gorm.DB
}

func (controller *AccessShortUrlController) HandleRequest(c *gin.Context) {
	slug := c.Param("slug")

	var shortUrl models.ShortUrl

	whereClause := models.ShortUrl{}
	whereClause.Slug = slug

	err := controller.DB.
		Where(&whereClause).
		First(&shortUrl).Error

	if err == nil {
		controller.DB.Model(&shortUrl).Association("Clicks").Append(&models.Click{})

		c.Writer.Header().Set("Location", shortUrl.LongUrl)
		c.Writer.Header().Set("Cache-Control", "private,max-age=0")
		c.Writer.WriteHeader(http.StatusMovedPermanently)
		return
	}

	status := http.StatusInternalServerError

	if errors.Is(err, gorm.ErrRecordNotFound) {
		status = http.StatusNotFound
	}

	c.Writer.WriteHeader(status)
}

func (controller *AccessShortUrlController) Register(r *gin.Engine) {
	r.GET("/:slug", controller.HandleRequest)
}
