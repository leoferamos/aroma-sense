package repository

import (
	"github.com/leoferamos/aroma-sense/internal/dto"
	"github.com/leoferamos/aroma-sense/internal/model"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(input dto.ProductFormDTO, imageURL string) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

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
