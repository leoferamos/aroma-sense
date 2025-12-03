package main

import (
	"log"

	"github.com/leoferamos/aroma-sense/internal/bootstrap"
	"github.com/leoferamos/aroma-sense/internal/db"
	"github.com/leoferamos/aroma-sense/internal/server"
	"github.com/leoferamos/aroma-sense/internal/storage"

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
	app := bootstrap.InitializeApp(db.DB, storageClient)

	// Start server (jobs, router, graceful shutdown centralized in server package)
	if err := server.StartServer(app, ":8080"); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
