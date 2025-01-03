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
	Log *logrus.Logger
}

func NewUserDTO(log *logrus.Logger) IUserDTO {
	return &UserDTO{
		Log: log,
	}
}

func UserDTOFactory(log *logrus.Logger) IUserDTO {
	return NewUserDTO(log)
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
		Roles:           nil,
	}
}

func (u *UserDTO) ConvertEntitiesToUserResponses(payload *[]entity.User) *[]response.UserResponse {
	var users []response.UserResponse
	for _, user := range *payload {
		users = append(users, response.UserResponse{
			ID:              user.ID,
			Email:           user.Email,
			Name:            user.Name,
			Username:        user.Username,
			EmailVerifiedAt: user.EmailVerifiedAt,
			Gender:          user.Gender,
			Status:          user.Status,
			Roles:           nil,
		})
	}
	return &users
}
