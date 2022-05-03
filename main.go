package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type ServerConfig struct {
	DB *sql.DB
}

type ShortUrl struct {
	Id        int64  `gorm:"primaryKey"`
	LongUrl   string `gorm:"index:uq_short_urls_long_url,unique" json:"long_url"`
	CreatedAt time.Time
}

func main() {
	config := ServerConfig{DB: nil}
	SetupServer(&config).Run()
}

func SetupServer(cfg *ServerConfig) *gin.Engine {
	r := gin.Default()

	gormDb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: cfg.DB,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		panic("Failed to connect to database")
	}

	gormDb.AutoMigrate(&ShortUrl{})

	// Access an existing short URL
	r.GET("/:slug", func(c *gin.Context) {
		c.Writer.WriteHeader(501)
	})

	// Create a new short URL
	r.POST("/shorturls", func(c *gin.Context) {
		var request ShortUrl
		if c.BindJSON(&request) == nil {
			// https://github.com/go-gorm/gorm/issues/4037
			if err := gormDb.Create(&request).Error; err != nil {
				var pgErr *pgconn.PgError

				if errors.As(err, &pgErr) {
					if pgErr.Code == pgerrcode.UniqueViolation {
						var existing ShortUrl

						gormDb.Where(&ShortUrl{LongUrl: request.LongUrl}).First(&existing)

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
	})

	// Get details about an existing short URL
	r.GET("/shorturls/:slug", func(c *gin.Context) {
		c.Writer.WriteHeader(501)
	})

	// Delete a short url
	r.DELETE("/shorturls/:slug", func(c *gin.Context) {
		c.Writer.WriteHeader(501)
	})

	return r
}
