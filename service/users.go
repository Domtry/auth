// service/user_service.go

package service

import (
	"auth/model"
	"auth/repository"
	"auth/utils"
	"errors"
	"time"
)

type User struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required"`
	Sername  string `json:"sername" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// UserService gère la logique métier liée aux utilisateurs
type UserService struct {
	userRepo *repository.UserRepository
	roleRepo *repository.RoleRepository
}

// NewUserService crée une nouvelle instance de UserService
func NewUserService(
	userRepo *repository.UserRepository, roleRepo *repository.RoleRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

// GetAllUsers récupère tous les utilisateurs
func (s *UserService) GetAllUsers() ([]model.User, error) {
	users, err := s.userRepo.GetAllUsers()
	if err != nil {
		return []model.User{}, err
	}

	for i := 0; i < len(users); i++ {
		userItem := users[i]
		roleResponse, _ := s.roleRepo.GetRoleById(userItem.RoleId)

		userItem.RoleId = roleResponse.Name
		users[i] = userItem
	}

	return users, nil
}

// GetUserById récupère un utilisateur par son ID
func (s *UserService) GetUserById(id string) (model.User, error) {

	userResponse, err := s.userRepo.GetUserById(id)
	if err != nil {
		return model.User{}, err
	}

	roleResponse, _ := s.roleRepo.GetRoleById(userResponse.RoleId)

	userResponse.RoleId = roleResponse.Name

	return userResponse, nil
}

// CreateUser crée un nouvel utilisateur avec un rôle spécifié
func (s *UserService) CreateUser(newUser User) (model.User, error) {

	_, err := s.userRepo.GetUserByEmail(newUser.Email)
	if err == nil {
		return model.User{}, errors.New("adresse email a déjà été utilisé")
	}

	_, err = s.userRepo.GetUserByUsername(newUser.Username)
	if err == nil {
		return model.User{}, errors.New("le nom d'utilisateur a déjà été utilisé")
	}

	roleResponse, err := s.roleRepo.GetRoleByName("client")
	if err != nil {
		return model.User{}, errors.New("aucun rôle ne correspond à cette valeur")
	}

	hashPassword, err := utils.GenerateHashPassword(newUser.Password)
	if err != nil {
		return model.User{}, errors.New("erreur de cryptage du mot de passe utilisateur")
	}

	userModel := model.User{
		Name:     newUser.Name,
		Email:    newUser.Email,
		RoleId:   roleResponse.Id,
		Sername:  newUser.Sername,
		Username: newUser.Username,
		Password: hashPassword,
	}

	userResponse, err := s.userRepo.CreateUser(userModel)
	if err != nil {
		return model.User{}, err
	}

	userResponse.RoleId = roleResponse.Name

	return userResponse, nil
}

// UpdateUser met à jour un utilisateur
func (s *UserService) UpdateUser(updatedUser User) (model.User, error) {

	existingUser, err := s.userRepo.GetUserByUsername(updatedUser.Username)
	if err != nil {
		return model.User{}, errors.New("utilisateur inconnu du système")
	}

	existingUser.Email = updatedUser.Email
	existingUser.Sername = updatedUser.Sername
	existingUser.Name = updatedUser.Name
	existingUser.UpdatedAt = time.Now()

	updateUser, err := s.userRepo.UpdateUser(existingUser)
	if err != nil {
		return model.User{}, errors.New("nous avons rencontré un problème durant la mise à jour")
	}

	roleResponse, _ := s.roleRepo.GetRoleById(updateUser.RoleId)

	updateUser.RoleId = roleResponse.Name

	return updateUser, nil
}

// DeleteUser supprime un utilisateur
func (s *UserService) DeleteUser(id string) error {
	existingUser, err := s.userRepo.GetUserById(id)
	if err != nil {
		return errors.New("utilisateur non trouvé")
	}

	existingUser.IsVisible = false
	existingUser.DeletedAt = time.Now()
	return s.userRepo.ArchivedUser(existingUser)
}

// GetUserByRole récupère tous les utilisateurs avec un rôle spécifié
func (s *UserService) GetUserByRole(role string) ([]model.User, error) {
	return s.userRepo.GetUserByRole(role)
}

// RemoveRole retire un rôle à un utilisateur
func (s *UserService) RemoveRole(userId string) (model.User, error) {
	user, err := s.userRepo.GetUserById(userId)
	if err != nil {
		return model.User{}, err
	}

	if user.RoleId == "" {
		return model.User{}, errors.New("l'utilisateur n'a pas de rôle à retirer")
	}

	user.RoleId = ""
	updatedUser, err := s.userRepo.UpdateUser(user)
	if err != nil {
		return model.User{}, err
	}

	return updatedUser, nil
}

// AssignRole assigne un rôle à un utilisateur
func (s *UserService) AssignRole(userId string, role string) (model.User, error) {
	user, err := s.userRepo.GetUserById(userId)
	if err != nil {
		return model.User{}, err
	}

	user.RoleId = role
	updatedUser, err := s.userRepo.UpdateUser(user)
	if err != nil {
		return model.User{}, err
	}

	return updatedUser, nil
}

func (s *UserService) GetUserByEmail(email string) (model.User, error) {

	userResponse, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return model.User{}, err
	}

	roleResponse, _ := s.roleRepo.GetRoleById(userResponse.RoleId)

	userResponse.RoleId = roleResponse.Name

	return userResponse, nil
}
