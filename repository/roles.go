package repository

import (
	"auth/model"
	"gorm.io/gorm"
)

type RoleRepository struct {
	db *gorm.DB
}

// NewRoleRepository crée une nouvelle instance de RoleRepository
func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{
		db: db,
	}
}

// GetRoleById récupère un role basé sur l'ID depuis la base de données
func (r *RoleRepository) GetRoleById(id string) (model.Role, error) {
	var role model.Role
	if err := r.db.Model(model.Role{}).First(&role, "id = ?", id).Error; err != nil {
		return model.Role{}, err
	}
	return role, nil
}
