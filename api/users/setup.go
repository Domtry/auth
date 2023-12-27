package users

import (
	"auth/repository"
	"auth/service"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// SetupUser config User
func SetupUser(apiGroup *echo.Group, db *gorm.DB) {
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	userService := service.NewUserService(userRepo, roleRepo)
	userHandler := NewUserHandler(userService)

	userGroup := apiGroup.Group("/users")
	RegisterUserRoutes(userGroup, userHandler)
}
