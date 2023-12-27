package users

import "time"

type UserIn struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required"`
	Sername  string `json:"sername" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UpdateUserIn struct {
	Name    string `json:"name" validate:"required"`
	Email   string `json:"email" validate:"required,email"`
	Sername string `json:"sername" validate:"required"`
}

type UserOut struct {
	Id        string    `json:"id" validate:"required"`
	Name      string    `json:"name" validate:"required"`
	Email     string    `json:"email" validate:"required,email"`
	Username  string    `json:"username" validate:"required"`
	Role      string    `json:"role" validate:"required,uuid4"`
	Sername   string    `json:"sername" validate:"required"`
	CreatedAt time.Time `json:"create_at,omitempty" validate:"-"`
	UpdatedAt time.Time `json:"updated_at,omitempty" validate:"-"`
}
