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

	// TODO:
	// * Look into custom binding for "slug"
	// * Refactor into generic handler
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

		switch createResult.Status {
		case enums.CreationResultCreated:
			status = http.StatusCreated
		case enums.CreationResultAlreadyExisted:
			status = http.StatusOK
		}

		c.JSON(status, gin.H{
			"short_url": fmt.Sprintf("http://%s/%s", c.Request.Host, createResult.Record.Slug),
		})
	} else {
		c.Writer.WriteHeader(http.StatusNotImplemented)
	}
}
