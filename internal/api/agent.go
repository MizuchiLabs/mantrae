package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"sync"

	"connectrpc.com/connect"

	agentv1 "github.com/MizuchiLabs/mantrae/agent/proto/gen/agent/v1"
	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/MizuchiLabs/mantrae/pkg/traefik"
	"github.com/MizuchiLabs/mantrae/pkg/util"
)

type AgentServer struct {
	mu sync.Mutex
}

func (s *AgentServer) HealthCheck(
	ctx context.Context,
	req *connect.Request[agentv1.HealthCheckRequest],
) (*connect.Response[agentv1.HealthCheckResponse], error) {
	if err := validate(req.Header()); err != nil {
		return nil, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	agent, err := db.Query.GetAgentByID(context.Background(), req.Msg.GetId())
	// No agent found, ignore since maybe a new agent wants to connect
	if err != nil {
		return connect.NewResponse(&agentv1.HealthCheckResponse{Ok: true}), nil
	}

	if agent.Deleted {
		if err := db.Query.DeleteAgentByID(context.Background(), agent.ID); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		util.Broadcast <- util.EventMessage{
			Type:    "agent_updated",
			Message: agent.ID,
		}
		return connect.NewResponse(&agentv1.HealthCheckResponse{Ok: false}), nil
	}

	return connect.NewResponse(&agentv1.HealthCheckResponse{Ok: true}), nil
}

func (s *AgentServer) RefreshToken(
	ctx context.Context,
	req *connect.Request[agentv1.RefreshTokenRequest],
) (*connect.Response[agentv1.RefreshTokenResponse], error) {
	if err := validate(req.Header()); err != nil {
		return nil, err
	}

	decoded, err := util.DecodeJWT(req.Msg.GetToken())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	token, err := util.EncodeAgentJWT(decoded.ProfileID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&agentv1.RefreshTokenResponse{
		Token: token,
	}), nil
}

func (s *AgentServer) GetContainer(
	ctx context.Context,
	req *connect.Request[agentv1.GetContainerRequest],
) (*connect.Response[agentv1.GetContainerResponse], error) {
	if err := validate(req.Header()); err != nil {
		return nil, err
	}
	decoded, err := util.DecodeJWT(req.Msg.GetToken())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
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
	lastSeen := req.Msg.GetLastSeen().AsTime()
	if _, err := db.Query.UpsertAgent(context.Background(), db.UpsertAgentParams{
		ID:         req.Msg.GetId(),
		ProfileID:  decoded.ProfileID,
		Hostname:   req.Msg.GetHostname(),
		PublicIp:   &req.Msg.PublicIp,
		PrivateIps: privateIpsJSON,
		Containers: containersJSON,
		LastSeen:   &lastSeen,
	}); err != nil {
		s.mu.Unlock()
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	util.Broadcast <- util.EventMessage{
		Type:    "agent_updated",
		Message: req.Msg.GetHostname(),
	}
	s.mu.Unlock()

	err = traefik.DecodeFromLabels(req.Msg.GetId(), containersJSON)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&agentv1.GetContainerResponse{}), nil
}

func validate(header http.Header) error {
	auth := header.Get("authorization")
	if len(auth) == 0 {
		return connect.NewError(connect.CodeInvalidArgument, errors.New("missing authorization"))
	}
	if !strings.HasPrefix(auth, "Bearer ") {
		return connect.NewError(connect.CodeUnauthenticated, errors.New("missing bearer prefix"))
	}
	if strings.TrimSpace(util.App.Secret) != strings.TrimPrefix(auth, "Bearer ") {
		return connect.NewError(
			connect.CodeUnauthenticated,
			errors.New("failed to validate token"),
		)
	}

	return nil
}
