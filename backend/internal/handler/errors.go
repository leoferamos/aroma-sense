package handler

import (
	"errors"
	"net/http"

	"github.com/leoferamos/aroma-sense/internal/service"
)

// errorMapping centralizes mapping of service-layer errors to HTTP status and client-facing messages.
var errorMapping = map[error]struct {
	Status  int
	Message string
}{
	service.ErrUnauthorized:        {Status: http.StatusUnauthorized, Message: "unauthorized"},
	service.ErrCartEmpty:           {Status: http.StatusBadRequest, Message: "cart is empty"},
	service.ErrInvalidPostalCode:   {Status: http.StatusBadRequest, Message: "invalid postal_code"},
	service.ErrOriginNotConfigured: {Status: http.StatusInternalServerError, Message: "shipping origin not configured"},
	service.ErrProviderUnavailable: {Status: http.StatusServiceUnavailable, Message: "shipping temporarily unavailable"},
	service.ErrNoOptions:           {Status: http.StatusServiceUnavailable, Message: "no shipping options available"},
	// Reviews
	service.ErrReviewUnauthenticated:   {Status: http.StatusUnauthorized, Message: "unauthenticated"},
	service.ErrReviewProfileIncomplete: {Status: http.StatusForbidden, Message: "profile_incomplete"},
	service.ErrReviewNotDelivered:      {Status: http.StatusForbidden, Message: "not_delivered"},
	service.ErrReviewAlreadyReviewed:   {Status: http.StatusConflict, Message: "already_reviewed"},
	service.ErrReviewInvalidRating:     {Status: http.StatusBadRequest, Message: "invalid rating"},
	service.ErrReviewCommentTooLong:    {Status: http.StatusBadRequest, Message: "comment too long"},
	service.ErrReviewProductNotFound:   {Status: http.StatusNotFound, Message: "product not found"},
}

// mapServiceError checks if an error matches a known service error and returns the appropriate HTTP mapping.
func mapServiceError(err error) (status int, message string, ok bool) {
	for known, info := range errorMapping {
		if errors.Is(err, known) {
			return info.Status, info.Message, true
		}
	}
	return 0, "", false
}
