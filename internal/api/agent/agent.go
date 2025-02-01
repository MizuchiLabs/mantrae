package agent

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"connectrpc.com/connect"

	agentv1 "github.com/MizuchiLabs/mantrae/agent/proto/gen/agent/v1"
	"github.com/MizuchiLabs/mantrae/internal/config"
	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/MizuchiLabs/mantrae/internal/traefik"
	"github.com/MizuchiLabs/mantrae/internal/util"
)

type AgentServer struct {
	app *config.App
}

func NewAgentServer(app *config.App) *AgentServer {
	return &AgentServer{
		app: app,
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

	q := s.app.Conn.GetQuery()
	updatedAgent, err := q.UpdateAgent(context.Background(), params)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	util.Broadcast <- util.EventMessage{
		Type:    util.EventTypeUpdate,
		Message: "agent",
	}

	if err = traefik.DecodeAgentConfig(s.app.Conn.Get(), updatedAgent); err != nil {
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
	q := s.app.Conn.GetQuery()
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

	if _, err := DecodeJWT(dbAgent.Token, s.app.Config.Secret); err != nil {
		return nil, connect.NewError(
			connect.CodeUnauthenticated,
			errors.New("failed to decode token"),
		)
	}

	return &dbAgent, nil
}
