package service

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"connectrpc.com/connect"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mizuchilabs/mantrae/internal/api/middlewares"
	"github.com/mizuchilabs/mantrae/internal/config"
	"github.com/mizuchilabs/mantrae/internal/db"
	"github.com/mizuchilabs/mantrae/internal/settings"
	"github.com/mizuchilabs/mantrae/internal/traefik"
	"github.com/mizuchilabs/mantrae/internal/util"
	"github.com/mizuchilabs/mantrae/pkg/meta"
	agentv1 "github.com/mizuchilabs/mantrae/proto/gen/agent/v1"
)

type AgentService struct {
	app *config.App
}

func NewAgentService(app *config.App) *AgentService {
	return &AgentService{app: app}
}

func (s *AgentService) HealthCheck(
	ctx context.Context,
	req *connect.Request[agentv1.HealthCheckRequest],
) (*connect.Response[agentv1.HealthCheckResponse], error) {
	// Rotate Token
	token, err := s.updateToken(ctx, req.Header().Get(meta.HeaderAgentID))
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	util.Broadcast <- util.EventMessage{
		Type:     util.EventTypeUpdate,
		Category: util.EventCategoryAgent,
	}
	return connect.NewResponse(&agentv1.HealthCheckResponse{Ok: true, Token: *token}), nil
}

func (s *AgentService) GetContainer(
	ctx context.Context,
	req *connect.Request[agentv1.GetContainerRequest],
) (*connect.Response[agentv1.GetContainerResponse], error) {
	agent := middlewares.GetAgentContext(ctx)
	if agent == nil {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("agent context missing"),
		)
	}

	// Upsert agent
	params := db.UpdateAgentParams{
		ID:       agent.ID,
		Hostname: &req.Msg.Hostname,
		PublicIp: &req.Msg.PublicIp,
	}
	if agent.ActiveIp == nil {
		params.ActiveIp = &req.Msg.PublicIp
	}

	privateIPs := db.AgentPrivateIPs{IPs: make([]string, len(req.Msg.PrivateIps))}
	privateIPs.IPs = req.Msg.PrivateIps
	params.PrivateIps = &privateIPs

	var containers db.AgentContainers
	for _, container := range req.Msg.Containers {
		containers = append(containers, db.AgentContainer{
			ID:      container.Id,
			Name:    container.Name,
			Labels:  container.Labels,
			Image:   container.Image,
			Portmap: container.Portmap,
			Status:  container.Status,
			Created: container.Created.AsTime(),
		})
	}
	params.Containers = &containers

	q := s.app.Conn.GetQuery()
	updatedAgent, err := q.UpdateAgent(ctx, params)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// Update agent config
	if err = traefik.DecodeAgentConfig(s.app.Conn.Get(), updatedAgent); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	util.Broadcast <- util.EventMessage{
		Type:     util.EventTypeUpdate,
		Category: util.EventCategoryAgent,
	}
	return connect.NewResponse(&agentv1.GetContainerResponse{}), nil
}

func (s *AgentService) updateToken(ctx context.Context, id string) (*string, error) {
	q := s.app.Conn.GetQuery()
	agent, err := q.GetAgent(ctx, id)
	if err != nil {
		return nil, err
	}

	claims, err := DecodeJWT(agent.Token, s.app.Config.Secret)
	if err != nil {
		return nil, err
	}

	// Only update the token if it's close to expiring (less than 25%)
	lifetime := claims.ExpiresAt.Sub(claims.IssuedAt.Time)
	remaining := time.Until(claims.ExpiresAt.Time)
	if remaining > lifetime/4 {
		return &agent.Token, nil // Still valid
	}

	agentInterval, err := s.app.SM.Get(ctx, settings.KeyAgentCleanupInterval)
	if err != nil {
		return nil, err
	}

	token, err := claims.EncodeJWT(s.app.Config.Secret, agentInterval.Duration(time.Hour*72))
	if err != nil {
		return nil, err
	}

	err = q.UpdateAgentToken(ctx, db.UpdateAgentTokenParams{ID: agent.ID, Token: token})
	if err != nil {
		return nil, err
	}
	slog.Info("Rotating agent token", "agentID", agent.ID, "token", token)

	return &token, nil
}

// Helpers --------------------------------------------------------------------
type AgentClaims struct {
	AgentID   string `json:"agentId,omitempty"`
	ProfileID int64  `json:"profileId,omitempty"`
	ServerURL string `json:"serverUrl,omitempty"`
	jwt.RegisteredClaims
}

// EncodeJWT generates a JWT for agents
func (a *AgentClaims) EncodeJWT(secret string, expirationTime time.Duration) (string, error) {
	if a.ServerURL == "" || a.ProfileID == 0 {
		return "", errors.New("serverUrl and profileID cannot be empty")
	}

	if expirationTime == 0 {
		expirationTime = time.Hour * 24
	}

	claims := &AgentClaims{
		AgentID:   a.AgentID,
		ProfileID: a.ProfileID,
		ServerURL: a.ServerURL,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expirationTime)),
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
