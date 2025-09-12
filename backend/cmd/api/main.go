package main

import (
	"log"

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

func main() {
	db.Connect()

	storageClient, err := storage.NewSupabaseS3()
	if err != nil {
		log.Fatal("Failed to initialize storage client:", err)
	}

	// Initialize modules
	userHandler := bootstrap.InitUserModule(db.DB)
	productHandler := bootstrap.InitProductModule(db.DB, storageClient)

	r := router.SetupRouter(userHandler, productHandler)

	// Swagger docs route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start the server
	r.Run(":8080")
}
