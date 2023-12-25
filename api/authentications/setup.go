package authentications

import (
	"auth/repository"
	"auth/service"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// SetupAuthentication config User
func SetupAuthentication(apiGroup *echo.Group, db *gorm.DB) {
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	userService := service.NewAuthenticationService(userRepo, roleRepo)
	userHandler := NewAuthenticationHandler(userService)

	authGroup := apiGroup.Group("/auth")
	RegisterAuthRoutes(authGroup, userHandler)
}
