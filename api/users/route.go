package users

import "github.com/labstack/echo/v4"

func RegisterUserRoutes(apiGroup *echo.Group, handler *UserHandler) {
	apiGroup.GET("", handler.GetAllUsersHandler)
	apiGroup.POST("", handler.CreateUserHandler)
	apiGroup.PUT("/:id", handler.UpdateUserHandler)
	apiGroup.GET("/:id", handler.GetUserByIdHandler)
	apiGroup.DELETE("/:id", handler.DeleteUserHandler)
	apiGroup.PUT("/:id/remove-role", handler.RemoveRoleHandler)
	apiGroup.PUT("/:id/assign-role", handler.AssignRoleHandler)
}
