package bootstrap

import (
	"github.com/leoferamos/aroma-sense/internal/handler"
	"github.com/leoferamos/aroma-sense/internal/repository"
	"github.com/leoferamos/aroma-sense/internal/service"
	"github.com/leoferamos/aroma-sense/internal/storage"
	"gorm.io/gorm"
)

// InitProductModule initializes the product module components and returns a ProductHandler
func InitProductModule(db *gorm.DB, storage storage.ImageStorage) *handler.ProductHandler {
	repo := repository.NewProductRepository(db)
	service := service.NewProductService(repo, storage)
	return handler.NewProductHandler(service)
}
