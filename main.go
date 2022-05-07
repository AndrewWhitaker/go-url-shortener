package main

import (
	"database/sql"
	"fmt"
	"os"
	"url-shortener/controllers"
	"url-shortener/db"

	"github.com/gin-gonic/gin"
)

type ServerConfig struct {
	DB *sql.DB
}

func main() {
	postgresHost := os.Getenv("POSTGRES_HOST")
	postgresPort := os.Getenv("POSTGRES_PORT")
	postgresUser := os.Getenv("POSTGRES_USER")
	postgresPass := os.Getenv("POSTGRES_PASSWORD")
	postgresDatabase := os.Getenv("POSTGRES_DATABASE")

	url := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		postgresUser,
		postgresPass,
		postgresHost,
		postgresPort,
		postgresDatabase,
	)

	fmt.Println(url)

	fmt.Printf("Connecting to %s\n", url)
	db, err := sql.Open("pgx", url)

	if err != nil {
		panic(fmt.Sprintf("Unable to connect to postgres: %s", err))
	}

	config := ServerConfig{DB: db}
	SetupServer(&config).Run()
}

func SetupServer(cfg *ServerConfig) *gin.Engine {
	r := gin.Default()
	db, err := db.ConnectDatabase(cfg.DB)

	if err != nil {
		panic("Failed to connect to database")
	}

	shortUrlController := controllers.NewShortUrlController(db)

	// Access an existing short URL
	r.GET("/:slug", shortUrlController.GetShortUrl)
	r.POST("/shorturls", shortUrlController.CreateShortUrl)
	r.DELETE("/shorturls/:slug", shortUrlController.DeleteShortUrl)

	// Get details about an existing short URL
	r.GET("/shorturls/:slug", func(c *gin.Context) {
		c.Writer.WriteHeader(501)
	})

	return r
}
