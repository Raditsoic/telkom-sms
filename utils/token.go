package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret []byte

func init() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "YOU_MIGHT_WANT_TO_IMPLEMENT_THIS_BRO"
	}
	jwtSecret = []byte(secret)
}

type Claims struct {
	ID   string `json:"id"`
	jwt.RegisteredClaims
}

func GenerateJWT(id uint) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour)

	claims := &Claims{
		ID:   fmt.Sprintf("%d", id),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   "user_id",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}
	return tokenString, nil
}

func VerifyToken(token string) (*Claims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	if claims, ok := parsedToken.Claims.(*Claims); ok && parsedToken.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token: %v", err)
}
