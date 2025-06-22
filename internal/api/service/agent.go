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
	agent, err := s.app.Conn.GetQuery().GetAgent(ctx, req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// var containers []*mantraev1.Container
	// err = json.Unmarshal(containers, &agent.Containers)
	// if err != nil {
	// 	return nil, connect.NewError(connect.CodeInternal, err)
	// }

	return connect.NewResponse(&mantraev1.GetAgentResponse{
		Agent: &mantraev1.Agent{
			Id:        agent.ID,
			ProfileId: agent.ProfileID,
			Hostname:  SafeString(agent.Hostname),
			PublicIp:  SafeString(agent.PublicIp),
			PrivateIp: SafeString(agent.PrivateIp),
			ActiveIp:  SafeString(agent.ActiveIp),
			Token:     agent.Token,
			// Containers: containers,
			CreatedAt: SafeTimestamp(agent.CreatedAt),
			UpdatedAt: SafeTimestamp(agent.UpdatedAt),
		},
	}), nil
}

func (s *AgentService) CreateAgent(
	ctx context.Context,
	req *connect.Request[mantraev1.CreateAgentRequest],
) (*connect.Response[mantraev1.CreateAgentResponse], error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	serverUrl, err := s.app.Conn.GetQuery().GetSetting(ctx, settings.KeyServerURL)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	if serverUrl.Value == "" {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("server url is required, check your settings"),
		)
	}

	token, err := s.createToken(id.String(), req.Msg.ProfileId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	agent, err := s.app.Conn.GetQuery().CreateAgent(ctx, db.CreateAgentParams{
		ID:        id.String(),
		ProfileID: req.Msg.ProfileId,
		Token:     *token,
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&mantraev1.CreateAgentResponse{
		Agent: &mantraev1.Agent{
			Id:        agent.ID,
			ProfileId: agent.ProfileID,
			Token:     agent.Token,
		},
	}), nil
}

func (s *AgentService) UpdateAgentIP(
	ctx context.Context,
	req *connect.Request[mantraev1.UpdateAgentIPRequest],
) (*connect.Response[mantraev1.UpdateAgentIPResponse], error) {
	if err := s.app.Conn.GetQuery().UpdateAgentIP(ctx, db.UpdateAgentIPParams{
		ID:       req.Msg.Id,
		ActiveIp: &req.Msg.Ip,
	}); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	agent, err := s.app.Conn.GetQuery().GetAgent(ctx, req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&mantraev1.UpdateAgentIPResponse{
		Agent: &mantraev1.Agent{
			Id:        agent.ID,
			ProfileId: agent.ProfileID,
			Hostname:  SafeString(agent.Hostname),
			PublicIp:  SafeString(agent.PublicIp),
			PrivateIp: SafeString(agent.PrivateIp),
			ActiveIp:  SafeString(agent.ActiveIp),
			Token:     agent.Token,
			CreatedAt: SafeTimestamp(agent.CreatedAt),
			UpdatedAt: SafeTimestamp(agent.UpdatedAt),
		},
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
	var params db.ListAgentsParams
	params.ProfileID = req.Msg.ProfileId
	if req.Msg.Limit == nil {
		params.Limit = 100
	} else {
		params.Limit = *req.Msg.Limit
	}
	if req.Msg.Offset == nil {
		params.Offset = 0
	} else {
		params.Offset = *req.Msg.Offset
	}

	dbAgents, err := s.app.Conn.GetQuery().ListAgents(ctx, params)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	totalCount, err := s.app.Conn.GetQuery().CountAgents(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	var agents []*mantraev1.Agent
	for _, agent := range dbAgents {
		agents = append(agents, &mantraev1.Agent{
			Id:        agent.ID,
			ProfileId: agent.ProfileID,
			Hostname:  SafeString(agent.Hostname),
			PublicIp:  SafeString(agent.PublicIp),
			PrivateIp: SafeString(agent.PrivateIp),
			ActiveIp:  SafeString(agent.ActiveIp),
			Token:     agent.Token,
			CreatedAt: SafeTimestamp(agent.CreatedAt),
			UpdatedAt: SafeTimestamp(agent.UpdatedAt),
		})
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
	if _, err := s.updateToken(ctx, &agent); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// Update Agent
	var params db.UpdateAgentParams
	params.ID = agent.ID
	// if req.Msg.MachineId != "" {
	// 	params.MachineId = &req.Msg.MachineId
	// }
	if req.Msg.Hostname != "" {
		params.Hostname = &req.Msg.Hostname
	}
	if req.Msg.PublicIp != "" {
		params.PublicIp = &req.Msg.PublicIp
	}
	if req.Msg.PrivateIp != "" {
		params.PrivateIp = &req.Msg.PrivateIp
	}

	agentNew, err := s.app.Conn.GetQuery().UpdateAgent(ctx, params)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&mantraev1.HealthCheckResponse{Agent: &mantraev1.Agent{
		Id:        agentNew.ID,
		ProfileId: agentNew.ProfileID,
		ActiveIp:  SafeString(agentNew.ActiveIp),
		Token:     agentNew.Token,
	}}), nil
}

func (s *AgentService) RotateAgentToken(
	ctx context.Context,
	req *connect.Request[mantraev1.RotateAgentTokenRequest],
) (*connect.Response[mantraev1.RotateAgentTokenResponse], error) {
	agent, err := s.app.Conn.GetQuery().GetAgent(ctx, req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	_, err = s.updateToken(ctx, &agent)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&mantraev1.RotateAgentTokenResponse{}), nil
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
	serverUrl, ok := s.app.SM.Get(settings.KeyServerURL)
	if !ok {
		return nil, errors.New("failed to get server url setting")
	}

	agentInterval, ok := s.app.SM.Get(settings.KeyAgentCleanupInterval)
	if !ok {
		return nil, errors.New("failed to get agent cleanup interval setting")
	}

	token, err := meta.EncodeAgentToken(
		profileID,
		agentID,
		serverUrl,
		s.app.Secret,
		time.Now().Add(settings.AsDuration(agentInterval)),
	)
	if err != nil {
		return nil, err
	}

	return &token, nil
}
