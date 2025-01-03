package repository

import (
	"errors"

	"github.com/IlhamSetiaji/gift-redeem-be/internal/config"
	"github.com/IlhamSetiaji/gift-redeem-be/internal/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IUserRepository interface {
	FindByEmail(email string) (*entity.User, error)
	FindAllPaginated(page int, pageSize int, search string) (*[]entity.User, int64, error)
	FindById(id uuid.UUID) (*entity.User, error)
	CreateUser(user *entity.User, roleIDs []uuid.UUID) (*entity.User, error)
	UpdateUser(user *entity.User, roleIDs []uuid.UUID) (*entity.User, error)
	DeleteUser(id uuid.UUID) error
}

type UserRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewUserRepository(log *logrus.Logger, db *gorm.DB) IUserRepository {
	return &UserRepository{
		Log: log,
		DB:  db,
	}
}

func (r *UserRepository) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.DB.Preload("Roles").Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.Log.Warn("[UserRepository.FindByEmail] User not found")
			return nil, nil
		} else {
			r.Log.Error("[UserRepository.FindByEmail] " + err.Error())
			return nil, errors.New("[UserRepository.FindByEmail] " + err.Error())
		}
	}
	return &user, nil
}

func (r *UserRepository) FindAllPaginated(page int, pageSize int, search string) (*[]entity.User, int64, error) {
	var users []entity.User
	var total int64

	query := r.DB.Preload("Roles")

	if search != "" {
		query = query.Where("email ILIKE ?", "%"+search+"%").Or("name ILIKE ?", "%"+search+"%").Or("username ILIKE ?", "%"+search+"%")
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error; err != nil {
		r.Log.Error("[UserRepository.FindAllPaginated] " + err.Error())
		return nil, 0, errors.New("[UserRepository.FindAllPaginated] " + err.Error())
	}

	if err := query.Model(&entity.User{}).Count(&total).Error; err != nil {
		r.Log.Error("[UserRepository.FindAllPaginated] " + err.Error())
		return nil, 0, errors.New("[UserRepository.FindAllPaginated] " + err.Error())
	}

	return &users, total, nil
}

func (r *UserRepository) FindById(id uuid.UUID) (*entity.User, error) {
	var user entity.User
	err := r.DB.Preload("Roles").Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.Log.Warn("[UserRepository.FindById] User not found")
			return nil, nil
		} else {
			r.Log.Error("[UserRepository.FindById] " + err.Error())
			return nil, errors.New("[UserRepository.FindById] " + err.Error())
		}
	}
	return &user, nil
}

func (r *UserRepository) CreateUser(user *entity.User, roleIDs []uuid.UUID) (*entity.User, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return nil, errors.New("[UserRepository.CreateUser] failed to begin transaction: " + tx.Error.Error())
	}

	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[UserRepository.CreateUser] " + err.Error())
		return nil, errors.New("[UserRepository.CreateUser] " + err.Error())
	}

	for _, roleID := range roleIDs {
		var role entity.Role
		if err := tx.First(&role, "id = ?", roleID).Error; err != nil {
			tx.Rollback()
			r.Log.Error("[UserRepository.CreateUser] Role not found: " + err.Error())
			return nil, errors.New("[UserRepository.CreateUser] Role not found: " + err.Error())
		}

		var userRole = entity.UserRole{
			UserID: user.ID,
			RoleID: role.ID,
		}

		if err := tx.Create(&userRole).Error; err != nil {
			tx.Rollback()
			r.Log.Error("[UserRepository.AppendUser] " + err.Error())
			return nil, errors.New("[UserRepository.AppendUser] " + err.Error())
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Error("[UserRepository.CreateUser] failed to commit transaction: " + err.Error())
		return nil, errors.New("[UserRepository.CreateUser] failed to commit transaction: " + err.Error())
	}

	if err := r.DB.Preload("Roles").First(user, user.ID).Error; err != nil {
		r.Log.Error("[UserRepository.CreateUser] Failed to reload user: " + err.Error())
		return nil, errors.New("[UserRepository.CreateUser] Failed to reload user: " + err.Error())
	}

	return user, nil
}

func (r *UserRepository) UpdateUser(user *entity.User, roleIDs []uuid.UUID) (*entity.User, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return nil, errors.New("[UserRepository.UpdateUser] failed to begin transaction: " + tx.Error.Error())
	}

	if err := tx.Model(&user).Where("id = ?", user.ID).Updates(entity.User{
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
		Gender:   user.Gender,
		Status:   user.Status,
	}).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[UserRepository.UpdateUser] " + err.Error())
		return nil, errors.New("[UserRepository.UpdateUser] " + err.Error())
	}

	if len(roleIDs) > 0 {
		// delete user role by user id
		if err := tx.Where("user_id = ?", user.ID).Delete(&entity.UserRole{}).Error; err != nil {
			tx.Rollback()
			r.Log.Error("[UserRepository.UpdateUser] " + err.Error())
			return nil, errors.New("[UserRepository.UpdateUser] " + err.Error())
		}
		for _, roleID := range roleIDs {
			var role entity.Role
			if err := tx.First(&role, "id = ?", roleID).Error; err != nil {
				tx.Rollback()
				r.Log.Error("[UserRepository.UpdateUser] Role not found: " + err.Error())
				return nil, errors.New("[UserRepository.UpdateUser] Role not found: " + err.Error())
			}

			var userRole = entity.UserRole{
				UserID: user.ID,
				RoleID: role.ID,
			}

			if err := tx.Create(&userRole).Error; err != nil {
				tx.Rollback()
				r.Log.Error("[UserRepository.UpdateUser] " + err.Error())
				return nil, errors.New("[UserRepository.UpdateUser] " + err.Error())
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Error("[UserRepository.CreateUser] failed to commit transaction: " + err.Error())
		return nil, errors.New("[UserRepository.CreateUser] failed to commit transaction: " + err.Error())
	}

	if err := r.DB.Preload("Roles").First(user, user.ID).Error; err != nil {
		r.Log.Error("[UserRepository.CreateUser] Failed to reload user: " + err.Error())
		return nil, errors.New("[UserRepository.CreateUser] Failed to reload user: " + err.Error())
	}

	return user, nil
}

func (r *UserRepository) DeleteUser(id uuid.UUID) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return errors.New("[UserRepository.DeleteUser] failed to begin transaction: " + tx.Error.Error())
	}

	var user entity.User
	if err := tx.First(&user, "id = ?", id).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[UserRepository.DeleteUser] User not found: " + err.Error())
		return errors.New("[UserRepository.DeleteUser] User not found: " + err.Error())
	}

	if err := tx.Delete(&user).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[UserRepository.DeleteUser] " + err.Error())
		return errors.New("[UserRepository.DeleteUser] " + err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Error("[UserRepository.DeleteUser] failed to commit transaction: " + err.Error())
		return errors.New("[UserRepository.DeleteUser] failed to commit transaction: " + err.Error())
	}

	return nil
}

func UserRepositoryFactory(log *logrus.Logger) IUserRepository {
	db := config.NewDatabase()
	return NewUserRepository(log, db)
}
