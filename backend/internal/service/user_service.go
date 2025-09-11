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
	Login(input dto.LoginRequest) (string, *model.User, error)
}

type userService struct {
	repo repository.UserRepository
}

// NewUserService creates a new instance of UserService
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
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

	return s.repo.Create(&user)
}

// Login handles the business logic for user login
func (s *userService) Login(input dto.LoginRequest) (string, *model.User, error) {
	user, err := s.repo.FindByEmail(input.Email)
	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	token, err := auth.GenerateJWT(user.PublicID, user.Role)
	if err != nil {
		return "", nil, errors.New("failed to generate token")
	}

	return token, user, nil
}
