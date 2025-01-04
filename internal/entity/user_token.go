package entity

import "time"

type UserTokenType string

const (
	UserTokenVerification  UserTokenType = "VERIFICATION"
	UserTokenResetPassword UserTokenType = "RESET_PASSWORD"
)

type UserToken struct {
	Email     string        `json:"email"`
	Token     int           `json:"token"`
	TokenType UserTokenType `json:"token_type"`
	ExpiredAt time.Time     `json:"expired_at"`
	CreatedAt time.Time     `gorm:"autoCreateTime"`
	UpdatedAt time.Time     `gorm:"autoUpdateTime"`
}

func (UserToken) TableName() string {
	return "user_tokens"
}
