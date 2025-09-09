package main

import (
	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/db"
)

func main() {
	// Connect to PostgreSQL
	db.Connect()
	defer db.Close()

	// Initialize Gin router
	router := gin.Default()

	// Health check endpoint
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Start server
	router.Run(":8080")
}
