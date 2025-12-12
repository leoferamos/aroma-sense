package service_test

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/leoferamos/aroma-sense/internal/auth"
	"github.com/leoferamos/aroma-sense/internal/dto"
	"github.com/leoferamos/aroma-sense/internal/model"
	"github.com/leoferamos/aroma-sense/internal/repository"
	service "github.com/leoferamos/aroma-sense/internal/service/auth"
	cartservice "github.com/leoferamos/aroma-sense/internal/service/cart"
	logservice "github.com/leoferamos/aroma-sense/internal/service/log"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

// Additional mocks for auth service
type mockCartService struct {
	createCartForUserErr error
}

func (m *mockCartService) CreateCartForUser(userID string) error {
	return m.createCartForUserErr
}

func (m *mockCartService) GetCartByUserID(userID string) (*model.Cart, error) {
	return nil, nil
}

func (m *mockCartService) GetCartResponse(userID string) (*dto.CartResponse, error) {
	return nil, nil
}

func (m *mockCartService) AddItemToCart(userID string, productSlug string, quantity int) (*dto.CartResponse, error) {
	return nil, nil
}

func (m *mockCartService) UpdateItemQuantity(userID string, itemID uint, quantity int) (*dto.CartResponse, error) {
	return nil, nil
}

func (m *mockCartService) UpdateItemQuantityBySlug(userID string, productSlug string, quantity int) (*dto.CartResponse, error) {
	return nil, nil
}

func (m *mockCartService) RemoveItem(userID string, itemID uint) (*dto.CartResponse, error) {
	return nil, nil
}

func (m *mockCartService) RemoveItemBySlug(userID string, productSlug string) (*dto.CartResponse, error) {
	return nil, nil
}

func (m *mockCartService) ClearCart(userID string) (*dto.CartResponse, error) {
	return nil, nil
}

type mockUserRepo struct {
	createErr                  error
	findByEmailUser            *model.User
	findByEmailErr             error
	findByRefreshTokenHashUser *model.User
	findByRefreshTokenHashErr  error
	findByPublicIDUser         *model.User
	findByPublicIDErr          error
	updateRefreshTokenErr      error
}

func (m *mockUserRepo) Create(user *model.User) error {
	return m.createErr
}

func (m *mockUserRepo) FindByEmail(email string) (*model.User, error) {
	return m.findByEmailUser, m.findByEmailErr
}

func (m *mockUserRepo) FindByRefreshTokenHash(hash string) (*model.User, error) {
	return m.findByRefreshTokenHashUser, m.findByRefreshTokenHashErr
}

func (m *mockUserRepo) FindByPublicID(publicID string) (*model.User, error) {
	return m.findByPublicIDUser, m.findByPublicIDErr
}

func (m *mockUserRepo) Update(user *model.User) error {
	return nil
}

func (m *mockUserRepo) UpdateRefreshToken(userID uint, hash *string, expiresAt *time.Time) error {
	return m.updateRefreshTokenErr
}

func (m *mockUserRepo) ListUsers(limit int, offset int, filters map[string]interface{}) ([]*model.User, int64, error) {
	return nil, 0, nil
}

func (m *mockUserRepo) FindByID(id uint) (*model.User, error) {
	return nil, nil
}

func (m *mockUserRepo) UpdateRole(userID uint, newRole string) error {
	return nil
}

func (m *mockUserRepo) DeactivateUser(userID uint, adminPublicID string, deactivatedAt time.Time, reason string, notes string, suspensionUntil *time.Time) error {
	return nil
}

func (m *mockUserRepo) RequestAccountDeletion(publicID string, requestedAt time.Time) error {
	return nil
}

func (m *mockUserRepo) ConfirmAccountDeletion(publicID string, confirmedAt time.Time) error {
	return nil
}

func (m *mockUserRepo) HasActiveDependencies(publicID string) (bool, error) {
	return false, nil
}

func (m *mockUserRepo) AnonymizeUser(publicID string, anonymizedEmail string, anonymizedDisplayName string) error {
	return nil
}

func (m *mockUserRepo) FindExpiredUsersForAnonymization() ([]*model.User, error) {
	return nil, nil
}

func (m *mockUserRepo) FindUsersPendingAutoConfirm(cutoff time.Time) ([]*model.User, error) {
	return nil, nil
}

func (m *mockUserRepo) DeleteByPublicID(publicID string) error {
	return nil
}

type mockAuditLogService struct{}

func (m *mockAuditLogService) LogUserAction(actorID *uint, userID *uint, action model.AuditAction, details map[string]interface{}) error {
	return nil
}

func (m *mockAuditLogService) LogUserUpdate(actorID uint, userID uint, oldUser, newUser *model.User) error {
	return nil
}

func (m *mockAuditLogService) LogAdminAction(adminID uint, userID uint, action model.AuditAction, details map[string]interface{}) error {
	return nil
}

func (m *mockAuditLogService) LogSystemAction(action model.AuditAction, resource, resourceID string, details map[string]interface{}) error {
	return nil
}

func (m *mockAuditLogService) LogDataAccess(userID uint, resource string, resourceID string) error {
	return nil
}

func (m *mockAuditLogService) LogDeletionAction(actorID *uint, userID uint, action model.AuditAction, details map[string]interface{}) error {
	return nil
}

func (m *mockAuditLogService) ListAuditLogs(filter *model.AuditLogFilter) ([]*model.AuditLog, int64, error) {
	return nil, 0, nil
}

func (m *mockAuditLogService) GetAuditLogByID(id uint) (*model.AuditLog, error) {
	return nil, nil
}

func (m *mockAuditLogService) GetUserAuditLogs(userID uint, limit, offset int) ([]*model.AuditLog, int64, error) {
	return nil, 0, nil
}

func (m *mockAuditLogService) GetResourceAuditLogs(resource, resourceID string) ([]*model.AuditLog, error) {
	return nil, nil
}

func (m *mockAuditLogService) GetAuditSummary(startDate, endDate *time.Time) (*model.AuditLogSummary, error) {
	return nil, nil
}

func (m *mockAuditLogService) CleanupOldLogs(retentionDays int) error {
	return nil
}

func (m *mockAuditLogService) ConvertAuditLogToResponse(auditLog *model.AuditLog) dto.AuditLogResponse {
	return dto.AuditLogResponse{}
}

func (m *mockAuditLogService) ConvertAuditLogToResponseDetailed(auditLog *model.AuditLog) dto.AuditLogResponse {
	return dto.AuditLogResponse{}
}

func (m *mockAuditLogService) ConvertAuditLogsToResponse(auditLogs []*model.AuditLog) []dto.AuditLogResponse {
	return nil
}

func (m *mockAuditLogService) ConvertAuditLogSummaryToResponse(summary *model.AuditLogSummary) dto.AuditLogSummaryResponse {
	return dto.AuditLogSummaryResponse{}
}

// Tests for AuthService
func TestAuthService_RegisterUser(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key-for-jwt-generation")
	tests := []struct {
		name            string
		input           dto.CreateUserRequest
		mockSetup       func() (repository.UserRepository, cartservice.CartService, logservice.AuditLogService)
		expectedErr     bool
		expectedErrCode string
	}{
		{
			name: "success",
			input: dto.CreateUserRequest{
				Email:    "newuser@example.com",
				Password: "ValidPass123!",
			},
			mockSetup: func() (repository.UserRepository, cartservice.CartService, logservice.AuditLogService) {
				userRepo := &mockUserRepo{
					findByEmailErr: errors.New("not found"),
				}
				cartSvc := &mockCartService{}
				auditSvc := &mockAuditLogService{}
				return userRepo, cartSvc, auditSvc
			},
			expectedErr: false,
		},
		{
			name: "email already registered",
			input: dto.CreateUserRequest{
				Email:    "existing@example.com",
				Password: "ValidPass123!",
			},
			mockSetup: func() (repository.UserRepository, cartservice.CartService, logservice.AuditLogService) {
				userRepo := &mockUserRepo{
					findByEmailUser: &model.User{Email: "existing@example.com"},
				}
				cartSvc := &mockCartService{}
				auditSvc := &mockAuditLogService{}
				return userRepo, cartSvc, auditSvc
			},
			expectedErr:     true,
			expectedErrCode: "email already registered",
		},
		{
			name: "invalid password",
			input: dto.CreateUserRequest{
				Email:    "newuser@example.com",
				Password: "weak",
			},
			mockSetup: func() (repository.UserRepository, cartservice.CartService, logservice.AuditLogService) {
				userRepo := &mockUserRepo{
					findByEmailErr: errors.New("not found"),
				}
				cartSvc := &mockCartService{}
				auditSvc := &mockAuditLogService{}
				return userRepo, cartSvc, auditSvc
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo, cartSvc, auditSvc := tt.mockSetup()
			authSvc := service.NewAuthService(userRepo, cartSvc, auditSvc)

			err := authSvc.RegisterUser(tt.input)

			if tt.expectedErr {
				assert.Error(t, err)
				if tt.expectedErrCode != "" {
					assert.Contains(t, err.Error(), tt.expectedErrCode)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key-for-jwt-generation")
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("ValidPass123!"), bcrypt.DefaultCost)

	tests := []struct {
		name            string
		input           dto.LoginRequest
		mockSetup       func() (repository.UserRepository, cartservice.CartService, logservice.AuditLogService)
		expectedErr     bool
		expectedErrCode string
	}{
		{
			name: "success",
			input: dto.LoginRequest{
				Email:    "user@example.com",
				Password: "ValidPass123!",
			},
			mockSetup: func() (repository.UserRepository, cartservice.CartService, logservice.AuditLogService) {
				user := &model.User{
					ID:           1,
					PublicID:     "user123",
					Email:        "user@example.com",
					PasswordHash: string(hashedPassword),
					Role:         "client",
				}
				userRepo := &mockUserRepo{
					findByEmailUser: user,
				}
				cartSvc := &mockCartService{}
				auditSvc := &mockAuditLogService{}
				return userRepo, cartSvc, auditSvc
			},
			expectedErr: false,
		},
		{
			name: "invalid credentials - user not found",
			input: dto.LoginRequest{
				Email:    "nonexistent@example.com",
				Password: "ValidPass123!",
			},
			mockSetup: func() (repository.UserRepository, cartservice.CartService, logservice.AuditLogService) {
				userRepo := &mockUserRepo{
					findByEmailErr: errors.New("not found"),
				}
				cartSvc := &mockCartService{}
				auditSvc := &mockAuditLogService{}
				return userRepo, cartSvc, auditSvc
			},
			expectedErr:     true,
			expectedErrCode: "invalid credentials",
		},
		{
			name: "invalid password",
			input: dto.LoginRequest{
				Email:    "user@example.com",
				Password: "WrongPass123!",
			},
			mockSetup: func() (repository.UserRepository, cartservice.CartService, logservice.AuditLogService) {
				user := &model.User{
					ID:           1,
					PublicID:     "user123",
					Email:        "user@example.com",
					PasswordHash: string(hashedPassword),
					Role:         "client",
				}
				userRepo := &mockUserRepo{
					findByEmailUser: user,
				}
				cartSvc := &mockCartService{}
				auditSvc := &mockAuditLogService{}
				return userRepo, cartSvc, auditSvc
			},
			expectedErr:     true,
			expectedErrCode: "invalid credentials",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo, cartSvc, auditSvc := tt.mockSetup()
			authSvc := service.NewAuthService(userRepo, cartSvc, auditSvc)

			accessToken, refreshToken, user, err := authSvc.Login(tt.input)

			if tt.expectedErr {
				assert.Error(t, err)
				if tt.expectedErrCode != "" {
					assert.Contains(t, err.Error(), tt.expectedErrCode)
				}
				assert.Empty(t, accessToken)
				assert.Empty(t, refreshToken)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, accessToken)
				assert.NotEmpty(t, refreshToken)
				assert.NotNil(t, user)
			}
		})
	}
}

func TestAuthService_RefreshAccessToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key-for-jwt-generation")
	refreshToken, expiresAt, _ := auth.GenerateRefreshToken()
	hash := auth.HashRefreshToken(refreshToken)

	tests := []struct {
		name            string
		refreshToken    string
		mockSetup       func() repository.UserRepository
		expectedErr     bool
		expectedErrCode string
	}{
		{
			name:         "success",
			refreshToken: refreshToken,
			mockSetup: func() repository.UserRepository {
				user := &model.User{
					ID:                    1,
					PublicID:              "user123",
					RefreshTokenHash:      &hash,
					RefreshTokenExpiresAt: &expiresAt,
				}
				return &mockUserRepo{
					findByRefreshTokenHashUser: user,
				}
			},
			expectedErr: false,
		},
		{
			name:         "invalid refresh token",
			refreshToken: "invalid",
			mockSetup: func() repository.UserRepository {
				return &mockUserRepo{
					findByRefreshTokenHashErr: errors.New("not found"),
				}
			},
			expectedErr:     true,
			expectedErrCode: "invalid refresh token",
		},
		{
			name:         "expired token",
			refreshToken: refreshToken,
			mockSetup: func() repository.UserRepository {
				pastTime := time.Now().Add(-time.Hour)
				user := &model.User{
					ID:                    1,
					PublicID:              "user123",
					RefreshTokenHash:      &hash,
					RefreshTokenExpiresAt: &pastTime,
				}
				return &mockUserRepo{
					findByRefreshTokenHashUser: user,
				}
			},
			expectedErr:     true,
			expectedErrCode: "refresh token expired",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := tt.mockSetup()
			cartSvc := &mockCartService{}
			auditSvc := &mockAuditLogService{}
			authSvc := service.NewAuthService(userRepo, cartSvc, auditSvc)

			accessToken, newRefreshToken, user, err := authSvc.RefreshAccessToken(tt.refreshToken)

			if tt.expectedErr {
				assert.Error(t, err)
				if tt.expectedErrCode != "" {
					assert.Contains(t, err.Error(), tt.expectedErrCode)
				}
				assert.Empty(t, accessToken)
				assert.Empty(t, newRefreshToken)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, accessToken)
				assert.NotEmpty(t, newRefreshToken)
				assert.NotNil(t, user)
			}
		})
	}
}

func TestAuthService_Logout(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key-for-jwt-generation")
	tests := []struct {
		name         string
		refreshToken string
		mockSetup    func() repository.UserRepository
		expectedErr  bool
	}{
		{
			name:         "success",
			refreshToken: "valid_token",
			mockSetup: func() repository.UserRepository {
				user := &model.User{ID: 1, PublicID: "user123"}
				return &mockUserRepo{
					findByRefreshTokenHashUser: user,
				}
			},
			expectedErr: false,
		},
		{
			name:         "missing token",
			refreshToken: "",
			mockSetup: func() repository.UserRepository {
				return &mockUserRepo{}
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := tt.mockSetup()
			cartSvc := &mockCartService{}
			auditSvc := &mockAuditLogService{}
			authSvc := service.NewAuthService(userRepo, cartSvc, auditSvc)

			err := authSvc.Logout(tt.refreshToken)

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAuthService_InvalidateRefreshToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key-for-jwt-generation")
	tests := []struct {
		name            string
		refreshToken    string
		mockSetup       func() repository.UserRepository
		expectedErr     bool
		expectedErrCode string
	}{
		{
			name:         "success",
			refreshToken: "valid_token",
			mockSetup: func() repository.UserRepository {
				user := &model.User{ID: 1, PublicID: "user123"}
				return &mockUserRepo{
					findByRefreshTokenHashUser: user,
				}
			},
			expectedErr: false,
		},
		{
			name:         "missing token",
			refreshToken: "",
			mockSetup: func() repository.UserRepository {
				return &mockUserRepo{}
			},
			expectedErr:     true,
			expectedErrCode: "no refresh token provided",
		},
		{
			name:         "invalid token",
			refreshToken: "invalid",
			mockSetup: func() repository.UserRepository {
				return &mockUserRepo{
					findByRefreshTokenHashErr: errors.New("not found"),
				}
			},
			expectedErr:     true,
			expectedErrCode: "invalid refresh token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepo := tt.mockSetup()
			cartSvc := &mockCartService{}
			auditSvc := &mockAuditLogService{}
			authSvc := service.NewAuthService(userRepo, cartSvc, auditSvc)

			err := authSvc.InvalidateRefreshToken(tt.refreshToken)

			if tt.expectedErr {
				assert.Error(t, err)
				if tt.expectedErrCode != "" {
					assert.Contains(t, err.Error(), tt.expectedErrCode)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
