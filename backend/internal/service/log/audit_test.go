package service

import (
	"errors"
	"testing"
	"time"

	"github.com/leoferamos/aroma-sense/internal/dto"
	"github.com/leoferamos/aroma-sense/internal/model"
	"github.com/stretchr/testify/assert"
)

type mockAuditLogRepo struct {
	createErr         error
	listLogs          []*model.AuditLog
	listCount         int64
	listErr           error
	getByIDLog        *model.AuditLog
	getByIDErr        error
	getByUserIDLogs   []*model.AuditLog
	getByUserIDCount  int64
	getByUserIDErr    error
	getByResourceLogs []*model.AuditLog
	getByResourceErr  error
	getSummaryResult  *model.AuditLogSummary
	getSummaryErr     error
	deleteOldErr      error
}

func (m *mockAuditLogRepo) Create(auditLog *model.AuditLog) error {
	return m.createErr
}

func (m *mockAuditLogRepo) List(filter *model.AuditLogFilter) ([]*model.AuditLog, int64, error) {
	return m.listLogs, m.listCount, m.listErr
}

func (m *mockAuditLogRepo) GetByID(id uint) (*model.AuditLog, error) {
	return m.getByIDLog, m.getByIDErr
}

func (m *mockAuditLogRepo) GetByPublicID(publicID string) (*model.AuditLog, error) {
	return nil, nil
}

func (m *mockAuditLogRepo) GetByUserID(userID uint, limit int, offset int) ([]*model.AuditLog, int64, error) {
	return m.getByUserIDLogs, m.getByUserIDCount, m.getByUserIDErr
}

func (m *mockAuditLogRepo) GetByResource(resource string, resourceID string) ([]*model.AuditLog, error) {
	return m.getByResourceLogs, m.getByResourceErr
}

func (m *mockAuditLogRepo) GetSummary(startDate, endDate *time.Time) (*model.AuditLogSummary, error) {
	return m.getSummaryResult, m.getSummaryErr
}

func (m *mockAuditLogRepo) DeleteOldLogs(olderThan time.Time) error {
	return m.deleteOldErr
}

// --- Test helpers ---
func createTestAuditLog() *model.AuditLog {
	return &model.AuditLog{
		ID:         1,
		UserID:     &[]uint{1}[0],
		ActorID:    &[]uint{2}[0],
		ActorType:  "user",
		Action:     model.AuditActionUserLogin,
		Resource:   "user",
		ResourceID: &[]string{"123"}[0],
		Details:    `{"key": "value"}`,
		OldValues:  "",
		NewValues:  "",
		Timestamp:  time.Now(),
		Compliance: "LGPD",
		Severity:   "info",
	}
}

func createTestUser() *model.User {
	return &model.User{
		ID:          1,
		PublicID:    "user123",
		Email:       "user@example.com",
		Role:        "user",
		DisplayName: &[]string{"Test User"}[0],
		CreatedAt:   time.Now(),
	}
}

// --- Tests ---
func TestLogUserAction(t *testing.T) {
	svc := NewAuditLogService(&mockAuditLogRepo{})

	err := svc.LogUserAction(&[]uint{1}[0], &[]uint{2}[0], model.AuditActionUserLogin, map[string]interface{}{"key": "value"})
	assert.NoError(t, err)

	// error case
	svc = NewAuditLogService(&mockAuditLogRepo{createErr: errors.New("db error")})
	err = svc.LogUserAction(&[]uint{1}[0], &[]uint{2}[0], model.AuditActionUserLogin, map[string]interface{}{"key": "value"})
	assert.Error(t, err)
}

func TestLogUserUpdate(t *testing.T) {
	oldUser := createTestUser()
	newUser := createTestUser()
	newUser.Email = "new@example.com"

	svc := NewAuditLogService(&mockAuditLogRepo{})

	err := svc.LogUserUpdate(1, 2, oldUser, newUser)
	assert.NoError(t, err)

	// error case
	svc = NewAuditLogService(&mockAuditLogRepo{createErr: errors.New("db error")})
	err = svc.LogUserUpdate(1, 2, oldUser, newUser)
	assert.Error(t, err)
}

func TestLogAdminAction(t *testing.T) {
	svc := NewAuditLogService(&mockAuditLogRepo{})

	err := svc.LogAdminAction(1, 2, model.AuditActionUserDeactivated, map[string]interface{}{"reason": "violation"})
	assert.NoError(t, err)

	// error case
	svc = NewAuditLogService(&mockAuditLogRepo{createErr: errors.New("db error")})
	err = svc.LogAdminAction(1, 2, model.AuditActionUserDeactivated, map[string]interface{}{"reason": "violation"})
	assert.Error(t, err)
}

func TestLogSystemAction(t *testing.T) {
	svc := NewAuditLogService(&mockAuditLogRepo{})

	err := svc.LogSystemAction(model.AuditActionDataAnonymized, "system", "backup", map[string]interface{}{"status": "completed"})
	assert.NoError(t, err)

	// error case
	svc = NewAuditLogService(&mockAuditLogRepo{createErr: errors.New("db error")})
	err = svc.LogSystemAction(model.AuditActionDataAnonymized, "system", "backup", map[string]interface{}{"status": "completed"})
	assert.Error(t, err)
}

func TestLogDataAccess(t *testing.T) {
	svc := NewAuditLogService(&mockAuditLogRepo{})

	err := svc.LogDataAccess(1, "user", "123")
	assert.NoError(t, err)

	// error case
	svc = NewAuditLogService(&mockAuditLogRepo{createErr: errors.New("db error")})
	err = svc.LogDataAccess(1, "user", "123")
	assert.Error(t, err)
}

func TestLogDeletionAction(t *testing.T) {
	svc := NewAuditLogService(&mockAuditLogRepo{})

	err := svc.LogDeletionAction(&[]uint{1}[0], 2, model.AuditActionDeletionConfirmed, map[string]interface{}{"method": "user_request"})
	assert.NoError(t, err)

	// error case
	svc = NewAuditLogService(&mockAuditLogRepo{createErr: errors.New("db error")})
	err = svc.LogDeletionAction(&[]uint{1}[0], 2, model.AuditActionDeletionConfirmed, map[string]interface{}{"method": "user_request"})
	assert.Error(t, err)
}

func TestListAuditLogs(t *testing.T) {
	logs := []*model.AuditLog{createTestAuditLog()}
	svc := NewAuditLogService(&mockAuditLogRepo{listLogs: logs, listCount: 1})

	result, count, err := svc.ListAuditLogs(&model.AuditLogFilter{})
	assert.NoError(t, err)
	assert.Equal(t, logs, result)
	assert.Equal(t, int64(1), count)

	// error case
	svc = NewAuditLogService(&mockAuditLogRepo{listErr: errors.New("db error")})
	result, count, err = svc.ListAuditLogs(&model.AuditLogFilter{})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, int64(0), count)
}

func TestGetAuditLogByID(t *testing.T) {
	log := createTestAuditLog()
	svc := NewAuditLogService(&mockAuditLogRepo{getByIDLog: log})

	result, err := svc.GetAuditLogByID(1)
	assert.NoError(t, err)
	assert.Equal(t, log, result)

	// error case
	svc = NewAuditLogService(&mockAuditLogRepo{getByIDErr: errors.New("not found")})
	result, err = svc.GetAuditLogByID(1)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetUserAuditLogs(t *testing.T) {
	logs := []*model.AuditLog{createTestAuditLog()}
	svc := NewAuditLogService(&mockAuditLogRepo{getByUserIDLogs: logs, getByUserIDCount: 1})

	result, count, err := svc.GetUserAuditLogs(1, 10, 0)
	assert.NoError(t, err)
	assert.Equal(t, logs, result)
	assert.Equal(t, int64(1), count)

	// error case
	svc = NewAuditLogService(&mockAuditLogRepo{getByUserIDErr: errors.New("db error")})
	result, count, err = svc.GetUserAuditLogs(1, 10, 0)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, int64(0), count)
}

func TestGetResourceAuditLogs(t *testing.T) {
	logs := []*model.AuditLog{createTestAuditLog()}
	svc := NewAuditLogService(&mockAuditLogRepo{getByResourceLogs: logs})

	result, err := svc.GetResourceAuditLogs("user", "123")
	assert.NoError(t, err)
	assert.Equal(t, logs, result)

	// error case
	svc = NewAuditLogService(&mockAuditLogRepo{getByResourceErr: errors.New("db error")})
	result, err = svc.GetResourceAuditLogs("user", "123")
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetAuditSummary(t *testing.T) {
	summary := &model.AuditLogSummary{TotalActions: 100}
	svc := NewAuditLogService(&mockAuditLogRepo{getSummaryResult: summary})

	start := time.Now().Add(-24 * time.Hour)
	end := time.Now()
	result, err := svc.GetAuditSummary(&start, &end)
	assert.NoError(t, err)
	assert.Equal(t, summary, result)

	// error case
	svc = NewAuditLogService(&mockAuditLogRepo{getSummaryErr: errors.New("db error")})
	result, err = svc.GetAuditSummary(&start, &end)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestCleanupOldLogs(t *testing.T) {
	svc := NewAuditLogService(&mockAuditLogRepo{})

	err := svc.CleanupOldLogs(30)
	assert.NoError(t, err)

	// error case
	svc = NewAuditLogService(&mockAuditLogRepo{deleteOldErr: errors.New("db error")})
	err = svc.CleanupOldLogs(30)
	assert.Error(t, err)
}

func TestConvertAuditLogToResponse(t *testing.T) {
	log := createTestAuditLog()
	svc := NewAuditLogService(&mockAuditLogRepo{})

	response := svc.ConvertAuditLogToResponse(log)
	assert.IsType(t, dto.AuditLogResponse{}, response)
}

func TestConvertAuditLogToResponseDetailed(t *testing.T) {
	log := createTestAuditLog()
	svc := NewAuditLogService(&mockAuditLogRepo{})

	response := svc.ConvertAuditLogToResponseDetailed(log)
	assert.IsType(t, dto.AuditLogResponse{}, response)
}

func TestConvertAuditLogsToResponse(t *testing.T) {
	logs := []*model.AuditLog{createTestAuditLog()}
	svc := NewAuditLogService(&mockAuditLogRepo{})

	responses := svc.ConvertAuditLogsToResponse(logs)
	assert.Len(t, responses, 1)
	assert.IsType(t, []dto.AuditLogResponse{}, responses)
}

func TestConvertAuditLogSummaryToResponse(t *testing.T) {
	summary := &model.AuditLogSummary{TotalActions: 100}
	svc := NewAuditLogService(&mockAuditLogRepo{})

	response := svc.ConvertAuditLogSummaryToResponse(summary)
	assert.IsType(t, dto.AuditLogSummaryResponse{}, response)
}
