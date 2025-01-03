package repository

import (
	"errors"

	"github.com/IlhamSetiaji/gift-redeem-be/internal/config"
	"github.com/IlhamSetiaji/gift-redeem-be/internal/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IRoleRepository interface {
	GetAllRoles() (*[]entity.Role, error)
	FindAllPaginated(page int, pageSize int, search string) (*[]entity.Role, int64, error)
	FindById(id uuid.UUID) (*entity.Role, error)
	StoreRole(role *entity.Role) (*entity.Role, error)
	UpdateRole(role *entity.Role) (*entity.Role, error)
	GetAllRolesNotInUserID(userID uuid.UUID) (*[]entity.Role, error)
	GetAllRolesInUserID(userID uuid.UUID) (*[]entity.Role, error)
	DeleteRole(id uuid.UUID) error
}

type RoleRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewRoleRepository(log *logrus.Logger, db *gorm.DB) IRoleRepository {
	return &RoleRepository{
		Log: log,
		DB:  db,
	}
}

func RoleRepositoryFactory(log *logrus.Logger) IRoleRepository {
	db := config.NewDatabase()
	return NewRoleRepository(log, db)
}

func (r *RoleRepository) GetAllRoles() (*[]entity.Role, error) {
	var roles []entity.Role
	if err := r.DB.Preload("Users").Find(&roles).Error; err != nil {
		r.Log.Error(err)
		return nil, err
	}
	return &roles, nil
}

func (r *RoleRepository) FindAllPaginated(page int, pageSize int, search string) (*[]entity.Role, int64, error) {
	var roles []entity.Role
	var total int64

	query := r.DB.Preload("Users")

	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&roles).Error; err != nil {
		r.Log.Error("[RoleRepository.FindAllPaginated] " + err.Error())
		return nil, 0, errors.New("[RoleRepository.FindAllPaginated] " + err.Error())
	}

	if err := query.Count(&total).Error; err != nil {
		r.Log.Error("[RoleRepository.FindAllPaginated] " + err.Error())
		return nil, 0, errors.New("[RoleRepository.FindAllPaginated] " + err.Error())
	}

	return &roles, total, nil
}

func (r *RoleRepository) FindById(id uuid.UUID) (*entity.Role, error) {
	var role entity.Role
	if err := r.DB.Preload("Users").Where("id = ?", id).First(&role).Error; err != nil {
		r.Log.Error(err)
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) StoreRole(role *entity.Role) (*entity.Role, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return nil, errors.New("[UserRepository.CreateUser] failed to begin transaction: " + tx.Error.Error())
	}

	if err := tx.Create(role).Error; err != nil {
		tx.Rollback()
		r.Log.Error(err)
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Error("[RoleRepository.CreateRole] failed to commit transaction: " + err.Error())
		return nil, errors.New("[RoleRepository.CreateRole] failed to commit transaction: " + err.Error())
	}

	return role, nil
}

func (r *RoleRepository) UpdateRole(role *entity.Role) (*entity.Role, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return nil, errors.New("[UserRepository.CreateUser] failed to begin transaction: " + tx.Error.Error())
	}

	if err := tx.Model(&role).Where("id = ?", role.ID).Updates(role).Error; err != nil {
		r.Log.Error(err)
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Error("[RoleRepository.UpdateRole] failed to commit transaction: " + err.Error())
		return nil, errors.New("[RoleRepository.UpdateRole] failed to commit transaction: " + err.Error())
	}
	return role, nil
}

func (r *RoleRepository) DeleteRole(id uuid.UUID) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return errors.New("[UserRepository.DeleteUser] failed to begin transaction: " + tx.Error.Error())
	}
	if err := tx.Where("id = ?", id).Delete(&entity.Role{}).Error; err != nil {
		r.Log.Error(err)
		return err
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Error("[RoleRepository.DeleteRole] failed to commit transaction: " + err.Error())
		return errors.New("[RoleRepository.DeleteRole] failed to commit transaction: " + err.Error())
	}
	return nil
}

func (r *RoleRepository) GetAllRolesNotInUserID(userID uuid.UUID) (*[]entity.Role, error) {
	var roles []entity.Role

	if err := r.DB.Preload("Users").Where("id NOT IN (?)", r.DB.Table("user_roles").Select("role_id").Where("user_id = ?", userID)).Find(&roles).Error; err != nil {
		r.Log.Error(err)
		return nil, err
	}

	return &roles, nil
}

func (r *RoleRepository) GetAllRolesInUserID(userID uuid.UUID) (*[]entity.Role, error) {
	var roles []entity.Role

	if err := r.DB.Preload("Users").Where("id IN (?)", r.DB.Table("user_roles").Select("role_id").Where("user_id = ?", userID)).Find(&roles).Error; err != nil {
		r.Log.Error(err)
		return nil, err
	}

	return &roles, nil
}
