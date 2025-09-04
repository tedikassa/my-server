package utils

import (
	"time"

	"example.com/ecomerce/model"
	"github.com/golang-jwt/jwt/v5"
)

var JwtKey = []byte("your-secret-key")

// GenerateToken creates JWT token for user
func GenerateToken(id int, name, role string) (string, error) {
	expiration := time.Now().Add(24 * time.Hour)

	claims := &model.Claims{
		ID:   id,
		Name: name,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
		},
	}
	// create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// sign with secret key
	return token.SignedString(JwtKey)
}


