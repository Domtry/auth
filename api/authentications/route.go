package authentications

import (
	"auth/middlewares"
	"github.com/labstack/echo/v4"
)

func RegisterAuthRoutes(apiGroup *echo.Group, handler *AuthenticationHandler) {
	apiGroup.POST("/login", handler.LoginHandler)
	apiGroup.GET("/me", handler.UserProfil, middlewares.IsAuthorizedMiddle)
}
