package usecase

import (
	"errors"

	"github.com/IlhamSetiaji/gift-redeem-be/internal/entity"
	"github.com/IlhamSetiaji/gift-redeem-be/internal/http/dto"
	"github.com/IlhamSetiaji/gift-redeem-be/internal/http/messaging"
	"github.com/IlhamSetiaji/gift-redeem-be/internal/http/request"
	"github.com/IlhamSetiaji/gift-redeem-be/internal/http/response"
	"github.com/IlhamSetiaji/gift-redeem-be/internal/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/rand"
)

type IUserUseCase interface {
	Login(payload *request.UserLoginRequest) (*response.UserResponse, error)
	Register(payload *request.UserRegisterRequest) (*response.UserResponse, error)
	FindByID(id uuid.UUID) (*response.UserResponse, error)
}

type UserUseCase struct {
	Log         *logrus.Logger
	Repository  repository.IUserRepository
	DTO         dto.IUserDTO
	MailMessage messaging.IMailMessage
}

func NewUserUseCase(
	log *logrus.Logger,
	repository repository.IUserRepository,
	dto dto.IUserDTO,
	mailMessage messaging.IMailMessage,
) IUserUseCase {
	return &UserUseCase{
		Log:         log,
		Repository:  repository,
		DTO:         dto,
		MailMessage: mailMessage,
	}
}

func UserUseCaseFactory(log *logrus.Logger) IUserUseCase {
	repository := repository.UserRepositoryFactory(log)
	dto := dto.UserDTOFactory(log)
	mailMessage := messaging.MailMessageFactory(log)
	return NewUserUseCase(log, repository, dto, mailMessage)
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

func (u *UserUseCase) Register(payload *request.UserRegisterRequest) (*response.UserResponse, error) {
	user, err := u.Repository.FindByEmail(payload.Email)
	if err != nil {
		u.Log.Error("[UserUseCase.Register] " + err.Error())
		return nil, err
	}

	if user != nil {
		u.Log.Warn("[UserUseCase.Register] User already registered")
		return nil, errors.New("user already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		u.Log.Error("[UserUseCase.Register] " + err.Error())
		return nil, err
	}

	user = &entity.User{
		Username: payload.Username,
		Email:    payload.Email,
		Name:     payload.Name,
		Password: string(hashedPassword),
		Gender:   payload.Gender,
		Status:   entity.USER_PENDING,
	}

	if _, err := u.Repository.CreateUser(user, payload.RoleIDs); err != nil {
		u.Log.Error("[UserUseCase.Register] " + err.Error())
		return nil, err
	}

	randomIntToken := rand.Intn(100000)
	if err := u.Repository.CreateUserToken(payload.Email, randomIntToken); err != nil {
		u.Log.Error("[UserUseCase.Register] " + err.Error())
		return nil, err
	}

	if _, err := u.MailMessage.SendMail(&request.MailRequest{
		Email:   payload.Email,
		Subject: "Email Verification",
		Body:    "Your verification code is " + string(randomIntToken),
		From:    "ilham.ahmadz18@gmail.com",
		To:      payload.Email,
	}); err != nil {
		u.Log.Error("[UserUseCase.Register] " + err.Error())
		return nil, err
	}

	return u.DTO.ConvertEntityToUserResponse(user), nil
}
