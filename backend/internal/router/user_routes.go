package router

import (
	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/auth"
	"github.com/leoferamos/aroma-sense/internal/handler"
)

// UserRoutes sets up the user-related routes
func UserRoutes(r *gin.Engine, userHandler *handler.UserHandler, resetHandler *handler.PasswordResetHandler) {
	userGroup := r.Group("/users")
	{
		userGroup.POST("/register", userHandler.RegisterUser)
		userGroup.POST("/login", userHandler.LoginUser)
		userGroup.POST("/refresh", userHandler.RefreshToken)
		userGroup.POST("/logout", userHandler.LogoutUser)
		userGroup.POST("/reset/request", resetHandler.RequestReset)
		userGroup.POST("/reset/confirm", resetHandler.ConfirmReset)

		// Authenticated user endpoints
		authGroup := userGroup.Group("")
		authGroup.Use(auth.JWTAuthMiddleware())
		{
			authGroup.GET("/me", userHandler.GetProfile)
			authGroup.PATCH("/me/profile", userHandler.UpdateProfile)
			authGroup.GET("/me/export", userHandler.ExportUserData)
			authGroup.POST("/me/deletion/confirm", userHandler.ConfirmAccountDeletion)
			authGroup.POST("/me/deletion/cancel", userHandler.CancelAccountDeletion)
			authGroup.POST("/me/deletion", userHandler.RequestAccountDeletion)
			authGroup.POST("/me/contest", userHandler.RequestContestation)
		}
	}
}
