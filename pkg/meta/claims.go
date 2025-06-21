package meta

import (
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
