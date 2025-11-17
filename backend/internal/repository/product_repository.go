package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/leoferamos/aroma-sense/internal/dto"
	"github.com/leoferamos/aroma-sense/internal/model"
	"github.com/leoferamos/aroma-sense/utils"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(input dto.ProductFormDTO, imageURL string, thumbnailURL string) error
	FindAll(limit int) ([]model.Product, error)
	FindByID(id uint) (model.Product, error)
	SearchProducts(ctx context.Context, query string, limit int, offset int, sort string) ([]model.Product, int, error)
	Update(product *model.Product) error
	Delete(id uint) error
	DecrementStock(productID uint, quantity int) error
	EnsureUniqueSlug(base string) (string, error)
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

// Create inserts a new product into the database
func (r *productRepository) Create(input dto.ProductFormDTO, imageURL string, thumbnailURL string) error {
	// Generate unique slug from brand + name
	base := utils.Slugify(input.Brand, input.Name)
	slug, err := r.uniqueSlug(base)
	if err != nil {
		return err
	}

	product := model.Product{
		Name:          input.Name,
		Brand:         input.Brand,
		Weight:        input.Weight,
		Description:   input.Description,
		Price:         input.Price,
		ImageURL:      imageURL,
		ThumbnailURL:  thumbnailURL,
		Slug:          slug,
		Category:      input.Category,
		StockQuantity: input.StockQuantity,
		Accords:       input.Accords,
		Occasions:     input.Occasions,
		Seasons:       input.Seasons,
		Intensity:     input.Intensity,
		Gender:        input.Gender,
		PriceRange:    input.PriceRange,
		NotesTop:      input.NotesTop,
		NotesHeart:    input.NotesHeart,
		NotesBase:     input.NotesBase,
	}
	return r.db.Create(&product).Error
}

// FindAll retrieves all products, limited by the specified number
func (r *productRepository) FindAll(limit int) ([]model.Product, error) {
	var products []model.Product
	query := r.db.Order("created_at desc")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&products).Error
	return products, err
}

// FindByID retrieves a product by its ID
func (r *productRepository) FindByID(id uint) (model.Product, error) {
	var product model.Product
	err := r.db.First(&product, id).Error
	return product, err
}

// Update updates an existing product in the database
func (r *productRepository) Update(product *model.Product) error {
	return r.db.Save(product).Error
}

// Delete removes a product from the database by its ID
func (r *productRepository) Delete(id uint) error {
	return r.db.Delete(&model.Product{}, id).Error
}

// DecrementStock decreases the stock quantity of a product
func (r *productRepository) DecrementStock(productID uint, quantity int) error {
	return r.db.Model(&model.Product{}).
		Where("id = ? AND stock_quantity >= ?", productID, quantity).
		UpdateColumn("stock_quantity", gorm.Expr("stock_quantity - ?", quantity)).Error
}

// SearchProducts performs a search with pagination and sort.
func (r *productRepository) SearchProducts(ctx context.Context, query string, limit int, offset int, sort string) ([]model.Product, int, error) {
	var products []model.Product

	// Build SQL depending on sort preference
	var selectSQL string
	var args []interface{}

	if sort == "latest" {
		selectSQL = `
		SELECT p.*
		FROM products p
		WHERE p.search_vector @@ websearch_to_tsquery('portuguese', unaccent(?))
		ORDER BY p.created_at DESC
		LIMIT ? OFFSET ?
		`
		args = []interface{}{query, limit, offset}
	} else {
		// relevance (default)
		selectSQL = `
		SELECT p.*
		FROM products p
		WHERE p.search_vector @@ websearch_to_tsquery('portuguese', unaccent(?))
		ORDER BY ts_rank_cd(p.search_vector, websearch_to_tsquery('portuguese', unaccent(?))) DESC, p.created_at DESC
		LIMIT ? OFFSET ?
		`
		args = []interface{}{query, query, limit, offset}
	}

	if err := r.db.WithContext(ctx).Raw(selectSQL, args...).Scan(&products).Error; err != nil {
		return nil, 0, err
	}

	// Count total matches
	var total int64
	countSQL := `SELECT COUNT(*) FROM products p WHERE p.search_vector @@ websearch_to_tsquery('portuguese', unaccent(?))`
	if err := r.db.WithContext(ctx).Raw(countSQL, query).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	return products, int(total), nil
}

// uniqueSlug ensures the provided base slug is unique.
func (r *productRepository) uniqueSlug(base string) (string, error) {
	candidate := base
	var count int64
	// Fast path
	if err := r.db.Model(&model.Product{}).Where("slug = ?", candidate).Count(&count).Error; err != nil {
		return "", err
	}
	if count == 0 {
		return candidate, nil
	}
	// Try with numeric suffixes
	for i := 2; i < 1000; i++ {
		candidate = fmt.Sprintf("%s-%d", base, i)
		if err := r.db.Model(&model.Product{}).Where("slug = ?", candidate).Count(&count).Error; err != nil {
			return "", err
		}
		if count == 0 {
			return candidate, nil
		}
	}
	// As a last resort append a timestamp-ish suffix
	suffix := time.Now().Unix() % 100000
	candidate = fmt.Sprintf("%s-%d", base, suffix)
	return candidate, nil
}

// EnsureUniqueSlug exposes slug uniqueness to callers outside this package.
func (r *productRepository) EnsureUniqueSlug(base string) (string, error) {
	return r.uniqueSlug(base)
}
