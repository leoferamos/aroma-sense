package router

import (
	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/auth"
	"github.com/leoferamos/aroma-sense/internal/handler"
)

// AdminRoutes sets up the admin-related routes
func AdminRoutes(r *gin.Engine, adminUserHandler *handler.AdminUserHandler,
	productHandler *handler.ProductHandler, orderHandler *handler.OrderHandler) {
	adminGroup := r.Group("/admin")
	adminGroup.Use(auth.JWTAuthMiddleware(), auth.AdminOnly())
	{
		// User management
		adminGroup.GET("/users", adminUserHandler.AdminListUsers)
		adminGroup.GET("/users/:id", adminUserHandler.AdminGetUser)
		adminGroup.PATCH("/users/:id/role", adminUserHandler.AdminUpdateUserRole)
		adminGroup.POST("/users/:id/deactivate", adminUserHandler.AdminDeactivateUser)

		// Product management
		adminGroup.POST("/products", productHandler.CreateProduct)
		adminGroup.GET("/products/:id", productHandler.GetProduct)
		adminGroup.PATCH("/products/:id", productHandler.UpdateProduct)
		adminGroup.DELETE("/products/:id", productHandler.DeleteProduct)

		// Order management
		adminGroup.GET("/orders", orderHandler.ListOrders)
	}
}
