package middleware

import (
	"github.com/golang-jwt/jwt"
	"time"
)

var jwtKey = []byte("kitty-go-123")

func GenerateJWT(userID int, userRole string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID:   userID,
		UserRole: userRole,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	return tokenString, err
}

type Claims struct {
	UserID   int    `json:"user_id"`
	UserRole string `json:"user_role"`
	jwt.StandardClaims
}
