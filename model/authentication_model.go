package model

import "time"

type Authentication struct {
	UserId string `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	Token  Token  `json:"token"`
}

type Token struct {
	ExpiresAt    time.Time `json:"expires_at"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
}
