package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/leoferamos/aroma-sense/internal/dto"
	"github.com/leoferamos/aroma-sense/internal/model"
	"github.com/leoferamos/aroma-sense/internal/utils"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(input dto.ProductFormDTO, imageURL string, thumbnailURL string) (uint, error)
	FindAll(limit int) ([]model.Product, error)
	FindAllPaginated(limit int, offset int) ([]model.Product, int, error)
	FindByID(id uint) (model.Product, error)
	SearchProducts(ctx context.Context, query string, limit int, offset int, sort string) ([]model.Product, int, error)
	Update(product *model.Product) error
	Delete(id uint) error
	DecrementStock(productID uint, quantity int) error
	EnsureUniqueSlug(base string) (string, error)
	UpsertProductEmbedding(productID uint, embedding []float32) error
	HasProductEmbedding(productID uint) (bool, error)
	FindSimilarProductsByEmbedding(ctx context.Context, embedding []float32, limit int) ([]model.Product, error)
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

// Create inserts a new product into the database
func (r *productRepository) Create(input dto.ProductFormDTO, imageURL string, thumbnailURL string) (uint, error) {
	// Generate unique slug from brand + name
	base := utils.Slugify(input.Brand, input.Name)
	slug, err := r.uniqueSlug(base)
	if err != nil {
		return 0, err
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
	if err := r.db.Create(&product).Error; err != nil {
		return 0, err
	}
	return product.ID, nil
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

// FindAllPaginated retrieves products with pagination support
func (r *productRepository) FindAllPaginated(limit int, offset int) ([]model.Product, int, error) {
	var products []model.Product
	var total int64

	// Count total products
	if err := r.db.Model(&model.Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch paginated products
	query := r.db.Order("created_at desc")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	err := query.Find(&products).Error
	return products, int(total), err
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

// UpsertProductEmbedding stores the embedding for a product.
func (r *productRepository) UpsertProductEmbedding(productID uint, embedding []float32) error {
	// Marshal embedding to JSON
	b, err := json.Marshal(embedding)
	if err != nil {
		return err
	}
	sql := `INSERT INTO product_embeddings (product_id, embedding) VALUES (?, ?::jsonb)
		ON CONFLICT (product_id) DO UPDATE SET embedding = EXCLUDED.embedding`
	return r.db.Exec(sql, productID, string(b)).Error
}

// HasProductEmbedding returns true if an embedding row exists for the given product_id.
func (r *productRepository) HasProductEmbedding(productID uint) (bool, error) {
	var exists bool
	sql := `SELECT EXISTS (SELECT 1 FROM product_embeddings WHERE product_id = ?)`
	if err := r.db.Raw(sql, productID).Scan(&exists).Error; err != nil {
		return false, err
	}
	return exists, nil
}

// FindSimilarProductsByEmbedding finds top-k products similar to the given embedding using cosine similarity.
func (r *productRepository) FindSimilarProductsByEmbedding(ctx context.Context, embedding []float32, limit int) ([]model.Product, error) {
	if len(embedding) == 0 {
		return []model.Product{}, nil
	}

	// Fetch all embeddings
	var rows []struct {
		ProductID uint   `json:"product_id"`
		Embedding string `json:"embedding"`
	}
	if err := r.db.WithContext(ctx).Table("product_embeddings").Find(&rows).Error; err != nil {
		return nil, err
	}

	type scoredProduct struct {
		product model.Product
		score   float32
	}

	var candidates []scoredProduct
	for _, row := range rows {
		var emb []float32
		if err := json.Unmarshal([]byte(row.Embedding), &emb); err != nil {
			continue
		}
		if len(emb) != len(embedding) {
			continue
		}
		score := cosineSimilarity(embedding, emb)
		if score > 0 {
			var prod model.Product
			if err := r.db.Raw("SELECT * FROM products WHERE id = ?", row.ProductID).Scan(&prod).Error; err != nil {
				continue
			}
			candidates = append(candidates, scoredProduct{product: prod, score: score})
		}
	}

	// Sort by score descending and take top limit
	for i := 0; i < len(candidates)-1; i++ {
		for j := i + 1; j < len(candidates); j++ {
			if candidates[i].score < candidates[j].score {
				candidates[i], candidates[j] = candidates[j], candidates[i]
			}
		}
	}

	var results []model.Product
	for i := 0; i < len(candidates) && i < limit; i++ {
		results = append(results, candidates[i].product)
	}

	return results, nil
}

// cosineSimilarity computes cosine similarity between two vectors.
func cosineSimilarity(a, b []float32) float32 {
	var dot, normA, normB float32
	for i := range a {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dot / (float32(math.Sqrt(float64(normA))) * float32(math.Sqrt(float64(normB))))
}
