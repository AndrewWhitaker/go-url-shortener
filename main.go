package main

import (
	"database/sql"
	"url-shortener/controllers"
	"url-shortener/db"

	"github.com/gin-gonic/gin"
)

type ServerConfig struct {
	DB *sql.DB
}

func main() {
	config := ServerConfig{DB: nil}
	SetupServer(&config).Run()
}

func SetupServer(cfg *ServerConfig) *gin.Engine {
	r := gin.Default()
	err := db.ConnectDatabase(cfg.DB)

	if err != nil {
		panic("Failed to connect to database")
	}

	// Access an existing short URL
	r.GET("/:slug", func(c *gin.Context) {
		c.Writer.WriteHeader(501)
	})

	// Create a new short URL
	r.POST("/shorturls", controllers.CreateShortUrl)

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
