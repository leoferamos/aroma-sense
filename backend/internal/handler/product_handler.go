package handler

import (
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/leoferamos/aroma-sense/internal/dto"
	"github.com/leoferamos/aroma-sense/internal/service"
)

type ProductHandler struct {
	productService service.ProductService
	reviewService  service.ReviewService
	userService    service.UserProfileService
}

func NewProductHandler(ps service.ProductService, rs service.ReviewService, us service.UserProfileService) *ProductHandler {
	return &ProductHandler{productService: ps, reviewService: rs, userService: us}
}

// CreateProduct handles admin product creation
//
// @Summary      Create a new product
// @Description  Creates a new product with image upload (Admin only)
// @Tags         admin
// @Accept       multipart/form-data
// @Produce      json
// @Param        name           formData  string   true   "Product name"
// @Param        brand          formData  string   true   "Product brand"
// @Param        weight         formData  number   true   "Product weight in ml"
// @Param        description    formData  string   false  "Product description"
// @Param        price          formData  number   true   "Product price"
// @Param        category       formData  string   true   "Product category"
// @Param        notes          formData  array    true   "Product notes (fragrance notes)"
// @Param        stock_quantity formData  integer  true   "Stock quantity"
// @Param        image          formData  file     true   "Product image"
// @Success      201  {object}  dto.MessageResponse  "Product created successfully"
// @Failure      400  {object}  dto.ErrorResponse    "Invalid request data or missing image"
// @Failure      401  {object}  dto.ErrorResponse    "Unauthorized"
// @Failure      403  {object}  dto.ErrorResponse    "Forbidden - Admin only"
// @Router       /admin/products [post]
// @Security     BearerAuth
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var form dto.ProductFormDTO
	if err := c.ShouldBindWith(&form, binding.FormMultipart); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	file, fileHeader, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "image is required"})
		return
	}
	defer file.Close()

	// Convert multipart.File to FileUpload abstraction
	fileUpload := dto.FileUpload{
		Content:     file,
		Name:        fileHeader.Filename,
		Size:        fileHeader.Size,
		ContentType: fileHeader.Header.Get("Content-Type"),
	}

	if fileUpload.ContentType == "" {
		// Read first 512 bytes to detect content type
		buf := make([]byte, 512)
		n, err := file.Read(buf)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "failed to read file"})
			return
		}

		// Reset file position
		if seeker, ok := file.(io.Seeker); ok {
			seeker.Seek(0, io.SeekStart)
		}

		fileUpload.ContentType = http.DetectContentType(buf[:n])
	}

	if err := h.productService.CreateProduct(c.Request.Context(), form, fileUpload); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, dto.MessageResponse{Message: "Product created successfully"})
}

// GetProduct handles fetching a product by its ID or slug
//
// @Summary      Get product by ID or slug
// @Description  Retrieves a specific product by its ID or slug
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id             path    string  true  "Product ID or slug"
// @Success      200  {object}  dto.ProductResponse  "Product details"
// @Failure      404  {object}  dto.ErrorResponse    "Product not found"
// @Router       /products/{id} [get]
func (h *ProductHandler) GetProduct(c *gin.Context) {
	param := c.Param("id")

	// Try to parse as ID first
	if id, err := strconv.Atoi(param); err == nil {
		product, err := h.productService.GetProductByID(c.Request.Context(), uint(id))
		if err == nil {
			c.JSON(http.StatusOK, product)
			return
		}
	}

	// If not a valid ID or ID lookup failed, try as slug
	product, err := h.productService.GetProductBySlug(c.Request.Context(), param)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// GetProductBySlug handles fetching a product by its slug (clean URLs)
//
// @Summary      Get product by slug
// @Description  Retrieves a specific product by its slug for clean URLs
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        slug           path    string  true  "Product slug"
// @Success      200  {object}  dto.ProductResponse  "Product details"
// @Failure      404  {object}  dto.ErrorResponse    "Product not found"
// @Router       /product/{slug} [get]
func (h *ProductHandler) GetProductBySlug(c *gin.Context) {
	slug := c.Param("slug")

	product, err := h.productService.GetProductBySlug(c.Request.Context(), slug)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// GetLatestProducts returns the latest products or performs a search when the
// `query` parameter is present.
//
// @Summary      List or search products
// @Description  If `query` is provided, returns a paginated search envelope; otherwise returns the latest products.
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        query  query    string  false  "Search term (websearch syntax)"
// @Param        page   query    int     false  "Page number (1-based)"  default(1)
// @Param        limit  query    int     false  "Items per page (default 10, max 100)"  default(10)
// @Param        sort   query    string  false  "Sort order: relevance|latest"  default(relevance)
// @Success      200  {array}   dto.ProductResponse        "List of latest products (when query is absent and page=1)"
// @Success      200  {object}  dto.ProductListResponse   "Search results envelope (when query is present) or paginated latest (when query absent and page>1)"
// @Failure      400  {object}  dto.ErrorResponse         "Invalid request parameters"
// @Failure      500  {object}  dto.ErrorResponse         "Internal server error"
// @Router       /products [get]
func (h *ProductHandler) GetLatestProducts(c *gin.Context) {
	const maxLimit = 100

	query := strings.TrimSpace(c.Query("query"))
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid page parameter"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid limit parameter"})
		return
	}
	if limit > maxLimit {
		limit = maxLimit
	}

	if query == "" {
		// return latest products
		products, total, err := h.productService.GetLatestProducts(c.Request.Context(), page, limit)
		if err != nil {
			log.Printf("GetLatestProducts: latest products error: %v", err)
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal server error"})
			return
		}

		if page == 1 {
			c.JSON(http.StatusOK, products)
		} else {
			resp := dto.ProductListResponse{
				Items: products,
				Total: total,
				Page:  page,
				Limit: limit,
			}
			c.JSON(http.StatusOK, resp)
		}
		return
	}

	// parse and validate sort
	sort := c.DefaultQuery("sort", "relevance")
	if sort != "relevance" && sort != "latest" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid sort parameter"})
		return
	}

	items, total, err := h.productService.SearchProducts(c.Request.Context(), query, page, limit, sort)
	if err != nil {
		log.Printf("GetLatestProducts: search error (query=%q, page=%d, limit=%d, sort=%s): %v", query, page, limit, sort, err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal server error"})
		return
	}

	resp := dto.ProductListResponse{
		Items: items,
		Total: total,
		Page:  page,
		Limit: limit,
	}
	c.JSON(http.StatusOK, resp)
}

// GetProductByID handles fetching a single product by ID
//
// @Summary      Get product by ID
// @Description  Retrieves detailed information about a specific product
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Product ID"
// @Success      200  {object}  dto.ProductResponse  "Product details"
// @Failure      400  {object}  dto.ErrorResponse    "Invalid product ID"
// @Failure      404  {object}  dto.ErrorResponse    "Product not found"
// @Router       /products/{id} [get]
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid product ID"})
		return
	}

	product, err := h.productService.GetProductByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "Product not found"})
		return
	}

	// Check if user can review this product
	if rawUserID, exists := c.Get("userID"); exists && rawUserID != "" && h.reviewService != nil && h.userService != nil {
		publicID := rawUserID.(string)
		if userModel, err := h.userService.GetByPublicID(publicID); err == nil {
			can, reason, canErr := h.reviewService.CanUserReview(c.Request.Context(), userModel, product.ID)
			if can && canErr == nil {
				trueVal := true
				product.CanReview = &trueVal
			} else if canErr == nil {
				falseVal := false
				product.CanReview = &falseVal
				if reason == "profile_incomplete" || reason == "not_delivered" || reason == "already_reviewed" {
					product.CannotReviewReason = &reason
				}
			}
		}
	}

	c.JSON(http.StatusOK, product)
}

// UpdateProduct handles updating an existing product
//
// @Summary      Update product
// @Description  Updates an existing product (Admin only)
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        id             path    int                         true  "Product ID"
// @Param        product        body    dto.UpdateProductRequest    true  "Product update data"
// @Success      200  {object}  dto.MessageResponse  "Product updated successfully"
// @Failure      400  {object}  dto.ErrorResponse    "Invalid product ID or request data"
// @Failure      401  {object}  dto.ErrorResponse    "Unauthorized"
// @Failure      403  {object}  dto.ErrorResponse    "Forbidden - Admin only"
// @Failure      500  {object}  dto.ErrorResponse    "Internal server error"
// @Router       /admin/products/{id} [patch]
// @Security     BearerAuth
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid product ID"})
		return
	}

	var input dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.productService.UpdateProduct(c.Request.Context(), uint(id), input); err != nil {
		log.Printf("UpdateProduct: service error: %v", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal server error"})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "Product updated successfully"})
}

// DeleteProduct handles deleting an existing product
//
// @Summary      Delete product
// @Description  Deletes an existing product (Admin only)
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        id             path    int     true  "Product ID"
// @Success      200  {object}  dto.MessageResponse  "Product deleted successfully"
// @Failure      400  {object}  dto.ErrorResponse    "Invalid product ID"
// @Failure      401  {object}  dto.ErrorResponse    "Unauthorized"
// @Failure      403  {object}  dto.ErrorResponse    "Forbidden - Admin only"
// @Failure      500  {object}  dto.ErrorResponse    "Internal server error"
// @Router       /admin/products/{id} [delete]
// @Security     BearerAuth
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid product ID"})
		return
	}

	if err := h.productService.DeleteProduct(c.Request.Context(), uint(id)); err != nil {
		log.Printf("DeleteProduct: service error: %v", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal server error"})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "Product deleted successfully"})
}
