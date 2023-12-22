package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AbstractModel struct {
	Id        string    `json:"id,omitempty" gorm:"primarykey" validate:"-"`
	CreatedAt time.Time `json:"create_at,omitempty" gorm:"autoCreateTime:nano" validate:"-"`
	UpdatedAt time.Time `json:"updated_at,omitempty" gorm:"autoUpdateTime:nano" validate:"-"`
	DeletedAt time.Time `json:"deleted_at,omitempty" validate:"-"`
	IsVisible bool      `json:"is_visible,omitempty" gorm:"default:true" validate:"-"`
}

func (a *AbstractModel) BeforeCreate(tx *gorm.DB) (err error) {
	a.Id = uuid.New().String()
	return
}
