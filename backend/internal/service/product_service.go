package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/leoferamos/aroma-sense/internal/dto"
	"github.com/leoferamos/aroma-sense/internal/repository"
)

type ProductService interface {
	CreateProduct(input dto.ProductFormDTO, file multipart.File, fileHeader *multipart.FileHeader) error
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) CreateProduct(input dto.ProductFormDTO, file multipart.File, fileHeader *multipart.FileHeader) error {
	if fileHeader.Size > 5*1024*1024 {
		return fmt.Errorf("image too large (max 5MB)")
	}

	// Validates image type
	allowedTypes := []string{"image/jpeg", "image/png"}
	buf := make([]byte, 512)
	_, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return fmt.Errorf("failed to read image: %w", err)
	}
	filetype := http.DetectContentType(buf)
	isValidType := false
	for _, t := range allowedTypes {
		if filetype == t {
			isValidType = true
			break
		}
	}
	if !isValidType {
		return fmt.Errorf("invalid image type: %s", filetype)
	}
	file.Seek(0, io.SeekStart)

	// Generate a unique name for the image
	uuidStr := uuid.New().String()
	ext := ""
	if strings.HasSuffix(fileHeader.Filename, ".png") {
		ext = ".png"
	} else {
		ext = ".jpg"
	}
	imageName := fmt.Sprintf("product-%s%s", uuidStr, ext)
	imageURL := fmt.Sprintf("https://bucket-url/%s", imageName)

	return s.repo.Create(input, imageURL)
}
