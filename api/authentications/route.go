package authentications

import (
	"auth/middlewares"
	"github.com/labstack/echo/v4"
)

func RegisterAuthRoutes(apiGroup *echo.Group, handler *AuthenticationHandler) {
	apiGroup.POST("/login", handler.LoginHandler)                                                 //OK
	apiGroup.GET("/me", handler.UserProfilHandler, middlewares.IsAuthorizedMiddle)                //OK
	apiGroup.PUT("/reset_password", handler.ResetPasswordHandler, middlewares.IsAuthorizedMiddle) //OK
	apiGroup.PUT("/forget_password", handler.ForgetPasswordHandler)                               //>Process non finalisé
	apiGroup.PUT("/init_password", handler.InitPasswordHandler, middlewares.IsAdminMiddle, middlewares.GetPermission)
	apiGroup.PUT("/refresh_token", handler.RefreshTokenHandler, middlewares.IsRefreshTokenMiddle) //OK
}
