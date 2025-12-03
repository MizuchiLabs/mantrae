package service

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"connectrpc.com/connect"

	"github.com/google/uuid"
	"github.com/mizuchilabs/mantrae/pkg/meta"
	"github.com/mizuchilabs/mantrae/pkg/util"
	mantraev1 "github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1"
	"github.com/mizuchilabs/mantrae/server/internal/config"
	"github.com/mizuchilabs/mantrae/server/internal/store/db"
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
	params := &db.CreateAgentParams{
		ID:        uuid.NewString(),
		ProfileID: req.Msg.ProfileId,
	}

	params.Token = util.GenerateAgentToken(
		strconv.Itoa(int(params.ProfileID)),
		params.ID,
	)

	result, err := s.app.Conn.GetQuery().CreateAgent(ctx, params)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	s.app.Event.Broadcast(&mantraev1.EventStreamResponse{
		Action: mantraev1.EventAction_EVENT_ACTION_CREATED,
		Data: &mantraev1.EventStreamResponse_Agent{
			Agent: result.ToProto(),
		},
	})
	return connect.NewResponse(&mantraev1.CreateAgentResponse{
		Agent: result.ToProto(),
	}), nil
}

func (s *AgentService) UpdateAgent(
	ctx context.Context,
	req *connect.Request[mantraev1.UpdateAgentRequest],
) (*connect.Response[mantraev1.UpdateAgentResponse], error) {
	params := &db.UpdateAgentParams{
		ID:       req.Msg.Id,
		ActiveIp: req.Msg.Ip,
	}

	if req.Msg.RotateToken != nil && *req.Msg.RotateToken {
		agent, err := s.app.Conn.GetQuery().GetAgent(ctx, params.ID)
		if err != nil {
			return nil, err
		}
		token := util.GenerateAgentToken(strconv.Itoa(int(agent.ProfileID)), agent.ID)
		params.Token = &token
	}

	result, err := s.app.Conn.GetQuery().UpdateAgent(ctx, params)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	s.app.Event.Broadcast(&mantraev1.EventStreamResponse{
		Action: mantraev1.EventAction_EVENT_ACTION_UPDATED,
		Data: &mantraev1.EventStreamResponse_Agent{
			Agent: result.ToProto(),
		},
	})
	return connect.NewResponse(&mantraev1.UpdateAgentResponse{
		Agent: result.ToProto(),
	}), nil
}

func (s *AgentService) DeleteAgent(
	ctx context.Context,
	req *connect.Request[mantraev1.DeleteAgentRequest],
) (*connect.Response[mantraev1.DeleteAgentResponse], error) {
	agent, err := s.app.Conn.GetQuery().GetAgent(ctx, req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if err := s.app.Conn.GetQuery().DeleteAgent(ctx, req.Msg.Id); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	s.app.Event.Broadcast(&mantraev1.EventStreamResponse{
		Action: mantraev1.EventAction_EVENT_ACTION_DELETED,
		Data: &mantraev1.EventStreamResponse_Agent{
			Agent: agent.ToProto(),
		},
	})
	return connect.NewResponse(&mantraev1.DeleteAgentResponse{}), nil
}

func (s *AgentService) ListAgents(
	ctx context.Context,
	req *connect.Request[mantraev1.ListAgentsRequest],
) (*connect.Response[mantraev1.ListAgentsResponse], error) {
	params := &db.ListAgentsParams{
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

	// Update Agent
	params := &db.UpdateAgentParams{
		ID:        agent.ID,
		Hostname:  &req.Msg.Hostname,
		PublicIp:  &req.Msg.PublicIp,
		PrivateIp: &req.Msg.PrivateIp,
	}
	// Update ip if it's not set
	if agent.ActiveIp == nil && req.Msg.PrivateIp != "" {
		params.ActiveIp = &req.Msg.PrivateIp
	}

	if req.Msg.Containers != nil {
		jsonBytes, err := json.Marshal(req.Msg.Containers)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		containers := string(jsonBytes)
		params.Containers = &containers
	}

	result, err := s.app.Conn.GetQuery().UpdateAgent(ctx, params)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	s.app.Event.Broadcast(&mantraev1.EventStreamResponse{
		Action: mantraev1.EventAction_EVENT_ACTION_UPDATED,
		Data: &mantraev1.EventStreamResponse_Agent{
			Agent: result.ToProto(),
		},
	})
	return connect.NewResponse(&mantraev1.HealthCheckResponse{
		Agent: result.ToProto(),
	}), nil
}
