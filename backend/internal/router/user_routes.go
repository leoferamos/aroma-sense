package router

import (
	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/handler"
)

// RegisterUserRoutes sets up the user-related routes
func RegisterUserRoutes(r *gin.Engine, handler *handler.UserHandler) {
	userGroup := r.Group("/users")
	{
		userGroup.POST("/register", handler.RegisterUser)

	}
}
