package client

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	ServerURL string `json:"serverUrl,omitempty"`
	ProfileID int64  `json:"profileId,omitempty"`
	Secret    string `json:"secret,omitempty"`
	jwt.RegisteredClaims
}

func DecodeJWT(tokenString string) (*Claims, error) {
	// Decode without verifying to extract the secret
	claims := &Claims{}
	parser := &jwt.Parser{}
	_, _, err := parser.ParseUnverified(tokenString, claims)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JWT without verification: %w", err)
	}

	if claims.Secret == "" {
		return nil, fmt.Errorf("missing or invalid secret in token claims")
	}

	// Re-parse the token with signature verification
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(claims.Secret), nil
		},
	)
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("failed to decode and verify JWT: %w", err)
	}

	return claims, nil
}
