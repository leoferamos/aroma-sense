package review

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/dto"
	handlererrors "github.com/leoferamos/aroma-sense/internal/handler/errors"
	"github.com/leoferamos/aroma-sense/internal/model"
	"github.com/leoferamos/aroma-sense/internal/rate"
	logservice "github.com/leoferamos/aroma-sense/internal/service/log"
	productservice "github.com/leoferamos/aroma-sense/internal/service/product"
	reviewservice "github.com/leoferamos/aroma-sense/internal/service/review"
	userservice "github.com/leoferamos/aroma-sense/internal/service/user"
)

type ReviewHandler struct {
	service        reviewservice.ReviewService
	userService    userservice.UserProfileService
	productService productservice.ProductService
	auditService   logservice.AuditLogService
	reportService  reviewservice.ReviewReportService
	rateLimiter    rate.RateLimiter
}

func NewReviewHandler(s reviewservice.ReviewService, reportService reviewservice.ReviewReportService, userService userservice.UserProfileService, productService productservice.ProductService, auditService logservice.AuditLogService, limiter rate.RateLimiter) *ReviewHandler {
	return &ReviewHandler{service: s, reportService: reportService, userService: userService, productService: productService, auditService: auditService, rateLimiter: limiter}
}

// Create review handles the creation of a product review
func (h *ReviewHandler) CreateReview(c *gin.Context) {
	slug := c.Param("slug")

	productID, err := h.productService.GetProductIDBySlug(c.Request.Context(), slug)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "product_not_found"})
		return
	}

	rawUserID, exists := c.Get("userID")
	if !exists || rawUserID == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthenticated"})
		return
	}
	publicID := rawUserID.(string)
	userModel, err := h.userService.GetByPublicID(publicID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthenticated"})
		return
	}

	var req dto.ReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request"})
		return
	}

	if userModel.DisplayName == nil || strings.TrimSpace(getPtrVal(userModel.DisplayName)) == "" {
		c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "profile_incomplete"})
		return
	}

	review, err := h.service.CreateReview(c.Request.Context(), userModel, productID, req.Rating, req.Comment)
	if err != nil {
		if status, code, ok := handlererrors.MapServiceError(err); ok {
			c.JSON(status, dto.ErrorResponse{Error: code})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal_error"})
		return
	}

	resp := dto.ReviewResponse{
		ID:            review.ID,
		Rating:        review.Rating,
		Comment:       review.Comment,
		AuthorID:      userModel.PublicID,
		AuthorDisplay: getPtrVal(userModel.DisplayName),
		CreatedAt:     review.CreatedAt,
	}
	c.JSON(http.StatusCreated, resp)
}

// ListReviews handles the product reviews listing
func (h *ReviewHandler) ListReviews(c *gin.Context) {
	slug := c.Param("slug")

	productID, err := h.productService.GetProductIDBySlug(c.Request.Context(), slug)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "product_not_found"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	reviews, total, err := h.service.ListReviews(c.Request.Context(), productID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal_error"})
		return
	}

	items := make([]dto.ReviewResponse, 0, len(reviews))
	for _, r := range reviews {
		display := ""
		if r.User != nil && r.User.DisplayName != nil {
			display = getPtrVal(r.User.DisplayName)
		}
		authorID := ""
		if r.User != nil {
			authorID = r.User.PublicID
		}
		items = append(items, dto.ReviewResponse{
			ID:            r.ID,
			Rating:        r.Rating,
			Comment:       r.Comment,
			AuthorID:      authorID,
			AuthorDisplay: display,
			CreatedAt:     r.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, dto.ReviewListResponse{Items: items, Total: total, Page: page, Limit: limit})
}

// GetSummary handles the product review summary
func (h *ReviewHandler) GetSummary(c *gin.Context) {
	slug := c.Param("slug")

	productID, err := h.productService.GetProductIDBySlug(c.Request.Context(), slug)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "product_not_found"})
		return
	}

	avg, count, err := h.service.GetAverage(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal_error"})
		return
	}

	dist := map[int]int{}
	reviews, _, err := h.service.ListReviews(c.Request.Context(), productID, 1, 1000)
	if err == nil {
		for _, r := range reviews {
			dist[r.Rating]++
		}
	}
	c.JSON(http.StatusOK, dto.ReviewSummary{Average: avg, Count: count, Distribution: dist})
}

// DeleteReview handles the deletion of a user's own review
func (h *ReviewHandler) DeleteReview(c *gin.Context) {
	reviewID := c.Param("reviewID")

	rawUserID, exists := c.Get("userID")
	if !exists || rawUserID == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthenticated"})
		return
	}
	publicID := rawUserID.(string)
	userModel, err := h.userService.GetByPublicID(publicID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthenticated"})
		return
	}

	err = h.service.DeleteOwnReview(c.Request.Context(), reviewID, publicID)
	if err != nil {
		if status, code, ok := handlererrors.MapServiceError(err); ok {
			c.JSON(status, dto.ErrorResponse{Error: code})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal_error"})
		return
	}

	if h.auditService != nil {
		h.auditService.LogDeletionAction(&userModel.ID, userModel.ID, model.AuditActionReviewDeleted, map[string]interface{}{
			"review_id":      reviewID,
			"user_public_id": publicID,
		})
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "review deleted successfully"})
}

// ReportReview handles reporting an abusive or inappropriate review
func (h *ReviewHandler) ReportReview(c *gin.Context) {
	reviewID := c.Param("reviewID")

	if h.rateLimiter != nil {
		bucket := "review_report:" + c.ClientIP() + ":" + c.GetString("userID")
		allowed, _, _, err := h.rateLimiter.Allow(c.Request.Context(), bucket, 5, time.Hour)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal_error"})
			return
		}
		if !allowed {
			c.JSON(http.StatusTooManyRequests, dto.ErrorResponse{Error: "rate_limited"})
			return
		}
	}

	rawUserID, exists := c.Get("userID")
	if !exists || rawUserID == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthenticated"})
		return
	}
	reporterID := rawUserID.(string)

	var req dto.ReviewReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_request"})
		return
	}

	if err := h.reportService.Report(c.Request.Context(), reviewID, reporterID, req.Category, req.Reason); err != nil {
		if status, code, ok := handlererrors.MapServiceError(err); ok {
			c.JSON(status, dto.ErrorResponse{Error: code})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal_error"})
		return
	}

	c.JSON(http.StatusCreated, dto.MessageResponse{Message: "review reported successfully"})
}

func getPtrVal(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}
