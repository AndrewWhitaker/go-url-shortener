package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"url-shortener/enums"
	"url-shortener/models"
	"url-shortener/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type ShortUrlController struct {
	DB                    *gorm.DB
	CreateShortUrlService *services.CreateShortUrlService
}

func NewShortUrlController(db *gorm.DB) *ShortUrlController {
	return &ShortUrlController{
		DB: db,
		CreateShortUrlService: &services.CreateShortUrlService{
			DB: db,
		},
	}
}

type CreateResponse struct {
	ShortUrl string `json:"short_url"`
}

func (controller *ShortUrlController) CreateShortUrl(c *gin.Context) {
	var request models.ShortUrl
	var err error

	if err = c.ShouldBindJSON(&request); err != nil {
		var verr validator.ValidationErrors

		if errors.As(err, &verr) {
			c.JSON(http.StatusBadRequest, gin.H{"errors": FormatErrors(verr)})
		} else {
			c.Writer.WriteHeader(http.StatusBadRequest)
		}

		return
	}

	if createResult := controller.CreateShortUrlService.Create(&request); createResult.Error == nil {
		var status int
		var body map[string]interface{}

		switch createResult.Status {
		case enums.CreationResultCreated:
			status = http.StatusCreated
			body = gin.H{"short_url": fmt.Sprintf("http://%s/%s", c.Request.Host, createResult.Record.Slug)}
		case enums.CreationResultAlreadyExists:
			status = http.StatusOK
			body = gin.H{"short_url": fmt.Sprintf("http://%s/%s", c.Request.Host, createResult.Record.Slug)}
		case enums.CreationResultDuplicateSlug:
			status = http.StatusConflict
			body = gin.H{"errors": []gin.H{gin.H{"field": "Slug", "reason": "must be unique"}}}
		}

		c.JSON(status, body)
	} else {
		c.Writer.WriteHeader(http.StatusNotImplemented)
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
