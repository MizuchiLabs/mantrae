// Package meta provides functionality for handling JWT claims.
package meta

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	UserID string `json:"user_id,omitempty"`
	jwt.RegisteredClaims
}

type AgentClaims struct {
	AgentID   string `json:"agent_id,omitempty"`
	ProfileID int64  `json:"profile_id,omitempty"`
	ServerURL string `json:"server_url,omitempty"`
	jwt.RegisteredClaims
}

func (u *UserClaims) Valid() error {
	if u.UserID == "" {
		return errors.New("user id is required")
	}
	return nil
}

func (u *UserClaims) IsExpired() bool {
	return u.ExpiresAt.Before(time.Now())
}

func DecodeUserToken(tokenStr string, secret string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&UserClaims{},
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		},
	)
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}

func EncodeUserToken(
	userID, secret string,
	expirationTime time.Time,
) (string, error) {
	claims := &UserClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	if err := claims.Valid(); err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func (a *AgentClaims) Valid() error {
	if a.AgentID == "" {
		return errors.New("agent id is required")
	}
	if a.ProfileID == 0 {
		return errors.New("profile id is required")
	}
	if a.ServerURL == "" {
		return errors.New("server url is required")
	}
	return nil
}

func (a *AgentClaims) IsExpired() bool {
	return a.ExpiresAt.Before(time.Now())
}

func DecodeAgentToken(tokenStr string, secret string) (*AgentClaims, error) {
	if tokenStr == "" {
		return nil, errors.New("token is required")
	}
	if secret == "" {
		claims := &AgentClaims{}
		parser := &jwt.Parser{}
		_, _, err := parser.ParseUnverified(tokenStr, claims)
		if err != nil {
			return nil, err
		}
		return claims, nil
	}
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&AgentClaims{},
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(secret), nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*AgentClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}

func EncodeAgentToken(
	profileID int64,
	agentID, serverURL, secret string,
	expirationTime time.Time,
) (string, error) {
	claims := &AgentClaims{
		AgentID:   agentID,
		ProfileID: profileID,
		ServerURL: serverURL,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	if err := claims.Valid(); err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
