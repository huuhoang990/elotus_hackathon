package common

import (
	"errors"
	"go_api/cmd/internal/models"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomJWTClaims struct {
	ID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateJWTToken(user models.User) (string, error) {
	userClaims := CustomJWTClaims{
		ID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Token valid for 24 hours
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)
	signedAccessToken, err := accessToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}
	return signedAccessToken, nil
}

func ParseJWTAccessToken(signedAccessToken string) (*CustomJWTClaims, error) {
	parsedJwtAccessToken, err := jwt.ParseWithClaims(signedAccessToken, &CustomJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, err
	} else if claims, ok := parsedJwtAccessToken.Claims.(*CustomJWTClaims); ok {
		return claims, nil
	} else {
		return nil, errors.New("invalid JWT access token")
	}
}

func IsClaimExpired(claims *CustomJWTClaims) bool {
	currentTime := jwt.NewNumericDate(time.Now())
	return claims.ExpiresAt.Time.Before(currentTime.Time)
}
