package router

import (
	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/auth"
	"github.com/leoferamos/aroma-sense/internal/handler"
)

// AdminRoutes sets up the admin-related routes
func AdminRoutes(r *gin.Engine, adminUserHandler *handler.AdminUserHandler,
	productHandler *handler.ProductHandler, orderHandler *handler.OrderHandler,
	auditLogHandler *handler.AuditLogHandler) {
	adminGroup := r.Group("/admin")
	adminGroup.Use(auth.JWTAuthMiddleware(), auth.AdminOnly())
	{
		// User management
		adminGroup.GET("/users/:id/audit-logs", auditLogHandler.GetUserAuditLogs)
		adminGroup.POST("/users/:id/reactivate", adminUserHandler.AdminReactivateUser)
		adminGroup.POST("/users/:id/deactivate", adminUserHandler.AdminDeactivateUser)
		adminGroup.PATCH("/users/:id/role", adminUserHandler.AdminUpdateUserRole)
		adminGroup.GET("/users/:id", adminUserHandler.AdminGetUser)
		adminGroup.GET("/users", adminUserHandler.AdminListUsers)

		// Product management
		adminGroup.POST("/products", productHandler.CreateProduct)
		adminGroup.GET("/products/:id", productHandler.GetProduct)
		adminGroup.PATCH("/products/:id", productHandler.UpdateProduct)
		adminGroup.DELETE("/products/:id", productHandler.DeleteProduct)

		// Order management
		adminGroup.GET("/orders", orderHandler.ListOrders)

		// Audit logs
		adminGroup.GET("/audit-logs/:id/detailed", auditLogHandler.GetAuditLogDetailed)
		adminGroup.GET("/audit-logs/:id", auditLogHandler.GetAuditLog)
		adminGroup.GET("/audit-logs/summary", auditLogHandler.GetAuditSummary)
		adminGroup.POST("/audit-logs/cleanup", auditLogHandler.CleanupOldLogs)
		adminGroup.GET("/audit-logs", auditLogHandler.ListAuditLogs)
	}
}
