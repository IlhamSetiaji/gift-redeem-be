package dto

import (
	"github.com/IlhamSetiaji/gift-redeem-be/internal/entity"
	"github.com/IlhamSetiaji/gift-redeem-be/internal/http/response"
	"github.com/sirupsen/logrus"
)

type IRoleDTO interface {
	ConvertEntityToRoleResponse(payload *entity.Role) *response.RoleResponse
	ConvertEntitiesToRoleResponses(payload *[]entity.Role) *[]response.RoleResponse
}

type RoleDTO struct {
	Log     *logrus.Logger
	UserDTO IUserDTO
}

func NewRoleDTO(log *logrus.Logger, userDTO IUserDTO) IRoleDTO {
	return &RoleDTO{
		Log:     log,
		UserDTO: userDTO,
	}
}

func RoleDTOFactory(log *logrus.Logger) IRoleDTO {
	userDTO := UserDTOFactory(log)
	return NewRoleDTO(log, userDTO)
}

func (r *RoleDTO) ConvertEntityToRoleResponse(payload *entity.Role) *response.RoleResponse {
	return &response.RoleResponse{
		ID:        payload.ID,
		Name:      payload.Name,
		GuardName: payload.GuardName,
		Status:    payload.Status,
		CreatedAt: payload.CreatedAt,
		UpdatedAt: payload.UpdatedAt,
		Users: func() *[]response.UserResponse {
			var users []response.UserResponse
			if len(payload.Users) == 0 || payload.Users == nil {
				return nil
			}
			for _, user := range payload.Users {
				users = append(users, *r.UserDTO.ConvertEntityToUserResponse(&user))
			}
			return &users
		}(),
	}
}

func (r *RoleDTO) ConvertEntitiesToRoleResponses(payload *[]entity.Role) *[]response.RoleResponse {
	var roles []response.RoleResponse
	for _, role := range *payload {
		roles = append(roles, *r.ConvertEntityToRoleResponse(&role))
	}
	return &roles
}
