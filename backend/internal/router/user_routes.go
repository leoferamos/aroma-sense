package router

import (
	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/handler"
)

// UserRoutes sets up the user-related routes
func UserRoutes(r *gin.Engine, handler *handler.UserHandler) {
	userGroup := r.Group("/users")
	{
		userGroup.POST("/register", handler.RegisterUser)
		userGroup.POST("/login", handler.LoginUser)
		userGroup.POST("/logout", handler.LogoutUser)
	}
}
