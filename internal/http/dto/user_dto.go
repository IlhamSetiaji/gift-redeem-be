package dto

import (
	"github.com/IlhamSetiaji/gift-redeem-be/internal/entity"
	"github.com/IlhamSetiaji/gift-redeem-be/internal/http/response"
	"github.com/sirupsen/logrus"
)

type IUserDTO interface {
	ConvertEntityToUserResponse(payload *entity.User) *response.UserResponse
	ConvertEntitiesToUserResponses(payload *[]entity.User) *[]response.UserResponse
}

type UserDTO struct {
	Log     *logrus.Logger
	RoleDTO IRoleDTO
}

func NewUserDTO(log *logrus.Logger, roleDTO IRoleDTO) IUserDTO {
	return &UserDTO{
		Log: log,
		RoleDTO: roleDTO,
	}
}

func UserDTOFactory(log *logrus.Logger) IUserDTO {
	roleDTO := RoleDTOFactory(log)
	return NewUserDTO(log, roleDTO)
}

func (u *UserDTO) ConvertEntityToUserResponse(payload *entity.User) *response.UserResponse {
	return &response.UserResponse{
		ID:              payload.ID,
		Email:           payload.Email,
		Name:            payload.Name,
		Username:        payload.Username,
		EmailVerifiedAt: payload.EmailVerifiedAt,
		Gender:          payload.Gender,
		Status:          payload.Status,
		CreatedAt:       payload.CreatedAt,
		UpdatedAt:       payload.UpdatedAt,
		Roles: func() *[]response.RoleResponse {
			var roles []response.RoleResponse
			if len(payload.Roles) == 0 || payload.Roles == nil {
				return nil
			}
			for _, role := range payload.Roles {
				roles = append(roles, *u.RoleDTO.ConvertEntityToRoleResponse(&role))
			}
			return &roles
		}(),
	}
}

func (u *UserDTO) ConvertEntitiesToUserResponses(payload *[]entity.User) *[]response.UserResponse {
	var users []response.UserResponse
	for _, user := range *payload {
		users = append(users, *u.ConvertEntityToUserResponse(&user))
	}
	return &users
}
