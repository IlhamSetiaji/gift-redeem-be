package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserStatus string
type UserGender string

const (
	USER_ACTIVE   UserStatus = "ACTIVE"
	USER_INACTIVE UserStatus = "INACTIVE"
	USER_PENDING  UserStatus = "PENDING"
)

const (
	MALE   UserGender = "MALE"
	FEMALE UserGender = "FEMALE"
)

type User struct {
	gorm.Model      `json:"-"`
	ID              uuid.UUID  `json:"id" gorm:"type:char(36);primaryKey"`
	Username        string     `json:"username" gorm:"unique;not null"`
	Email           string     `json:"email" gorm:"unique;not null"`
	Name            string     `json:"name"`
	Password        string     `json:"password"`
	Gender          UserGender `json:"gender" gorm:"not null"`
	EmailVerifiedAt time.Time  `json:"email_verified_at" gorm:"default:null"`
	Status          UserStatus `json:"status" gorm:"default:PENDING"`
	Roles           []Role     `json:"roles" gorm:"many2many:user_roles;"`

	RedeemedGifts []Redemption `json:"redeemed_gifts" gorm:"many2many:redemptions;constraint:onDelete:CASCADE;"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.ID = uuid.New()
	return nil
}

func (user *User) BeforeUpdate(tx *gorm.DB) (err error) {
	return nil
}

func (user *User) BeforeDelete(tx *gorm.DB) (err error) {
	if user.DeletedAt.Valid {
		return nil
	}

	randomString := uuid.New().String()

	user.Email = user.Email + "_deleted_" + randomString
	tx.Model(&User{}).Where("id = ?", user.ID).Updates(map[string]interface{}{
		"email": user.Email,
	})
	return nil
}

func (User) TableName() string {
	return "users"
}
