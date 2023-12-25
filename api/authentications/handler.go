package authentications

import (
	"auth/service"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

// AuthenticationHandler gère les requêtes liées aux utilisateurs
type AuthenticationHandler struct {
	authService *service.AuthenticationService
}

// NewAuthenticationHandler crée une nouvelle instance de UserHandler
func NewAuthenticationHandler(authService *service.AuthenticationService) *AuthenticationHandler {
	return &AuthenticationHandler{
		authService: authService,
	}
}

// LoginHandler authentifier un utlisateur
// @Summary Authentifier un utlisateur
// @Description Authentifier un utlisateur
// @Tags Authentications
// @Param user body AuthIn true "Détails de l'utilisateur"
// @Success 201 {object} AuthOut
// @Produce json
// @Router /auth/login [post]
func (h *AuthenticationHandler) LoginHandler(ctx echo.Context) error {

	var payload AuthIn

	if err := ctx.Bind(&payload); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Données JSON invalides")
	}

	// Validation des données
	validate := validator.New()
	if err := validate.Struct(payload); err != nil {
		var validationErrors []string

		// Parcourez les erreurs de validation
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, formatValidationError(err))
		}

		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"message":    "Validation failed",
			"success":    false,
			"code_error": http.StatusBadRequest,
			"data": map[string]interface{}{
				"details": validationErrors,
			},
		})
	}

	authResponse, err := h.authService.Login(payload.Username, payload.Password)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message":    err.Error(),
			"success":    false,
			"code_error": http.StatusBadRequest,
			"data":       nil,
		})
	}

	return ctx.JSON(http.StatusOK, authResponse)
}

// UserProfil Voir le profil d'un utilisateur
// @Summary Voir le profil d'un utilisateur
// @Description Voir le profil d'un utilisateur
// @Tags Authentications
// @Param user body model.User true "Détails de l'utilisateur"
// @Success 201 {object} model.User
// @Produce json
// @Router /auth/me [get]
func (h *AuthenticationHandler) UserProfil(ctx echo.Context) error {
	userId := ctx.Get("userId")

	userProfil, err := h.authService.UserProfil(fmt.Sprintf("%v", userId))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"message":    err.Error(),
			"success":    false,
			"code_error": http.StatusBadRequest,
			"data":       nil,
		})
	}

	return ctx.JSON(http.StatusOK, userProfil)
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
