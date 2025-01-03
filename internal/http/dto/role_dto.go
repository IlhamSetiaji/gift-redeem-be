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
	Log *logrus.Logger
}

func NewRoleDTO(log *logrus.Logger) IRoleDTO {
	return &RoleDTO{
		Log: log,
	}
}

func RoleDTOFactory(log *logrus.Logger) IRoleDTO {
	return NewRoleDTO(log)
}

func (r *RoleDTO) ConvertEntityToRoleResponse(payload *entity.Role) *response.RoleResponse {
	return &response.RoleResponse{
		ID:        payload.ID,
		Name:      payload.Name,
		GuardName: payload.GuardName,
		Status:    payload.Status,
		CreatedAt: payload.CreatedAt,
		UpdatedAt: payload.UpdatedAt,
	}
}

func (r *RoleDTO) ConvertEntitiesToRoleResponses(payload *[]entity.Role) *[]response.RoleResponse {
	var roles []response.RoleResponse
	for _, role := range *payload {
		roles = append(roles, *r.ConvertEntityToRoleResponse(&role))
	}
	return &roles
}
