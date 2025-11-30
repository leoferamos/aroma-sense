package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/dto"
)

// AccountStatusMiddleware blocks access to protected endpoints when the user is deactivated or has a pending/confirmed deletion.
func AccountStatusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		rawUserID, exists := c.Get("userID")
		if !exists {
			c.Next()
			return
		}
		publicID, ok := rawUserID.(string)
		if !ok || publicID == "" {
			c.Next()
			return
		}

		user, err := userProfileSvc.GetByPublicID(publicID)
		if err != nil || user == nil {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthenticated"})
			c.Abort()
			return
		}

		// Admin deactivation
		if user.DeactivatedAt != nil {
			// Respond with structured details about deactivation
			payload := gin.H{
				"error":                 "account_deactivated",
				"message":               "Your account has been deactivated.",
				"deactivated_at":        user.DeactivatedAt,
				"deactivated_by":        user.DeactivatedBy,
				"deactivation_reason":   user.DeactivationReason,
				"deactivation_notes":    user.DeactivationNotes,
				"suspension_until":      user.SuspensionUntil,
				"contestation_deadline": user.ContestationDeadline,
			}
			c.JSON(http.StatusForbidden, payload)
			c.Abort()
			return
		}

		// User requested deletion and still in cooling-off period
		if user.DeletionRequestedAt != nil && user.DeletionConfirmedAt == nil {
			payload := gin.H{
				"error":                 "deletion_requested",
				"message":               "Account in cooling-off period. You can cancel the request.",
				"deletion_requested_at": user.DeletionRequestedAt,
			}
			c.JSON(http.StatusForbidden, payload)
			c.Abort()
			return
		}

		// Deletion already confirmed — block everything
		if user.DeletionConfirmedAt != nil {
			payload := gin.H{
				"error":                 "deletion_confirmed",
				"message":               "Deletion confirmed — data retention in progress.",
				"deletion_confirmed_at": user.DeletionConfirmedAt,
			}
			c.JSON(http.StatusForbidden, payload)
			c.Abort()
			return
		}

		c.Next()
	}
}
