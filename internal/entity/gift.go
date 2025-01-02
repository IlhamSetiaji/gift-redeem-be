package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Gift struct {
	gorm.Model  `json:"-"`
	ID          uuid.UUID `json:"id" gorm:"type:char(36);primaryKey"`
	RedeemCode  string    `json:"redeem_code" gorm:"unique;not null"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description" gorm:"type:text;default:null"`
	Price       int       `json:"price" gorm:"default:0"`
	Stock       int       `json:"stock" gorm:"default:0"`
	ExpiredAt   string    `json:"expired_at" gorm:"not null"`

	RedeemedUsers []Redemption `json:"redeemed_users" gorm:"many2many:redemptions;constraint:onDelete:CASCADE;"`
}

func (gift *Gift) BeforeCreate(tx *gorm.DB) (err error) {
	gift.ID = uuid.New()
	return nil
}

func (gift *Gift) BeforeUpdate(tx *gorm.DB) (err error) {
	return nil
}

func (gift *Gift) BeforeDelete(tx *gorm.DB) (err error) {
	if gift.DeletedAt.Valid {
		return nil
	}

	randomString := uuid.New().String()

	gift.RedeemCode = gift.RedeemCode + "_deleted_" + randomString
	tx.Model(&Gift{}).Where("id = ?", gift.ID).Updates(map[string]interface{}{
		"redeem_code": gift.RedeemCode,
	})
	return nil
}

func (Gift) TableName() string {
	return "gifts"
}
