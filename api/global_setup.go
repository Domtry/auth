package api

import (
	"auth/api/users"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func GlobalSetup(ech *echo.Echo, db *gorm.DB) {
	apiGroup := ech.Group("/api/v1")
	users.SetupUser(apiGroup, db)
}
