package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTUtils struct {
	secret []byte
}

func NewJWTUtils() *JWTUtils {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "YOU_MIGHT_WANT_TO_IMPLEMENT_THIS_BRO"
	}
	return &JWTUtils{secret: []byte(secret)}
}

type Claims struct {
	ID string `json:"id"`
	jwt.RegisteredClaims
}

func (j *JWTUtils) GenerateJWT(id uint) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour)

	claims := &Claims{
		ID: fmt.Sprintf("%d", id),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   "user_id",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(j.secret)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	return tokenString, nil
}

func (j *JWTUtils) VerifyToken(token string) (*Claims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return j.secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if claims, ok := parsedToken.Claims.(*Claims); ok && parsedToken.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
