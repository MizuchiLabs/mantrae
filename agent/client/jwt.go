package client

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	AgentID   string `json:"agentId,omitempty"`
	ProfileID int64  `json:"profileId,omitempty"`
	ServerURL string `json:"serverUrl,omitempty"`
	jwt.RegisteredClaims
}

func DecodeJWT(tokenString string) (*Claims, error) {
	// Decode without verifying
	claims := &Claims{}
	parser := &jwt.Parser{}
	_, _, err := parser.ParseUnverified(tokenString, claims)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JWT without verification: %w", err)
	}

	return claims, nil
}
