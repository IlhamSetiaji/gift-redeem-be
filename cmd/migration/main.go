package main

import (
	"time"

	"github.com/IlhamSetiaji/gift-redeem-be/internal/config"
	"github.com/IlhamSetiaji/gift-redeem-be/internal/entity"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	viper := config.NewViper()
	log := config.NewLogrus(viper)
	db := config.NewDatabase()

	// migrate the schema
	err := db.AutoMigrate(&entity.Role{}, &entity.User{}, &entity.UserToken{}, &entity.UserRole{}, &entity.Gift{}, &entity.Redemption{}, &entity.Rating{})
	if err != nil {
		log.Fatal(err)
	} else {
		log.Info("Migration success")
	}

	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte("changeme"), bcrypt.DefaultCost)

	// seed superadmin role data
	superAdminRole := entity.Role{
		Name:      "superadmin",
		GuardName: "api",
		Status:    entity.ROLE_ACTIVE,
	}

	err = db.Create(&superAdminRole).Error

	if err != nil {
		log.Fatal(err)
	}

	// seed user role data
	userRole := entity.Role{
		Name:      "user",
		GuardName: "api",
		Status:    entity.ROLE_ACTIVE,
	}

	err = db.Create(&userRole).Error

	if err != nil {
		log.Fatal(err)
	}

	// seed superadmin user data
	superAdminUser := entity.User{
		Name:            "Super Admin",
		Username:        "superadmin",
		Email:           "superadmin@test.test",
		Password:        string(hashedPasswordBytes),
		Gender:          entity.MALE,
		EmailVerifiedAt: time.Now(),
		Status:          entity.USER_ACTIVE,
	}

	err = db.Create(&superAdminUser).Error

	if err != nil {
		log.Fatal(err)
	}

	err = db.Create(&entity.UserRole{
		UserID: superAdminUser.ID,
		RoleID: superAdminRole.ID,
	}).Error

	if err != nil {
		log.Fatal(err)
	}

	// seed user data
	user := entity.User{
		Name:            "User",
		Username:        "user",
		Email:           "user@test.test",
		Password:        string(hashedPasswordBytes),
		Gender:          entity.FEMALE,
		EmailVerifiedAt: time.Now(),
		Status:          entity.USER_ACTIVE,
	}

	err = db.Create(&user).Error

	if err != nil {
		log.Fatal(err)
	}

	err = db.Create(&entity.UserRole{
		UserID: user.ID,
		RoleID: userRole.ID,
	}).Error

	if err != nil {
		log.Fatal(err)
	}

	// seed gifts data
	gifts := []entity.Gift{
		{
			Name:        "Red Dead Redemption 2 Gift Card",
			RedeemCode:  "REDEMPTION2",
			Description: "Kalo main wajib pake kuda",
			Price:       10000,
			Stock:       10,
		},
		{
			Name:        "The Witcher 3 Gift Card",
			RedeemCode:  "WITCHER3",
			Description: "Kalo main wajib pake Geralt",
			Price:       20000,
			Stock:       20,
		},
		{
			Name:        "God of War 2018 Gift Card",
			RedeemCode:  "GODOFWAR",
			Description: "Kalo main wajib pake Kratos",
			Price:       30000,
			Stock:       30,
		},
	}

	for _, gift := range gifts {
		err = db.Create(&gift).Error
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Info("Seed success")
}
