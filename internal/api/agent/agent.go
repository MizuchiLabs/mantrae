package agent

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"sync"

	"connectrpc.com/connect"

	agentv1 "github.com/MizuchiLabs/mantrae/agent/proto/gen/agent/v1"
	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/MizuchiLabs/mantrae/internal/util"
)

type AgentServer struct {
	db *sql.DB
	mu sync.Mutex
}

func (s *AgentServer) HealthCheck(
	ctx context.Context,
	req *connect.Request[agentv1.HealthCheckRequest],
) (*connect.Response[agentv1.HealthCheckResponse], error) {
	if _, err := s.validate(req.Header(), req.Msg.GetId()); err != nil {
		return nil, err
	}

	return connect.NewResponse(&agentv1.HealthCheckResponse{Ok: true}), nil
}

func (s *AgentServer) GetContainer(
	ctx context.Context,
	req *connect.Request[agentv1.GetContainerRequest],
) (*connect.Response[agentv1.GetContainerResponse], error) {
	_, err := s.validate(req.Header(), req.Msg.GetId())
	if err != nil {
		return nil, err
	}

	// Upsert agent
	s.mu.Lock()
	privateIpsJSON, err := json.Marshal(req.Msg.GetPrivateIps())
	if err != nil {
		s.mu.Unlock()
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	containersJSON, err := json.Marshal(req.Msg.GetContainers())
	if err != nil {
		s.mu.Unlock()
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	q := db.New(s.db)
	if err := q.UpdateAgent(context.Background(), db.UpdateAgentParams{
		ID:         req.Msg.GetId(),
		Hostname:   &req.Msg.Hostname,
		PublicIp:   &req.Msg.PublicIp,
		PrivateIps: privateIpsJSON,
		Containers: containersJSON,
	}); err != nil {
		s.mu.Unlock()
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	util.Broadcast <- util.EventMessage{
		Type:    "agent_updated",
		Message: req.Msg.GetHostname(),
	}
	s.mu.Unlock()

	// err = traefik.DecodeFromLabels(req.Msg.GetId(), containersJSON)
	// if err != nil {
	// 	return nil, connect.NewError(connect.CodeInternal, err)
	// }
	return connect.NewResponse(&agentv1.GetContainerResponse{}), nil
}

func (s *AgentServer) validate(header http.Header, id string) (*db.Agent, error) {
	auth := header.Get("authorization")
	if len(auth) == 0 {
		return nil, connect.NewError(
			connect.CodeInvalidArgument,
			errors.New("missing authorization"),
		)
	}
	if !strings.HasPrefix(auth, "Bearer ") {
		return nil, connect.NewError(
			connect.CodeUnauthenticated,
			errors.New("missing bearer prefix"),
		)
	}

	// Check if agent exists
	q := db.New(s.db)
	agent, err := q.GetAgent(context.Background(), id)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("agent not found"))
	}

	// Check if token is valid
	if agent.Token != strings.TrimPrefix(auth, "Bearer ") {
		return nil, connect.NewError(
			connect.CodeUnauthenticated,
			errors.New("failed to validate token"),
		)
	}

	// if _, err := util.DecodeAgentJWT(agent.Token); err != nil {
	// 	return nil, connect.NewError(
	// 		connect.CodeUnauthenticated,
	// 		errors.New("failed to decode token"),
	// 	)
	// }

	return &agent, nil
}
