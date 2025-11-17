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
