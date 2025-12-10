package router

import (
	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/auth"
	admin "github.com/leoferamos/aroma-sense/internal/handler/admin"
	loghandler "github.com/leoferamos/aroma-sense/internal/handler/log"
	orderhandler "github.com/leoferamos/aroma-sense/internal/handler/order"
	product "github.com/leoferamos/aroma-sense/internal/handler/product"
)

// AdminRoutes sets up the admin-related routes
func AdminRoutes(r *gin.Engine, adminUserHandler *admin.AdminUserHandler,
	productHandler *product.ProductHandler, orderHandler *orderhandler.OrderHandler,
	auditLogHandler *loghandler.AuditLogHandler,
	adminContestationHandler *admin.AdminContestationHandler,
	adminReviewReportHandler *admin.AdminReviewReportHandler) {
	adminGroup := r.Group("/admin")
	adminGroup.Use(auth.JWTAuthMiddleware(), auth.AdminOnly())

	superAdminGroup := r.Group("/admin")
	superAdminGroup.Use(auth.JWTAuthMiddleware(), auth.SuperAdminOnly())
	{
		superAdminGroup.POST("/users", adminUserHandler.AdminCreateAdmin)

		// User management
		adminGroup.GET("/users/:id/audit-logs", auditLogHandler.GetUserAuditLogs)
		adminGroup.POST("/users/:id/reactivate", adminUserHandler.AdminReactivateUser)
		adminGroup.POST("/users/:id/deactivate", adminUserHandler.AdminDeactivateUser)
		adminGroup.PATCH("/users/:id/role", adminUserHandler.AdminUpdateUserRole)
		adminGroup.GET("/users/:id", adminUserHandler.AdminGetUser)
		adminGroup.GET("/users", adminUserHandler.AdminListUsers)

		// Product management
		adminGroup.GET("/products", productHandler.AdminListProducts)
		adminGroup.POST("/products", productHandler.CreateProduct)
		adminGroup.GET("/products/:id", productHandler.GetProductByID)
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

		// User contestations
		adminGroup.GET("/contestations", adminContestationHandler.ListPendingContestions)
		adminGroup.POST("/contestations/:id/approve", adminContestationHandler.ApproveContestation)
		adminGroup.POST("/contestations/:id/reject", adminContestationHandler.RejectContestation)

		// Review reports
		adminGroup.GET("/review-reports", adminReviewReportHandler.ListReports)
		adminGroup.POST("/review-reports/:id/resolve", adminReviewReportHandler.ResolveReport)
	}
}
