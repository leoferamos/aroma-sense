package admin_test

import (
	"context"
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
	reviewservice "github.com/leoferamos/aroma-sense/internal/service/review"
	"github.com/stretchr/testify/assert"
)

type mockReviewReportService struct {
	listReports []model.ReviewReport
	listTotal   int64
	listErr     error
	resolveErr  error
}

func (m *mockReviewReportService) Report(ctx context.Context, reviewID string, reporterID string, category string, reason string) error {
	return nil
}

func (m *mockReviewReportService) List(ctx context.Context, status string, limit, offset int) ([]model.ReviewReport, int64, error) {
	return m.listReports, m.listTotal, m.listErr
}

func (m *mockReviewReportService) Resolve(ctx context.Context, reportID string, action string, deactivateUser bool, adminPublicID string, suspensionUntil *time.Time, notes *string) error {
	return m.resolveErr
}

// --- Test helpers ---
func createTestReviewReport() model.ReviewReport {
	now := time.Now()
	return model.ReviewReport{
		ID:             "report-123",
		ReviewID:       "review-123",
		ReportedBy:     "user-123",
		ReasonCategory: "spam",
		ReasonText:     "This is spam content",
		Status:         "pending",
		CreatedAt:      now,
		Review: &model.Review{
			ID:      "review-123",
			Comment: "Great product!",
			Rating:  5,
			UserID:  "reviewer-123",
		},
		Reporter: &model.User{
			PublicID: "user-123",
			Email:    "reporter@example.com",
		},
	}
}

func setupAdminReviewReportRouter(svc reviewservice.ReviewReportService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Add middleware to simulate authentication
	r.Use(func(c *gin.Context) {
		c.Set("userID", "admin-123")
		c.Next()
	})

	handler := admin.NewAdminReviewReportHandler(svc)
	r.GET("/admin/review-reports", handler.ListReports)
	r.POST("/admin/review-reports/:id/resolve", handler.ResolveReport)
	return r
}

// --- Tests ---
func TestAdminReviewReportHandler_ListReports(t *testing.T) {
	t.Run("success with pending status", func(t *testing.T) {
		reports := []model.ReviewReport{createTestReviewReport()}
		svc := &mockReviewReportService{
			listReports: reports,
			listTotal:   1,
		}
		r := setupAdminReviewReportRouter(svc)

		req, _ := http.NewRequest("GET", "/admin/review-reports?status=pending&limit=10&offset=0", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response dto.ReviewReportAdminResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Items, 1)
		assert.Equal(t, int64(1), response.Total)
		assert.Equal(t, 10, response.Limit)
		assert.Equal(t, 0, response.Offset)
	})

	t.Run("success with default params", func(t *testing.T) {
		reports := []model.ReviewReport{createTestReviewReport()}
		svc := &mockReviewReportService{
			listReports: reports,
			listTotal:   1,
		}
		r := setupAdminReviewReportRouter(svc)

		req, _ := http.NewRequest("GET", "/admin/review-reports", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response dto.ReviewReportAdminResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 20, response.Limit) // default limit
		assert.Equal(t, 0, response.Offset) // default offset
	})

	t.Run("empty result", func(t *testing.T) {
		svc := &mockReviewReportService{
			listReports: []model.ReviewReport{},
			listTotal:   0,
		}
		r := setupAdminReviewReportRouter(svc)

		req, _ := http.NewRequest("GET", "/admin/review-reports", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response dto.ReviewReportAdminResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Items, 0)
		assert.Equal(t, int64(0), response.Total)
	})

	t.Run("service error", func(t *testing.T) {
		svc := &mockReviewReportService{
			listErr: errors.New("db error"),
		}
		r := setupAdminReviewReportRouter(svc)

		req, _ := http.NewRequest("GET", "/admin/review-reports", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "internal_error", response.Error)
	})
}

func TestAdminReviewReportHandler_ResolveReport(t *testing.T) {
	t.Run("success accept", func(t *testing.T) {
		svc := &mockReviewReportService{}
		r := setupAdminReviewReportRouter(svc)

		reqBody := `{
			"action": "accept",
			"deactivate_user": false,
			"notes": "Report accepted and review hidden"
		}`
		req, _ := http.NewRequest("POST", "/admin/review-reports/report-123/resolve", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response dto.MessageResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "report resolved", response.Message)
	})

	t.Run("success reject", func(t *testing.T) {
		svc := &mockReviewReportService{}
		r := setupAdminReviewReportRouter(svc)

		reqBody := `{
			"action": "reject",
			"deactivate_user": true,
			"suspension_until": "2025-12-31T23:59:59Z",
			"notes": "Report rejected, user suspended"
		}`
		req, _ := http.NewRequest("POST", "/admin/review-reports/report-123/resolve", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response dto.MessageResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "report resolved", response.Message)
	})

	t.Run("success minimal payload", func(t *testing.T) {
		svc := &mockReviewReportService{}
		r := setupAdminReviewReportRouter(svc)

		reqBody := `{"action": "accept"}`
		req, _ := http.NewRequest("POST", "/admin/review-reports/report-123/resolve", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("invalid action", func(t *testing.T) {
		svc := &mockReviewReportService{}
		r := setupAdminReviewReportRouter(svc)

		reqBody := `{"action": "invalid"}`
		req, _ := http.NewRequest("POST", "/admin/review-reports/report-123/resolve", strings.NewReader(reqBody))
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
		svc := &mockReviewReportService{}
		r := setupAdminReviewReportRouter(svc)

		reqBody := `{"action": "accept", "invalid": json}`
		req, _ := http.NewRequest("POST", "/admin/review-reports/report-123/resolve", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response dto.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid_request", response.Error)
	})

	t.Run("invalid suspension date", func(t *testing.T) {
		svc := &mockReviewReportService{}
		r := setupAdminReviewReportRouter(svc)

		reqBody := `{
			"action": "reject",
			"deactivate_user": true,
			"suspension_until": "invalid-date"
		}`
		req, _ := http.NewRequest("POST", "/admin/review-reports/report-123/resolve", strings.NewReader(reqBody))
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
		svc := &mockReviewReportService{
			resolveErr: errors.New("db error"),
		}
		r := setupAdminReviewReportRouter(svc)

		reqBody := `{"action": "accept"}`
		req, _ := http.NewRequest("POST", "/admin/review-reports/report-123/resolve", strings.NewReader(reqBody))
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
