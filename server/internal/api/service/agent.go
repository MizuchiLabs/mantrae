package service

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"connectrpc.com/connect"

	"github.com/google/uuid"
	"github.com/mizuchilabs/mantrae/server/internal/config"
	"github.com/mizuchilabs/mantrae/server/internal/settings"
	"github.com/mizuchilabs/mantrae/server/internal/store/db"
	"github.com/mizuchilabs/mantrae/pkg/meta"
	mantraev1 "github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1"
)

type AgentService struct {
	app *config.App
}

func NewAgentService(app *config.App) *AgentService {
	return &AgentService{app: app}
}

func (s *AgentService) GetAgent(
	ctx context.Context,
	req *connect.Request[mantraev1.GetAgentRequest],
) (*connect.Response[mantraev1.GetAgentResponse], error) {
	result, err := s.app.Conn.GetQuery().GetAgent(ctx, req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&mantraev1.GetAgentResponse{
		Agent: result.ToProto(),
	}), nil
}

func (s *AgentService) CreateAgent(
	ctx context.Context,
	req *connect.Request[mantraev1.CreateAgentRequest],
) (*connect.Response[mantraev1.CreateAgentResponse], error) {
	params := db.CreateAgentParams{
		ID:        uuid.NewString(),
		ProfileID: req.Msg.ProfileId,
	}

	token, err := s.createToken(params.ID, params.ProfileID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	params.Token = *token

	result, err := s.app.Conn.GetQuery().CreateAgent(ctx, params)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&mantraev1.CreateAgentResponse{
		Agent: result.ToProto(),
	}), nil
}

func (s *AgentService) UpdateAgent(
	ctx context.Context,
	req *connect.Request[mantraev1.UpdateAgentRequest],
) (*connect.Response[mantraev1.UpdateAgentResponse], error) {
	params := db.UpdateAgentParams{
		ID:       req.Msg.Id,
		ActiveIp: req.Msg.Ip,
	}

	if req.Msg.RotateToken != nil && *req.Msg.RotateToken {
		agent, err := s.app.Conn.GetQuery().GetAgent(ctx, params.ID)
		if err != nil {
			return nil, err
		}
		token, err := s.createToken(params.ID, agent.ProfileID)
		if err != nil {
			return nil, err
		}
		params.Token = token
	}

	result, err := s.app.Conn.GetQuery().UpdateAgent(ctx, params)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&mantraev1.UpdateAgentResponse{
		Agent: result.ToProto(),
	}), nil
}

func (s *AgentService) DeleteAgent(
	ctx context.Context,
	req *connect.Request[mantraev1.DeleteAgentRequest],
) (*connect.Response[mantraev1.DeleteAgentResponse], error) {
	if err := s.app.Conn.GetQuery().DeleteAgent(ctx, req.Msg.Id); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&mantraev1.DeleteAgentResponse{}), nil
}

func (s *AgentService) ListAgents(
	ctx context.Context,
	req *connect.Request[mantraev1.ListAgentsRequest],
) (*connect.Response[mantraev1.ListAgentsResponse], error) {
	params := db.ListAgentsParams{
		ProfileID: req.Msg.ProfileId,
		Limit:     req.Msg.Limit,
		Offset:    req.Msg.Offset,
	}

	result, err := s.app.Conn.GetQuery().ListAgents(ctx, params)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	totalCount, err := s.app.Conn.GetQuery().CountAgents(ctx, req.Msg.ProfileId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	agents := make([]*mantraev1.Agent, 0, len(result))
	for _, a := range result {
		agents = append(agents, a.ToProto())
	}
	return connect.NewResponse(&mantraev1.ListAgentsResponse{
		Agents:     agents,
		TotalCount: totalCount,
	}), nil
}

func (s *AgentService) HealthCheck(
	ctx context.Context,
	req *connect.Request[mantraev1.HealthCheckRequest],
) (*connect.Response[mantraev1.HealthCheckResponse], error) {
	agentID := req.Header().Get(meta.HeaderAgentID)
	if agentID == "" {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("agent id is required"),
		)
	}
	agent, err := s.app.Conn.GetQuery().GetAgent(ctx, agentID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// Rotate Token if it's close to expiring
	if _, err = s.updateToken(ctx, &agent); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// Update Agent
	params := db.UpdateAgentParams{
		ID:        agent.ID,
		Hostname:  &req.Msg.Hostname,
		PublicIp:  &req.Msg.PublicIp,
		PrivateIp: &req.Msg.PrivateIp,
	}

	// Update ActiveIp if it's not set
	if agent.ActiveIp == nil && req.Msg.PrivateIp != "" {
		params.ActiveIp = &req.Msg.PrivateIp
	}

	result, err := s.app.Conn.GetQuery().UpdateAgent(ctx, params)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&mantraev1.HealthCheckResponse{
		Agent: result.ToProto(),
	}), nil
}

func (s *AgentService) updateToken(ctx context.Context, agent *db.Agent) (*string, error) {
	claims, err := meta.DecodeAgentToken(agent.Token, s.app.Secret)
	if err != nil {
		return nil, err
	}

	// Only update the token if it's close to expiring (less than 25%)
	lifetime := claims.ExpiresAt.Sub(claims.IssuedAt.Time)
	remaining := time.Until(claims.ExpiresAt.Time)
	if remaining > lifetime/4 {
		return &agent.Token, nil // Token is still valid
	}

	token, err := s.createToken(agent.ID, agent.ProfileID)
	if err != nil {
		return nil, err
	}

	if _, err = s.app.Conn.GetQuery().UpdateAgent(ctx, db.UpdateAgentParams{
		ID:    agent.ID,
		Token: token,
	}); err != nil {
		return nil, err
	}

	slog.Info("Rotating agent token", "agentID", agent.ID, "token", token)
	return token, nil
}

func (s *AgentService) createToken(agentID string, profileID int64) (*string, error) {
	serverURL, ok := s.app.SM.Get(context.Background(), settings.KeyServerURL)
	if !ok || serverURL == "" {
		return nil, errors.New("server url is required, check your settings")
	}
	cleanupInterval, ok := s.app.SM.Get(context.Background(), settings.KeyAgentCleanupInterval)
	if !ok {
		return nil, errors.New("agent cleanup interval is required, check your settings")
	}

	token, err := meta.EncodeAgentToken(
		profileID,
		agentID,
		serverURL,
		s.app.Secret,
		time.Now().Add(settings.AsDuration(cleanupInterval)),
	)
	if err != nil {
		return nil, err
	}

	return &token, nil
}
