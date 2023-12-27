package users

import (
	"auth/middlewares"

	"github.com/labstack/echo/v4"
)

func RegisterUserRoutes(apiGroup *echo.Group, handler *UserHandler) {
	apiGroup.POST("", handler.CreateUserHandler)                                                                       //OK
	apiGroup.PUT("", handler.UpdateUserHandler, middlewares.IsAuthorizedMiddle, middlewares.GetPermission)             //OK
	apiGroup.GET("/profile", handler.GetUserProfileHandler, middlewares.IsAuthorizedMiddle, middlewares.GetPermission) //OK

	//Admin method
	apiGroup.GET("", handler.GetAllUsersHandler, middlewares.IsAdminMiddle, middlewares.GetPermission)        //OK
	apiGroup.GET("/:id", handler.GetUserByIdHandler, middlewares.IsAdminMiddle, middlewares.GetPermission)    //OK
	apiGroup.PUT("/:id", handler.UpdateUserByIdHandler, middlewares.IsAdminMiddle, middlewares.GetPermission) //OK

	apiGroup.DELETE("/:id", handler.DeleteUserHandler, middlewares.IsAdminMiddle, middlewares.GetPermission) //OK
	apiGroup.PUT("/:id/remove-role", handler.RemoveRoleHandler, middlewares.IsAdminMiddle, middlewares.GetPermission)
	apiGroup.PUT("/:id/assign-role", handler.AssignRoleHandler, middlewares.IsAdminMiddle, middlewares.GetPermission)
}
