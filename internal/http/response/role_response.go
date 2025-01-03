package response

import (
	"time"

	"github.com/IlhamSetiaji/gift-redeem-be/internal/entity"
	"github.com/google/uuid"
)

type RoleResponse struct {
	ID        uuid.UUID         `json:"id"`
	Name      string            `json:"name"`
	GuardName string            `json:"guard_name"`
	Status    entity.RoleStatus `json:"status"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	Users     *[]UserResponse   `json:"users"`
}
