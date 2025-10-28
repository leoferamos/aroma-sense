package repository

import (
	"github.com/leoferamos/aroma-sense/internal/model"
	"gorm.io/gorm"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(user *model.User) error
	FindByEmail(email string) (*model.User, error)
	FindByRefreshTokenHash(hash string) (*model.User, error)
	Update(user *model.User) error
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
	if err := r.db.Where("refresh_token_hash = ?", hash).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Update saves changes to an existing user
func (r *userRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}
