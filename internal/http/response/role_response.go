package response

import "time"

type RoleResponse struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	GuardName string          `json:"guard_name"`
	Status    string          `json:"status"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	Users     *[]UserResponse `json:"users"`
}
