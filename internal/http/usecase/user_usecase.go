package usecase

import (
	"errors"

	"github.com/IlhamSetiaji/gift-redeem-be/internal/http/dto"
	"github.com/IlhamSetiaji/gift-redeem-be/internal/http/request"
	"github.com/IlhamSetiaji/gift-redeem-be/internal/http/response"
	"github.com/IlhamSetiaji/gift-redeem-be/internal/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type IUserUseCase interface {
	Login(payload *request.UserLoginRequest) (*response.UserResponse, error)
	FindByID(id uuid.UUID) (*response.UserResponse, error)
}

type UserUseCase struct {
	Log        *logrus.Logger
	Repository repository.IUserRepository
	DTO        dto.IUserDTO
}

func NewUserUseCase(
	log *logrus.Logger,
	repository repository.IUserRepository,
	dto dto.IUserDTO,
) IUserUseCase {
	return &UserUseCase{
		Log:        log,
		Repository: repository,
		DTO:        dto,
	}
}

func UserUseCaseFactory(log *logrus.Logger) IUserUseCase {
	repository := repository.UserRepositoryFactory(log)
	dto := dto.UserDTOFactory(log)
	return NewUserUseCase(log, repository, dto)
}

func (u *UserUseCase) Login(payload *request.UserLoginRequest) (*response.UserResponse, error) {
	user, err := u.Repository.FindByEmail(payload.Email)
	if err != nil {
		u.Log.Error("[UserUseCase.Login] " + err.Error())
		return nil, err
	}

	if user == nil {
		u.Log.Warn("[UserUseCase.Login] User not found")
		return nil, nil
	}

	if user.EmailVerifiedAt.IsZero() {
		u.Log.Warn("[UserUseCase.Login] User email not verified")
		return nil, errors.New("email not verified")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		u.Log.Error("Password not match")
		return nil, errors.New("email or password is incorrect")
	}

	return u.DTO.ConvertEntityToUserResponse(user), nil
}

func (u *UserUseCase) FindByID(id uuid.UUID) (*response.UserResponse, error) {
	user, err := u.Repository.FindById(id)
	if err != nil {
		u.Log.Error("[UserUseCase.FindByID] " + err.Error())
		return nil, err
	}

	if user == nil {
		u.Log.Warn("[UserUseCase.FindByID] User not found")
		return nil, nil
	}

	return u.DTO.ConvertEntityToUserResponse(user), nil
}
