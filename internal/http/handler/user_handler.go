package handler

import (
	"net/http"

	"github.com/IlhamSetiaji/gift-redeem-be/internal/config"
	"github.com/IlhamSetiaji/gift-redeem-be/internal/http/middleware"
	"github.com/IlhamSetiaji/gift-redeem-be/internal/http/request"
	"github.com/IlhamSetiaji/gift-redeem-be/internal/http/usecase"
	"github.com/IlhamSetiaji/gift-redeem-be/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IUserHandler interface {
	Login(ctx *gin.Context)
	UserMe(ctx *gin.Context)
}

type UserHandler struct {
	Log      *logrus.Logger
	Viper    *viper.Viper
	Validate *validator.Validate
	UseCase  usecase.IUserUseCase
}

func NewUserHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	validate *validator.Validate,
	useCase usecase.IUserUseCase,
) IUserHandler {
	return &UserHandler{
		Log:      log,
		Viper:    viper,
		Validate: validate,
		UseCase:  useCase,
	}
}

func UserHandlerFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) IUserHandler {
	useCase := usecase.UserUseCaseFactory(log)
	validate := config.NewValidator(viper)
	return NewUserHandler(log, viper, validate, useCase)
}

func (u *UserHandler) Login(ctx *gin.Context) {
	var payload = new(request.UserLoginRequest)
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		u.Log.Error("[UserHandler.Login] " + err.Error())
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}

	err := u.Validate.Struct(payload)
	if err != nil {
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		u.Log.Errorf("Error when validating request: %v", err)
		return
	}

	user, err := u.UseCase.Login(payload)
	if err != nil {
		u.Log.Error("[UserHandler.Login] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	if user == nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "error", "User not found")
		return
	}

	token, err := utils.GenerateToken(user)
	if err != nil {
		u.Log.Error("[UserHandler.Login] " + err.Error())
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "error", err.Error())
		return
	}

	var data = map[string]interface{}{
		"token":      token,
		"token_type": "Bearer",
		"user":       user,
	}

	utils.SuccessResponse(ctx, http.StatusOK, "success", data)
}

func (u *UserHandler) UserMe(ctx *gin.Context) {
	user, err := middleware.GetUser(ctx)
	if err != nil {
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		u.Log.Errorf("Error when getting user: %v", err)
		return
	}
	if user == nil {
		utils.ErrorResponse(ctx, 404, "error", "User not found")
		u.Log.Errorf("User not found")
		return
	}

	me, err := u.UseCase.FindByID(uuid.MustParse(user["id"].(string)))
	if err != nil {
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		u.Log.Errorf("Error when getting user: %v", err)
		return
	}

	if me == nil {
		utils.ErrorResponse(ctx, 404, "error", "User not found")
		u.Log.Errorf("User not found")
		return
	}

	utils.SuccessResponse(ctx, 200, "success", me)
}
