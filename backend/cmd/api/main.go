package main

import (
	"log"
	"os"

	"github.com/leoferamos/aroma-sense/internal/bootstrap"
	"github.com/leoferamos/aroma-sense/internal/db"
	"github.com/leoferamos/aroma-sense/internal/router"
	"github.com/leoferamos/aroma-sense/internal/storage"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/leoferamos/aroma-sense/docs"
)

// @title           Aroma Sense API
// @version         1.0
// @description     REST API for user Aroma Sense Eccomerce

// @contact.name   API Support
// @contact.url    https://github.com/leoferamos/aroma-sense

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	db.Connect()

	storageClient, err := storage.NewSupabaseS3()
	if err != nil {
		log.Fatal("Failed to initialize storage client:", err)
	}

	// Initialize all modules with proper dependency injection
	handlers := bootstrap.InitializeApp(db.DB, storageClient)

	// Setup router with all handlers
	r := router.SetupRouter(handlers)

	// Swagger docs route
	if os.Getenv("ENABLE_SWAGGER") == "true" {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// Start the server
	r.Run(":8080")
}
