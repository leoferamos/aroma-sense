package handler

import (
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/leoferamos/aroma-sense/internal/dto"
	"github.com/leoferamos/aroma-sense/internal/service"
)

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(s service.ProductService) *ProductHandler {
	return &ProductHandler{productService: s}
}

// CreateProduct handles admin product creation
//
// @Summary      Create a new product
// @Description  Creates a new product with image upload (Admin only)
// @Tags         admin
// @Accept       multipart/form-data
// @Produce      json
// @Param        Authorization  header  string  true  "Bearer JWT token"
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
	if err := c.ShouldBind(&form); err != nil {
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

// GetProduct handles fetching a product by its ID
//
// @Summary      Get product by ID
// @Description  Retrieves a specific product by its ID (Admin only)
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        Authorization  header  string  true  "Bearer JWT token"
// @Param        id             path    int     true  "Product ID"
// @Success      200  {object}  dto.ProductResponse  "Product details"
// @Failure      400  {object}  dto.ErrorResponse    "Invalid product ID"
// @Failure      401  {object}  dto.ErrorResponse    "Unauthorized"
// @Failure      403  {object}  dto.ErrorResponse    "Forbidden - Admin only"
// @Failure      404  {object}  dto.ErrorResponse    "Product not found"
// @Router       /admin/products/{id} [get]
// @Security     BearerAuth
func (h *ProductHandler) GetProduct(c *gin.Context) {
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

	c.JSON(http.StatusOK, product)
}

// GetLatestProducts handles fetching the latest products with an optional limit
//
// @Summary      List latest products
// @Description  Retrieves a list of the latest products with optional limit (Public endpoint)
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        limit  query    int  false  "Maximum number of products to return"  default(10)
// @Success      200  {array}   dto.ProductResponse  "List of latest products"
// @Failure      400  {object}  dto.ErrorResponse    "Invalid limit parameter"
// @Failure      500  {object}  dto.ErrorResponse    "Internal server error"
// @Router       /products [get]
func (h *ProductHandler) GetLatestProducts(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid limit parameter"})
		return
	}

	products, err := h.productService.GetLatestProducts(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
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

	c.JSON(http.StatusOK, product)
}

// UpdateProduct handles updating an existing product
//
// @Summary      Update product
// @Description  Updates an existing product (Admin only)
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        Authorization  header  string                      true  "Bearer JWT token"
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
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
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
// @Param        Authorization  header  string  true  "Bearer JWT token"
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
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "Product deleted successfully"})
}
