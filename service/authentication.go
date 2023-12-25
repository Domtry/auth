package service

import (
	"auth/model"
	"auth/repository"
	"auth/utils"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var jwtKey = []byte("my_secret_key")

type TokenResponse struct {
	token     string
	expiredAt time.Time
}

// AuthenticationService gère la logique métier liée aux utilisateurs
type AuthenticationService struct {
	userRepo *repository.UserRepository
	roleRepo *repository.RoleRepository
}

// NewAuthenticationService create new AuthenticationService instance
func NewAuthenticationService(
	userRepo *repository.UserRepository, roleRepo *repository.RoleRepository) *AuthenticationService {
	return &AuthenticationService{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

// Login permit to authenticated user
func (a *AuthenticationService) Login(username string, password string) (model.Authentication, error) {

	var existingUser model.User
	existingUser, err := a.userRepo.GetUserByUsername(username)
	if err != nil {
		return model.Authentication{}, fmt.Errorf("username or passwword is invalid")
	}

	errHash := utils.CompareHashPassword(password, existingUser.Password)
	if !errHash {
		return model.Authentication{}, fmt.Errorf("username or passwword is invalid")
	}

	userRole, err := a.roleRepo.GetRoleById(existingUser.RoleId)
	if err != nil {
		return model.Authentication{}, fmt.Errorf(
			"error system : user role has not found please call admin system to resolve this problem")
	}

	accessToken, _ := a.generateToken("access_token", userRole.Name, existingUser)
	refreshToken, _ := a.generateToken("refresh_token", userRole.Name, existingUser)

	authModel := model.Authentication{
		UserId: existingUser.Id,
		Email:  existingUser.Email,
		Name:   existingUser.Name,
		Token: model.Token{
			AccessToken:  accessToken.token,
			RefreshToken: refreshToken.token,
			ExpiresAt:    accessToken.expiredAt,
		},
	}
	return authModel, nil
}

func (a *AuthenticationService) Logout() {
}

func (a *AuthenticationService) UserProfil(userId string) (model.User, error) {
	response, err := a.userRepo.GetUserById(userId)
	if err != nil {
		return model.User{}, err
	}

	return response, nil
}

func (a *AuthenticationService) generateToken(source string, roleName string, user model.User) (TokenResponse, error) {

	var expirationTime time.Time

	switch source {
	case "refresh_token":
		expirationTime = time.Now().Add(6 * time.Hour)
	case "access_token":
		expirationTime = time.Now().Add(1 * time.Hour)
	}

	claimsAccessToken := model.Claims{
		Role:   roleName,
		Source: source,
		StandardClaims: jwt.StandardClaims{
			Id:        user.Id,
			Subject:   user.Email,
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsAccessToken)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return TokenResponse{}, fmt.Errorf(
			"error system : failled to signed token")
	}
	return TokenResponse{
		token:     tokenString,
		expiredAt: expirationTime,
	}, nil
}
