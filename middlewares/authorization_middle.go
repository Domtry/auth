package middlewares

import (
	"auth/config"
	"auth/utils"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
)

func getDBInstance() *gorm.DB {
	dbConfig, err := config.LoadDBonfig()
	if err != nil {
		panic("File not found")
	}

	db, err := config.GetDB(dbConfig)

	return db
}

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
				"message":    "Access token has not valid",
				"success":    false,
				"code_error": http.StatusUnauthorized,
				"data":       nil,
			})
		}

		ctx.Set("userId", claims.Id)

		return next(ctx)
	}
}

func IsRefreshTokenMiddle(next echo.HandlerFunc) echo.HandlerFunc {
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
		if err != nil || claims.Source == "access_token" {
			return ctx.JSON(http.StatusUnauthorized, map[string]interface{}{
				"message":    "Refresh token has not valid",
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

func GetPermission(next echo.HandlerFunc) echo.HandlerFunc {
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

		claims, _ := utils.ParseToken(extractToken)
		roles := utils.LoadPermissionByRoleName(claims.Role)

		ctx.Set("roles", roles)

		return next(ctx)
	}
}
