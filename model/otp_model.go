package model

import (
	"time"
)

type Otp struct {
	AbstractModel
	Code      string    `json:"code"`
	IsUsed    bool      `json:"is_used"`
	UserId    string    `json:"user_id"`
	ExpireHas time.Time `json:"expire_has"`
}
