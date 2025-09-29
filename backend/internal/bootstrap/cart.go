package bootstrap

import (
	"github.com/leoferamos/aroma-sense/internal/handler"
	"github.com/leoferamos/aroma-sense/internal/repository"
	"github.com/leoferamos/aroma-sense/internal/service"
	"gorm.io/gorm"
)

// InitCartModule initializes the cart module components and returns a CartHandler
func InitCartModule(db *gorm.DB) *handler.CartHandler {
	cartRepo := repository.NewCartRepository(db)
	cartService := service.NewCartService(cartRepo)
	return handler.NewCartHandler(cartService)
}
