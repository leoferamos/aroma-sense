package main

import (
	"github.com/leoferamos/aroma-sense/internal/bootstrap"
	"github.com/leoferamos/aroma-sense/internal/db"
	"github.com/leoferamos/aroma-sense/internal/router"

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

	// Initialize user module
	userHandler := bootstrap.InitUserModule(db.DB)

	r := router.SetupRouter(userHandler)

	// Swagger docs route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start the server
	r.Run(":8080")
}
