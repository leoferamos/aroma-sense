package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"

	"github.com/leoferamos/aroma-sense/internal/dto"
	"github.com/leoferamos/aroma-sense/internal/repository"
	"github.com/leoferamos/aroma-sense/internal/storage"
)

// ProductService defines the interface for product-related business logic
type ProductService interface {
	CreateProduct(ctx context.Context, input dto.ProductFormDTO, file dto.FileUpload) error
}

type productService struct {
	repo    repository.ProductRepository
	storage storage.ImageStorage
}

func NewProductService(repo repository.ProductRepository, storage storage.ImageStorage) ProductService {
	return &productService{repo: repo, storage: storage}
}

func (s *productService) CreateProduct(ctx context.Context, input dto.ProductFormDTO, file dto.FileUpload) error {
	// Validate the file upload
	if err := file.Validate(); err != nil {
		return err
	}

	// Read first 512 bytes to detect actual content type
	buf := make([]byte, 512)
	n, err := file.Content.Read(buf)
	if err != nil && err != io.EOF {
		return fmt.Errorf("failed to read image: %w", err)
	}
	detectedType := http.DetectContentType(buf[:n])

	// Verify the detected type matches the provided content type
	if detectedType != file.ContentType {
		return fmt.Errorf("content type mismatch: detected %s, provided %s", detectedType, file.ContentType)
	}

	// Generate a unique name for the image
	uuidStr := uuid.New().String()
	var ext string
	switch file.ContentType {
	case "image/jpeg":
		ext = ".jpg"
	case "image/png":
		ext = ".png"
	default:
		ext = ""
	}
	imageName := fmt.Sprintf("product-%s%s", uuidStr, ext)

	combinedReader := io.MultiReader(
		bytes.NewReader(buf[:n]),
		file.Content,
	)

	// Upload the image to storage
	imageURL, err := s.storage.UploadImage(ctx, imageName, combinedReader, file.Size, file.ContentType)
	if err != nil {
		return fmt.Errorf("failed to upload image: %w", err)
	}

	// Call the repository to save to database
	return s.repo.Create(input, imageURL)
}
