package service

import (
	"context"
	"errors"
	"time"

	"connectrpc.com/connect"

	"github.com/google/uuid"
	"github.com/mizuchilabs/mantrae/internal/config"
	"github.com/mizuchilabs/mantrae/internal/settings"
	"github.com/mizuchilabs/mantrae/internal/store/db"
	mantraev1 "github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1"
)

type AgentManagementService struct {
	app *config.App
}

func NewAgentManagementService(app *config.App) *AgentManagementService {
	return &AgentManagementService{app: app}
}

func (s *AgentManagementService) GetAgent(
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
			Id:         agent.ID,
			ProfileId:  agent.ProfileID,
			Hostname:   SafeString(agent.Hostname),
			PublicIp:   SafeString(agent.PublicIp),
			ActiveIp:   SafeString(agent.ActiveIp),
			Token:      agent.Token,
			PrivateIps: agent.PrivateIps.IPs,
			// Containers: containers,
			CreatedAt: SafeTimestamp(agent.CreatedAt),
			UpdatedAt: SafeTimestamp(agent.UpdatedAt),
		},
	}), nil
}

func (s *AgentManagementService) CreateAgent(
	ctx context.Context,
	req *connect.Request[mantraev1.CreateAgentRequest],
) (*connect.Response[mantraev1.CreateAgentResponse], error) {
	if req.Msg.ProfileId == 0 {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("profile id is required"),
		)
	}

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

	claims := &AgentClaims{
		AgentID:   id.String(),
		ProfileID: req.Msg.ProfileId,
		ServerURL: serverUrl.Value,
	}
	token, err := claims.EncodeJWT(s.app.Secret, time.Hour*72)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	agent, err := s.app.Conn.GetQuery().CreateAgent(ctx, db.CreateAgentParams{
		ID:        claims.AgentID,
		ProfileID: claims.ProfileID,
		Token:     token,
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

func (s *AgentManagementService) UpdateAgentIP(
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
			Id:         agent.ID,
			ProfileId:  agent.ProfileID,
			Hostname:   SafeString(agent.Hostname),
			PublicIp:   SafeString(agent.PublicIp),
			ActiveIp:   SafeString(agent.ActiveIp),
			Token:      agent.Token,
			PrivateIps: agent.PrivateIps.IPs,
			CreatedAt:  SafeTimestamp(agent.CreatedAt),
			UpdatedAt:  SafeTimestamp(agent.UpdatedAt),
		},
	}), nil
}

func (s *AgentManagementService) DeleteAgent(
	ctx context.Context,
	req *connect.Request[mantraev1.DeleteAgentRequest],
) (*connect.Response[mantraev1.DeleteAgentResponse], error) {
	if err := s.app.Conn.GetQuery().DeleteAgent(ctx, req.Msg.Id); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&mantraev1.DeleteAgentResponse{}), nil
}

func (s *AgentManagementService) ListAgents(
	ctx context.Context,
	req *connect.Request[mantraev1.ListAgentsRequest],
) (*connect.Response[mantraev1.ListAgentsResponse], error) {
	if req.Msg.ProfileId == 0 {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("profile id is required"),
		)
	}

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
			ActiveIp:  SafeString(agent.ActiveIp),
			Token:     agent.Token,
			CreatedAt: SafeTimestamp(agent.CreatedAt),
			UpdatedAt: SafeTimestamp(agent.UpdatedAt),
		})
		if agent.PrivateIps != nil {
			for _, ip := range agent.PrivateIps.IPs {
				agents[len(agents)-1].PrivateIps = append(agents[len(agents)-1].PrivateIps, ip)
			}
		}
	}
	return connect.NewResponse(&mantraev1.ListAgentsResponse{
		Agents:     agents,
		TotalCount: totalCount,
	}), nil
}

func (s *AgentManagementService) RotateAgentToken(
	ctx context.Context,
	req *connect.Request[mantraev1.RotateAgentTokenRequest],
) (*connect.Response[mantraev1.RotateAgentTokenResponse], error) {
	agent, err := s.app.Conn.GetQuery().GetAgent(ctx, req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	serverUrl, err := s.app.Conn.GetQuery().GetSetting(ctx, settings.KeyServerURL)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	claims := &AgentClaims{
		AgentID:   agent.ID,
		ProfileID: agent.ProfileID,
		ServerURL: serverUrl.Value,
	}
	token, err := claims.EncodeJWT(s.app.Secret, time.Hour*72)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if err := s.app.Conn.GetQuery().UpdateAgentToken(ctx, db.UpdateAgentTokenParams{
		ID:    agent.ID,
		Token: token,
	}); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&mantraev1.RotateAgentTokenResponse{}), nil
}

func (s *AgentManagementService) BootstrapAgent(
	ctx context.Context,
	req *connect.Request[mantraev1.BootstrapAgentRequest],
) (*connect.Response[mantraev1.BootstrapAgentResponse], error) {
	enabled, ok := s.app.SM.Get(settings.KeyAgentBootstrapEnabled)
	if !ok {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("failed to get agent bootstrap enabled setting"),
		)
	}
	if enabled != "true" {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("agent bootstrap is disabled, check your settings"),
		)
	}
	if req.Msg.Token == "" {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("token is required"),
		)
	}

	// Check if token is valid
	bootstrapToken, ok := s.app.SM.Get(settings.KeyAgentBootstrapToken)
	if !ok {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("failed to get agent bootstrap token setting"),
		)
	}
	if bootstrapToken != req.Msg.Token {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("invalid token"),
		)
	}

	// Toke is valid, create agent
	agent, err := s.app.Conn.GetQuery().CreateAgent(ctx, db.CreateAgentParams{
		ProfileID: req.Msg.ProfileId,
	})
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

	claims := &AgentClaims{
		AgentID:   agent.ID,
		ProfileID: agent.ProfileID,
		ServerURL: serverUrl.Value,
	}
	token, err := claims.EncodeJWT(s.app.Secret, time.Hour*72)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if err := s.app.Conn.GetQuery().UpdateAgentToken(ctx, db.UpdateAgentTokenParams{
		ID:    agent.ID,
		Token: token,
	}); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&mantraev1.BootstrapAgentResponse{Token: token}), nil
}
