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
	GetLatestProducts(ctx context.Context, limit int) ([]dto.ProductResponse, error)
	UpdateProduct(ctx context.Context, id uint, input dto.UpdateProductRequest) error
	DeleteProduct(ctx context.Context, id uint) error
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

// GetLatestProducts retrieves the latest products up to the specified limit
func (s *productService) GetLatestProducts(ctx context.Context, limit int) ([]dto.ProductResponse, error) {
	products, err := s.repo.FindAll(limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}

	var response []dto.ProductResponse
	for _, p := range products {
		response = append(response, dto.ProductResponse{
			ID:            p.ID,
			Name:          p.Name,
			Brand:         p.Brand,
			Weight:        p.Weight,
			Description:   p.Description,
			Price:         p.Price,
			ImageURL:      p.ImageURL,
			Category:      p.Category,
			Notes:         p.Notes,
			StockQuantity: p.StockQuantity,
			CreatedAt:     p.CreatedAt,
		})
	}

	return response, nil
}

// UpdateProduct updates an existing product with the provided details
func (s *productService) UpdateProduct(ctx context.Context, id uint, input dto.UpdateProductRequest) error {
	product, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("product not found: %w", err)
	}

	if input.Name != nil {
		product.Name = *input.Name
	}
	if input.Brand != nil {
		product.Brand = *input.Brand
	}
	if input.Weight != nil {
		product.Weight = *input.Weight
	}
	if input.Description != nil {
		product.Description = *input.Description
	}
	if input.Price != nil {
		product.Price = *input.Price
	}
	if input.Category != nil {
		product.Category = *input.Category
	}
	if input.StockQuantity != nil {
		product.StockQuantity = *input.StockQuantity
	}
	if len(input.Notes) > 0 {
		notes := input.Notes[0]
		if len(input.Notes) > 1 {
			for _, n := range input.Notes[1:] {
				notes += ", " + n
			}
		}
		product.Notes = notes
	}

	return s.repo.Update(&product)
}

func (s *productService) DeleteProduct(ctx context.Context, id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("product not found: %w", err)
	}
	return s.repo.Delete(id)
}
