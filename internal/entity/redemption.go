package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Redemption struct {
	gorm.Model `json:"-"`
	ID         uuid.UUID `json:"id" gorm:"type:char(36);primaryKey"`
	UserID     uuid.UUID `json:"user_id" gorm:"type:char(36);not null"`
	GiftID     uuid.UUID `json:"gift_id" gorm:"type:char(36);not null"`
	RedeemedAt time.Time `json:"redeemed_at" gorm:"default:CURRENT_TIMESTAMP"`

	User User `json:"user" gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	Gift Gift `json:"gift" gorm:"foreignKey:GiftID;references:ID;constraint:OnDelete:CASCADE"`

	Rating *Rating `json:"rating" gorm:"foreignKey:RedemptionID;references:ID;constraint:OnDelete:CASCADE"`
}

func (redemption *Redemption) BeforeCreate(tx *gorm.DB) (err error) {
	redemption.ID = uuid.New()
	return nil
}

func (redemption *Redemption) BeforeUpdate(tx *gorm.DB) (err error) {
	return nil
}

func (Redemption) TableName() string {
	return "redemptions"
}
