package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"url-shortener/e"
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

type CreateResponse struct {
	ShortUrl string `json:"short_url"`
}

func (controller *ShortUrlController) CreateShortUrl(c *gin.Context) {
	var request models.ShortUrl
	var err error

	if err = c.ShouldBindJSON(&request); err != nil {
		var verr validator.ValidationErrors

		if errors.As(err, &verr) {
			c.JSON(http.StatusBadRequest, e.ErrorResponse{
				Errors: e.FormatErrors(verr),
			})
		} else {
			c.Writer.WriteHeader(http.StatusBadRequest)
		}

		return
	}

	if createResult := controller.CreateShortUrlService.Create(&request); createResult.Error == nil {
		var status int
		var body interface{}

		switch createResult.Status {
		case enums.CreationResultCreated:
			status = http.StatusCreated
			body = gin.H{"short_url": fmt.Sprintf("http://%s/%s", c.Request.Host, createResult.Record.Slug)}
		case enums.CreationResultAlreadyExists:
			status = http.StatusOK
			body = gin.H{"short_url": fmt.Sprintf("http://%s/%s", c.Request.Host, createResult.Record.Slug)}
		case enums.CreationResultDuplicateSlug:
			status = http.StatusConflict
			body = e.ErrorResponse{
				Errors: []e.ValidationError{
					e.ValidationError{
						Field:  "Slug",
						Reason: "must be unique",
					},
				},
			}
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
