package controllers

import (
	"errors"
	"net/http"

	"url-shortener/e"
	"url-shortener/enums"
	"url-shortener/models"
	"url-shortener/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ShortUrlController struct {
	DB                    *gorm.DB
	CreateShortUrlService *services.CreateShortUrlService
	DeleteShortUrlService *services.DeleteShortUrlService
}

func NewShortUrlController(db *gorm.DB) *ShortUrlController {
	return &ShortUrlController{
		DB: db,
		CreateShortUrlService: &services.CreateShortUrlService{
			DB: db,
		},
		DeleteShortUrlService: &services.DeleteShortUrlService{
			DB: db,
		},
	}
}

func (controller *ShortUrlController) GetShortUrl(c *gin.Context) {
	slug := c.Param("slug")

	var shortUrl models.ShortUrl

	err := controller.DB.
		Where(&models.ShortUrl{Slug: slug}).
		First(&shortUrl).Error

	if err == nil {
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

func (controller *ShortUrlController) DeleteShortUrl(c *gin.Context) {
	slug := c.Param("slug")

	result := controller.DeleteShortUrlService.Delete(slug)

	if result.Status == enums.DeleteResultSuccessful {
		c.Writer.WriteHeader(http.StatusNoContent)
	} else if result.Status == enums.DeleteResultNotFound {
		c.JSON(http.StatusNotFound, e.ErrorResponse{
			Errors: []e.ValidationError{
				e.ValidationError{
					Field:  "Slug",
					Reason: "not found",
				},
			},
		})
	} else {
		c.Writer.WriteHeader(http.StatusInternalServerError)
	}
}
