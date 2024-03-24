package authentications

import "time"

type AuthOut struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
	Token Token  `json:"token"`
}

type AuthResponse[T any] struct {
	Username        string `json:"username"`
	IsAuthenticated bool   `json:"is_authenticated"`
	UseOTP          bool   `json:"use_otp"`
	Content         T      `json:"content"`
}

type AuthIn struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type TwoFactorIn struct {
	SessionId string `json:"session_id"`
	Otp       string `json:"otp"`
}

type RefreshTokenOut struct {
	Id    string `json:"id"`
	Token Token  `json:"token"`
}

type ResetPasswordIn struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required"`
}

type InitUserPasswordIn struct {
	UserId   string `json:"user_id" validate:"required,uuid"`
	Password string `json:"password" validate:"required"`
}

type Token struct {
	ExpiresAt    time.Time `json:"expires_at"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
}
