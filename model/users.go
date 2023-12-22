package model

type User struct {
	AbstractModel
	Name     string `json:"name" gorm:"size:80" validate:"required"`
	Email    string `json:"email" gorm:"unique" validate:"required,email"`
	Username string `json:"username" gorm:"unique" validate:"required"`
	RoleId   string `json:"role_id" validate:"required,uuid4"`
	Sername  string `json:"sername" validate:"required"`
	Password string `json:"password" validate:"required"`
}
