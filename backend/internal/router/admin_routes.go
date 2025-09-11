package router

import (
	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/auth"
	"github.com/leoferamos/aroma-sense/internal/handler"
)

// AdminRoutes sets up the admin-related routes
func AdminRoutes(r *gin.Engine, handler *handler.UserHandler) {
	adminGroup := r.Group("/admin")
	adminGroup.Use(auth.JWTAuthMiddleware(), auth.AdminOnly())
	{

	}
}
