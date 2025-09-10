package main

import (
	"github.com/leoferamos/aroma-sense/internal/bootstrap"
	"github.com/leoferamos/aroma-sense/internal/db"
	"github.com/leoferamos/aroma-sense/internal/router"
)

func main() {

	db.Connect()

	// Initialize user module
	userHandler := bootstrap.InitUserModule(db.DB)

	r := router.SetupRouter(userHandler)

	// Start the server
	r.Run(":8080")
}
