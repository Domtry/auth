package authentications

import (
	"auth/api/users"
	"auth/service"
	"auth/utils"
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
// @Success 201 {object} utils.HttpResponse[AuthOut]
// @Produce json
// @Router /auth/login [post]
func (h *AuthenticationHandler) LoginHandler(ctx echo.Context) error {

	var payload AuthIn

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

		jsonResponse := utils.HttpResponse[map[string]interface{}]{
			Message:   "Validation failed",
			Success:   false,
			CodeError: http.StatusBadRequest,
			Data: map[string]interface{}{
				"details": validationErrors,
			},
		}

		return ctx.JSON(http.StatusBadRequest, jsonResponse)
	}

	authResponse, err := h.authService.Login(payload.Username, payload.Password)
	if err != nil {
		jsonResponse := utils.HttpResponse[any]{
			Message:   err.Error(),
			Success:   false,
			CodeError: http.StatusBadRequest,
			Data:      nil,
		}
		return ctx.JSON(http.StatusBadRequest, jsonResponse)
	}

	jsonResponse := utils.HttpResponse[AuthOut]{
		Message:   "Connexion succès",
		Success:   true,
		CodeError: http.StatusOK,
		Data: AuthOut{
			Id:    authResponse.UserId,
			Name:  authResponse.Name,
			Email: authResponse.Email,
			Role:  authResponse.Role,
			Token: Token{
				AccessToken:  authResponse.Token.AccessToken,
				RefreshToken: authResponse.Token.RefreshToken,
				ExpiresAt:    authResponse.Token.ExpiresAt,
			},
		},
	}
	return ctx.JSON(http.StatusOK, jsonResponse)
}

// UserProfilHandler Voir le profil d'un utilisateur
// @Summary Voir le profil d'un utilisateur
// @Description Voir le profil d'un utilisateur
// @Tags Authentications
// @Success 201 {object} utils.HttpResponse[users.UserOut]
// @Produce json
// @Router /auth/me [get]
func (h *AuthenticationHandler) UserProfilHandler(ctx echo.Context) error {

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

	userProfil, err := h.authService.UserProfil(fmt.Sprintf("%v", userId))
	if err != nil {
		jsonResponse := utils.HttpResponse[any]{
			Message:   err.Error(),
			Success:   false,
			CodeError: http.StatusBadRequest,
			Data:      nil,
		}
		return ctx.JSON(http.StatusBadRequest, jsonResponse)
	}

	jsonResponse := utils.HttpResponse[users.UserOut]{
		Message:   "Données récupérée avec succès",
		Success:   true,
		CodeError: http.StatusOK,
		Data: users.UserOut{
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

// RefreshTokenHandler Rafraichir le token
// @Summary Rafraichir le token
// @Description Rafraichir le token
// @Tags Authentications
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @Success 202 {object} utils.HttpResponse[RefreshTokenOut]
// @Produce json
// @Router /auth/refresh_token [put]
func (h *AuthenticationHandler) RefreshTokenHandler(ctx echo.Context) error {

	userId := ctx.Get("userId")

	authResponse, err := h.authService.RefreshToken(fmt.Sprintf("%s", userId))
	if err != nil {
		jsonResponse := utils.HttpResponse[any]{
			Message:   err.Error(),
			Success:   false,
			CodeError: http.StatusBadRequest,
			Data:      nil,
		}
		return ctx.JSON(http.StatusBadRequest, jsonResponse)
	}

	jsonResponse := utils.HttpResponse[RefreshTokenOut]{
		Message:   "Token has been refresh",
		Success:   true,
		CodeError: http.StatusAccepted,
		Data: RefreshTokenOut{
			Id: authResponse.UserId,
			Token: Token{
				AccessToken:  authResponse.Token.AccessToken,
				RefreshToken: authResponse.Token.RefreshToken,
				ExpiresAt:    authResponse.Token.ExpiresAt,
			},
		},
	}
	return ctx.JSON(http.StatusAccepted, jsonResponse)
}

// ForgetPasswordHandler Mot de passe oublié
// @Summary Voir le profil d'un utilisateur
// @Description Voir le profil d'un utilisateur
// @Tags Authentications
// @Produce json
// @Router /auth/forget_password [put]
func (h *AuthenticationHandler) ForgetPasswordHandler(ctx echo.Context) error {

	h.authService.ForgetPassword("")
	return nil
}

// ResetPasswordHandler Changer de mot de passe
// @Summary Changer de mot de passe
// @Description Changer de mot de passe
// @Tags Authentications
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @Param reset_password body ResetPasswordIn true "Détails de l'utilisateur"
// @Success 202 {object} utils.HttpResponse[AuthOut]
// @Produce json
// @Router /auth/reset_password [put]
func (h *AuthenticationHandler) ResetPasswordHandler(ctx echo.Context) error {

	userId := ctx.Get("userId")

	var payload ResetPasswordIn

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

		jsonResponse := utils.HttpResponse[map[string]interface{}]{
			Message:   "Validation failed",
			Success:   false,
			CodeError: http.StatusBadRequest,
			Data: map[string]interface{}{
				"details": validationErrors,
			},
		}

		return ctx.JSON(http.StatusBadRequest, jsonResponse)
	}

	if payload.OldPassword == payload.NewPassword {
		jsonResponse := utils.HttpResponse[any]{
			Message:   "Vous devez saisir un mot de passe different de l'actuel",
			Success:   false,
			CodeError: http.StatusBadRequest,
			Data:      nil,
		}
		return ctx.JSON(http.StatusBadRequest, jsonResponse)
	}

	restPasswordResponse, err := h.authService.ChangePassword(fmt.Sprintf("%s", userId), payload.NewPassword)
	if err != nil {
		jsonResponse := utils.HttpResponse[any]{
			Message:   err.Error(),
			Success:   false,
			CodeError: http.StatusBadRequest,
			Data:      nil,
		}
		return ctx.JSON(http.StatusBadRequest, jsonResponse)
	}

	jsonResponse := utils.HttpResponse[AuthOut]{
		Message:   "Mot de passe à bien mise à jour avec succès",
		Success:   true,
		CodeError: http.StatusAccepted,
		Data: AuthOut{
			Id:    restPasswordResponse.UserId,
			Name:  restPasswordResponse.Name,
			Email: restPasswordResponse.Email,
			Role:  restPasswordResponse.Role,
			Token: Token{
				AccessToken:  restPasswordResponse.Token.AccessToken,
				RefreshToken: restPasswordResponse.Token.RefreshToken,
				ExpiresAt:    restPasswordResponse.Token.ExpiresAt,
			},
		},
	}
	return ctx.JSON(http.StatusAccepted, jsonResponse)
}

// InitPasswordHandler Initialisation du mot de passe utilisateur
// @Summary Initialisation du mot de passe utilisateur
// @Description Initialisation du mot de passe utilisateur
// @Tags Authentications
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @Param init_password body InitUserPasswordIn true "Body data"
// @Success 202 {object} utils.HttpResponse[any]
// @Produce json
// @Router /auth/init_password [put]
func (h *AuthenticationHandler) InitPasswordHandler(ctx echo.Context) error {

	var payload InitUserPasswordIn

	err := utils.VerifyPermission(ctx, "init_user_password")
	if err != nil {
		jsonResponse := utils.HttpResponse[any]{
			Message:   "Vous n'avez pas suffisament les droits pour continuer cette opération",
			Success:   false,
			CodeError: http.StatusUnauthorized,
			Data:      nil,
		}
		return ctx.JSON(http.StatusBadRequest, jsonResponse)
	}

	if err = ctx.Bind(&payload); err != nil {
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

		jsonResponse := utils.HttpResponse[map[string]interface{}]{
			Message:   "Validation failed",
			Success:   false,
			CodeError: http.StatusBadRequest,
			Data: map[string]interface{}{
				"details": validationErrors,
			},
		}

		return ctx.JSON(http.StatusBadRequest, jsonResponse)
	}

	err = h.authService.InitPassword(payload.UserId, payload.Password)
	if err != nil {
		jsonResponse := utils.HttpResponse[any]{
			Message:   "Nous avons rencontré un problème durant l'initialisation du mot de passe",
			Success:   false,
			CodeError: http.StatusBadRequest,
			Data:      nil,
		}

		return ctx.JSON(http.StatusBadRequest, jsonResponse)
	}

	jsonResponse := utils.HttpResponse[any]{
		Message:   "Le nouveau mot de passe de l'utilisateur à bien été initialisé",
		Success:   true,
		CodeError: http.StatusAccepted,
		Data:      nil,
	}
	return ctx.JSON(http.StatusAccepted, jsonResponse)
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
