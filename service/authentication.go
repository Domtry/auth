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
	Token     string
	ExpiredAt time.Time
}

type CustomResponse struct {
	Data interface{}
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
		Role:   userRole.Name,
		Token: model.Token{
			AccessToken:  accessToken.Token,
			RefreshToken: refreshToken.Token,
			ExpiresAt:    accessToken.ExpiredAt,
		},
	}
	return authModel, nil
}

func (a *AuthenticationService) RefreshToken(userId string) (model.Authentication, error) {

	var existingUser model.User
	existingUser, err := a.userRepo.GetUserById(userId)
	if err != nil {
		return model.Authentication{}, fmt.Errorf("user is not exist")
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
		Role:   userRole.Name,
		Token: model.Token{
			AccessToken:  accessToken.Token,
			RefreshToken: refreshToken.Token,
			ExpiresAt:    accessToken.ExpiredAt,
		},
	}

	return authModel, nil
}

func (a *AuthenticationService) Logout() {
}

func (a *AuthenticationService) ForgetPassword(email string) {
	searchResponse, err := a.userRepo.GetUserByEmail(email)
	if err != nil {
		return
	}

	fmt.Println("searchResponse :::> ", searchResponse.Email)

}

func (a *AuthenticationService) ChangePassword(userId string, password string) (model.Authentication, error) {

	var existingUser model.User

	existingUser, err := a.userRepo.GetUserById(userId)
	if err != nil {
		return model.Authentication{}, fmt.Errorf("user is not exist")
	}

	isDuplicatePassword := utils.CompareHashPassword(password, existingUser.Password)
	if isDuplicatePassword {
		return model.Authentication{}, fmt.Errorf("vous ne pouvez pas utiliser le même mot de passe")
	}

	newPassword, _ := utils.GenerateHashPassword(password)
	existingUser.Password = newPassword
	existingUser.UpdatedAt = time.Now()

	updateUserRespone, err := a.userRepo.UpdateUser(existingUser)
	if err != nil {
		return model.Authentication{}, fmt.Errorf(
			"nous avons rencontré un problème durant la mise à du mot de passe. merci de réessayer")
	}

	userRole, err := a.roleRepo.GetRoleById(updateUserRespone.RoleId)
	if err != nil {
		return model.Authentication{}, fmt.Errorf(
			"error system : user role has not found please call admin system to resolve this problem")
	}

	accessToken, _ := a.generateToken("access_token", userRole.Name, updateUserRespone)
	refreshToken, _ := a.generateToken("refresh_token", userRole.Name, updateUserRespone)

	authModel := model.Authentication{
		UserId: updateUserRespone.Id,
		Email:  updateUserRespone.Email,
		Name:   updateUserRespone.Name,
		Role:   updateUserRespone.Name,
		Token: model.Token{
			AccessToken:  accessToken.Token,
			RefreshToken: refreshToken.Token,
			ExpiresAt:    accessToken.ExpiredAt,
		},
	}
	return authModel, nil
}

func (a *AuthenticationService) UserProfil(userId string) (model.User, error) {
	response, err := a.userRepo.GetUserById(userId)
	if err != nil {
		return model.User{}, err
	}

	userRole, err := a.roleRepo.GetRoleById(response.RoleId)
	if err != nil {
		return model.User{}, err
	}
	response.RoleId = userRole.Name
	return response, nil
}

func (a *AuthenticationService) InitPassword(userId string, password string) error {
	currentUser, err := a.userRepo.GetUserById(userId)
	if err != nil {
		return fmt.Errorf("utilisateur inexistant")
	}

	newPassword, _ := utils.GenerateHashPassword(password)
	currentUser.Password = newPassword
	currentUser.UpdatedAt = time.Now()

	_, err = a.userRepo.UpdateUser(currentUser)
	if err != nil {
		return fmt.Errorf("une erreur est survenue durant la mise à des données client")
	}

	//Send Email or SMS to submit new password
	return nil
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
		Token:     tokenString,
		ExpiredAt: expirationTime,
	}, nil
}
