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

	token, err := util.EncodeAgentJWT(decoded.ServerURL)
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
		Hostname:   req.Msg.GetHostname(),
		PublicIp:   &req.Msg.PublicIp,
		PrivateIps: privateIpsJSON,
		Containers: containersJSON,
		LastSeen:   &lastSeen,
	}); err != nil {
		s.mu.Unlock()
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	s.mu.Unlock()

	_, err = traefik.DecodeFromLabels(req.Msg.GetId())
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
}
