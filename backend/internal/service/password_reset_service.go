package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/leoferamos/aroma-sense/internal/email"
	"github.com/leoferamos/aroma-sense/internal/model"
	"github.com/leoferamos/aroma-sense/internal/repository"
	"github.com/leoferamos/aroma-sense/internal/validation"
	"github.com/leoferamos/aroma-sense/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	// ResetTokenExpiration defines how long a reset token is valid
	ResetTokenExpiration = 10 * time.Minute
)

// PasswordResetService defines the interface for password reset operations
type PasswordResetService interface {
	RequestReset(email string) error
	ConfirmReset(email, code, newPassword string) error
}

type passwordResetService struct {
	resetTokenRepo repository.ResetTokenRepository
	userRepo       repository.UserRepository
	emailService   email.EmailService
}

// NewPasswordResetService creates a new instance of PasswordResetService
func NewPasswordResetService(
	resetTokenRepo repository.ResetTokenRepository,
	userRepo repository.UserRepository,
	emailService email.EmailService,
) PasswordResetService {
	return &passwordResetService{
		resetTokenRepo: resetTokenRepo,
		userRepo:       userRepo,
		emailService:   emailService,
	}
}

// RequestReset generates a reset code and sends it via email.
func (s *passwordResetService) RequestReset(email string) error {
	// Check if user exists
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		// User doesn't exist
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return fmt.Errorf("failed to check user: %w", err)
	}

	// Generate 6-digit OTP code
	code, err := utils.GenerateOTP()
	if err != nil {
		return fmt.Errorf("failed to generate OTP: %w", err)
	}

	// Delete any existing tokens for this email
	if err := s.resetTokenRepo.DeleteByEmail(email); err != nil {
		return fmt.Errorf("failed to delete old tokens: %w", err)
	}

	// Create new reset token
	token := &model.ResetToken{
		Email:     email,
		Code:      code,
		ExpiresAt: time.Now().Add(ResetTokenExpiration),
	}

	if err := s.resetTokenRepo.Create(token); err != nil {
		return fmt.Errorf("failed to create reset token: %w", err)
	}

	// Send email with code
	if err := s.emailService.SendPasswordResetCode(user.Email, code); err != nil {
		return fmt.Errorf("failed to send reset email: %w", err)
	}

	return nil
}

// ConfirmReset validates the code and resets the user's password.
func (s *passwordResetService) ConfirmReset(email, code, newPassword string) error {
	// Validate new password first
	if err := validation.ValidatePassword(newPassword, email); err != nil {
		return err
	}

	// Find valid token
	token, err := s.resetTokenRepo.FindByEmailAndCode(email, code)
	if err != nil {
		// Token not found or expired
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("invalid or expired reset code")
		}
		return fmt.Errorf("failed to validate reset token: %w", err)
	}

	// Double-check expiration
	if token.IsExpired() {
		return errors.New("invalid or expired reset code")
	}

	// Find user
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("invalid or expired reset code")
		}
		return fmt.Errorf("failed to find user: %w", err)
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update user password
	user.PasswordHash = string(hashedPassword)
	if err := s.userRepo.Update(user); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Delete used token
	if err := s.resetTokenRepo.DeleteByEmail(email); err != nil {
		fmt.Printf("Failed to delete used reset token for %s: %v\n", email, err)
	}

	return nil
}
