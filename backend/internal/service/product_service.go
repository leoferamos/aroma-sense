package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/leoferamos/aroma-sense/internal/dto"
	"github.com/leoferamos/aroma-sense/internal/integrations/ai/embeddings"
	"github.com/leoferamos/aroma-sense/internal/repository"
	"github.com/leoferamos/aroma-sense/internal/storage"
	"github.com/leoferamos/aroma-sense/internal/utils"
)

// ProductService defines the interface for product-related business logic
type ProductService interface {
	CreateProduct(ctx context.Context, input dto.ProductFormDTO, file dto.FileUpload) error
	GetProductByID(ctx context.Context, id uint) (dto.ProductResponse, error)
	GetLatestProducts(ctx context.Context, page int, limit int) ([]dto.ProductResponse, int, error)
	SearchProducts(ctx context.Context, query string, page int, limit int, sort string) ([]dto.ProductResponse, int, error)
	UpdateProduct(ctx context.Context, id uint, input dto.UpdateProductRequest) error
	DeleteProduct(ctx context.Context, id uint) error
}

type productService struct {
	repo       repository.ProductRepository
	storage    storage.ImageStorage
	embeddings embeddings.Provider
}

func NewProductService(repo repository.ProductRepository, storage storage.ImageStorage, embProvider embeddings.Provider) ProductService {
	return &productService{repo: repo, storage: storage, embeddings: embProvider}
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

	// Upload the image and thumbnail to storage
	origURL, thumbURL, err := s.storage.UploadImageWithThumbnail(ctx, imageName, combinedReader, file.Size, file.ContentType, 256, 256)
	if err != nil {
		return fmt.Errorf("failed to upload image: %w", err)
	}

	// Call the repository to save to database
	productID, err := s.repo.Create(input, origURL, thumbURL)
	if err != nil {
		return err
	}

	// Generate embedding asynchronously
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("PANIC in embedding generation for product %d: %v\n", productID, r)
			}
		}()

		text := s.buildProductText(input)
		emb, embErr := s.embeddings.Embed([]string{text})
		if embErr == nil && len(emb) > 0 && len(emb[0]) > 0 {
			_ = s.repo.UpsertProductEmbedding(productID, emb[0])
		}
	}()

	return nil
}

// buildProductText creates a text representation of the product for embedding.
func (s *productService) buildProductText(input dto.ProductFormDTO) string {
	var parts []string
	parts = append(parts, input.Name)
	parts = append(parts, input.Brand)
	parts = append(parts, input.Description)
	if len(input.Accords) > 0 {
		parts = append(parts, "Accords: "+strings.Join(input.Accords, ", "))
	}
	if len(input.Occasions) > 0 {
		parts = append(parts, "Occasions: "+strings.Join(input.Occasions, ", "))
	}
	if len(input.Seasons) > 0 {
		parts = append(parts, "Seasons: "+strings.Join(input.Seasons, ", "))
	}
	if input.Intensity != "" {
		parts = append(parts, "Intensity: "+input.Intensity)
	}
	if len(input.NotesTop) > 0 {
		parts = append(parts, "Top notes: "+strings.Join(input.NotesTop, ", "))
	}
	if len(input.NotesHeart) > 0 {
		parts = append(parts, "Heart notes: "+strings.Join(input.NotesHeart, ", "))
	}
	if len(input.NotesBase) > 0 {
		parts = append(parts, "Base notes: "+strings.Join(input.NotesBase, ", "))
	}
	return strings.Join(parts, ". ")
}

// GetProductByID retrieves a product by its ID and maps it to a DTO
func (s *productService) GetProductByID(ctx context.Context, id uint) (dto.ProductResponse, error) {
	product, err := s.repo.FindByID(id)
	if err != nil {
		return dto.ProductResponse{}, fmt.Errorf("failed to get product: %w", err)
	}

	return dto.ProductResponse{
		ID:            product.ID,
		Name:          product.Name,
		Brand:         product.Brand,
		Weight:        product.Weight,
		Description:   product.Description,
		Price:         product.Price,
		ImageURL:      product.ImageURL,
		ThumbnailURL:  product.ThumbnailURL,
		Slug:          product.Slug,
		Accords:       product.Accords,
		Occasions:     product.Occasions,
		Seasons:       product.Seasons,
		Intensity:     product.Intensity,
		Gender:        product.Gender,
		PriceRange:    product.PriceRange,
		NotesTop:      product.NotesTop,
		NotesHeart:    product.NotesHeart,
		NotesBase:     product.NotesBase,
		Category:      product.Category,
		StockQuantity: product.StockQuantity,
		CreatedAt:     product.CreatedAt,
		UpdatedAt:     product.UpdatedAt,
	}, nil
}

// GetLatestProducts retrieves the latest products with pagination
func (s *productService) GetLatestProducts(ctx context.Context, page int, limit int) ([]dto.ProductResponse, int, error) {
	if page < 1 {
		page = 1
	}
	const maxLimit = 100
	if limit <= 0 {
		limit = 10
	}
	if limit > maxLimit {
		limit = maxLimit
	}

	offset := (page - 1) * limit

	products, total, err := s.repo.FindAllPaginated(limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get products: %w", err)
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
			ThumbnailURL:  p.ThumbnailURL,
			Slug:          p.Slug,
			Accords:       p.Accords,
			Occasions:     p.Occasions,
			Seasons:       p.Seasons,
			Intensity:     p.Intensity,
			Gender:        p.Gender,
			PriceRange:    p.PriceRange,
			NotesTop:      p.NotesTop,
			NotesHeart:    p.NotesHeart,
			NotesBase:     p.NotesBase,
			Category:      p.Category,
			StockQuantity: p.StockQuantity,
			CreatedAt:     p.CreatedAt,
		})
	}

	return response, total, nil
}

// UpdateProduct updates an existing product with the provided details
func (s *productService) UpdateProduct(ctx context.Context, id uint, input dto.UpdateProductRequest) error {
	product, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("product not found: %w", err)
	}

	nameChanged := false
	brandChanged := false
	if input.Name != nil {
		product.Name = *input.Name
		nameChanged = true
	}
	if input.Brand != nil {
		product.Brand = *input.Brand
		brandChanged = true
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
	if input.Accords != nil {
		product.Accords = *input.Accords
	}
	if input.Occasions != nil {
		product.Occasions = *input.Occasions
	}
	if input.Seasons != nil {
		product.Seasons = *input.Seasons
	}
	if input.Intensity != nil {
		product.Intensity = *input.Intensity
	}
	if input.Gender != nil {
		product.Gender = *input.Gender
	}
	if input.PriceRange != nil {
		product.PriceRange = *input.PriceRange
	}
	if input.NotesTop != nil {
		product.NotesTop = *input.NotesTop
	}
	if input.NotesHeart != nil {
		product.NotesHeart = *input.NotesHeart
	}
	if input.NotesBase != nil {
		product.NotesBase = *input.NotesBase
	}

	// If name or brand changed, regenerate slug.
	if nameChanged || brandChanged {
		base := utils.Slugify(product.Brand, product.Name)
		if slug, err := s.repo.EnsureUniqueSlug(base); err == nil {
			product.Slug = slug
		} else {
			product.Slug = base
		}
	}

	return s.repo.Update(&product)
}

// DeleteProduct removes a product by its ID
func (s *productService) DeleteProduct(ctx context.Context, id uint) error {
	product, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("product not found: %w", err)
	}

	// Delete images from storage first
	if product.ImageURL != "" {
		// Extract image name from URL
		imageName := extractImageNameFromURL(product.ImageURL)
		if imageName != "" {
			if err := s.storage.DeleteImage(ctx, imageName); err != nil {
				fmt.Printf("Warning: failed to delete image %s: %v\n", imageName, err)
			}
		}
	}

	// Also delete thumbnail if it exists and is different from main image
	if product.ThumbnailURL != "" && product.ThumbnailURL != product.ImageURL {
		thumbName := extractImageNameFromURL(product.ThumbnailURL)
		if thumbName != "" {
			if err := s.storage.DeleteImage(ctx, thumbName); err != nil {
				fmt.Printf("Warning: failed to delete thumbnail %s: %v\n", thumbName, err)
			}
		}
	}

	return s.repo.Delete(id)
}

// SearchProducts performs a product search with pagination and sorting.
func (s *productService) SearchProducts(ctx context.Context, query string, page int, limit int, sort string) ([]dto.ProductResponse, int, error) {
	if page < 1 {
		page = 1
	}
	// Enforce sensible limits (defense in depth)
	const maxLimit = 100
	if limit <= 0 {
		limit = 10
	}
	if limit > maxLimit {
		limit = maxLimit
	}

	offset := (page - 1) * limit

	products, total, err := s.repo.SearchProducts(ctx, query, limit, offset, sort)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search products: %w", err)
	}

	var resp []dto.ProductResponse
	for _, p := range products {
		resp = append(resp, dto.ProductResponse{
			ID:            p.ID,
			Name:          p.Name,
			Brand:         p.Brand,
			Weight:        p.Weight,
			Description:   p.Description,
			Price:         p.Price,
			ImageURL:      p.ImageURL,
			ThumbnailURL:  p.ThumbnailURL,
			Slug:          p.Slug,
			Accords:       p.Accords,
			Occasions:     p.Occasions,
			Seasons:       p.Seasons,
			Intensity:     p.Intensity,
			Gender:        p.Gender,
			PriceRange:    p.PriceRange,
			NotesTop:      p.NotesTop,
			NotesHeart:    p.NotesHeart,
			NotesBase:     p.NotesBase,
			Category:      p.Category,
			StockQuantity: p.StockQuantity,
			CreatedAt:     p.CreatedAt,
			UpdatedAt:     p.UpdatedAt,
		})
	}

	return resp, total, nil
}

// extractImageNameFromURL extracts the image name from a Supabase storage URL
func extractImageNameFromURL(imageURL string) string {
	// URL format: https://domain.com/storage/v1/object/public/bucket/image-name.jpg
	// We need to extract "image-name.jpg"
	parts := strings.Split(imageURL, "/")
	if len(parts) >= 2 {
		return parts[len(parts)-1]
	}
	return ""
}
