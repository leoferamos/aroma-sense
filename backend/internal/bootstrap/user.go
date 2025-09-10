package bootstrap

import (
	"github.com/leoferamos/aroma-sense/internal/handler"
	"github.com/leoferamos/aroma-sense/internal/repository"
	"github.com/leoferamos/aroma-sense/internal/service"
	"gorm.io/gorm"
)

// InitUserModule initializes the user module components and returns a UserHandler
func InitUserModule(db *gorm.DB) *handler.UserHandler {
	repo := repository.NewUserRepository(db)
	service := service.NewUserService(repo)
	return handler.NewUserHandler(service)
}
