package repository

import (
	"auth/model"
	"fmt"
	"gorm.io/gorm"
	"log"
	"time"
)

type OtpRepository struct {
	db *gorm.DB
}

// NewOtpRepository crée une nouvelle instance de OtpRepository
func NewOtpRepository(db *gorm.DB) *OtpRepository {
	return &OtpRepository{
		db: db,
	}
}

// CreateUser crée un nouvel utilisateur dans la base de données
func (r *OtpRepository) CreateOtp(newOtp model.Otp) (model.Otp, error) {
	// Générez la date actuelle pour les champs créés et modifiés
	currentTime := time.Now()
	newOtp.CreatedAt = currentTime
	newOtp.UpdatedAt = currentTime

	// Insertion dans la base de données
	if err := r.db.Model(model.Otp{}).Create(&newOtp).Error; err != nil {
		return model.Otp{}, err
	}

	msg := fmt.Sprintf("Created new otp: %+v", newOtp)
	log.Println(msg)
	return newOtp, nil
}

// GetOtpById recuperate un utilisateur basé sur l'ID depuis la base de données
func (r *OtpRepository) GetOtpById(id string) (model.Otp, error) {
	var otp model.Otp
	tx := r.db.Model(model.Otp{}).
		First(&otp, "id = ? and is_used = ?", id, false)
	if tx.Error != nil {
		return model.Otp{}, tx.Error
	}
	return otp, nil
}

// UpdateUser met à jour un utilisateur dans la base de données
func (r *OtpRepository) UpdateOtp(otp model.Otp) (model.Otp, error) {
	otp.UpdatedAt = time.Now()
	tx := r.db.Model(model.User{}).
		Where("id = ? and is_visible = ?", otp.Id, true).Save(&otp)
	if tx.Error != nil {
		return model.Otp{}, tx.Error
	}

	msg := fmt.Sprintf("Updated otp with ID %v", otp.Id)
	log.Println(msg)

	return otp, nil
}
