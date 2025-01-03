package response

import (
	"time"

	"github.com/IlhamSetiaji/gift-redeem-be/internal/entity"
	"github.com/google/uuid"
)

type UserResponse struct {
	ID              uuid.UUID         `json:"id"`
	Email           string            `json:"email"`
	Name            string            `json:"name"`
	Username        string            `json:"username"`
	EmailVerifiedAt time.Time         `json:"email_verified_at"`
	Gender          entity.UserGender `json:"gender"`
	Status          entity.UserStatus `json:"status"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
	Roles           *[]RoleResponse   `json:"roles"`
}
