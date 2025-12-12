package service

import (
	"errors"
	"testing"
	"time"

	"github.com/leoferamos/aroma-sense/internal/dto"
	"github.com/leoferamos/aroma-sense/internal/model"
	"github.com/stretchr/testify/assert"
)

type mockUserRepo struct {
	user            *model.User
	err             error
	hasDep          bool
	hasDepErr       error
	reqDelErr       error
	confDelErr      error
	updateErr       error
	anonymErr       error
	usersPending    []*model.User
	usersPendingErr error
	usersExpired    []*model.User
	usersExpiredErr error
}

func (m *mockUserRepo) FindByPublicID(publicID string) (*model.User, error) {
	return m.user, m.err
}
func (m *mockUserRepo) FindByEmail(email string) (*model.User, error)           { return m.user, m.err }
func (m *mockUserRepo) FindByRefreshTokenHash(hash string) (*model.User, error) { return m.user, m.err }
func (m *mockUserRepo) FindByID(id uint) (*model.User, error)                   { return m.user, m.err }
func (m *mockUserRepo) HasActiveDependencies(publicID string) (bool, error) {
	return m.hasDep, m.hasDepErr
}
func (m *mockUserRepo) RequestAccountDeletion(publicID string, t time.Time) error { return m.reqDelErr }
func (m *mockUserRepo) ConfirmAccountDeletion(publicID string, t time.Time) error {
	return m.confDelErr
}
func (m *mockUserRepo) Update(user *model.User) error { return m.updateErr }
func (m *mockUserRepo) UpdateRefreshToken(userID uint, hash *string, expiresAt *time.Time) error {
	return nil
}
func (m *mockUserRepo) UpdateRole(userID uint, newRole string) error { return nil }
func (m *mockUserRepo) ListUsers(limit int, offset int, filters map[string]interface{}) ([]*model.User, int64, error) {
	return nil, 0, nil
}
func (m *mockUserRepo) AnonymizeUser(publicID, email, displayName string) error { return m.anonymErr }
func (m *mockUserRepo) FindUsersPendingAutoConfirm(cutoff time.Time) ([]*model.User, error) {
	return m.usersPending, m.usersPendingErr
}
func (m *mockUserRepo) FindExpiredUsersForAnonymization() ([]*model.User, error) {
	return m.usersExpired, m.usersExpiredErr
}
func (m *mockUserRepo) Create(user *model.User) error { return nil }
func (m *mockUserRepo) DeactivateUser(userID uint, adminPublicID string, deactivatedAt time.Time, reason string, notes string, suspensionUntil *time.Time) error {
	return nil
}
func (m *mockUserRepo) DeleteByPublicID(publicID string) error { return nil }

type mockUserContestationRepo struct {
	createErr error
}

func (m *mockUserContestationRepo) Create(contest *model.UserContestation) error { return m.createErr }
func (m *mockUserContestationRepo) FindByID(id uint) (*model.UserContestation, error) {
	return nil, nil
}
func (m *mockUserContestationRepo) FindByUserID(userID uint) ([]model.UserContestation, error) {
	return nil, nil
}
func (m *mockUserContestationRepo) ListPending(limit, offset int) ([]model.UserContestation, int64, error) {
	return nil, 0, nil
}
func (m *mockUserContestationRepo) Update(contest *model.UserContestation) error { return nil }

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
func (m *mockAuditLogService) GetAuditLogByID(id uint) (*model.AuditLog, error) { return nil, nil }
func (m *mockAuditLogService) GetUserAuditLogs(userID uint, limit, offset int) ([]*model.AuditLog, int64, error) {
	return nil, 0, nil
}
func (m *mockAuditLogService) GetResourceAuditLogs(resource, resourceID string) ([]*model.AuditLog, error) {
	return nil, nil
}
func (m *mockAuditLogService) GetAuditSummary(startDate, endDate *time.Time) (*model.AuditLogSummary, error) {
	return nil, nil
}
func (m *mockAuditLogService) CleanupOldLogs(retentionDays int) error { return nil }
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

type mockNotifier struct {
	err error
}

func (m *mockNotifier) SendPasswordResetCode(to, code string) error               { return m.err }
func (m *mockNotifier) SendWelcomeEmail(to, name string) error                    { return m.err }
func (m *mockNotifier) SendOrderConfirmation(to string, order *model.Order) error { return m.err }
func (m *mockNotifier) SendAccountDeactivated(to, reason string, contestationDeadline string) error {
	return m.err
}
func (m *mockNotifier) SendContestationReceived(to string) error { return m.err }
func (m *mockNotifier) SendContestationResult(to string, approved bool, reason string) error {
	return m.err
}
func (m *mockNotifier) SendDeletionRequested(to string, cancelLink string) error { return m.err }
func (m *mockNotifier) SendDeletionAutoConfirmed(to string) error                { return m.err }
func (m *mockNotifier) SendDeletionCancelled(to string) error                    { return m.err }
func (m *mockNotifier) SendDataAnonymized(to string) error                       { return m.err }
func (m *mockNotifier) SendPromotional(to, subject, htmlBody string) error       { return m.err }

// --- Test helpers: create a base user for tests ---
func baseUser() *model.User {
	now := time.Now().Add(-10 * 24 * time.Hour)
	displayName := "Test User"
	return &model.User{
		ID:          1,
		PublicID:    "publicid",
		Email:       "test@example.com",
		Role:        "user",
		DisplayName: &displayName,
		CreatedAt:   now,
		LastLoginAt: &now,
	}
}

// --- Tests: covers all public methods and error branches ---
func TestExportUserData(t *testing.T) {
	svc := NewLgpdService(&mockUserRepo{user: baseUser()}, &mockUserContestationRepo{}, &mockAuditLogService{}, &mockNotifier{})
	resp, err := svc.ExportUserData("publicid")
	assert.NoError(t, err)
	assert.Equal(t, "publicid", resp.PublicID)
}

func TestRequestAccountDeletion(t *testing.T) {
	user := baseUser()
	svc := NewLgpdService(&mockUserRepo{user: user}, &mockUserContestationRepo{}, &mockAuditLogService{}, &mockNotifier{})
	err := svc.RequestAccountDeletion("publicid")
	assert.NoError(t, err)
	// error: empty publicID
	err = svc.RequestAccountDeletion("")
	assert.Error(t, err)
	// error: user has active dependencies
	svc = NewLgpdService(&mockUserRepo{user: user, hasDep: true}, &mockUserContestationRepo{}, &mockAuditLogService{}, &mockNotifier{})
	err = svc.RequestAccountDeletion("publicid")
	assert.Error(t, err)
	// error: failed to check dependencies
	svc = NewLgpdService(&mockUserRepo{user: user, hasDepErr: errors.New("fail")}, &mockUserContestationRepo{}, &mockAuditLogService{}, &mockNotifier{})
	err = svc.RequestAccountDeletion("publicid")
	assert.Error(t, err)
	// error: deletion already requested
	u2 := baseUser()
	now := time.Now()
	u2.DeletionRequestedAt = &now
	svc = NewLgpdService(&mockUserRepo{user: u2}, &mockUserContestationRepo{}, &mockAuditLogService{}, &mockNotifier{})
	err = svc.RequestAccountDeletion("publicid")
	assert.Error(t, err)
	// error: user not found
	svc = NewLgpdService(&mockUserRepo{err: errors.New("fail")}, &mockUserContestationRepo{}, &mockAuditLogService{}, &mockNotifier{})
	err = svc.RequestAccountDeletion("publicid")
	assert.Error(t, err)
	// error: failed to request deletion
	svc = NewLgpdService(&mockUserRepo{user: user, reqDelErr: errors.New("fail")}, &mockUserContestationRepo{}, &mockAuditLogService{}, &mockNotifier{})
	err = svc.RequestAccountDeletion("publicid")
	assert.Error(t, err)
}

func TestConfirmAccountDeletion(t *testing.T) {
	now := time.Now().Add(-8 * 24 * time.Hour)
	user := baseUser()
	user.DeletionRequestedAt = &now
	svc := NewLgpdService(&mockUserRepo{user: user}, &mockUserContestationRepo{}, &mockAuditLogService{}, &mockNotifier{})
	err := svc.ConfirmAccountDeletion("publicid")
	assert.NoError(t, err)
	// error: empty publicID
	err = svc.ConfirmAccountDeletion("")
	assert.Error(t, err)
	// error: user not found
	svc = NewLgpdService(&mockUserRepo{err: errors.New("fail")}, &mockUserContestationRepo{}, &mockAuditLogService{}, &mockNotifier{})
	err = svc.ConfirmAccountDeletion("publicid")
	assert.Error(t, err)
	// error: deletion not requested
	svc = NewLgpdService(&mockUserRepo{user: baseUser()}, &mockUserContestationRepo{}, &mockAuditLogService{}, &mockNotifier{})
	err = svc.ConfirmAccountDeletion("publicid")
	assert.Error(t, err)
	// error: cooling off period not expired
	n2 := time.Now()
	u2 := baseUser()
	u2.DeletionRequestedAt = &n2
	svc = NewLgpdService(&mockUserRepo{user: u2}, &mockUserContestationRepo{}, &mockAuditLogService{}, &mockNotifier{})
	err = svc.ConfirmAccountDeletion("publicid")
	assert.Error(t, err)
	// error: failed to confirm deletion
	n3 := time.Now().Add(-8 * 24 * time.Hour)
	u3 := baseUser()
	u3.DeletionRequestedAt = &n3
	svc = NewLgpdService(&mockUserRepo{user: u3, confDelErr: errors.New("fail")}, &mockUserContestationRepo{}, &mockAuditLogService{}, &mockNotifier{})
	err = svc.ConfirmAccountDeletion("publicid")
	assert.Error(t, err)
}

func TestCancelAccountDeletion(t *testing.T) {
	now := time.Now()
	user := baseUser()
	user.DeletionRequestedAt = &now
	svc := NewLgpdService(&mockUserRepo{user: user}, &mockUserContestationRepo{}, &mockAuditLogService{}, &mockNotifier{})
	err := svc.CancelAccountDeletion("publicid")
	assert.NoError(t, err)
	// error: empty publicID
	err = svc.CancelAccountDeletion("")
	assert.Error(t, err)
	// error: user not found
	svc = NewLgpdService(&mockUserRepo{err: errors.New("fail")}, &mockUserContestationRepo{}, &mockAuditLogService{}, &mockNotifier{})
	err = svc.CancelAccountDeletion("publicid")
	assert.Error(t, err)
	// error: deletion not requested
	svc = NewLgpdService(&mockUserRepo{user: baseUser()}, &mockUserContestationRepo{}, &mockAuditLogService{}, &mockNotifier{})
	err = svc.CancelAccountDeletion("publicid")
	assert.Error(t, err)
	// error: failed to update user
	u2 := baseUser()
	u2.DeletionRequestedAt = &now
	svc = NewLgpdService(&mockUserRepo{user: u2, updateErr: errors.New("fail")}, &mockUserContestationRepo{}, &mockAuditLogService{}, &mockNotifier{})
	err = svc.CancelAccountDeletion("publicid")
	assert.Error(t, err)
}

func TestAnonymizeExpiredUser(t *testing.T) {
	now := time.Now().Add(-6 * 365 * 24 * time.Hour)
	user := baseUser()
	user.DeletionConfirmedAt = &now
	svc := NewLgpdService(&mockUserRepo{user: user}, &mockUserContestationRepo{}, &mockAuditLogService{}, &mockNotifier{})
	err := svc.AnonymizeExpiredUser("publicid")
	assert.NoError(t, err)
	// error: user not found
	svc = NewLgpdService(&mockUserRepo{err: errors.New("fail")}, &mockUserContestationRepo{}, &mockAuditLogService{}, &mockNotifier{})
	err = svc.AnonymizeExpiredUser("publicid")
	assert.Error(t, err)
	// error: deletion not confirmed
	svc = NewLgpdService(&mockUserRepo{user: baseUser()}, &mockUserContestationRepo{}, &mockAuditLogService{}, &mockNotifier{})
	err = svc.AnonymizeExpiredUser("publicid")
	assert.Error(t, err)
	// error: retention period not expired
	n2 := time.Now().Add(-2 * 365 * 24 * time.Hour)
	u2 := baseUser()
	u2.DeletionConfirmedAt = &n2
	svc = NewLgpdService(&mockUserRepo{user: u2}, &mockUserContestationRepo{}, &mockAuditLogService{}, &mockNotifier{})
	err = svc.AnonymizeExpiredUser("publicid")
	assert.Error(t, err)
	// error: failed to anonymize user
	n3 := time.Now().Add(-6 * 365 * 24 * time.Hour)
	u3 := baseUser()
	u3.DeletionConfirmedAt = &n3
	svc = NewLgpdService(&mockUserRepo{user: u3, anonymErr: errors.New("fail")}, &mockUserContestationRepo{}, &mockAuditLogService{}, &mockNotifier{})
	err = svc.AnonymizeExpiredUser("publicid")
	assert.Error(t, err)
}

func TestRequestContestation(t *testing.T) {
	now := time.Now().Add(-2 * 24 * time.Hour)
	user := baseUser()
	user.DeactivatedAt = &now
	svc := NewLgpdService(&mockUserRepo{user: user}, &mockUserContestationRepo{}, &mockAuditLogService{}, &mockNotifier{})
	err := svc.RequestContestation("publicid", "motivo")
	assert.NoError(t, err)
	// error: empty publicID
	err = svc.RequestContestation("", "motivo")
	assert.Error(t, err)
	// error: user not found
	svc = NewLgpdService(&mockUserRepo{err: errors.New("fail")}, &mockUserContestationRepo{}, &mockAuditLogService{}, &mockNotifier{})
	err = svc.RequestContestation("publicid", "motivo")
	assert.Error(t, err)
	// error: user is not deactivated
	svc = NewLgpdService(&mockUserRepo{user: baseUser()}, &mockUserContestationRepo{}, &mockAuditLogService{}, &mockNotifier{})
	err = svc.RequestContestation("publicid", "motivo")
	assert.Error(t, err)
	// error: contestation deadline expired
	n2 := time.Now().Add(-10 * 24 * time.Hour)
	u2 := baseUser()
	u2.DeactivatedAt = &n2
	d := time.Now().Add(-2 * 24 * time.Hour)
	u2.ContestationDeadline = &d
	svc = NewLgpdService(&mockUserRepo{user: u2}, &mockUserContestationRepo{}, &mockAuditLogService{}, &mockNotifier{})
	err = svc.RequestContestation("publicid", "motivo")
	assert.Error(t, err)
	// error: reactivation already requested
	n3 := time.Now().Add(-2 * 24 * time.Hour)
	u3 := baseUser()
	u3.DeactivatedAt = &n3
	u3.ReactivationRequested = true
	svc = NewLgpdService(&mockUserRepo{user: u3}, &mockUserContestationRepo{}, &mockAuditLogService{}, &mockNotifier{})
	err = svc.RequestContestation("publicid", "motivo")
	assert.Error(t, err)
	// error: failed to create contestation
	n4 := time.Now().Add(-2 * 24 * time.Hour)
	u4 := baseUser()
	u4.DeactivatedAt = &n4
	svc = NewLgpdService(&mockUserRepo{user: u4}, &mockUserContestationRepo{createErr: errors.New("fail")}, &mockAuditLogService{}, &mockNotifier{})
	err = svc.RequestContestation("publicid", "motivo")
	assert.Error(t, err)
}

func TestProcessPendingDeletions(t *testing.T) {
	user := baseUser()
	now := time.Now().Add(-8 * 24 * time.Hour)
	user.DeletionRequestedAt = &now
	svc := NewLgpdService(&mockUserRepo{user: user, usersPending: []*model.User{user}}, &mockUserContestationRepo{}, &mockAuditLogService{}, &mockNotifier{})
	err := svc.ProcessPendingDeletions()
	assert.NoError(t, err)
	// error: failed to find users for pending deletions
	svc = NewLgpdService(&mockUserRepo{usersPendingErr: errors.New("fail")}, &mockUserContestationRepo{}, &mockAuditLogService{}, &mockNotifier{})
	err = svc.ProcessPendingDeletions()
	assert.Error(t, err)
}

func TestProcessExpiredAnonymizations(t *testing.T) {
	user := baseUser()
	now := time.Now().Add(-6 * 365 * 24 * time.Hour)
	user.DeletionConfirmedAt = &now
	svc := NewLgpdService(&mockUserRepo{user: user, usersExpired: []*model.User{user}}, &mockUserContestationRepo{}, &mockAuditLogService{}, &mockNotifier{})
	err := svc.ProcessExpiredAnonymizations()
	assert.NoError(t, err)
	// error: failed to find users for anonymization
	svc = NewLgpdService(&mockUserRepo{usersExpiredErr: errors.New("fail")}, &mockUserContestationRepo{}, &mockAuditLogService{}, &mockNotifier{})
	err = svc.ProcessExpiredAnonymizations()
	assert.Error(t, err)
}
