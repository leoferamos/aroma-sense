package service

import (
	"crypto/subtle"
	"errors"
	"strings"
	"time"

	"github.com/leoferamos/aroma-sense/internal/auth"
	"github.com/leoferamos/aroma-sense/internal/dto"
	"github.com/leoferamos/aroma-sense/internal/model"
	"github.com/leoferamos/aroma-sense/internal/repository"
	"github.com/leoferamos/aroma-sense/internal/validation"
	"golang.org/x/crypto/bcrypt"
)

// UserService defines the interface for user-related business logic
type UserService interface {
	RegisterUser(input dto.CreateUserRequest) error
	Login(input dto.LoginRequest) (accessToken string, refreshToken string, user *model.User, err error)
	RefreshAccessToken(refreshToken string) (accessToken string, newRefreshToken string, user *model.User, err error)
	InvalidateRefreshToken(refreshToken string) error
	GetByPublicID(publicID string) (*model.User, error)
	UpdateDisplayName(publicID string, displayName string) (*model.User, error)
}

type userService struct {
	repo        repository.UserRepository
	cartService CartService
}

// NewUserService creates a new instance of UserService
func NewUserService(repo repository.UserRepository, cartService CartService) UserService {
	return &userService{repo: repo, cartService: cartService}
}

// RegisterUser handles the business logic for user registration
func (s *userService) RegisterUser(input dto.CreateUserRequest) error {
	if err := validation.ValidatePassword(input.Password, input.Email); err != nil {
		return err
	}
	// Check if email already exists
	if _, err := s.repo.FindByEmail(input.Email); err == nil {
		return errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	user := model.User{
		Email:        input.Email,
		PasswordHash: string(hashedPassword),
	}

	// Create the user
	if err := s.repo.Create(&user); err != nil {
		return err
	}

	// Create a cart for the new user
	if err := s.cartService.CreateCartForUser(user.PublicID); err != nil {
		return errors.New("failed to create cart for user")
	}

	return nil
}

// Login handles the business logic for user login
func (s *userService) Login(input dto.LoginRequest) (string, string, *model.User, error) {
	user, err := s.repo.FindByEmail(input.Email)
	if err != nil {
		return "", "", nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return "", "", nil, errors.New("invalid credentials")
	}

	// Ensure user has a cart
	if err := s.cartService.CreateCartForUser(user.PublicID); err != nil {
		return "", "", nil, errors.New("failed to ensure cart exists")
	}

	// Generate access token
	accessToken, err := auth.GenerateJWT(user.PublicID, user.Role)
	if err != nil {
		return "", "", nil, errors.New("failed to generate access token")
	}

	// Generate refresh token
	refreshToken, expiresAt, err := auth.GenerateRefreshToken()
	if err != nil {
		return "", "", nil, errors.New("failed to generate refresh token")
	}

	// Save refresh token hash in DB
	refreshTokenHash := auth.HashRefreshToken(refreshToken)
	if err := s.repo.UpdateRefreshToken(user.ID, &refreshTokenHash, &expiresAt); err != nil {
		return "", "", nil, errors.New("failed to save refresh token")
	}

	return accessToken, refreshToken, user, nil
}

// RefreshAccessToken validates refresh token and generates new access token
func (s *userService) RefreshAccessToken(refreshToken string) (string, string, *model.User, error) {
	// Hash the refresh token to compare with DB
	refreshTokenHash := auth.HashRefreshToken(refreshToken)

	// Find user by refresh token hash
	user, err := s.repo.FindByRefreshTokenHash(refreshTokenHash)
	if err != nil {
		return "", "", nil, errors.New("invalid refresh token")
	}

	// Constant-time double-check to mitigate timing attacks
	if user.RefreshTokenHash == nil {
		return "", "", nil, errors.New("invalid refresh token")
	}
	if subtle.ConstantTimeCompare([]byte(*user.RefreshTokenHash), []byte(refreshTokenHash)) != 1 {
		return "", "", nil, errors.New("invalid refresh token")
	}

	// Check if refresh token is expired
	if user.RefreshTokenExpiresAt == nil || user.RefreshTokenExpiresAt.Before(time.Now()) {
		return "", "", nil, errors.New("refresh token expired")
	}

	// Generate new access token
	accessToken, err := auth.GenerateJWT(user.PublicID, user.Role)
	if err != nil {
		return "", "", nil, errors.New("failed to generate access token")
	}

	// Rotate refresh token
	newRefreshToken, newExpiresAt, err := auth.GenerateRefreshToken()
	if err != nil {
		return "", "", nil, errors.New("failed to generate refresh token")
	}
	newHash := auth.HashRefreshToken(newRefreshToken)
	if err := s.repo.UpdateRefreshToken(user.ID, &newHash, &newExpiresAt); err != nil {
		return "", "", nil, errors.New("failed to save refresh token")
	}

	return accessToken, newRefreshToken, user, nil
}

// InvalidateRefreshToken clears the stored refresh token for the owning user
func (s *userService) InvalidateRefreshToken(refreshToken string) error {
	if refreshToken == "" {
		return errors.New("missing refresh token")
	}
	hash := auth.HashRefreshToken(refreshToken)
	user, err := s.repo.FindByRefreshTokenHash(hash)
	if err != nil {
		return err
	}
	// constant-time double-check before invalidating
	if user.RefreshTokenHash == nil {
		return errors.New("invalid refresh token")
	}
	if subtle.ConstantTimeCompare([]byte(*user.RefreshTokenHash), []byte(hash)) != 1 {
		return errors.New("invalid refresh token")
	}
	return s.repo.UpdateRefreshToken(user.ID, nil, nil)
}

// GetByPublicID returns the user by public id
func (s *userService) GetByPublicID(publicID string) (*model.User, error) {
	if publicID == "" {
		return nil, errors.New("unauthenticated")
	}
	return s.repo.FindByPublicID(publicID)
}

// UpdateDisplayName updates the user's public display name
func (s *userService) UpdateDisplayName(publicID string, displayName string) (*model.User, error) {
	if publicID == "" {
		return nil, errors.New("unauthenticated")
	}
	trimmed := strings.TrimSpace(displayName)
	if len(trimmed) < 2 {
		return nil, errors.New("display_name too short")
	}
	if len(trimmed) > 50 {
		return nil, errors.New("display_name too long")
	}
	user, err := s.repo.FindByPublicID(publicID)
	if err != nil {
		return nil, err
	}
	dn := strings.TrimSpace(displayName)
	user.DisplayName = &dn
	if err := s.repo.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}
