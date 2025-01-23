package agent

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AgentClaims struct {
	ServerURL string `json:"serverUrl,omitempty"`
	ProfileID int64  `json:"profileId,omitempty"`
	jwt.RegisteredClaims
}

// EncodeJWT generates a JWT for agents
func EncodeJWT(
	serverurl string,
	profileid int64,
	secret string,
) (string, error) {
	if serverurl == "" || profileid == 0 {
		return "", errors.New("serverUrl and profileID cannot be empty")
	}
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &AgentClaims{
		ServerURL: serverurl,
		ProfileID: profileid,
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
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		},
	)

	if err != nil || !token.Valid {
		return nil, err
	}
	return claims, nil
}
