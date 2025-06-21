package util

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// EncodeJWT generates a token string from a claims struct
func EncodeJWT[T jwt.Claims](claims T, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// DecodeJWT parses and validates a JWT string
func DecodeJWT[T jwt.Claims](tokenStr, secret string) (T, error) {
	var claims T

	token, err := jwt.ParseWithClaims(
		tokenStr,
		claims,
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		},
	)
	if err != nil {
		return claims, err
	}

	typedClaims, ok := token.Claims.(T)
	if !ok || !token.Valid {
		return claims, errors.New("invalid token")
	}

	return typedClaims, nil
}

func DecodeUnsafeJWT[T jwt.Claims](tokenStr string) (T, error) {
	var claims T
	parser := &jwt.Parser{}
	_, _, err := parser.ParseUnverified(tokenStr, claims)
	if err != nil {
		return claims, fmt.Errorf("failed to decode JWT without verification: %w", err)
	}

	return claims, nil
}
