package users

import (
	"auth/middlewares"

	"github.com/labstack/echo/v4"
)

func RegisterUserRoutes(apiGroup *echo.Group, handler *UserHandler) {
	apiGroup.GET("", handler.GetAllUsersHandler, middlewares.IsAdminMiddle)
	apiGroup.POST("", handler.CreateUserHandler, middlewares.IsAdminMiddle)
	apiGroup.PUT("/:id", handler.UpdateUserHandler, middlewares.IsAdminMiddle)
	apiGroup.GET("/:id", handler.GetUserByIdHandler, middlewares.IsAdminMiddle)
	apiGroup.DELETE("/:id", handler.DeleteUserHandler, middlewares.IsAdminMiddle)
	apiGroup.PUT("/:id/remove-role", handler.RemoveRoleHandler, middlewares.IsAdminMiddle)
	apiGroup.PUT("/:id/assign-role", handler.AssignRoleHandler, middlewares.IsAdminMiddle)
}
