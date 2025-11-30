package router

import (
	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/auth"
	"github.com/leoferamos/aroma-sense/internal/handler"
	"github.com/leoferamos/aroma-sense/internal/middleware"
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
		authGroup.Use(auth.JWTAuthMiddleware(), middleware.AccountStatusMiddleware())
		{
			authGroup.GET("/me", userHandler.GetProfile)
			authGroup.PATCH("/me/profile", userHandler.UpdateProfile)
			authGroup.GET("/me/export", userHandler.ExportUserData)
			authGroup.POST("/me/deletion/confirm", userHandler.ConfirmAccountDeletion)
			authGroup.POST("/me/deletion", userHandler.RequestAccountDeletion)
		}
		// Routes that require authentication but must remain callable while the account is suspended or in cooling-off period.
		authNoStatus := userGroup.Group("")
		authNoStatus.Use(auth.JWTAuthMiddleware())
		{
			authNoStatus.POST("/me/deletion/cancel", userHandler.CancelAccountDeletion)
			authNoStatus.POST("/me/contest", userHandler.RequestContestation)
		}
	}
}
