package token

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// TokenClaims struct untuk JWT
type TokenClaims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateAccessToken berlaku 60 menit

func GenerateAccessToken(userID uint, role string) (string, error) {
	claims := &TokenClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(60 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("SECRET_KEY")
	fmt.Println("SECRET_KEY GENERATE:", os.Getenv("SECRET_KEY"))
	return token.SignedString([]byte(secret))
}

// GenerateRefreshToken berlaku 7 hari
func GenerateRefreshToken(userID uint, role string) (string, error) {
	claims := &TokenClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("REFRESH_SECRET")
	return token.SignedString([]byte(secret))
}

// ValidateAccessToken: cek access token
func ValidateAccessToken(tokenString string) (*TokenClaims, error) {
	secret := os.Getenv("SECRET_KEY")
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid access token")
	}

	return claims, nil
}

// ValidateRefreshToken: cek refresh token
func ValidateRefreshToken(tokenString string) (*TokenClaims, error) {
	secret := os.Getenv("REFRESH_SECRET")
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	return claims, nil
}
