package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/IlhamSetiaji/gift-redeem-be/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

func NewAuth(viper *viper.Viper) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "error", "No Authorization header provided")
			c.Abort()
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 {
			utils.ErrorResponse(c, http.StatusUnauthorized, "error", "Invalid Authorization header format")
			c.Abort()
			return
		}

		token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(viper.GetString("jwt.secret")), nil
		})

		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "error", err.Error())
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("auth", claims)
		} else {
			utils.ErrorResponse(c, http.StatusUnauthorized, "error", "Invalid token")
			c.Abort()
			return
		}

		c.Next()
	}
}

func GetUser(c *gin.Context) (jwt.MapClaims, error) {
	auth, exists := c.Get("auth")
	if !exists {
		return nil, errors.New("auth key not found in context")
	}

	claims, ok := auth.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid auth claims type")
	}

	return claims, nil
}
