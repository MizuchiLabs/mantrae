package service

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"connectrpc.com/connect"

	"github.com/google/uuid"
	"github.com/mizuchilabs/mantrae/internal/config"
	mantraev1 "github.com/mizuchilabs/mantrae/internal/gen/mantrae/v1"
	"github.com/mizuchilabs/mantrae/internal/meta"
	"github.com/mizuchilabs/mantrae/internal/store/db"
	"github.com/mizuchilabs/mantrae/internal/util"
)

type AgentService struct {
	app *config.App
}

func NewAgentService(app *config.App) *AgentService {
	return &AgentService{app: app}
}

func (s *AgentService) GetAgent(
	ctx context.Context,
	req *mantraev1.GetAgentRequest,
) (*mantraev1.GetAgentResponse, error) {
	result, err := s.app.Conn.Q.GetAgent(ctx, req.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return &mantraev1.GetAgentResponse{Agent: result.ToProto()}, nil
}

func (s *AgentService) CreateAgent(
	ctx context.Context,
	req *mantraev1.CreateAgentRequest,
) (*mantraev1.CreateAgentResponse, error) {
	params := &db.CreateAgentParams{
		ID:        uuid.NewString(),
		ProfileID: req.ProfileId,
	}

	params.Token = util.GenerateAgentToken(
		strconv.Itoa(int(params.ProfileID)),
		params.ID,
	)

	result, err := s.app.Conn.Q.CreateAgent(ctx, params)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	s.app.Event.Broadcast(&mantraev1.EventStreamResponse{
		Action: mantraev1.EventAction_EVENT_ACTION_CREATED,
		Data: &mantraev1.EventStreamResponse_Agent{
			Agent: result.ToProto(),
		},
	})
	return &mantraev1.CreateAgentResponse{Agent: result.ToProto()}, nil
}

func (s *AgentService) UpdateAgent(
	ctx context.Context,
	req *mantraev1.UpdateAgentRequest,
) (*mantraev1.UpdateAgentResponse, error) {
	params := &db.UpdateAgentParams{
		ID:       req.Id,
		ActiveIp: req.Ip,
	}

	if req.RotateToken != nil && *req.RotateToken {
		agent, err := s.app.Conn.Q.GetAgent(ctx, params.ID)
		if err != nil {
			return nil, err
		}
		token := util.GenerateAgentToken(strconv.Itoa(int(agent.ProfileID)), agent.ID)
		params.Token = &token
	}

	result, err := s.app.Conn.Q.UpdateAgent(ctx, params)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	s.app.Event.Broadcast(&mantraev1.EventStreamResponse{
		Action: mantraev1.EventAction_EVENT_ACTION_UPDATED,
		Data: &mantraev1.EventStreamResponse_Agent{
			Agent: result.ToProto(),
		},
	})
	return &mantraev1.UpdateAgentResponse{Agent: result.ToProto()}, nil
}

func (s *AgentService) DeleteAgent(
	ctx context.Context,
	req *mantraev1.DeleteAgentRequest,
) (*mantraev1.DeleteAgentResponse, error) {
	agent, err := s.app.Conn.Q.GetAgent(ctx, req.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if err := s.app.Conn.Q.DeleteAgent(ctx, req.Id); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	s.app.Event.Broadcast(&mantraev1.EventStreamResponse{
		Action: mantraev1.EventAction_EVENT_ACTION_DELETED,
		Data: &mantraev1.EventStreamResponse_Agent{
			Agent: agent.ToProto(),
		},
	})
	return &mantraev1.DeleteAgentResponse{}, nil
}

func (s *AgentService) ListAgents(
	ctx context.Context,
	req *mantraev1.ListAgentsRequest,
) (*mantraev1.ListAgentsResponse, error) {
	params := &db.ListAgentsParams{
		ProfileID: req.ProfileId,
		Limit:     req.Limit,
		Offset:    req.Offset,
	}

	result, err := s.app.Conn.Q.ListAgents(ctx, params)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	totalCount, err := s.app.Conn.Q.CountAgents(ctx, req.ProfileId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	agents := make([]*mantraev1.Agent, 0, len(result))
	for _, a := range result {
		agents = append(agents, a.ToProto())
	}
	return &mantraev1.ListAgentsResponse{
		Agents:     agents,
		TotalCount: totalCount,
	}, nil
}

func (s *AgentService) HealthCheck(
	ctx context.Context,
	req *mantraev1.HealthCheckRequest,
) (*mantraev1.HealthCheckResponse, error) {
	ci, ok := connect.CallInfoForHandlerContext(ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeInternal, errors.New("failed to get call info"))
	}

	agentID := ci.RequestHeader().Get(meta.HeaderAgentID)
	if agentID == "" {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("agent id is required"),
		)
	}
	agent, err := s.app.Conn.Q.GetAgent(ctx, agentID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// Update Agent
	params := &db.UpdateAgentParams{
		ID:        agent.ID,
		Hostname:  &req.Hostname,
		PublicIp:  &req.PublicIp,
		PrivateIp: &req.PrivateIp,
	}
	// Update ip if it's not set
	if agent.ActiveIp == nil && req.PrivateIp != "" {
		params.ActiveIp = &req.PrivateIp
	}

	if req.Containers != nil {
		jsonBytes, err := json.Marshal(req.Containers)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		containers := string(jsonBytes)
		params.Containers = &containers
	}

	result, err := s.app.Conn.Q.UpdateAgent(ctx, params)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	s.app.Event.Broadcast(&mantraev1.EventStreamResponse{
		Action: mantraev1.EventAction_EVENT_ACTION_UPDATED,
		Data: &mantraev1.EventStreamResponse_Agent{
			Agent: result.ToProto(),
		},
	})
	return &mantraev1.HealthCheckResponse{Agent: result.ToProto()}, nil
}
