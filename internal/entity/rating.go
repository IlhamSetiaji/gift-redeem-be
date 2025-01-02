package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Rating struct {
	gorm.Model   `json:"-"`
	ID           uuid.UUID  `json:"id" gorm:"type:char(36);primaryKey"`
	RedemptionID *uuid.UUID `json:"redemption_id" gorm:"unique;type:char(36);not null"`
	Rating       float64    `json:"rating" gorm:"not null"`
	Comment      string     `json:"comment" gorm:"type:text;default:null"`

	Redemption *Redemption `json:"redemption" gorm:"foreignKey:RedemptionID;references:ID;constraint:OnDelete:CASCADE"`
}

func (rating *Rating) BeforeCreate(tx *gorm.DB) (err error) {
	rating.ID = uuid.New()
	return nil
}

func (rating *Rating) BeforeUpdate(tx *gorm.DB) (err error) {
	return nil
}

func (Rating) TableName() string {
	return "ratings"
}
