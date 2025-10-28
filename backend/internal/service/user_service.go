package service

import (
	"errors"

	"github.com/leoferamos/aroma-sense/internal/auth"
	"github.com/leoferamos/aroma-sense/internal/dto"
	"github.com/leoferamos/aroma-sense/internal/model"
	"github.com/leoferamos/aroma-sense/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// UserService defines the interface for user-related business logic
type UserService interface {
	RegisterUser(input dto.CreateUserRequest) error
	Login(input dto.LoginRequest) (accessToken string, refreshToken string, user *model.User, err error)
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
	user.RefreshTokenHash = &refreshTokenHash
	user.RefreshTokenExpiresAt = &expiresAt
	if err := s.repo.Update(user); err != nil {
		return "", "", nil, errors.New("failed to save refresh token")
	}

	return accessToken, refreshToken, user, nil
}
