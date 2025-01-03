package utils

import (
	"time"

	"github.com/IlhamSetiaji/gift-redeem-be/internal/http/response"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func GenerateToken(user *response.UserResponse) (string, error) {
	viper := viper.New()
	logger := logrus.New()

	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("./")
	err := viper.ReadInConfig()

	if err != nil {
		logger.Fatalf("Fatal error config file: %v", err)
	}

	roles := make([]map[string]interface{}, len(*user.Roles))
	for i, role := range *user.Roles {
		roles[i] = map[string]interface{}{
			"name": role.Name,
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"name":     user.Name,
		"username": user.Username,
		"email":    user.Email,
		"roles":    roles,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte(viper.GetString("jwt.secret")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateTokenForOAuth2(data *map[string]interface{}) (string, error) {
	viper := viper.New()
	logger := logrus.New()

	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("./")
	err := viper.ReadInConfig()

	if err != nil {
		logger.Fatalf("Fatal error config file: %v", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"data": data,
		"exp":  time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte(viper.GetString("jwt.secret")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
