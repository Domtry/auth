package users

import (
	"auth/service"
	"auth/utils"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// UserHandler gère les requêtes liées aux utilisateurs
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler crée une nouvelle instance de UserHandler
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetAllUsersHandler gère la requête pour récupérer tous les utilisateurs
// @Summary Récupère tous les utilisateurs
// @Description Récupère la liste de tous les utilisateurs.
// @Tags Users
// @Produce json
// @Success 200 {array} model.User
// @Router /users [get]
func (h *UserHandler) GetAllUsersHandler(ctx echo.Context) error {

	users, err := h.userService.GetAllUsers()
	if err != nil {
		jsonResponse := utils.HttpResponse[any]{
			Message:   "Aucun utilisateur trouvé",
			Success:   false,
			CodeError: http.StatusBadRequest,
			Data:      nil,
		}
		return ctx.JSON(http.StatusBadRequest, jsonResponse)
	}

	var userList []UserOut

	for i := 0; i < len(users); i++ {
		userItem := users[i]
		currentData := UserOut{
			Id:        userItem.Id,
			Name:      userItem.Name,
			Email:     userItem.Email,
			Role:      userItem.RoleId,
			Sername:   userItem.Sername,
			Username:  userItem.Username,
			CreatedAt: userItem.CreatedAt,
			UpdatedAt: userItem.UpdatedAt,
		}

		userList = append(userList, currentData)
	}

	if len(userList) == 0 {
		userList = []UserOut{}
	}

	jsonResponse := utils.HttpResponse[[]UserOut]{
		Message:   "Données récupérée avec succès",
		Success:   true,
		CodeError: http.StatusOK,
		Data:      userList,
	}
	return ctx.JSON(http.StatusOK, jsonResponse)
}

// GetUserByIdHandler gère la requête pour récupérer un utilisateur par son ID
// @Summary Récupère un utilisateur par ID
// @Description Récupère un utilisateur en fonction de son ID.
// @Tags Users
// @Produce json
// @Param id path string true "ID de l'utilisateur"
// @Success 200 {object} model.User
// @Router /users/{id} [get]
func (h *UserHandler) GetUserByIdHandler(ctx echo.Context) error {

	userId := ctx.Param("id")
	if userId == "" {
		jsonResponse := utils.HttpResponse[any]{
			Message:   "Invalid user Id",
			Success:   false,
			CodeError: http.StatusBadRequest,
			Data:      nil,
		}
		return ctx.JSON(http.StatusBadRequest, jsonResponse)
	}

	userProfil, err := h.userService.GetUserById(fmt.Sprintf("%s", userId))
	if err != nil {
		jsonResponse := utils.HttpResponse[any]{
			Message:   "Ce compte n'existe dans le système",
			Success:   false,
			CodeError: http.StatusBadRequest,
			Data:      nil,
		}
		return ctx.JSON(http.StatusBadRequest, jsonResponse)
	}

	jsonResponse := utils.HttpResponse[UserOut]{
		Message:   "Données récupérée avec succès",
		Success:   true,
		CodeError: http.StatusOK,
		Data: UserOut{
			Id:        userProfil.Id,
			Name:      userProfil.Name,
			Email:     userProfil.Email,
			Role:      userProfil.RoleId,
			Sername:   userProfil.Sername,
			Username:  userProfil.Username,
			CreatedAt: userProfil.CreatedAt,
			UpdatedAt: userProfil.UpdatedAt,
		},
	}
	return ctx.JSON(http.StatusOK, jsonResponse)
}

// GetUserProfileHandler View user profile info
// @Summary User profile
// @Description view user profile info.
// @Tags Users
// @Produce json
// @Success 200 {object} utils.HttpResponse[users.UserOut]
// @Router /users/profile [get]
func (h *UserHandler) GetUserProfileHandler(ctx echo.Context) error {

	userId := ctx.Get("userId")

	err := utils.VerifyPermission(ctx, "view_user_profile")
	if err != nil {
		jsonResponse := utils.HttpResponse[any]{
			Message:   "Vous n'avez pas suffisament les droits pour continuer cette opération",
			Success:   false,
			CodeError: http.StatusUnauthorized,
			Data:      nil,
		}
		return ctx.JSON(http.StatusBadRequest, jsonResponse)
	}

	userProfil, err := h.userService.GetUserById(fmt.Sprintf("%s", userId))
	if err != nil {
		jsonResponse := utils.HttpResponse[any]{
			Message:   "Ce compte n'existe dans le système",
			Success:   false,
			CodeError: http.StatusBadRequest,
			Data:      nil,
		}
		return ctx.JSON(http.StatusBadRequest, jsonResponse)
	}

	jsonResponse := utils.HttpResponse[UserOut]{
		Message:   "Données récupérée avec succès",
		Success:   true,
		CodeError: http.StatusOK,
		Data: UserOut{
			Id:        userProfil.Id,
			Name:      userProfil.Name,
			Email:     userProfil.Email,
			Role:      userProfil.RoleId,
			Sername:   userProfil.Sername,
			Username:  userProfil.Username,
			CreatedAt: userProfil.CreatedAt,
			UpdatedAt: userProfil.UpdatedAt,
		},
	}
	return ctx.JSON(http.StatusOK, jsonResponse)
}

// CreateUserHandler gère la requête pour créer un nouvel utilisateur
// @Summary Crée un nouvel utilisateur
// @Description Crée un nouvel utilisateur avec les détails fournis.
// @Tags Users
// @Accept json
// @Produce json
// @Param user body UserIn true "Body data"
// @Success 201 {object} utils.HttpResponse[UserOut]
// @Router /users [post]
func (h *UserHandler) CreateUserHandler(ctx echo.Context) error {

	var payload UserIn

	/** Verify user permission
	err := utils.VerifyPermission(ctx, "create_user")
	if err != nil {
		jsonResponse := utils.HttpResponse[any]{
			Message:   "Vous n'avez pas suffisament les droits pour continuer cette opération",
			Success:   false,
			CodeError: http.StatusUnauthorized,
			Data:      nil,
		}
		return ctx.JSON(http.StatusBadRequest, jsonResponse)
	}
	*/

	if err := ctx.Bind(&payload); err != nil {
		jsonResponse := utils.HttpResponse[any]{
			Message:   "Données JSON invalides",
			Success:   false,
			CodeError: http.StatusBadRequest,
			Data:      nil,
		}
		return ctx.JSON(http.StatusBadRequest, jsonResponse)
	}

	// Validation des données
	validate := validator.New()
	if err := validate.Struct(payload); err != nil {
		var validationErrors []string

		// Parcourez les erreurs de validation
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, formatValidationError(err))
		}

		jsonResponse := utils.HttpResponse[any]{
			Message:   "Données JSON invalides",
			Success:   false,
			CodeError: http.StatusBadRequest,
			Data:      validationErrors,
		}
		return ctx.JSON(http.StatusBadRequest, jsonResponse)
	}

	newUser := service.User(payload)

	createdUser, err := h.userService.CreateUser(newUser)
	if err != nil {
		jsonResponse := utils.HttpResponse[any]{
			Message:   err.Error(),
			Success:   false,
			CodeError: http.StatusUnprocessableEntity,
			Data:      nil,
		}
		return ctx.JSON(http.StatusUnprocessableEntity, jsonResponse)
	}

	jsonResponse := utils.HttpResponse[UserOut]{
		Message:   "User has been created",
		Success:   true,
		CodeError: http.StatusCreated,
		Data: UserOut{
			Id:        createdUser.Id,
			Username:  createdUser.Username,
			Email:     createdUser.Email,
			Role:      createdUser.RoleId,
			Sername:   createdUser.Sername,
			Name:      createdUser.Name,
			CreatedAt: createdUser.CreatedAt,
		},
	}

	return ctx.JSON(http.StatusCreated, jsonResponse)
}

// UpdateUserHandler gère la requête pour mettre à jour un utilisateur
// @Summary Met à jour un utilisateur
// @Description Met à jour un utilisateur en fonction de son ID avec les détails fournis.
// @Tags Users
// @Accept json
// @Produce json
// @Param user body UserIn true "Détails de l'utilisateur à mettre à jour"
// @Success 200 {object} utils.HttpResponse[UserOut]
// @Router /users [put]
func (h *UserHandler) UpdateUserHandler(ctx echo.Context) error {

	userId := ctx.Get("userId")

	err := utils.VerifyPermission(ctx, "update_user_profile")
	if err != nil {
		jsonResponse := utils.HttpResponse[any]{
			Message:   "Vous n'avez pas suffisament les droits pour continuer cette opération",
			Success:   false,
			CodeError: http.StatusUnauthorized,
			Data:      nil,
		}
		return ctx.JSON(http.StatusBadRequest, jsonResponse)
	}

	var payload UpdateUserIn
	if err := ctx.Bind(&payload); err != nil {
		jsonResponse := utils.HttpResponse[any]{
			Message:   "Données JSON invalides",
			Success:   false,
			CodeError: http.StatusBadRequest,
			Data:      nil,
		}
		return ctx.JSON(http.StatusBadRequest, jsonResponse)
	}

	// Validation des données
	validate := validator.New()
	if err := validate.Struct(payload); err != nil {
		var validationErrors []string

		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, formatValidationError(err))
		}

		jsonResponse := utils.HttpResponse[any]{
			Message:   "Données JSON invalides",
			Success:   false,
			CodeError: http.StatusBadRequest,
			Data:      validationErrors,
		}
		return ctx.JSON(http.StatusBadRequest, jsonResponse)
	}

	existingUser, err := h.userService.GetUserById(fmt.Sprintf("%s", userId))
	if err != nil {
		jsonResponse := utils.HttpResponse[any]{
			Message:   "Ce compte n'existe dans le système",
			Success:   false,
			CodeError: http.StatusBadRequest,
			Data:      nil,
		}
		return ctx.JSON(http.StatusBadRequest, jsonResponse)
	}

	if payload.Name != existingUser.Name {
		existingUser.Name = payload.Name
	}

	if payload.Sername != existingUser.Sername {
		existingUser.Sername = payload.Sername
	}

	if payload.Email != existingUser.Email {
		_, err = h.userService.GetUserByEmail(payload.Email)
		if err == nil {
			jsonResponse := utils.HttpResponse[any]{
				Message:   "Ce email a déjà été utilisé",
				Success:   false,
				CodeError: http.StatusBadRequest,
				Data:      nil,
			}
			return ctx.JSON(http.StatusBadRequest, jsonResponse)
		}

		existingUser.Name = payload.Name
	}

	newUser := service.User{
		Name:     existingUser.Name,
		Email:    existingUser.Email,
		Sername:  existingUser.Sername,
		Username: existingUser.Username,
		Password: existingUser.Password,
	}

	updateUserResponse, err := h.userService.UpdateUser(newUser)
	if err != nil {
		jsonResponse := utils.HttpResponse[any]{
			Message:   err.Error(),
			Success:   false,
			CodeError: http.StatusUnprocessableEntity,
			Data:      nil,
		}
		return ctx.JSON(http.StatusUnprocessableEntity, jsonResponse)
	}

	jsonResponse := utils.HttpResponse[UserOut]{
		Message:   "Compte utilisateur mise à jour avec succès",
		Success:   true,
		CodeError: http.StatusAccepted,
		Data: UserOut{
			Id:        updateUserResponse.Id,
			Name:      updateUserResponse.Name,
			Email:     updateUserResponse.Email,
			Sername:   updateUserResponse.Sername,
			Username:  updateUserResponse.Username,
			CreatedAt: updateUserResponse.CreatedAt,
			UpdatedAt: updateUserResponse.UpdatedAt,
			Role:      updateUserResponse.RoleId,
		},
	}
	return ctx.JSON(http.StatusAccepted, jsonResponse)
}

// UpdateUserByIdHandler gère la requête pour mettre à jour un utilisateur
// @Summary Met à jour un utilisateur
// @Description Met à jour un utilisateur en fonction de son ID avec les détails fournis.
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "ID de l'utilisateur"
// @Param user body UserIn true "Body data"
// @Success 200 {object} utils.HttpResponse[UserOut]
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUserByIdHandler(ctx echo.Context) error {

	userId := ctx.Param("id")

	err := utils.VerifyPermission(ctx, "update_user_profile")
	if err != nil {
		jsonResponse := utils.HttpResponse[any]{
			Message:   "Vous n'avez pas suffisament les droits pour continuer cette opération",
			Success:   false,
			CodeError: http.StatusUnauthorized,
			Data:      nil,
		}
		return ctx.JSON(http.StatusBadRequest, jsonResponse)
	}

	if userId == "" {
		jsonResponse := utils.HttpResponse[any]{
			Message:   "Invalid user Id",
			Success:   false,
			CodeError: http.StatusBadRequest,
			Data:      nil,
		}
		return ctx.JSON(http.StatusBadRequest, jsonResponse)
	}

	var payload UpdateUserIn
	if err := ctx.Bind(&payload); err != nil {
		jsonResponse := utils.HttpResponse[any]{
			Message:   "Données JSON invalides",
			Success:   false,
			CodeError: http.StatusBadRequest,
			Data:      nil,
		}
		return ctx.JSON(http.StatusBadRequest, jsonResponse)
	}

	// Validation des données
	validate := validator.New()
	if err := validate.Struct(payload); err != nil {
		var validationErrors []string

		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, formatValidationError(err))
		}

		jsonResponse := utils.HttpResponse[any]{
			Message:   "Données JSON invalides",
			Success:   false,
			CodeError: http.StatusBadRequest,
			Data:      validationErrors,
		}
		return ctx.JSON(http.StatusBadRequest, jsonResponse)
	}

	existingUser, err := h.userService.GetUserById(fmt.Sprintf("%s", userId))
	if err != nil {
		jsonResponse := utils.HttpResponse[any]{
			Message:   "Ce compte n'existe dans le système",
			Success:   false,
			CodeError: http.StatusBadRequest,
			Data:      nil,
		}
		return ctx.JSON(http.StatusBadRequest, jsonResponse)
	}

	if payload.Name != existingUser.Name {
		existingUser.Name = payload.Name
	}

	if payload.Sername != existingUser.Sername {
		existingUser.Sername = payload.Sername
	}

	if payload.Email != existingUser.Email {
		_, err = h.userService.GetUserByEmail(payload.Email)
		if err == nil {
			jsonResponse := utils.HttpResponse[any]{
				Message:   "Ce email a déjà été utilisé",
				Success:   false,
				CodeError: http.StatusBadRequest,
				Data:      nil,
			}
			return ctx.JSON(http.StatusBadRequest, jsonResponse)
		}

		existingUser.Email = payload.Email
	}

	newUser := service.User{
		Name:     existingUser.Name,
		Email:    existingUser.Email,
		Sername:  existingUser.Sername,
		Username: existingUser.Username,
		Password: existingUser.Password,
	}

	updateUserResponse, err := h.userService.UpdateUser(newUser)
	if err != nil {
		jsonResponse := utils.HttpResponse[any]{
			Message:   err.Error(),
			Success:   false,
			CodeError: http.StatusUnprocessableEntity,
			Data:      nil,
		}
		return ctx.JSON(http.StatusUnprocessableEntity, jsonResponse)
	}

	jsonResponse := utils.HttpResponse[UserOut]{
		Message:   "Compte utilisateur mise à jour avec succès",
		Success:   true,
		CodeError: http.StatusAccepted,
		Data: UserOut{
			Id:        updateUserResponse.Id,
			Name:      updateUserResponse.Name,
			Email:     updateUserResponse.Email,
			Sername:   updateUserResponse.Sername,
			Username:  updateUserResponse.Username,
			CreatedAt: updateUserResponse.CreatedAt,
			UpdatedAt: updateUserResponse.UpdatedAt,
			Role:      updateUserResponse.RoleId,
		},
	}
	return ctx.JSON(http.StatusAccepted, jsonResponse)
}

// DeleteUserHandler gère la requête pour supprimer un utilisateur
// @Summary Supprime un utilisateur
// @Description Supprime un utilisateur en fonction de son ID.
// @Tags Users
// @Param id path int true "ID de l'utilisateur"
// @Success 204 "Aucun contenu"
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUserHandler(ctx echo.Context) error {

	userId := ctx.Param("id")

	err := utils.VerifyPermission(ctx, "delete_user")
	if err != nil {
		jsonResponse := utils.HttpResponse[any]{
			Message:   "Vous n'avez pas suffisament les droits pour continuer cette opération",
			Success:   false,
			CodeError: http.StatusUnauthorized,
			Data:      nil,
		}
		return ctx.JSON(http.StatusBadRequest, jsonResponse)
	}

	if userId == "" {
		jsonResponse := utils.HttpResponse[any]{
			Message:   "Invalid user Id",
			Success:   false,
			CodeError: http.StatusBadRequest,
			Data:      nil,
		}
		return ctx.JSON(http.StatusBadRequest, jsonResponse)
	}

	err = h.userService.DeleteUser(userId)
	if err != nil {
		jsonResponse := utils.HttpResponse[any]{
			Message:   "Invalid user Id",
			Success:   false,
			CodeError: http.StatusBadRequest,
			Data:      nil,
		}
		return ctx.JSON(http.StatusBadRequest, jsonResponse)
	}

	jsonResponse := utils.HttpResponse[any]{
		Message:   "User has been deleted",
		Success:   true,
		CodeError: http.StatusAccepted,
		Data:      nil,
	}
	return ctx.JSON(http.StatusAccepted, jsonResponse)
}

// AssignRoleHandler gère la requête pour assigner un rôle à un utilisateur
// @Summary Assigner un rôle à un utilisateur
// @Description Assigner un rôle spécifié à un utilisateur en fonction de son ID.
// @Tags Users
// @Param id path int true "ID de l'utilisateur"
// @Param role path string true "Rôle à assigner"
// @Success 204 "Aucun contenu"
// @Router /users/{id}/assign-role/{role} [post]
func (h *UserHandler) AssignRoleHandler(c echo.Context) error {
	userId := c.Param("id")
	if userId != "" {
		return c.JSON(http.StatusBadRequest, "ID invalide")
	}

	var roleRequest struct {
		Role string `json:"role"`
	}

	if err := c.Bind(&roleRequest); err != nil {
		return c.JSON(http.StatusBadRequest, "Données JSON invalides")
	}

	updatedUser, err := h.userService.AssignRole(userId, roleRequest.Role)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, updatedUser)
}

// RemoveRoleHandler gère la requête pour supprimer un rôle d'un utilisateur
// @Summary Supprimer un rôle d'un utilisateur
// @Description Supprimer un rôle spécifié d'un utilisateur en fonction de son ID.
// @Tags Users
// @Param id path int true "ID de l'utilisateur"
// @Param role path string true "Rôle à supprimer"
// @Success 204 "Aucun contenu"
// @Router /users/{id}/remove-role/{role} [post]
func (h *UserHandler) RemoveRoleHandler(c echo.Context) error {
	userId := c.Param("id")
	if userId != "" {
		return c.JSON(http.StatusBadRequest, "ID invalide")
	}

	updatedUser, err := h.userService.RemoveRole(userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, updatedUser)
}

// Fonction pour formater les erreurs de validation
func formatValidationError(err validator.FieldError) string {
	fieldName := strings.ToLower(err.Field())
	switch err.Tag() {
	case "required":
		return fieldName + " is required"
	case "min":
		return fieldName + " must be at least " + err.Param()
	case "max":
		return fieldName + " must be at most " + err.Param()
	case "email":
		return fieldName + " must be a valid email address"
	case "uuid4":
		return fieldName + " must be a valid UUID (version 4)"
	default:
		return fieldName + " is invalid"
	}
}
