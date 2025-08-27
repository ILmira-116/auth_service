package jwt

import (
	"auth-service/internal/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func NewToken(user model.User, app model.App, secret string, ttl time.Duration) (string, error) {

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(ttl).Unix(), // срок действия токена
		"iat":     time.Now().Unix(),          // время выпуска
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
