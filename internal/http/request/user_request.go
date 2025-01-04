package request

import (
	"time"

	"github.com/IlhamSetiaji/gift-redeem-be/internal/entity"
	"github.com/google/uuid"
)

type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserRequest struct {
	ID                   uuid.UUID         `json:"id" validate:"omitempty,uuid"`
	Username             string            `json:"username" validate:"omitempty"`
	Email                string            `json:"email" validate:"omitempty,email"`
	Name                 string            `json:"name" validate:"omitempty"`
	Password             string            `json:"password" validate:"omitempty"`
	PasswordConfirmation string            `json:"password_confirmation" validate:"omitempty,eqfield=Password"`
	Gender               entity.UserGender `json:"gender" validate:"omitempty,UserGenderValidation"`
	EmailVerifiedAt      time.Time         `json:"email_verified_at" validate:"omitempty"`
	Status               entity.UserStatus `json:"status" validate:"omitempty,UserStatusValidation"`
	RoleIDs              []uuid.UUID       `json:"role_ids" validate:"omitempty,dive,uuid"`
}

type UserRegisterRequest struct {
	Username             string            `json:"username" validate:"required"`
	Email                string            `json:"email" validate:"required,email"`
	Name                 string            `json:"name" validate:"required"`
	Password             string            `json:"password" validate:"required"`
	PasswordConfirmation string            `json:"password_confirmation" validate:"required,eqfield=Password"`
	Gender               entity.UserGender `json:"gender" validate:"required,UserGenderValidation"`
	RoleIDs              []uuid.UUID       `json:"role_ids" validate:"required,dive,uuid"`
}
