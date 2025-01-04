package config

import (
	"github.com/IlhamSetiaji/gift-redeem-be/internal/http/request"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

func NewValidator(viper *viper.Viper) *validator.Validate {
	validate := validator.New()
	validate.RegisterValidation("UserStatusValidation", request.UserStatusValidation)
	validate.RegisterValidation("UserGenderValidation", request.UserGenderValidation)
	return validate
}
