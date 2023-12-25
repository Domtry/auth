package middlewares

import (
	"auth/utils"
	"github.com/labstack/echo/v4"
	"net/http"
)

func IsAuthorizedMiddle(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		token := ctx.Request().Header.Get("Authorization")
		extractToken := utils.ExtractToken(token)
		if len(extractToken) == 0 {
			return ctx.JSON(http.StatusUnauthorized, map[string]interface{}{
				"message":    "Vous n'etes pas autorisé à executer cette route",
				"success":    false,
				"code_error": http.StatusUnauthorized,
				"data":       nil,
			})
		}

		claims, err := utils.ParseToken(extractToken)
		if err != nil || claims.Source == "refresh_token" {
			return ctx.JSON(http.StatusUnauthorized, map[string]interface{}{
				"message":    "Token non valid",
				"success":    false,
				"code_error": http.StatusUnauthorized,
				"data":       nil,
			})
		}

		ctx.Set("userId", claims.Id)

		return next(ctx)
	}
}

func IsAdminMiddle(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		token := ctx.Request().Header.Get("Authorization")
		extractToken := utils.ExtractToken(token)
		if len(extractToken) == 0 {
			return ctx.JSON(http.StatusUnauthorized, map[string]interface{}{
				"message":    "Vous n'etes pas autorisé à executer cette route",
				"success":    false,
				"code_error": http.StatusUnauthorized,
				"data":       nil,
			})
		}

		claims, err := utils.ParseToken(extractToken)
		if err != nil || claims.Source == "refresh_token" {
			return ctx.JSON(http.StatusUnauthorized, map[string]interface{}{
				"message":    "Token non valid",
				"success":    false,
				"code_error": http.StatusUnauthorized,
				"data":       nil,
			})
		}

		if claims.Role != "admin" {
			return ctx.JSON(http.StatusUnauthorized, map[string]interface{}{
				"message":    "Vous devez être un admin pour executer cette methode",
				"success":    false,
				"code_error": http.StatusUnauthorized,
				"data":       nil,
			})
		}

		ctx.Set("userId", claims.Id)

		return next(ctx)
	}
}
