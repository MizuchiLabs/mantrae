package util

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const CookieName = "auth_token"

type UserClaims struct {
	UserID string `json:"user_id,omitempty"`
	jwt.RegisteredClaims
}

// EncodeUserJWT generates a JWT for user login
func EncodeUserJWT(userID, secret string, expirationTime time.Time) (string, error) {
	if userID == "" {
		return "", errors.New("username cannot be empty")
	}
	if expirationTime.IsZero() {
		expirationTime = time.Now().Add(24 * time.Hour)
	}
	claims := &UserClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// DecodeUserJWT decodes the user token and returns claims if valid
func DecodeUserJWT(tokenString, secret string) (*UserClaims, error) {
	claims := &UserClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (any, error) {
			return []byte(secret), nil
		},
	)

	if err != nil || !token.Valid {
		return nil, err
	}
	return claims, nil
}
