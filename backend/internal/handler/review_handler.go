package handler

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
	userservice "github.com/leoferamos/aroma-sense/internal/service"
	logservice "github.com/leoferamos/aroma-sense/internal/service/log"
	productservice "github.com/leoferamos/aroma-sense/internal/service/product"
	reviewservice "github.com/leoferamos/aroma-sense/internal/service/review"
)

type ReviewHandler struct {
	reviewService  reviewservice.ReviewService
	userService    userservice.UserProfileService
	productService productservice.ProductService
	auditService   logservice.AuditLogService
	reportService  reviewservice.ReviewReportService
	rateLimiter    rate.RateLimiter
}

func NewReviewHandler(s reviewservice.ReviewService, reportService reviewservice.ReviewReportService, userService userservice.UserProfileService, productService productservice.ProductService, auditService logservice.AuditLogService, limiter rate.RateLimiter) *ReviewHandler {
	return &ReviewHandler{reviewService: s, reportService: reportService, userService: userService, productService: productService, auditService: auditService, rateLimiter: limiter}
}

// Create review handles the creation of a product review
//
// @Summary      Create product review
// @Description  Creates a review for a delivered product. Requires authentication, a display_name set, and at least one delivered order containing the product. One review per user/product.
// @Tags         reviews
// @Accept       json
// @Produce      json
// @Param        slug    path     string            true  "Product slug"
// @Param        review  body     dto.ReviewRequest true  "Review payload"
// @Success      201  {object}  dto.ReviewResponse       "Review created"
// @Failure      400  {object}  dto.ErrorResponse        "Error code: invalid_request"
// @Failure      401  {object}  dto.ErrorResponse        "Error code: unauthenticated"
// @Failure      403  {object}  dto.ErrorResponse        "Error code: profile_incomplete or not_delivered"
// @Failure      404  {object}  dto.ErrorResponse        "Error code: product_not_found"
// @Failure      409  {object}  dto.ErrorResponse        "Error code: already_reviewed"
// @Failure      500  {object}  dto.ErrorResponse        "Error code: internal_error"
// @Router       /products/{slug}/reviews [post]
// @Security     BearerAuth
func (h *ReviewHandler) CreateReview(c *gin.Context) {
	slug := c.Param("slug")

	// Get product ID by slug
	productID, err := h.productService.GetProductIDBySlug(c.Request.Context(), slug)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "product_not_found"})
		return
	}

	// Get authenticated user ID from context (set by JWT middleware)
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

	review, err := h.reviewService.CreateReview(c.Request.Context(), userModel, productID, req.Rating, req.Comment)
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
//
// @Summary      List product reviews
// @Description  Returns published reviews for a product in descending creation order.
// @Tags         reviews
// @Accept       json
// @Produce      json
// @Param        slug   path   string true   "Product slug"
// @Param        page   query  int    false  "Page number"      default(1)
// @Param        limit  query  int    false  "Items per page"   default(10)
// @Success      200  {object}  dto.ReviewListResponse   "Paginated reviews"
// @Failure      400  {object}  dto.ErrorResponse        "Error code: invalid_request"
// @Failure      404  {object}  dto.ErrorResponse        "Error code: product_not_found"
// @Failure      500  {object}  dto.ErrorResponse        "Error code: internal_error"
// @Router       /products/{slug}/reviews [get]
func (h *ReviewHandler) ListReviews(c *gin.Context) {
	slug := c.Param("slug")

	// Get product ID by slug
	productID, err := h.productService.GetProductIDBySlug(c.Request.Context(), slug)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "product_not_found"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	reviews, total, err := h.reviewService.ListReviews(c.Request.Context(), productID, page, limit)
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
//
// @Summary      Get product review summary
// @Description  Returns average rating, total review count, and rating distribution for a product.
// @Tags         reviews
// @Accept       json
// @Produce      json
// @Param        slug path  string  true  "Product slug"
// @Success      200  {object}  dto.ReviewSummary       "Review summary"
// @Failure      400  {object}  dto.ErrorResponse       "Error code: invalid_request"
// @Failure      404  {object}  dto.ErrorResponse       "Error code: product_not_found"
// @Failure      500  {object}  dto.ErrorResponse       "Error code: internal_error"
// @Router       /products/{slug}/reviews/summary [get]
func (h *ReviewHandler) GetSummary(c *gin.Context) {
	slug := c.Param("slug")

	// Get product ID by slug
	productID, err := h.productService.GetProductIDBySlug(c.Request.Context(), slug)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "product_not_found"})
		return
	}

	avg, count, err := h.reviewService.GetAverage(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal_error"})
		return
	}

	dist := map[int]int{}
	reviews, _, err := h.reviewService.ListReviews(c.Request.Context(), productID, 1, 1000)
	if err == nil {
		for _, r := range reviews {
			dist[r.Rating]++
		}
	}
	c.JSON(http.StatusOK, dto.ReviewSummary{Average: avg, Count: count, Distribution: dist})
}

// DeleteReview handles the deletion of a user's own review
//
// @Summary      Delete product review
// @Description  Soft deletes a review created by the authenticated user. Only the review author can delete their own review. The review data is retained for compliance purposes but marked as deleted and no longer visible.
// @Tags         reviews
// @Accept       json
// @Produce      json
// @Param        reviewID  path     string  true  "Review ID (UUID)"
// @Success      200  {object}  dto.MessageResponse  "Review deleted successfully"
// @Failure      401  {object}  dto.ErrorResponse    "Error code: unauthenticated"
// @Failure      403  {object}  dto.ErrorResponse    "Error code: unauthorized"
// @Failure      404  {object}  dto.ErrorResponse    "Error code: review_not_found"
// @Failure      500  {object}  dto.ErrorResponse    "Error code: internal_error"
// @Router       /reviews/{reviewID} [delete]
// @Security     BearerAuth
func (h *ReviewHandler) DeleteReview(c *gin.Context) {
	reviewID := c.Param("reviewID")

	// Get authenticated user ID from context
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

	err = h.reviewService.DeleteOwnReview(c.Request.Context(), reviewID, publicID)
	if err != nil {
		if status, code, ok := handlererrors.MapServiceError(err); ok {
			c.JSON(status, dto.ErrorResponse{Error: code})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal_error"})
		return
	}

	// Log the deletion
	if h.auditService != nil {
		h.auditService.LogDeletionAction(&userModel.ID, userModel.ID, model.AuditActionReviewDeleted, map[string]interface{}{
			"review_id":      reviewID,
			"user_public_id": publicID,
		})
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "review deleted successfully"})
}

// ReportReview handles reporting an abusive or inappropriate review
//
// @Summary      Report a review
// @Description  Allows an authenticated user to report a review for abuse (spam/offensive/improper/other). Prevents self-reporting and duplicate reports per user.
// @Tags         reviews
// @Accept       json
// @Produce      json
// @Param        reviewID  path     string                     true  "Review ID (UUID)"
// @Param        report    body     dto.ReviewReportRequest    true  "Report payload"
// @Success      201  {object}  dto.MessageResponse  "Review reported successfully"
// @Failure      400  {object}  dto.ErrorResponse    "Error code: invalid_request or invalid_category or reason_too_long"
// @Failure      401  {object}  dto.ErrorResponse    "Error code: unauthenticated"
// @Failure      403  {object}  dto.ErrorResponse    "Error code: cannot_report_own_review"
// @Failure      404  {object}  dto.ErrorResponse    "Error code: review_not_found"
// @Failure      409  {object}  dto.ErrorResponse    "Error code: already_reported"
// @Failure      500  {object}  dto.ErrorResponse    "Error code: internal_error"
// @Router       /reviews/{reviewID}/report [post]
// @Security     BearerAuth
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

	// Get authenticated user ID from context
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
