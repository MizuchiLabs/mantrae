package agent

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strings"
	"sync"

	"connectrpc.com/connect"

	agentv1 "github.com/MizuchiLabs/mantrae/agent/proto/gen/agent/v1"
	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/MizuchiLabs/mantrae/internal/traefik"
	"github.com/MizuchiLabs/mantrae/internal/util"
)

type AgentServer struct {
	db *sql.DB
	mu sync.Mutex
}

func NewAgentServer(db *sql.DB) *AgentServer {
	return &AgentServer{
		db: db,
	}
}

func (s *AgentServer) HealthCheck(
	ctx context.Context,
	req *connect.Request[agentv1.HealthCheckRequest],
) (*connect.Response[agentv1.HealthCheckResponse], error) {
	if _, err := s.validate(req.Header(), req.Msg.GetAgentId()); err != nil {
		return nil, err
	}

	return connect.NewResponse(&agentv1.HealthCheckResponse{Ok: true}), nil
}

func (s *AgentServer) GetContainer(
	ctx context.Context,
	req *connect.Request[agentv1.GetContainerRequest],
) (*connect.Response[agentv1.GetContainerResponse], error) {
	agent, err := s.validate(req.Header(), req.Msg.GetAgentId())
	if err != nil {
		return nil, err
	}

	// Upsert agent
	s.mu.Lock()
	params := db.UpdateAgentParams{
		ID:       req.Msg.GetAgentId(),
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

	q := db.New(s.db)
	updatedAgent, err := q.UpdateAgent(context.Background(), params)
	if err != nil {
		s.mu.Unlock()
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	util.Broadcast <- util.EventMessage{
		Type:    "agent_updated",
		Message: req.Msg.GetHostname(),
	}
	s.mu.Unlock()

	if err = traefik.DecodeAgentConfig(s.db, updatedAgent); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
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
	dbAgent, err := q.GetAgent(context.Background(), id)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("agent not found"))
	}

	// Check if token is valid
	if dbAgent.Token != strings.TrimPrefix(auth, "Bearer ") {
		return nil, connect.NewError(
			connect.CodeUnauthenticated,
			errors.New("failed to validate token"),
		)
	}

	// if _, err := agent.DecodeJWT(dbAgent.Token); err != nil {
	// 	return nil, connect.NewError(
	// 		connect.CodeUnauthenticated,
	// 		errors.New("failed to decode token"),
	// 	)
	// }

	return &dbAgent, nil
}
