package repository

import (
	"time"

	"github.com/leoferamos/aroma-sense/internal/model"
	"gorm.io/gorm"
)

// ResetTokenRepository defines the interface for reset token data operations
type ResetTokenRepository interface {
	Create(token *model.ResetToken) error
	FindByEmailAndCode(email, code string) (*model.ResetToken, error)
	DeleteByEmail(email string) error
	DeleteExpired() error
}

type resetTokenRepository struct {
	db *gorm.DB
}

// NewResetTokenRepository creates a new instance of ResetTokenRepository
func NewResetTokenRepository(db *gorm.DB) ResetTokenRepository {
	return &resetTokenRepository{db: db}
}

// Create saves a new reset token in the database
func (r *resetTokenRepository) Create(token *model.ResetToken) error {
	return r.db.Create(token).Error
}

// FindByEmailAndCode retrieves a valid (non-expired) token by email and code
func (r *resetTokenRepository) FindByEmailAndCode(email, code string) (*model.ResetToken, error) {
	var token model.ResetToken
	if err := r.db.Where("email = ? AND code = ? AND expires_at > ?", email, code, time.Now()).
		Order("created_at DESC").
		First(&token).Error; err != nil {
		return nil, err
	}
	return &token, nil
}

// DeleteByEmail removes all reset tokens for a given email
func (r *resetTokenRepository) DeleteByEmail(email string) error {
	return r.db.Where("email = ?", email).Delete(&model.ResetToken{}).Error
}

// DeleteExpired removes all expired reset tokens
func (r *resetTokenRepository) DeleteExpired() error {
	return r.db.Where("expires_at <= ?", time.Now()).Delete(&model.ResetToken{}).Error
}
