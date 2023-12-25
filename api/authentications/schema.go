package authentications

import "time"

type AuthOut struct {
	UserId string `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	Token  Token  `json:"token"`
}

type AuthIn struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Token struct {
	ExpiresAt    time.Time `json:"expires_at"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
}
