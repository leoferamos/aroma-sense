package handlererrors

import (
	"errors"
	"net/http"

	"github.com/leoferamos/aroma-sense/internal/apperror"
)

// domainCodeStatus maps domain error codes to HTTP statuses.
var domainCodeStatus = map[string]int{
	"unauthenticated":                http.StatusUnauthorized,
	"unauthorized":                   http.StatusUnauthorized,
	"invalid_request":                http.StatusBadRequest,
	"profile_incomplete":             http.StatusForbidden,
	"not_delivered":                  http.StatusForbidden,
	"already_reviewed":               http.StatusConflict,
	"invalid_postal_code":            http.StatusBadRequest,
	"cart_empty":                     http.StatusBadRequest,
	"origin_not_configured":          http.StatusInternalServerError,
	"provider_unavailable":           http.StatusServiceUnavailable,
	"no_shipping_options":            http.StatusServiceUnavailable,
	"product_not_found":              http.StatusNotFound,
	"review_not_found":               http.StatusNotFound,
	"insufficient_stock":             http.StatusBadRequest,
	"cart_item_not_found":            http.StatusNotFound,
	"cart_update_failed":             http.StatusInternalServerError,
	"stock_update_failed":            http.StatusInternalServerError,
	"invalid_shipping_selection":     http.StatusBadRequest,
	"cart_clear_failed":              http.StatusInternalServerError,
	"invalid_rating":                 http.StatusBadRequest,
	"comment_too_long":               http.StatusBadRequest,
	"invalid_category":               http.StatusBadRequest,
	"reason_too_long":                http.StatusBadRequest,
	"cannot_report_own_review":       http.StatusForbidden,
	"already_reported":               http.StatusConflict,
	"invalid_status":                 http.StatusBadRequest,
	"invalid_action":                 http.StatusBadRequest,
	"report_not_found":               http.StatusNotFound,
	"report_already_resolved":        http.StatusConflict,
	"active_orders_block_deletion":   http.StatusBadRequest,
	"deletion_already_requested":     http.StatusConflict,
	"deletion_not_requested":         http.StatusNotFound,
	"cooling_off_not_expired":        http.StatusBadRequest,
	"deletion_not_confirmed":         http.StatusBadRequest,
	"retention_not_expired":          http.StatusBadRequest,
	"account_not_deactivated":        http.StatusBadRequest,
	"contestation_deadline_expired":  http.StatusBadRequest,
	"reactivation_already_requested": http.StatusConflict,
	"email_already_registered":       http.StatusConflict,
	"password_hash_failed":           http.StatusInternalServerError,
	"rate_limited":                   http.StatusTooManyRequests,
	"cart_create_failed":             http.StatusInternalServerError,
	"invalid_credentials":            http.StatusUnauthorized,
	"access_token_failed":            http.StatusInternalServerError,
	"refresh_token_failed":           http.StatusInternalServerError,
	"refresh_token_save_failed":      http.StatusInternalServerError,
	"invalid_refresh_token":          http.StatusUnauthorized,
	"refresh_token_expired":          http.StatusUnauthorized,
	"refresh_token_missing":          http.StatusBadRequest,
	"reset_generation_failed":        http.StatusInternalServerError,
	"reset_cleanup_failed":           http.StatusInternalServerError,
	"reset_token_create_failed":      http.StatusInternalServerError,
	"reset_email_failed":             http.StatusInternalServerError,
	"reset_code_invalid":             http.StatusBadRequest,
	"password_update_failed":         http.StatusInternalServerError,
	"display_name_too_short":         http.StatusBadRequest,
	"display_name_too_long":          http.StatusBadRequest,
	"current_password_incorrect":     http.StatusUnauthorized,
	"new_password_same":              http.StatusBadRequest,
	"password_process_failed":        http.StatusInternalServerError,
	"topic_restricted":               http.StatusBadRequest,
	"invalid_role":                   http.StatusBadRequest,
	"cannot_change_own_role":         http.StatusForbidden,
	"suspension_until_past":          http.StatusBadRequest,
	"user_not_deactivated":           http.StatusBadRequest,
	"invalid_webhook":                http.StatusBadRequest,
	"internal_error":                 http.StatusInternalServerError,
}

// mapServiceError checks if an error matches a known service error and returns the appropriate HTTP mapping.
func mapServiceError(err error) (status int, message string, ok bool) {
	if err == nil {
		return 0, "", false
	}

	var de *apperror.DomainError
	if errors.As(err, &de) {
		code := de.Code
		status, exists := domainCodeStatus[code]
		if !exists {
			status = http.StatusInternalServerError
			if code == "" {
				code = "internal_error"
			}
		}
		return status, code, true
	}

	return 0, "", false
}

// MapServiceError exposes error-to-HTTP mapping for handler subpackages.
func MapServiceError(err error) (status int, message string, ok bool) {
	return mapServiceError(err)
}
