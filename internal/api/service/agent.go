package service

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"connectrpc.com/connect"

	"github.com/google/uuid"
	"github.com/mizuchilabs/mantrae/internal/config"
	"github.com/mizuchilabs/mantrae/internal/settings"
	"github.com/mizuchilabs/mantrae/internal/store/db"
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

	serverURL, ok := s.app.SM.Get(ctx, settings.KeyServerURL)
	if !ok || serverURL == "" {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("server url is required, check your settings"),
		)
	}

	token, err := s.createToken(params.ID, req.Msg.ProfileId)
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

func (s *AgentService) UpdateAgentIP(
	ctx context.Context,
	req *connect.Request[mantraev1.UpdateAgentIPRequest],
) (*connect.Response[mantraev1.UpdateAgentIPResponse], error) {
	params := db.UpdateAgentIPParams{
		ID:       req.Msg.Id,
		ActiveIp: &req.Msg.Ip,
	}

	if err := s.app.Conn.GetQuery().UpdateAgentIP(ctx, params); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	result, err := s.app.Conn.GetQuery().GetAgent(ctx, req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&mantraev1.UpdateAgentIPResponse{
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

func (s *AgentService) RotateAgentToken(
	ctx context.Context,
	req *connect.Request[mantraev1.RotateAgentTokenRequest],
) (*connect.Response[mantraev1.RotateAgentTokenResponse], error) {
	result, err := s.app.Conn.GetQuery().GetAgent(ctx, req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	token, err := s.createToken(result.ID, result.ProfileID)
	if err != nil {
		return nil, err
	}
	if err = s.app.Conn.GetQuery().UpdateAgentToken(ctx, db.UpdateAgentTokenParams{
		ID:    result.ID,
		Token: *token,
	}); err != nil {
		return nil, err
	}

	return connect.NewResponse(&mantraev1.RotateAgentTokenResponse{
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

	if err = s.app.Conn.GetQuery().UpdateAgentToken(ctx, db.UpdateAgentTokenParams{
		ID:    agent.ID,
		Token: *token,
	}); err != nil {
		return nil, err
	}

	slog.Info("Rotating agent token", "agentID", agent.ID, "token", token)
	return token, nil
}

func (s *AgentService) createToken(agentID string, profileID int64) (*string, error) {
	sets := s.app.SM.GetMany(
		context.Background(),
		[]string{settings.KeyServerURL, settings.KeyAgentCleanupInterval},
	)

	token, err := meta.EncodeAgentToken(
		profileID,
		agentID,
		sets[settings.KeyServerURL],
		s.app.Secret,
		time.Now().Add(settings.AsDuration(sets[settings.KeyAgentCleanupInterval])),
	)
	if err != nil {
		return nil, err
	}

	return &token, nil
}
