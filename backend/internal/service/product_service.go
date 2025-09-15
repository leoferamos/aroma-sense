package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/google/uuid"

	"github.com/leoferamos/aroma-sense/internal/dto"
	"github.com/leoferamos/aroma-sense/internal/repository"
	"github.com/leoferamos/aroma-sense/internal/storage"
)

// ProductService defines the interface for product-related business logic
type ProductService interface {
	CreateProduct(ctx context.Context, input dto.ProductFormDTO, file multipart.File, fileHeader *multipart.FileHeader) error
}

type productService struct {
	repo    repository.ProductRepository
	storage storage.ImageStorage
}

func NewProductService(repo repository.ProductRepository, storage storage.ImageStorage) ProductService {
	return &productService{repo: repo, storage: storage}
}

func (s *productService) CreateProduct(ctx context.Context, input dto.ProductFormDTO, file multipart.File, fileHeader *multipart.FileHeader) error {
	if fileHeader.Size > 5*1024*1024 {
		return fmt.Errorf("image too large (max 5MB)")
	}

	// Validates image type
	allowedTypes := []string{"image/jpeg", "image/png"}
	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return fmt.Errorf("failed to read image: %w", err)
	}
	filetype := http.DetectContentType(buf[:n])
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
	var ext string
	switch filetype {
	case "image/jpeg":
		ext = ".jpg"
	case "image/png":
		ext = ".png"
	default:
		ext = ""
	}
	imageName := fmt.Sprintf("product-%s%s", uuidStr, ext)

	// Uploads the image to storage
	imageURL, err := s.storage.UploadImage(ctx, imageName, file, fileHeader)
	if err != nil {
		return fmt.Errorf("failed to upload image: %w", err)
	}
	
	// Call the repository to save to database
	return s.repo.Create(input, imageURL)
}
