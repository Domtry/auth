package users

import (
	"auth/model"
	"auth/service"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

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
func (h *UserHandler) GetAllUsersHandler(c echo.Context) error {

	users, err := h.userService.GetAllUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, users)
}

// GetUserByIdHandler gère la requête pour récupérer un utilisateur par son ID
// @Summary Récupère un utilisateur par ID
// @Description Récupère un utilisateur en fonction de son ID.
// @Tags Users
// @Produce json
// @Param id path string true "ID de l'utilisateur"
// @Success 200 {object} model.User
// @Router /users/{id} [get]
func (h *UserHandler) GetUserByIdHandler(c echo.Context) error {

	userId := c.Param("id")
	if userId == "" {
		return c.JSON(http.StatusBadRequest, "ID invalid")
	}

	if _, err := uuid.Parse(userId); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	user, err := h.userService.GetUserById(userId)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, user)
}

// CreateUserHandler gère la requête pour créer un nouvel utilisateur
// @Summary Crée un nouvel utilisateur
// @Description Crée un nouvel utilisateur avec les détails fournis.
// @Tags Users
// @Accept json
// @Produce json
// @Param user body model.User true "Détails de l'utilisateur"
// @Success 201 {object} model.User
// @Router /users [post]
func (h *UserHandler) CreateUserHandler(c echo.Context) error {

	var newUser model.User
	if err := c.Bind(&newUser); err != nil {
		return c.JSON(http.StatusBadRequest, "Données JSON invalides")
	}

	// Validation des données
	validate := validator.New()
	if err := validate.Struct(newUser); err != nil {
		var validationErrors []string

		// Parcourez les erreurs de validation
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, formatValidationError(err))
		}

		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Validation failed",
			"details": validationErrors,
		})
	}

	createdUser, err := h.userService.CreateUser(newUser)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, createdUser)
}

// UpdateUserHandler gère la requête pour mettre à jour un utilisateur
// @Summary Met à jour un utilisateur
// @Description Met à jour un utilisateur en fonction de son ID avec les détails fournis.
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "ID de l'utilisateur"
// @Param user body model.User true "Détails de l'utilisateur à mettre à jour"
// @Success 200 {object} model.User
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUserHandler(c echo.Context) error {
	userId := c.Param("id")
	if userId != "" {
		return c.JSON(http.StatusBadRequest, "ID invalide")
	}

	if _, err := uuid.Parse(userId); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var updatedUser model.User
	if err := c.Bind(&updatedUser); err != nil {
		return c.JSON(http.StatusBadRequest, "Données JSON invalides")
	}

	updatedUser.Id = string(userId)

	updatedUser, err := h.userService.UpdateUser(updatedUser)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, updatedUser)
}

// DeleteUserHandler gère la requête pour supprimer un utilisateur
// @Summary Supprime un utilisateur
// @Description Supprime un utilisateur en fonction de son ID.
// @Tags Users
// @Param id path int true "ID de l'utilisateur"
// @Success 204 "Aucun contenu"
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUserHandler(c echo.Context) error {
	userId := c.Param("id")
	if userId != "" {
		return c.JSON(http.StatusBadRequest, "ID invalide")
	}

	err := h.userService.DeleteUser(userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
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
