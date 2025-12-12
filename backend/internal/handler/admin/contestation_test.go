package admin_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/dto"
	"github.com/leoferamos/aroma-sense/internal/handler/admin"
	"github.com/leoferamos/aroma-sense/internal/model"
	userservice "github.com/leoferamos/aroma-sense/internal/service/user"
	"github.com/stretchr/testify/assert"
)

type mockUserContestationService struct {
	listPendingContestations []model.UserContestation
	listPendingTotal         int64
	listPendingErr           error
	approveErr               error
	rejectErr                error
}

func (m *mockUserContestationService) Create(userID uint, reason string) error {
	return nil
}

func (m *mockUserContestationService) ListPending(limit, offset int) ([]model.UserContestation, int64, error) {
	return m.listPendingContestations, m.listPendingTotal, m.listPendingErr
}

func (m *mockUserContestationService) Approve(id uint, adminPublicID string, notes *string) error {
	return m.approveErr
}

func (m *mockUserContestationService) Reject(id uint, adminPublicID string, notes *string) error {
	return m.rejectErr
}

// --- Test helpers ---
func createTestContestation() model.UserContestation {
	now := time.Now()
	return model.UserContestation{
		ID:          1,
		UserID:      123,
		Reason:      "Test reason",
		Status:      "pending",
		RequestedAt: now,
	}
}

func setupAdminContestationRouter(svc userservice.UserContestationService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Add middleware to simulate authentication
	r.Use(func(c *gin.Context) {
		c.Set("userID", "admin123")
		c.Next()
	})

	handler := admin.NewAdminContestationHandler(svc)
	r.GET("/admin/contestations", handler.ListPendingContestions)
	r.POST("/admin/contestations/:id/approve", handler.ApproveContestation)
	r.POST("/admin/contestations/:id/reject", handler.RejectContestation)
	return r
}

// --- Tests ---
func TestAdminContestationHandler_ListPendingContestions(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		contestations := []model.UserContestation{createTestContestation()}
		svc := &mockUserContestationService{
			listPendingContestations: contestations,
			listPendingTotal:         1,
		}
		r := setupAdminContestationRouter(svc)

		req, _ := http.NewRequest("GET", "/admin/contestations?limit=10&offset=0", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		data, ok := response["data"].([]interface{})
		assert.True(t, ok)
		assert.Len(t, data, 1)

		total, ok := response["total"].(float64)
		assert.True(t, ok)
		assert.Equal(t, float64(1), total)
	})

	t.Run("success with default params", func(t *testing.T) {
		contestations := []model.UserContestation{createTestContestation()}
		svc := &mockUserContestationService{
			listPendingContestations: contestations,
			listPendingTotal:         1,
		}
		r := setupAdminContestationRouter(svc)

		req, _ := http.NewRequest("GET", "/admin/contestations", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("empty result", func(t *testing.T) {
		svc := &mockUserContestationService{
			listPendingContestations: []model.UserContestation{},
			listPendingTotal:         0,
		}
		r := setupAdminContestationRouter(svc)

		req, _ := http.NewRequest("GET", "/admin/contestations", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// When empty, data can be null or empty array
		data := response["data"]
		if data != nil {
			dataArray, ok := data.([]interface{})
			assert.True(t, ok)
			assert.Len(t, dataArray, 0)
		}

		total, ok := response["total"].(float64)
		assert.True(t, ok)
		assert.Equal(t, float64(0), total)
	})

	t.Run("service error", func(t *testing.T) {
		svc := &mockUserContestationService{
			listPendingErr: errors.New("db error"),
		}
		r := setupAdminContestationRouter(svc)

		req, _ := http.NewRequest("GET", "/admin/contestations", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "internal_error", response.Error)
	})
}

func TestAdminContestationHandler_ApproveContestation(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &mockUserContestationService{}
		r := setupAdminContestationRouter(svc)

		reqBody := `{"notes": "Approved for valid reason"}`
		req, _ := http.NewRequest("POST", "/admin/contestations/1/approve", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response dto.MessageResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "contestation approved", response.Message)
	})

	t.Run("success without notes", func(t *testing.T) {
		svc := &mockUserContestationService{}
		r := setupAdminContestationRouter(svc)

		reqBody := `{}`
		req, _ := http.NewRequest("POST", "/admin/contestations/1/approve", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := &mockUserContestationService{}
		r := setupAdminContestationRouter(svc)

		reqBody := `{}`
		req, _ := http.NewRequest("POST", "/admin/contestations/invalid/approve", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid_request", response.Error)
	})

	t.Run("invalid json", func(t *testing.T) {
		svc := &mockUserContestationService{}
		r := setupAdminContestationRouter(svc)

		reqBody := `{"notes": invalid}`
		req, _ := http.NewRequest("POST", "/admin/contestations/1/approve", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid_request", response.Error)
	})

	t.Run("service error", func(t *testing.T) {
		svc := &mockUserContestationService{
			approveErr: errors.New("db error"),
		}
		r := setupAdminContestationRouter(svc)

		reqBody := `{}`
		req, _ := http.NewRequest("POST", "/admin/contestations/1/approve", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "internal_error", response.Error)
	})
}

func TestAdminContestationHandler_RejectContestation(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &mockUserContestationService{}
		r := setupAdminContestationRouter(svc)

		reqBody := `{"notes": "Rejected for invalid reason"}`
		req, _ := http.NewRequest("POST", "/admin/contestations/1/reject", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response dto.MessageResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "contestation rejected", response.Message)
	})

	t.Run("success without notes", func(t *testing.T) {
		svc := &mockUserContestationService{}
		r := setupAdminContestationRouter(svc)

		reqBody := `{}`
		req, _ := http.NewRequest("POST", "/admin/contestations/1/reject", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := &mockUserContestationService{}
		r := setupAdminContestationRouter(svc)

		reqBody := `{}`
		req, _ := http.NewRequest("POST", "/admin/contestations/invalid/reject", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid_request", response.Error)
	})

	t.Run("invalid json", func(t *testing.T) {
		svc := &mockUserContestationService{}
		r := setupAdminContestationRouter(svc)

		reqBody := `{"notes": invalid}`
		req, _ := http.NewRequest("POST", "/admin/contestations/1/reject", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid_request", response.Error)
	})

	t.Run("service error", func(t *testing.T) {
		svc := &mockUserContestationService{
			rejectErr: errors.New("db error"),
		}
		r := setupAdminContestationRouter(svc)

		reqBody := `{}`
		req, _ := http.NewRequest("POST", "/admin/contestations/1/reject", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "internal_error", response.Error)
	})
}
