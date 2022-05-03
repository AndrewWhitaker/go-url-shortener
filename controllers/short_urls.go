package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"url-shortener/db"
	"url-shortener/models"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

// POST /shorturls
func CreateShortUrl(c *gin.Context) {
	gormDb := db.DB
	var request models.ShortUrl

	if c.BindJSON(&request) == nil {
		// https://github.com/go-gorm/gorm/issues/4037
		if err := gormDb.Create(&request).Error; err != nil {
			var pgErr *pgconn.PgError

			if errors.As(err, &pgErr) {
				if pgErr.Code == pgerrcode.UniqueViolation {
					var existing models.ShortUrl

					gormDb.Where(&models.ShortUrl{LongUrl: request.LongUrl}).First(&existing)

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
	}
}