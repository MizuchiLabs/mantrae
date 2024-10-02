package api

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"sync"

	"connectrpc.com/connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	agentv1 "github.com/MizuchiLabs/mantrae/agent/proto/gen/agent/v1"
	"github.com/MizuchiLabs/mantrae/agent/proto/gen/agent/v1/agentv1connect"
)

type AgentServer struct {
	mu     sync.Mutex
	agents map[string]agentv1.GetContainerRequest
}

func (s *AgentServer) GetContainer(
	ctx context.Context,
	req *connect.Request[agentv1.GetContainerRequest],
) (*connect.Response[agentv1.GetContainerResponse], error) {
	if err := validate(req.Header()); err != nil {
		return nil, err
	}

	// Add agent
	s.mu.Lock()
	s.agents[req.Msg.GetId()] = agentv1.GetContainerRequest{
		Id:         req.Msg.GetId(),
		Hostname:   req.Msg.GetHostname(),
		Containers: req.Msg.GetContainers(),
		LastSeen:   req.Msg.GetLastSeen(),
	}
	s.mu.Unlock()
	return connect.NewResponse(&agentv1.GetContainerResponse{}), nil
}

func (s *AgentServer) deleteAgent(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.agents, id)
}

func validate(header http.Header) error {
	auth := header.Get("authorization")
	if len(auth) == 0 {
		return connect.NewError(connect.CodeInvalidArgument, errors.New("missing authorization"))
	}
	if !strings.HasPrefix(auth, "Bearer ") {
		return connect.NewError(connect.CodeUnauthenticated, errors.New("missing bearer prefix"))
	}

	token := "test"

	if strings.TrimSpace(string(token)) != strings.TrimPrefix(auth, "Bearer ") {
		return connect.NewError(
			connect.CodeUnauthenticated,
			errors.New("failed to validate token"),
		)
	}

	return nil
}

func Server() {
	agent := &AgentServer{}

	mux := http.NewServeMux()
	path, handler := agentv1connect.NewAgentServiceHandler(agent)
	mux.Handle(path, handler)

	http.ListenAndServe(
		":8080",
		h2c.NewHandler(mux, &http2.Server{}),
	)
}
