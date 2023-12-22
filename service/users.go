// service/user_service.go

package service

import (
	"auth/model"
	"auth/repository"
	"errors"
)

// UserService gère la logique métier liée aux utilisateurs
type UserService struct {
	userRepo *repository.UserRepository
}

// NewUserService crée une nouvelle instance de UserService
func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// GetAllUsers récupère tous les utilisateurs
func (s *UserService) GetAllUsers() ([]model.User, error) {
	return s.userRepo.GetAllUsers()
}

// GetUserById récupère un utilisateur par son ID
func (s *UserService) GetUserById(id string) (model.User, error) {
	return s.userRepo.GetUserById(id)
}

// CreateUser crée un nouvel utilisateur avec un rôle spécifié
func (s *UserService) CreateUser(newUser model.User) (model.User, error) {
	// Ajoutez ici une logique de validation si nécessaire

	// Par exemple, assurez-vous que le nom d'utilisateur n'est pas déjà pris
	existingUser, err := s.userRepo.GetUserByUsername(newUser.Username)
	if err == nil && existingUser.Id != "" {
		return model.User{}, errors.New("le nom d'utilisateur est déjà pris")
	}

	return s.userRepo.CreateUser(newUser)
}

// UpdateUser met à jour un utilisateur
func (s *UserService) UpdateUser(updatedUser model.User) (model.User, error) {
	// Ajoutez ici une logique de validation si nécessaire

	// Par exemple, assurez-vous que le nom d'utilisateur n'est pas déjà pris par un autre utilisateur
	existingUser, err := s.userRepo.GetUserByUsername(updatedUser.Username)
	if err == nil && existingUser.Id != "" && existingUser.Id != updatedUser.Id {
		return model.User{}, errors.New("le nom d'utilisateur est déjà pris")
	}

	return s.userRepo.UpdateUser(updatedUser)
}

// DeleteUser supprime un utilisateur
func (s *UserService) DeleteUser(id string) error {
	// Ajoutez ici une logique de validation si nécessaire

	// Par exemple, assurez-vous que l'utilisateur que vous essayez de supprimer existe réellement
	_, err := s.userRepo.GetUserById(id)
	if err != nil {
		return errors.New("utilisateur non trouvé")
	}

	return s.userRepo.DeleteUser(id)
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
