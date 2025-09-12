package bootstrap

import (
	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/handler"
	"github.com/leoferamos/aroma-sense/internal/router"
)

func InitAdminRoutes(r *gin.Engine, userHandler *handler.UserHandler, productHandler *handler.ProductHandler) {
	router.AdminRoutes(r, userHandler, productHandler)
}
