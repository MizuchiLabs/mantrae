package agent

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AgentClaims struct {
	AgentID   string `json:"agentId,omitempty"`
	ProfileID int64  `json:"profileId,omitempty"`
	ServerURL string `json:"serverUrl,omitempty"`
	jwt.RegisteredClaims
}

// EncodeJWT generates a JWT for agents
func (a *AgentClaims) EncodeJWT(secret string, expirationTime time.Time) (string, error) {
	if a.ServerURL == "" || a.ProfileID == 0 {
		return "", errors.New("serverUrl and profileID cannot be empty")
	}
	claims := &AgentClaims{
		AgentID:   a.AgentID,
		ProfileID: a.ProfileID,
		ServerURL: a.ServerURL,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// DecodeJWT decodes the agent token and returns claims if valid
func DecodeJWT(tokenString, secret string) (*AgentClaims, error) {
	claims := &AgentClaims{}
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
