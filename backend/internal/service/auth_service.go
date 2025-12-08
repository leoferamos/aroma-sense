package service

import (
	"crypto/subtle"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/leoferamos/aroma-sense/internal/apperror"
	"github.com/leoferamos/aroma-sense/internal/auth"
	"github.com/leoferamos/aroma-sense/internal/dto"
	"github.com/leoferamos/aroma-sense/internal/model"
	"github.com/leoferamos/aroma-sense/internal/repository"
	"github.com/leoferamos/aroma-sense/internal/utils"
	"github.com/leoferamos/aroma-sense/internal/validation"
)

// AuthService defines the interface for authentication-related business logic
type AuthService interface {
	RegisterUser(input dto.CreateUserRequest) error
	Login(input dto.LoginRequest) (accessToken, refreshToken string, user *model.User, err error)
	RefreshAccessToken(refreshToken string) (accessToken, newRefreshToken string, user *model.User, err error)
	Logout(refreshToken string) error
	InvalidateRefreshToken(refreshToken string) error
}

type authService struct {
	repo            repository.UserRepository
	cartService     CartService
	auditLogService AuditLogService
}

// NewAuthService cria uma nova instância de AuthService
func NewAuthService(repo repository.UserRepository, cartService CartService, auditLogService AuditLogService) AuthService {
	return &authService{repo: repo, cartService: cartService, auditLogService: auditLogService}
}

// RegisterUser handles the business logic for user registration
func (s *authService) RegisterUser(input dto.CreateUserRequest) error {
	if err := validation.ValidatePassword(input.Password, input.Email); err != nil {
		return err
	}
	// Check if email already exists
	if _, err := s.repo.FindByEmail(input.Email); err == nil {
		return apperror.NewCodeMessage("email_already_registered", "email already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return apperror.NewCodeMessage("password_hash_failed", "failed to hash password")
	}

	user := model.User{
		Email:        input.Email,
		PasswordHash: string(hashedPassword),
	}

	// Create the user
	if err := s.repo.Create(&user); err != nil {
		return err
	}

	// Log successful user creation
	if s.auditLogService != nil {
		s.auditLogService.LogUserAction(nil, &user.ID, model.AuditActionUserCreated,
			map[string]interface{}{
				"email": utils.MaskEmail(user.Email),
				"role":  user.Role,
			})
	}

	// Create a cart for the new user
	if err := s.cartService.CreateCartForUser(user.PublicID); err != nil {
		return apperror.NewCodeMessage("cart_create_failed", "failed to create cart for user")
	}

	return nil
}

// Login handles the business logic for user login
func (s *authService) Login(input dto.LoginRequest) (string, string, *model.User, error) {
	user, err := s.repo.FindByEmail(input.Email)
	if err != nil {
		// Log failed login attempt
		if s.auditLogService != nil {
			s.auditLogService.LogUserAction(nil, nil, model.AuditActionUserLogin,
				map[string]interface{}{
					"identifier_hash": utils.HashEmailForLogging(input.Email),
					"success":         false,
					"reason":          "invalid_credentials",
				})
		}
		return "", "", nil, apperror.NewCodeMessage("invalid_credentials", "invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		// Log failed login attempt
		if s.auditLogService != nil {
			s.auditLogService.LogUserAction(nil, &user.ID, model.AuditActionUserLogin,
				map[string]interface{}{
					"user_id": user.ID,
					"success": false,
					"reason":  "invalid_password",
				})
		}
		return "", "", nil, apperror.NewCodeMessage("invalid_credentials", "invalid credentials")
	}

	// Ensure user has a cart
	if err := s.cartService.CreateCartForUser(user.PublicID); err != nil {
		return "", "", nil, apperror.NewCodeMessage("cart_create_failed", "failed to ensure cart exists")
	}

	// Generate access token
	accessToken, err := auth.GenerateJWT(user.PublicID, user.Role)
	if err != nil {
		return "", "", nil, apperror.NewCodeMessage("access_token_failed", "failed to generate access token")
	}

	// Generate refresh token
	refreshToken, expiresAt, err := auth.GenerateRefreshToken()
	if err != nil {
		return "", "", nil, apperror.NewCodeMessage("refresh_token_failed", "failed to generate refresh token")
	}

	// Save refresh token hash in DB
	refreshTokenHash := auth.HashRefreshToken(refreshToken)
	if err := s.repo.UpdateRefreshToken(user.ID, &refreshTokenHash, &expiresAt); err != nil {
		return "", "", nil, apperror.NewCodeMessage("refresh_token_save_failed", "failed to save refresh token")
	}
	user.RefreshTokenExpiresAt = &expiresAt

	// Log successful login
	if s.auditLogService != nil {
		s.auditLogService.LogUserAction(nil, &user.ID, model.AuditActionUserLogin,
			map[string]interface{}{
				"email":   utils.MaskEmail(user.Email),
				"success": true,
			})
	}

	return accessToken, refreshToken, user, nil
}

// RefreshAccessToken validates refresh token and generates new access token
func (s *authService) RefreshAccessToken(refreshToken string) (string, string, *model.User, error) {
	refreshTokenHash := auth.HashRefreshToken(refreshToken)

	user, err := s.repo.FindByRefreshTokenHash(refreshTokenHash)
	if err != nil {
		return "", "", nil, apperror.NewCodeMessage("invalid_refresh_token", "invalid refresh token")
	}

	if user.RefreshTokenHash == nil {
		return "", "", nil, apperror.NewCodeMessage("invalid_refresh_token", "invalid refresh token")
	}
	if subtle.ConstantTimeCompare([]byte(*user.RefreshTokenHash), []byte(refreshTokenHash)) != 1 {
		return "", "", nil, apperror.NewCodeMessage("invalid_refresh_token", "invalid refresh token")
	}

	if user.RefreshTokenExpiresAt == nil || user.RefreshTokenExpiresAt.Before(time.Now()) {
		return "", "", nil, apperror.NewCodeMessage("refresh_token_expired", "refresh token expired")
	}

	accessToken, err := auth.GenerateJWT(user.PublicID, user.Role)
	if err != nil {
		return "", "", nil, apperror.NewCodeMessage("access_token_failed", "failed to generate access token")
	}

	newRefreshToken, newExpiresAt, err := auth.GenerateRefreshToken()
	if err != nil {
		return "", "", nil, apperror.NewCodeMessage("refresh_token_failed", "failed to generate refresh token")
	}
	newHash := auth.HashRefreshToken(newRefreshToken)
	if err := s.repo.UpdateRefreshToken(user.ID, &newHash, &newExpiresAt); err != nil {
		return "", "", nil, apperror.NewCodeMessage("refresh_token_save_failed", "failed to save refresh token")
	}

	return accessToken, newRefreshToken, user, nil
}

// Logout handles user logout by invalidating refresh token
func (s *authService) Logout(refreshToken string) error {
	if refreshToken == "" {
		return apperror.NewCodeMessage("refresh_token_missing", "no refresh token provided")
	}

	hash := auth.HashRefreshToken(refreshToken)
	user, err := s.repo.FindByRefreshTokenHash(hash)
	if err == nil {
		if s.auditLogService != nil {
			s.auditLogService.LogUserAction(nil, &user.ID, model.AuditActionUserLogout,
				map[string]interface{}{
					"email": utils.MaskEmail(user.Email),
				})
		}
	}

	if err := s.repo.UpdateRefreshToken(user.ID, nil, nil); err != nil {
		return apperror.NewDomain(err, "refresh_token_save_failed", "failed to save refresh token")
	}
	return nil
}

// InvalidateRefreshToken invalidates a refresh token by clearing it from the database
func (s *authService) InvalidateRefreshToken(refreshToken string) error {
	if refreshToken == "" {
		return apperror.NewCodeMessage("refresh_token_missing", "no refresh token provided")
	}

	hash := auth.HashRefreshToken(refreshToken)
	user, err := s.repo.FindByRefreshTokenHash(hash)
	if err != nil {
		return apperror.NewDomain(fmt.Errorf("invalid refresh token: %w", err), "invalid_refresh_token", "invalid refresh token")
	}

	if err := s.repo.UpdateRefreshToken(user.ID, nil, nil); err != nil {
		return apperror.NewDomain(err, "refresh_token_save_failed", "failed to save refresh token")
	}
	return nil
}
