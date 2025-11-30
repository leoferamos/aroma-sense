package middleware

import "github.com/leoferamos/aroma-sense/internal/model"

// minimal interface used by middleware to fetch user data
type userProfileReader interface {
	GetByPublicID(publicID string) (*model.User, error)
}

var userProfileSvc userProfileReader

// SetUserProfileService sets the service used by AccountStatusMiddleware
func SetUserProfileService(s userProfileReader) {
	userProfileSvc = s
}
