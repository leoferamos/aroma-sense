package bootstrap

import (
	"github.com/leoferamos/aroma-sense/internal/handler"
	"github.com/leoferamos/aroma-sense/internal/repository"
	"github.com/leoferamos/aroma-sense/internal/service"
	"gorm.io/gorm"
)

// InitUserModule initializes the user module components and returns a UserHandler
func InitUserModule(db *gorm.DB) *handler.UserHandler {
	userRepo := repository.NewUserRepository(db)
	cartRepo := repository.NewCartRepository(db)
	cartService := service.NewCartService(cartRepo)
	userService := service.NewUserService(userRepo, cartService)
	return handler.NewUserHandler(userService)
}
