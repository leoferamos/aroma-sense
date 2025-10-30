package repository

import (
	"github.com/leoferamos/aroma-sense/internal/dto"
	"github.com/leoferamos/aroma-sense/internal/model"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(input dto.ProductFormDTO, imageURL string) error
	FindAll(limit int) ([]model.Product, error)
	FindByID(id uint) (model.Product, error)
	Update(product *model.Product) error
	Delete(id uint) error
	DecrementStock(productID uint, quantity int) error
}
// DecrementStock decreases the stock quantity of a product
func (r *productRepository) DecrementStock(productID uint, quantity int) error {
	return r.db.Model(&model.Product{}).
		Where("id = ? AND stock_quantity >= ?", productID, quantity).
		UpdateColumn("stock_quantity", gorm.Expr("stock_quantity - ?", quantity)).Error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

// Create inserts a new product into the database
func (r *productRepository) Create(input dto.ProductFormDTO, imageURL string) error {
	notes := ""
	if len(input.Notes) > 0 {
		notes = input.Notes[0]
		if len(input.Notes) > 1 {
			for _, n := range input.Notes[1:] {
				notes += ", " + n
			}
		}
	}

	product := model.Product{
		Name:          input.Name,
		Brand:         input.Brand,
		Weight:        input.Weight,
		Description:   input.Description,
		Price:         input.Price,
		ImageURL:      imageURL,
		Category:      input.Category,
		Notes:         notes,
		StockQuantity: input.StockQuantity,
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
