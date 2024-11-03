package api

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"connectrpc.com/connect"

	agentv1 "github.com/MizuchiLabs/mantrae/agent/proto/gen/agent/v1"
	"github.com/MizuchiLabs/mantrae/agent/proto/gen/agent/v1/agentv1connect"
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

	decoded, err := util.DecodeJWT(req.Msg.GetToken())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	agent, err := db.Query.GetAgentByProfileID(context.Background(), db.GetAgentByProfileIDParams{
		ID:        req.Msg.GetId(),
		ProfileID: decoded.ProfileID,
	})
	// No agent found, ignore since a new agent wants to connect
	if err != nil {
		return connect.NewResponse(&agentv1.HealthCheckResponse{Ok: true}), nil
	}

	if agent.Deleted {
		if err := db.Query.DeleteAgentByID(context.Background(), agent.ID); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
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

	token, err := util.EncodeAgentJWT(decoded.ServerURL, decoded.ProfileID)
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
	if strings.TrimSpace(os.Getenv("SECRET")) != strings.TrimPrefix(auth, "Bearer ") {
		return connect.NewError(
			connect.CodeUnauthenticated,
			errors.New("failed to validate token"),
		)
	}

	return nil
}

func (s *AgentServer) cleanupAgents() {
	// Run cleanup every 5 minutes
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	enabled, err := db.Query.GetSettingByKey(context.Background(), "agent-cleanup-enabled")
	if err != nil {
		slog.Error("failed to get agent cleanup timeout", "error", err)
		return
	}

	if enabled.Value != "true" {
		return
	}

	// Timeout to delete old agents
	timeout, err := db.Query.GetSettingByKey(context.Background(), "agent-cleanup-timeout")
	if err != nil {
		slog.Error("failed to get agent cleanup timeout", "error", err)
		return
	}

	timeoutDuration, err := time.ParseDuration(timeout.Value)
	if err != nil {
		slog.Error("failed to parse timeout cleanup duration", "error", err)
	}

	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		agents, err := db.Query.ListAgents(context.Background())
		if err != nil {
			slog.Error("failed to query disconnected agents", "error", err)
			s.mu.Unlock()
			continue
		}

		for _, agent := range agents {
			if agent.LastSeen == nil {
				continue
			}

			if now.Sub(*agent.LastSeen) > timeoutDuration {
				if err := db.Query.DeleteAgentByID(context.Background(), agent.ID); err != nil {
					slog.Error("failed to delete disconnected agent", "id", agent.ID, "error", err)
				} else {
					slog.Info("Deleted disconnected agent", "id", agent.ID)
				}
			}

		}
		s.mu.Unlock()
	}
}

func Server(port string) {
	agent := &AgentServer{}

	mux := http.NewServeMux()
	path, handler := agentv1connect.NewAgentServiceHandler(agent)
	mux.Handle(path, handler)

	if port == "" {
		port = ":8090"
	} else if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}
	srv := &http.Server{
		Addr:              port,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1MB
	}

	slog.Info("gRPC server running on", "port", port)
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				slog.Info("gRPC server closed")
				return
			}
			slog.Error("gRPC server error", "err", err)
			return
		}
	}()

	go agent.cleanupAgents()
}
