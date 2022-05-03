package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()

	// Access an existing short URL
	r.GET("/:slug", func(c *gin.Context) {
		c.Writer.WriteHeader(501)
	})

	// Create a new short URL
	r.POST("/shorturls", func(c *gin.Context) {
		c.Writer.WriteHeader(501)
	})

	// Get details about an existing short URL
	r.GET("/shorturls/:slug", func(c *gin.Context) {
		c.Writer.WriteHeader(501)
	})

	// Delete a short url
	r.DELETE("/shorturls/:slug", func(c *gin.Context) {
		c.Writer.WriteHeader(501)
	})

	r.Run()
}
