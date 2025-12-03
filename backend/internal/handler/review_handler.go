package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/dto"
	"github.com/leoferamos/aroma-sense/internal/service"
)

type ReviewHandler struct {
	service     service.ReviewService
	userService service.UserProfileService
}

func NewReviewHandler(s service.ReviewService, userService service.UserProfileService) *ReviewHandler {
	return &ReviewHandler{service: s, userService: userService}
}

// Create review handles the creation of a product review
//
// @Summary      Create product review
// @Description  Creates a review for a delivered product. Requires authentication, a display_name set, and at least one delivered order containing the product. One review per user/product.
// @Tags         reviews
// @Accept       json
// @Produce      json
// @Param        id      path     int               true  "Product ID"
// @Param        review  body     dto.ReviewRequest true  "Review payload"
// @Success      201  {object}  dto.ReviewResponse       "Review created"
// @Failure      400  {object}  dto.ErrorResponse        "Validation error"
// @Failure      401  {object}  dto.ErrorResponse        "Unauthorized"
// @Failure      403  {object}  dto.ErrorResponse        "Forbidden (not delivered or profile incomplete)"
// @Failure      404  {object}  dto.ErrorResponse        "Product not found"
// @Failure      409  {object}  dto.ErrorResponse        "Already reviewed"
// @Failure      500  {object}  dto.ErrorResponse        "Internal error"
// @Router       /products/{id}/reviews [post]
// @Security     BearerAuth
func (h *ReviewHandler) CreateReview(c *gin.Context) {
	productID, err := parseProductID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid product id"})
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
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	if userModel.DisplayName == nil || strings.TrimSpace(getPtrVal(userModel.DisplayName)) == "" {
		c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "profile_incomplete"})
		return
	}

	review, err := h.service.CreateReview(c.Request.Context(), userModel, productID, req.Rating, req.Comment)
	if err != nil {
		if status, message, ok := mapServiceError(err); ok {
			// Preserve detailed validation messages when available
			if errors.Is(err, service.ErrReviewInvalidRating) || errors.Is(err, service.ErrReviewCommentTooLong) {
				c.JSON(status, dto.ErrorResponse{Error: err.Error()})
				return
			}
			c.JSON(status, dto.ErrorResponse{Error: message})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal error"})
		return
	}

	resp := dto.ReviewResponse{
		ID:            review.ID,
		Rating:        review.Rating,
		Comment:       review.Comment,
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
// @Param        id     path   int   true   "Product ID"
// @Param        page   query  int   false  "Page number"      default(1)
// @Param        limit  query  int   false  "Items per page"   default(10)
// @Success      200  {object}  dto.ReviewListResponse   "Paginated reviews"
// @Failure      400  {object}  dto.ErrorResponse        "Invalid product id"
// @Failure      500  {object}  dto.ErrorResponse        "Internal error"
// @Router       /products/{id}/reviews [get]
func (h *ReviewHandler) ListReviews(c *gin.Context) {
	productID, err := parseProductID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid product id"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	reviews, total, err := h.service.ListReviews(c.Request.Context(), productID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal error"})
		return
	}

	items := make([]dto.ReviewResponse, 0, len(reviews))
	for _, r := range reviews {
		display := ""
		if r.User != nil && r.User.DisplayName != nil {
			display = getPtrVal(r.User.DisplayName)
		}
		items = append(items, dto.ReviewResponse{
			ID:            r.ID,
			Rating:        r.Rating,
			Comment:       r.Comment,
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
// @Param        id   path  int  true  "Product ID"
// @Success      200  {object}  dto.ReviewSummary       "Review summary"
// @Failure      400  {object}  dto.ErrorResponse       "Invalid product id"
// @Failure      500  {object}  dto.ErrorResponse       "Internal error"
// @Router       /products/{id}/reviews/summary [get]
func (h *ReviewHandler) GetSummary(c *gin.Context) {
	productID, err := parseProductID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid product id"})
		return
	}

	avg, count, err := h.service.GetAverage(c.Request.Context(), productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal error"})
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

func parseProductID(raw string) (uint, error) {
	id, err := strconv.Atoi(raw)
	if err != nil || id < 1 {
		return 0, err
	}
	return uint(id), nil
}

func getPtrVal(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}
