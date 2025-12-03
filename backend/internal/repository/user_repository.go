package repository

import (
	"time"

	"github.com/leoferamos/aroma-sense/internal/model"
	"gorm.io/gorm"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(user *model.User) error
	FindByEmail(email string) (*model.User, error)
	FindByRefreshTokenHash(hash string) (*model.User, error)
	FindByPublicID(publicID string) (*model.User, error)
	Update(user *model.User) error
	UpdateRefreshToken(userID uint, hash *string, expiresAt *time.Time) error
	ListUsers(limit int, offset int, filters map[string]interface{}) ([]*model.User, int64, error)
	FindByID(id uint) (*model.User, error)
	UpdateRole(userID uint, newRole string) error
	DeactivateUser(userID uint, adminPublicID string, deactivatedAt time.Time, reason string, notes string, suspensionUntil *time.Time) error
	RequestAccountDeletion(publicID string, requestedAt time.Time) error
	ConfirmAccountDeletion(publicID string, confirmedAt time.Time) error
	HasActiveDependencies(publicID string) (bool, error)
	AnonymizeUser(publicID string, anonymizedEmail string, anonymizedDisplayName string) error
	FindExpiredUsersForAnonymization() ([]*model.User, error)
	FindUsersPendingAutoConfirm(cutoff time.Time) ([]*model.User, error)
	DeleteByPublicID(publicID string) error
}

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create saves a new user in the database
func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// FindByEmail retrieves a user by email
func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByRefreshTokenHash retrieves a user by refresh token hash
func (r *userRepository) FindByRefreshTokenHash(hash string) (*model.User, error) {
	var user model.User
	if err := r.db.Select("id, public_id, role, refresh_token_hash, refresh_token_expires_at").
		Where("refresh_token_hash = ? AND refresh_token_expires_at > ?", hash, time.Now()).
		First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByPublicID retrieves a user by their public UUID identifier
func (r *userRepository) FindByPublicID(publicID string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("public_id = ?", publicID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Update saves changes to an existing user
func (r *userRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

// UpdateRefreshToken updates only the refresh-token related columns for a user.
func (r *userRepository) UpdateRefreshToken(userID uint, hash *string, expiresAt *time.Time) error {
	return r.db.Model(&model.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"refresh_token_hash":       hash,
		"refresh_token_expires_at": expiresAt,
	}).Error
}

// ListUsers returns paginated list of users for admin (LGPD compliance)
func (r *userRepository) ListUsers(limit int, offset int, filters map[string]interface{}) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	query := r.db.Model(&model.User{})

	// Apply filters
	if role, exists := filters["role"]; exists && role != "" {
		query = query.Where("role = ?", role)
	}
	if status, exists := filters["status"]; exists && status != "" {
		switch status {
		case "active":
			query = query.Where("deleted_at IS NULL AND deactivated_at IS NULL")
		case "deactivated":
			query = query.Where("deactivated_at IS NOT NULL")
		case "deleted":
			query = query.Where("deleted_at IS NOT NULL")
		}
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := query.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// FindByID retrieves a user by database ID
func (r *userRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateRole updates user role
func (r *userRepository) UpdateRole(userID uint, newRole string) error {
	return r.db.Model(&model.User{}).Where("id = ?", userID).Update("role", newRole).Error
}

// DeactivateUser soft deletes a user account with enhanced LGPD compliance
func (r *userRepository) DeactivateUser(userID uint, adminPublicID string, deactivatedAt time.Time, reason string, notes string, suspensionUntil *time.Time) error {
	updates := map[string]interface{}{
		"deactivated_by":      adminPublicID,
		"deactivated_at":      deactivatedAt,
		"deactivation_reason": reason,
		"deactivation_notes":  notes,
	}

	if suspensionUntil != nil {
		updates["suspension_until"] = suspensionUntil
	}

	return r.db.Model(&model.User{}).Where("id = ?", userID).Updates(updates).Error
}

// DeleteByPublicID permanently deletes a user by public ID
func (r *userRepository) DeleteByPublicID(publicID string) error {
	return r.db.Unscoped().Where("public_id = ?", publicID).Delete(&model.User{}).Error
}

// RequestAccountDeletion marks account for deletion with retention period (LGPD compliance)
func (r *userRepository) RequestAccountDeletion(publicID string, requestedAt time.Time) error {
	return r.db.Model(&model.User{}).Where("public_id = ?", publicID).Update("deletion_requested_at", requestedAt).Error
}

// ConfirmAccountDeletion confirms account deletion after retention period
func (r *userRepository) ConfirmAccountDeletion(publicID string, confirmedAt time.Time) error {
	return r.db.Model(&model.User{}).Where("public_id = ?", publicID).Update("deletion_confirmed_at", confirmedAt).Error
}

// HasActiveDependencies checks if user has active orders or other dependencies that prevent deletion
func (r *userRepository) HasActiveDependencies(publicID string) (bool, error) {
	// Check for active (non-delivered) orders
	var activeOrdersCount int64
	if err := r.db.Table("orders").
		Joins("JOIN order_items ON orders.id = order_items.order_id").
		Where("orders.user_id = ? AND orders.status NOT IN ('delivered', 'cancelled')", publicID).
		Count(&activeOrdersCount).Error; err != nil {
		return false, err
	}

	return activeOrdersCount > 0, nil
}

// FindExpiredUsersForAnonymization finds users who have confirmed deletion and exceeded retention period
func (r *userRepository) FindExpiredUsersForAnonymization() ([]*model.User, error) {
	var users []*model.User

	// Calculate cutoff date (2 years ago from now)
	cutoffDate := time.Now().AddDate(-2, 0, 0)

	// Find users who confirmed deletion more than 2 years ago
	if err := r.db.Where("deletion_confirmed_at IS NOT NULL AND deletion_confirmed_at < ?", cutoffDate).
		Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

// FindUsersPendingAutoConfirm finds users who requested deletion and whose deletion_requested_at is older than cutoff
func (r *userRepository) FindUsersPendingAutoConfirm(cutoff time.Time) ([]*model.User, error) {
	var users []*model.User

	if err := r.db.Where("deletion_requested_at IS NOT NULL AND deletion_confirmed_at IS NULL AND deletion_requested_at <= ?", cutoff).
		Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

// AnonymizeUser anonymizes user personal data while keeping audit trail (LGPD compliance)
func (r *userRepository) AnonymizeUser(publicID string, anonymizedEmail string, anonymizedDisplayName string) error {
	return r.db.Model(&model.User{}).Where("public_id = ?", publicID).Updates(map[string]interface{}{
		"email":                    anonymizedEmail,
		"password_hash":            "",
		"display_name":             anonymizedDisplayName,
		"last_login_at":            nil,
		"refresh_token_hash":       nil,
		"refresh_token_expires_at": nil,
		"deactivated_at":           nil,
		"deactivated_by":           nil,
		"deactivation_reason":      nil,
		"deactivation_notes":       nil,
		"suspension_until":         nil,
		"reactivation_requested":   false,
		"contestation_deadline":    nil,
	}).Error
}
