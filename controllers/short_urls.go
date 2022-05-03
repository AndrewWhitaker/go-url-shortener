package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"url-shortener/models"

	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"gorm.io/gorm"
)

type ShortUrlController struct {
	DB *gorm.DB
}

func NewShortUrlController(db *gorm.DB) *ShortUrlController {
	return &ShortUrlController{DB: db}
}

func (controller *ShortUrlController) CreateShortUrl(c *gin.Context) {
	db := controller.DB
	var request models.ShortUrl

	if err := c.ShouldBindJSON(&request); err == nil {
		// https://github.com/go-gorm/gorm/issues/4037
		if err := db.Create(&request).Error; err != nil {
			var pgErr *pgconn.PgError

			if errors.As(err, &pgErr) {
				if pgErr.Code == pgerrcode.UniqueViolation {
					var existing models.ShortUrl

					db.Where(&models.ShortUrl{LongUrl: request.LongUrl}).First(&existing)

					c.JSON(http.StatusOK, gin.H{
						"short_url": fmt.Sprintf("http://%s/%d", c.Request.Host, existing.Id),
					})
				}
			} else {
				c.Writer.WriteHeader(http.StatusNotImplemented)
			}
		} else {
			c.JSON(http.StatusCreated, gin.H{
				"short_url": fmt.Sprintf("http://%s/%d", c.Request.Host, request.Id),
			})
		}
	} else {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			c.JSON(http.StatusBadRequest, gin.H{"errors": FormatErrors(verr)})
		}
	}
}

// POST /shorturls
func CreateShortUrl(c *gin.Context) {
}
