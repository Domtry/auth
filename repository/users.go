package repository

import (
	"auth/model"
	"errors"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository crée une nouvelle instance de UserRepository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// GetAllUsers récupère tous les utilisateurs depuis la base de données
func (r *UserRepository) GetAllUsers() ([]model.User, error) {
	var users []model.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// GetUserById récupère un utilisateur basé sur l'ID depuis la base de données
func (r *UserRepository) GetUserById(id string) (model.User, error) {
	var user model.User
	if err := r.db.First(&user, id).Error; err != nil {
		return model.User{}, err
	}
	return user, nil
}

// CreateUser crée un nouvel utilisateur dans la base de données
func (r *UserRepository) CreateUser(newUser model.User) (model.User, error) {
	// Générez la date actuelle pour les champs créés et modifiés
	currentTime := time.Now()
	newUser.CreatedAt = currentTime
	newUser.UpdatedAt = currentTime

	// Insertion dans la base de données
	if err := r.db.Create(&newUser).Error; err != nil {
		return model.User{}, err
	}

	msg := fmt.Sprintf("Created new user: %+v", newUser)
	log.Println(msg)
	return newUser, nil
}

// UpdateUser met à jour un utilisateur dans la base de données
func (r *UserRepository) UpdateUser(updatedUser model.User) (model.User, error) {
	// Générez la date actuelle pour le champ modifié
	updatedUser.UpdatedAt = time.Now()

	// Mise à jour dans la base de données
	if err := r.db.Save(&updatedUser).Error; err != nil {
		return model.User{}, err
	}

	msg := fmt.Sprintf("Updated user with ID %v and new role %v", updatedUser.Id, updatedUser.RoleId)
	log.Println(msg)

	return updatedUser, nil
}

// DeleteUser supprime un utilisateur de la base de données
func (r *UserRepository) DeleteUser(id string) error {
	// Suppression dans la base de données
	if err := r.db.Delete(&model.User{}, id).Error; err != nil {
		return err
	}

	return nil
}

// GetUserByRole récupère tous les utilisateurs avec un rôle spécifié depuis la base de données
func (r *UserRepository) GetUserByRole(role string) ([]model.User, error) {
	var users []model.User
	if err := r.db.Where("role = ?", role).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// GetUserByUsername récupère un utilisateur par son nom d'utilisateur depuis la base de données
func (r *UserRepository) GetUserByUsername(username string) (model.User, error) {
	var user model.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.User{}, errors.New("utilisateur non trouvé")
		}
		return model.User{}, err
	}
	return user, nil
}
