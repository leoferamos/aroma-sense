package service

import (
	"fmt"
	"time"

	"github.com/leoferamos/aroma-sense/internal/model"
	"github.com/leoferamos/aroma-sense/internal/notification"
	"github.com/leoferamos/aroma-sense/internal/repository"
	"github.com/leoferamos/aroma-sense/internal/utils"
	"github.com/leoferamos/aroma-sense/internal/validation"
	"golang.org/x/crypto/bcrypt"
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
	notifier       notification.NotificationService
}

// NewPasswordResetService creates a new instance of PasswordResetService
func NewPasswordResetService(
	resetTokenRepo repository.ResetTokenRepository,
	userRepo repository.UserRepository,
	notifier notification.NotificationService,
) PasswordResetService {
	return &passwordResetService{
		resetTokenRepo: resetTokenRepo,
		userRepo:       userRepo,
		notifier:       notifier,
	}
}

// RequestReset generates a reset code and sends it via email.
func (s *passwordResetService) RequestReset(email string) error {
	// Check if user exists
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil
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

	// Send email with code via notifier
	if s.notifier != nil {
		if err := s.notifier.SendPasswordResetCode(user.Email, code); err != nil {
			return fmt.Errorf("failed to send reset email: %w", err)
		}
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
		// Generic error message
		return fmt.Errorf("invalid or expired reset code")
	}

	// Double-check expiration
	if token.IsExpired() {
		return fmt.Errorf("invalid or expired reset code")
	}

	// Find user
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return fmt.Errorf("invalid or expired reset code")
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
