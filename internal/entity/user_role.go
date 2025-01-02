package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserRole struct {
	UserID    uuid.UUID `json:"user_id" gorm:"type:char(36);primaryKey"`
	RoleID    uuid.UUID `json:"role_id" gorm:"type:char(36);primaryKey"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	User User `json:"user" gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	Role Role `json:"role" gorm:"foreignKey:RoleID;references:ID;constraint:OnDelete:CASCADE"`
}

func (userRole *UserRole) BeforeCreate() (err error) {
	userRole.CreatedAt = time.Now()
	userRole.UpdatedAt = time.Now()
	return nil
}

func (userRole *UserRole) BeforeUpdate() (err error) {
	userRole.UpdatedAt = time.Now()
	return nil
}

func (UserRole) TableName() string {
	return "user_roles"
}