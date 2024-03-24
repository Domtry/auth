package repository

import (
	"auth/model"
	"fmt"
	"github.com/google/uuid"
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
	tx := r.db.Model(model.User{}).
		Where("is_visible = ?", true).Find(&users)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return users, nil
}

// GetUserById récupère un utilisateur basé sur l'ID depuis la base de données
func (r *UserRepository) GetUserById(id string) (model.User, error) {
	var user model.User
	tx := r.db.Model(model.User{}).
		First(&user, "id = ? and is_visible = ?", id, true)
	if tx.Error != nil {
		return model.User{}, tx.Error
	}
	return user, nil
}

// CreateUser crée un nouvel utilisateur dans la base de données
func (r *UserRepository) CreateUser(newUser model.User) (model.User, error) {

	newUser.Id = uuid.New().String()
	newUser.CreatedAt = time.Now()

	// Insertion dans la base de données
	if err := r.db.Model(model.User{}).Create(&newUser).Error; err != nil {
		return model.User{}, err
	}

	msg := fmt.Sprintf("Created new user: %+v", newUser)
	log.Println(msg)
	return newUser, nil
}

// UpdateUser met à jour un utilisateur dans la base de données
func (r *UserRepository) UpdateUser(updatedUser model.User) (model.User, error) {
	updatedUser.UpdatedAt = time.Now()
	tx := r.db.Model(model.User{}).
		Where("id = ? and is_visible = ?", updatedUser.Id, true).Save(&updatedUser)
	if tx.Error != nil {
		return model.User{}, tx.Error
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

// DeleteUser supprime un utilisateur de la base de données
func (r *UserRepository) ArchivedUser(user model.User) error {
	// Suppression dans la base de données
	_, err := r.UpdateUser(user)
	return err
}

// GetUserByRole récupère tous les utilisateurs avec un rôle spécifié
func (r *UserRepository) GetUserByRole(role string) ([]model.User, error) {
	var users []model.User
	tx := r.db.Model(model.User{}).Where(
		"role = ? and is_visible = ?", role, true).Find(&users)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return users, nil
}

// GetUserByUsername récupère un utilisateur par son nom
func (r *UserRepository) GetUserByUsername(username string) (model.User, error) {
	var user model.User
	tx := r.db.Model(model.User{}).
		First(&user, "username = ? and is_visible = ?", username, true)
	if tx.Error != nil {
		return model.User{}, tx.Error
	}
	return user, nil
}

// GetUserByEmail récupère un utilisateur par son email
func (r *UserRepository) GetUserByEmail(email string) (model.User, error) {
	var user model.User
	tx := r.db.Model(model.User{}).
		First(&user, "email = ? and is_visible = ?", email, true)
	if tx.Error != nil {
		return model.User{}, tx.Error
	}
	return user, nil
}
