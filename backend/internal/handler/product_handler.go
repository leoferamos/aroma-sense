package handler

import (
	"io"
	"net/http"

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
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var form dto.ProductFormDTO
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	file, fileHeader, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image is required"})
		return
	}
	defer file.Close()

	// Convert multipart.File to oFileUpload abstraction
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read file"})
			return
		}

		// Reset file position
		if seeker, ok := file.(io.Seeker); ok {
			seeker.Seek(0, io.SeekStart)
		}

		fileUpload.ContentType = http.DetectContentType(buf[:n])
	}

	if err := h.productService.CreateProduct(c.Request.Context(), form, fileUpload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Product created successfully"})
}
