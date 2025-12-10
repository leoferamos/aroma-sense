package review_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/apperror"
	"github.com/leoferamos/aroma-sense/internal/dto"
	handler "github.com/leoferamos/aroma-sense/internal/handler/review"
	"github.com/leoferamos/aroma-sense/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubReviewService struct {
	createFn func(ctx context.Context, user *model.User, productID uint, rating int, comment string) (*model.Review, error)
	listFn   func(ctx context.Context, productID uint, page, perPage int) ([]model.Review, int, error)
	avgFn    func(ctx context.Context, productID uint) (float64, int, error)
}

func (s stubReviewService) CanUserReview(ctx context.Context, user *model.User, productID uint) (bool, string, error) {
	return false, "", nil
}

func (s stubReviewService) CanUserReviewBySlug(ctx context.Context, user *model.User, slug string) (bool, string, error) {
	return false, "", nil
}

func (s stubReviewService) CreateReview(ctx context.Context, user *model.User, productID uint, rating int, comment string) (*model.Review, error) {
	return s.createFn(ctx, user, productID, rating, comment)
}

func (s stubReviewService) ListReviews(ctx context.Context, productID uint, page, perPage int) ([]model.Review, int, error) {
	if s.listFn == nil {
		return nil, 0, nil
	}
	return s.listFn(ctx, productID, page, perPage)
}

func (s stubReviewService) GetAverage(ctx context.Context, productID uint) (float64, int, error) {
	if s.avgFn == nil {
		return 0, 0, nil
	}
	return s.avgFn(ctx, productID)
}

func (s stubReviewService) DeleteOwnReview(ctx context.Context, reviewID string, userID string) error {
	return nil
}

type stubProductService struct {
	id  uint
	err error
}

func (s stubProductService) CreateProduct(ctx context.Context, input dto.ProductFormDTO, file dto.FileUpload) error {
	return nil
}
func (s stubProductService) GetProductByID(ctx context.Context, id uint) (dto.ProductResponse, error) {
	return dto.ProductResponse{}, nil
}
func (s stubProductService) GetProductBySlug(ctx context.Context, slug string) (dto.ProductResponse, error) {
	return dto.ProductResponse{}, nil
}
func (s stubProductService) GetProductIDBySlug(ctx context.Context, slug string) (uint, error) {
	return s.id, s.err
}
func (s stubProductService) GetLatestProducts(ctx context.Context, page int, limit int) ([]dto.ProductResponse, int, error) {
	return nil, 0, nil
}
func (s stubProductService) SearchProducts(ctx context.Context, query string, page int, limit int, sort string) ([]dto.ProductResponse, int, error) {
	return nil, 0, nil
}
func (s stubProductService) AdminListProducts(ctx context.Context, page int, limit int) ([]dto.ProductResponse, int, error) {
	return nil, 0, nil
}
func (s stubProductService) UpdateProduct(ctx context.Context, id uint, input dto.UpdateProductRequest) error {
	return nil
}
func (s stubProductService) DeleteProduct(ctx context.Context, id uint) error { return nil }

type stubUserProfileService struct {
	user *model.User
	err  error
}

func (s stubUserProfileService) GetByPublicID(publicID string) (*model.User, error) {
	return s.user, s.err
}
func (s stubUserProfileService) UpdateDisplayName(publicID string, displayName string) (*model.User, error) {
	return s.user, s.err
}
func (s stubUserProfileService) SetPasswordHash(publicID string, hashedPassword string) error {
	return nil
}
func (s stubUserProfileService) ChangePassword(publicID string, currentPassword string, newPassword string) error {
	return nil
}

type stubAuditLogService struct{}

func (stubAuditLogService) LogUserAction(actorID *uint, userID *uint, action model.AuditAction, details map[string]interface{}) error {
	return nil
}
func (stubAuditLogService) LogUserUpdate(actorID uint, userID uint, oldUser, newUser *model.User) error {
	return nil
}
func (stubAuditLogService) LogAdminAction(adminID uint, userID uint, action model.AuditAction, details map[string]interface{}) error {
	return nil
}
func (stubAuditLogService) LogSystemAction(action model.AuditAction, resource, resourceID string, details map[string]interface{}) error {
	return nil
}
func (stubAuditLogService) LogDataAccess(userID uint, resource string, resourceID string) error {
	return nil
}
func (stubAuditLogService) LogDeletionAction(actorID *uint, userID uint, action model.AuditAction, details map[string]interface{}) error {
	return nil
}
func (stubAuditLogService) ListAuditLogs(filter *model.AuditLogFilter) ([]*model.AuditLog, int64, error) {
	return nil, 0, nil
}
func (stubAuditLogService) GetAuditLogByID(id uint) (*model.AuditLog, error) { return nil, nil }
func (stubAuditLogService) GetUserAuditLogs(userID uint, limit, offset int) ([]*model.AuditLog, int64, error) {
	return nil, 0, nil
}
func (stubAuditLogService) GetResourceAuditLogs(resource, resourceID string) ([]*model.AuditLog, error) {
	return nil, nil
}
func (stubAuditLogService) GetAuditSummary(startDate, endDate *time.Time) (*model.AuditLogSummary, error) {
	return nil, nil
}
func (stubAuditLogService) CleanupOldLogs(retentionDays int) error { return nil }
func (stubAuditLogService) ConvertAuditLogToResponse(auditLog *model.AuditLog) dto.AuditLogResponse {
	return dto.AuditLogResponse{}
}
func (stubAuditLogService) ConvertAuditLogToResponseDetailed(auditLog *model.AuditLog) dto.AuditLogResponse {
	return dto.AuditLogResponse{}
}
func (stubAuditLogService) ConvertAuditLogsToResponse(auditLogs []*model.AuditLog) []dto.AuditLogResponse {
	return nil
}
func (stubAuditLogService) ConvertAuditLogSummaryToResponse(summary *model.AuditLogSummary) dto.AuditLogSummaryResponse {
	return dto.AuditLogSummaryResponse{}
}

type stubReviewReportService struct{ err error }

func (s stubReviewReportService) Report(ctx context.Context, reviewID string, reporterID string, category string, reason string) error {
	return s.err
}
func (s stubReviewReportService) List(ctx context.Context, status string, limit, offset int) ([]model.ReviewReport, int64, error) {
	return nil, 0, nil
}
func (s stubReviewReportService) Resolve(ctx context.Context, reportID string, action string, deactivateUser bool, adminPublicID string, suspensionUntil *time.Time, notes *string) error {
	return nil
}

func setupReviewRouter(handler *handler.ReviewHandler) *gin.Engine {
	r := gin.New()
	r.POST("/products/:slug/reviews", func(c *gin.Context) {
		c.Set("userID", "user-1")
		handler.CreateReview(c)
	})
	r.GET("/products/:slug/reviews", handler.ListReviews)
	r.GET("/products/:slug/reviews/summary", handler.GetSummary)
	r.DELETE("/reviews/:id", func(c *gin.Context) {
		c.Set("userID", "user-1")
		handler.DeleteReview(c)
	})
	r.POST("/reviews/:id/report", func(c *gin.Context) {
		c.Set("userID", "user-1")
		handler.ReportReview(c)
	})
	return r
}

func TestReviewHandler_CreateReview(t *testing.T) {
	t.Parallel()

	user := &model.User{PublicID: "user-1", DisplayName: ptr("Alice")}
	productSvc := stubProductService{id: 42}
	reviewSvc := stubReviewService{createFn: func(ctx context.Context, u *model.User, productID uint, rating int, comment string) (*model.Review, error) {
		return &model.Review{ID: "rev-1", Rating: rating, Comment: comment, CreatedAt: time.Unix(1, 0)}, nil
	}}
	userSvc := stubUserProfileService{user: user}
	h := handler.NewReviewHandler(reviewSvc, stubReviewReportService{}, userSvc, productSvc, stubAuditLogService{}, nil)
	r := setupReviewRouter(h)

	body, err := json.Marshal(dto.ReviewRequest{Rating: 5, Comment: "nice"})
	require.NoError(t, err)
	req := httptest.NewRequest(http.MethodPost, "/products/slug-1/reviews", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	require.Equal(t, http.StatusCreated, res.Code)
	var resp dto.ReviewResponse
	require.NoError(t, json.Unmarshal(res.Body.Bytes(), &resp))
	assert.Equal(t, "rev-1", resp.ID)
	assert.Equal(t, "Alice", resp.AuthorDisplay)

	t.Run("product not found returns 404", func(t *testing.T) {
		badProduct := stubProductService{err: assert.AnError}
		h := handler.NewReviewHandler(reviewSvc, stubReviewReportService{}, userSvc, badProduct, stubAuditLogService{}, nil)
		r := setupReviewRouter(h)

		req := httptest.NewRequest(http.MethodPost, "/products/slug-1/reviews", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusNotFound, res.Code)
	})

	t.Run("unauthenticated returns 401", func(t *testing.T) {
		rNoAuth := gin.New()
		rNoAuth.POST("/products/:slug/reviews", h.CreateReview)

		req := httptest.NewRequest(http.MethodPost, "/products/slug-1/reviews", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		rNoAuth.ServeHTTP(res, req)
		assert.Equal(t, http.StatusUnauthorized, res.Code)
	})

	t.Run("mapped service error", func(t *testing.T) {
		reviewSvcErr := stubReviewService{createFn: func(ctx context.Context, user *model.User, productID uint, rating int, comment string) (*model.Review, error) {
			return nil, apperror.NewCodeMessage("already_reviewed", "")
		}}
		hErr := handler.NewReviewHandler(reviewSvcErr, stubReviewReportService{}, userSvc, productSvc, stubAuditLogService{}, nil)
		rErr := setupReviewRouter(hErr)

		req := httptest.NewRequest(http.MethodPost, "/products/slug-1/reviews", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		rErr.ServeHTTP(res, req)
		assert.Equal(t, http.StatusConflict, res.Code)
	})
}

func TestReviewHandler_ListReviews(t *testing.T) {
	t.Parallel()

	reviews := []model.Review{
		{ID: "r1", Rating: 4, Comment: "ok", CreatedAt: time.Unix(1, 0), User: &model.User{PublicID: "u1", DisplayName: ptr("A")}},
		{ID: "r2", Rating: 5, Comment: "great", CreatedAt: time.Unix(2, 0)},
	}

	reviewSvc := stubReviewService{
		listFn: func(ctx context.Context, productID uint, page, perPage int) ([]model.Review, int, error) {
			return reviews, len(reviews), nil
		},
		avgFn: func(ctx context.Context, productID uint) (float64, int, error) { return 4.5, 2, nil },
	}
	productSvc := stubProductService{id: 99}
	userSvc := stubUserProfileService{user: &model.User{PublicID: "u1", DisplayName: ptr("A")}}
	h := handler.NewReviewHandler(reviewSvc, stubReviewReportService{}, userSvc, productSvc, stubAuditLogService{}, nil)
	r := setupReviewRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/products/slug-1/reviews?page=1&limit=10", nil)
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)
	require.Equal(t, http.StatusOK, res.Code)

	var listResp dto.ReviewListResponse
	require.NoError(t, json.Unmarshal(res.Body.Bytes(), &listResp))
	assert.Len(t, listResp.Items, 2)

	reqSummary := httptest.NewRequest(http.MethodGet, "/products/slug-1/reviews/summary", nil)
	resSummary := httptest.NewRecorder()
	r.ServeHTTP(resSummary, reqSummary)
	require.Equal(t, http.StatusOK, resSummary.Code)

	var summary dto.ReviewSummary
	require.NoError(t, json.Unmarshal(resSummary.Body.Bytes(), &summary))
	assert.Equal(t, 4.5, summary.Average)
	assert.Equal(t, 2, summary.Count)
}

func ptr(s string) *string { return &s }
