package client

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"sync"

	"connectrpc.com/connect"
	"github.com/mizuchilabs/mantrae/pkg/meta"
	agentv1 "github.com/mizuchilabs/mantrae/proto/gen/agent/v1"
	"github.com/mizuchilabs/mantrae/proto/gen/agent/v1/agentv1connect"
)

const tokenFile = "data/.mantrae-token"

type TokenSource struct {
	mu       sync.Mutex
	client   agentv1connect.AgentServiceClient
	token    string
	fallback bool
}

func NewTokenSource() *TokenSource {
	return &TokenSource{fallback: false}
}

// SetToken loads the token from disk or env
func (ts *TokenSource) SetToken(ctx context.Context) error {
	ts.mu.Lock()
	if ts.token != "" {
		ts.mu.Unlock()
		return nil
	}

	// Try to load from disk
	data, err := os.ReadFile(tokenFile)
	if err == nil {
		ts.token = strings.TrimSpace(string(data))
	}

	// Fallback to env
	if ts.token == "" {
		ts.token = strings.TrimSpace(os.Getenv("TOKEN"))
	}
	if ts.token == "" {
		ts.mu.Unlock()
		return errors.New("no token found in environment or file")
	}

	// Write it back
	_ = os.MkdirAll("data", 0o755)
	if err := os.WriteFile(tokenFile, []byte(ts.token), 0o600); err != nil {
		slog.Warn("could not write token file", "error", err)
	}
	ts.mu.Unlock()

	return ts.SetClient()
}

// SetClient initializes the client
func (ts *TokenSource) SetClient() error {
	ts.mu.Lock()
	if ts.token == "" {
		ts.mu.Unlock()
		return errors.New("no token")
	}

	claims, err := DecodeJWT(ts.token)
	if err != nil {
		ts.mu.Unlock()
		return err
	}

	ts.client = agentv1connect.NewAgentServiceClient(
		http.DefaultClient,
		claims.ServerURL,
		connect.WithInterceptors(ts.Interceptor()),
	)
	ts.mu.Unlock()

	return ts.Refresh(context.Background()) // Check health
}

// Refresh calls HealthCheck and handles token rotation
func (ts *TokenSource) Refresh(ctx context.Context) error {
	if ts.client == nil {
		return errors.New("no client")
	}

	req := connect.NewRequest(&agentv1.HealthCheckRequest{})
	req.Header().Set("Authorization", "Bearer "+ts.token)
	if claims, err := DecodeJWT(ts.token); err == nil {
		req.Header().Set(meta.HeaderAgentID, claims.AgentID)
	}

	resp, err := ts.client.HealthCheck(ctx, req)
	if err != nil {
		// Try fallback to env after removing token
		if connect.CodeOf(err) == connect.CodeUnauthenticated {
			if err := os.Remove(tokenFile); err != nil {
				return err
			}
			if !ts.fallback {
				ts.fallback = true
				return ts.SetToken(ctx)
			}
			return errors.New("unauthenticated and no fallback $TOKEN available")
		}
		return err
	}

	// Shutdown on agent deletion
	if !resp.Msg.Ok {
		return errors.New("agent deleted")
	}

	// Handle token rotation
	if newToken := resp.Msg.GetToken(); newToken != "" && newToken != ts.token {
		ts.mu.Lock()
		ts.token = newToken
		ts.fallback = false
		_ = os.WriteFile(tokenFile, []byte(newToken), 0o600)
		ts.mu.Unlock()
	}
	return nil
}

// Interceptor injects Authorization header, auto-refreshing on 401.
func (ts *TokenSource) Interceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			if err := ts.SetToken(ctx); err != nil {
				return nil, connect.NewError(connect.CodeUnauthenticated, err)
			}
			req.Header().Set("Authorization", "Bearer "+ts.token)
			if claims, err := DecodeJWT(ts.token); err == nil {
				req.Header().Set(meta.HeaderAgentID, claims.AgentID)
			}

			resp, err := next(ctx, req)
			if connect.CodeOf(err) == connect.CodeUnauthenticated {
				ts.mu.Lock()
				ts.token = ""
				ts.mu.Unlock()
			}
			return resp, err
		}
	}
}

func (ts *TokenSource) GetToken() string {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	return ts.token
}

func (ts *TokenSource) GetClient() agentv1connect.AgentServiceClient {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	return ts.client
}

func (ts *TokenSource) PrintConnection() {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	if ts.client != nil {
		claims, err := DecodeJWT(ts.token)
		if err == nil {
			slog.Info("Connected", "server", claims.ServerURL)
		}
	}
}
