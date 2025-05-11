package agent

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"connectrpc.com/connect"

	agentv1 "github.com/MizuchiLabs/mantrae/agent/proto/gen/agent/v1"
	"github.com/MizuchiLabs/mantrae/internal/api/middlewares"
	"github.com/MizuchiLabs/mantrae/internal/config"
	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/MizuchiLabs/mantrae/internal/settings"
	"github.com/MizuchiLabs/mantrae/internal/traefik"
	"github.com/MizuchiLabs/mantrae/internal/util"
	"github.com/MizuchiLabs/mantrae/pkg/meta"
)

type AgentServer struct {
	app *config.App
}

func NewAgentServer(app *config.App) *AgentServer {
	return &AgentServer{app: app}
}

func (s *AgentServer) HealthCheck(
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

func (s *AgentServer) GetContainer(
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

func (s *AgentServer) updateToken(ctx context.Context, id string) (*string, error) {
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
	lifetime := claims.ExpiresAt.Time.Sub(claims.IssuedAt.Time)
	remaining := claims.ExpiresAt.Time.Sub(time.Now())
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
