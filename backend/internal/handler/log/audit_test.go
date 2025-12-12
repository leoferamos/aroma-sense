package handler_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/dto"
	handler "github.com/leoferamos/aroma-sense/internal/handler/log"
	"github.com/leoferamos/aroma-sense/internal/model"
	logservice "github.com/leoferamos/aroma-sense/internal/service/log"
	"github.com/stretchr/testify/assert"
)

type mockAuditLogService struct {
	listAuditLogsResult                     []*model.AuditLog
	listAuditLogsTotal                      int64
	listAuditLogsErr                        error
	getAuditLogByIDResult                   *model.AuditLog
	getAuditLogByIDErr                      error
	getUserAuditLogsResult                  []*model.AuditLog
	getUserAuditLogsTotal                   int64
	getUserAuditLogsErr                     error
	getAuditSummaryResult                   *model.AuditLogSummary
	getAuditSummaryErr                      error
	cleanupOldLogsErr                       error
	convertAuditLogToResponseResult         dto.AuditLogResponse
	convertAuditLogToResponseDetailedResult dto.AuditLogResponse
	convertAuditLogsToResponseResult        []dto.AuditLogResponse
	convertAuditLogSummaryToResponseResult  dto.AuditLogSummaryResponse
}

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
	return m.listAuditLogsResult, m.listAuditLogsTotal, m.listAuditLogsErr
}

func (m *mockAuditLogService) GetAuditLogByID(id uint) (*model.AuditLog, error) {
	return m.getAuditLogByIDResult, m.getAuditLogByIDErr
}

func (m *mockAuditLogService) GetUserAuditLogs(userID uint, limit, offset int) ([]*model.AuditLog, int64, error) {
	return m.getUserAuditLogsResult, m.getUserAuditLogsTotal, m.getUserAuditLogsErr
}

func (m *mockAuditLogService) GetResourceAuditLogs(resource, resourceID string) ([]*model.AuditLog, error) {
	return nil, nil
}

func (m *mockAuditLogService) GetAuditSummary(startDate, endDate *time.Time) (*model.AuditLogSummary, error) {
	return m.getAuditSummaryResult, m.getAuditSummaryErr
}

func (m *mockAuditLogService) CleanupOldLogs(retentionDays int) error {
	return m.cleanupOldLogsErr
}

func (m *mockAuditLogService) ConvertAuditLogToResponse(auditLog *model.AuditLog) dto.AuditLogResponse {
	return m.convertAuditLogToResponseResult
}

func (m *mockAuditLogService) ConvertAuditLogToResponseDetailed(auditLog *model.AuditLog) dto.AuditLogResponse {
	return m.convertAuditLogToResponseDetailedResult
}

func (m *mockAuditLogService) ConvertAuditLogsToResponse(auditLogs []*model.AuditLog) []dto.AuditLogResponse {
	return m.convertAuditLogsToResponseResult
}

func (m *mockAuditLogService) ConvertAuditLogSummaryToResponse(summary *model.AuditLogSummary) dto.AuditLogSummaryResponse {
	return m.convertAuditLogSummaryToResponseResult
}

func setupAuditLogRouter(svc logservice.AuditLogService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Add middleware to simulate authentication
	r.Use(func(c *gin.Context) {
		c.Set("userID", "1")
		c.Next()
	})

	handler := handler.NewAuditLogHandler(svc)
	r.GET("/admin/audit-logs", handler.ListAuditLogs)
	r.GET("/admin/audit-logs/:id", handler.GetAuditLog)
	r.GET("/admin/audit-logs/:id/detailed", handler.GetAuditLogDetailed)
	r.GET("/admin/users/:id/audit-logs", handler.GetUserAuditLogs)
	r.GET("/admin/audit-logs/summary", handler.GetAuditSummary)
	r.POST("/admin/audit-logs/cleanup", handler.CleanupOldLogs)
	return r
}

func TestAuditLogHandler_ListAuditLogs(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		userID := uint(1)
		auditLogs := []*model.AuditLog{
			{ID: 1, UserID: &userID, Action: model.AuditActionUserLogin, CreatedAt: time.Now()},
		}
		svc := &mockAuditLogService{
			listAuditLogsResult: auditLogs,
			listAuditLogsTotal:  1,
			convertAuditLogsToResponseResult: []dto.AuditLogResponse{
				{ID: 1, UserID: &userID, Action: model.AuditActionUserLogin},
			},
		}
		r := setupAuditLogRouter(svc)

		req, _ := http.NewRequest("GET", "/admin/audit-logs", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response dto.AuditLogListResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.AuditLogs, 1)
		assert.Equal(t, int64(1), response.Total)
	})

	t.Run("service error", func(t *testing.T) {
		svc := &mockAuditLogService{
			listAuditLogsErr: errors.New("database error"),
		}
		r := setupAuditLogRouter(svc)

		req, _ := http.NewRequest("GET", "/admin/audit-logs", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "internal_error", response.Error)
	})
}

func TestAuditLogHandler_GetAuditLog(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		userID := uint(1)
		auditLog := &model.AuditLog{ID: 1, UserID: &userID, Action: model.AuditActionUserLogin}
		svc := &mockAuditLogService{
			getAuditLogByIDResult:           auditLog,
			convertAuditLogToResponseResult: dto.AuditLogResponse{ID: 1, UserID: &userID, Action: model.AuditActionUserLogin},
		}
		r := setupAuditLogRouter(svc)

		req, _ := http.NewRequest("GET", "/admin/audit-logs/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response dto.AuditLogResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, uint(1), response.ID)
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := &mockAuditLogService{}
		r := setupAuditLogRouter(svc)

		req, _ := http.NewRequest("GET", "/admin/audit-logs/invalid", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid_request", response.Error)
	})

	t.Run("not found", func(t *testing.T) {
		svc := &mockAuditLogService{
			getAuditLogByIDErr: errors.New("not found"),
		}
		r := setupAuditLogRouter(svc)

		req, _ := http.NewRequest("GET", "/admin/audit-logs/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "not_found", response.Error)
	})
}

func TestAuditLogHandler_GetAuditSummary(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		summary := &model.AuditLogSummary{TotalActions: 100, ActionsByType: map[string]int64{"user_login": 50}}
		svc := &mockAuditLogService{
			getAuditSummaryResult:                  summary,
			convertAuditLogSummaryToResponseResult: dto.AuditLogSummaryResponse{TotalActions: 100},
		}
		r := setupAuditLogRouter(svc)

		req, _ := http.NewRequest("GET", "/admin/audit-logs/summary", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response dto.AuditLogSummaryResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, int64(100), response.TotalActions)
	})

	t.Run("service error", func(t *testing.T) {
		svc := &mockAuditLogService{
			getAuditSummaryErr: errors.New("database error"),
		}
		r := setupAuditLogRouter(svc)

		req, _ := http.NewRequest("GET", "/admin/audit-logs/summary", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "internal_error", response.Error)
	})
}

func TestAuditLogHandler_CleanupOldLogs(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &mockAuditLogService{}
		r := setupAuditLogRouter(svc)

		req, _ := http.NewRequest("POST", "/admin/audit-logs/cleanup", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response dto.MessageResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Audit logs cleanup completed", response.Message)
	})

	t.Run("service error", func(t *testing.T) {
		svc := &mockAuditLogService{
			cleanupOldLogsErr: errors.New("cleanup failed"),
		}
		r := setupAuditLogRouter(svc)

		req, _ := http.NewRequest("POST", "/admin/audit-logs/cleanup", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "internal_error", response.Error)
	})
}
