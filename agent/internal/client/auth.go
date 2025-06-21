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
	"github.com/mizuchilabs/mantrae/agent/internal/collector"
	"github.com/mizuchilabs/mantrae/pkg/meta"
	"github.com/mizuchilabs/mantrae/pkg/util"
	mantraev1 "github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1"
	"github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1/mantraev1connect"
)

const tokenFile = "data/.mantrae-token"

type TokenSource struct {
	mu       sync.Mutex
	client   mantraev1connect.AgentServiceClient
	token    string
	claims   *meta.AgentClaims
	activeIP string
	fallback bool
}

func NewTokenSource() *TokenSource {
	return &TokenSource{fallback: false}
}

// SetToken loads the token from disk or env
func (t *TokenSource) SetToken(ctx context.Context) error {
	t.mu.Lock()
	if t.token != "" {
		t.mu.Unlock()
		return nil
	}

	// Try to load from disk
	data, err := os.ReadFile(tokenFile)
	if err == nil {
		t.token = strings.TrimSpace(string(data))
	}

	// Fallback to env
	if t.token == "" {
		t.token = strings.TrimSpace(os.Getenv("TOKEN"))
	}
	if t.token == "" {
		t.mu.Unlock()
		return errors.New("no token found in environment or file")
	}

	// Write it back
	_ = os.MkdirAll("data", 0o755)
	if err := os.WriteFile(tokenFile, []byte(t.token), 0o600); err != nil {
		slog.Warn("could not write token file", "error", err)
	}
	t.mu.Unlock()

	return t.SetClient()
}

// SetClient initializes the client
func (t *TokenSource) SetClient() error {
	t.mu.Lock()
	if t.token == "" {
		t.mu.Unlock()
		return errors.New("no token")
	}

	claims, err := util.DecodeUnsafeJWT[*meta.AgentClaims](t.token)
	if err != nil {
		t.mu.Unlock()
		return err
	}
	t.claims = claims

	t.client = mantraev1connect.NewAgentServiceClient(
		http.DefaultClient,
		claims.ServerURL,
		connect.WithInterceptors(t.Interceptor()),
	)
	t.mu.Unlock()

	return t.Refresh(context.Background()) // Check health
}

// Refresh calls HealthCheck and handles token rotation
func (t *TokenSource) Refresh(ctx context.Context) error {
	if t.client == nil {
		return errors.New("no client")
	}

	info := collector.GetMachineInfo()

	req := connect.NewRequest(&mantraev1.HealthCheckRequest{
		PublicIp:  info.PublicIPs.IPv4,
		PrivateIp: info.PrivateIPs.IPv4,
	})
	req.Header().Set("Authorization", "Bearer "+t.token)
	req.Header().Set(meta.HeaderAgentID, t.claims.AgentID)

	resp, err := t.client.HealthCheck(ctx, req)
	if err != nil {
		// Try fallback to env after removing token
		if connect.CodeOf(err) == connect.CodeUnauthenticated {
			if err := os.Remove(tokenFile); err != nil {
				return err
			}
			if !t.fallback {
				t.fallback = true
				return t.SetToken(ctx)
			}
			return errors.New("unauthenticated and no fallback $TOKEN available")
		}
		return err
	}

	// Shutdown on agent deletion
	if resp.Msg.Agent == nil {
		return errors.New("agent deleted")
	}

	// Handle token rotation
	newToken := resp.Msg.Agent.Token
	if newToken != "" && newToken != t.token {
		t.mu.Lock()
		t.token = newToken
		t.fallback = false
		_ = os.WriteFile(tokenFile, []byte(newToken), 0o600)
		t.mu.Unlock()
	}

	// Handle active IP rotation
	newIP := resp.Msg.Agent.ActiveIp
	if newIP != "" && newIP != t.activeIP {
		t.mu.Lock()
		t.activeIP = newIP
		t.mu.Unlock()
	}
	return nil
}

// Interceptor injects Authorization header, auto-refreshing on 401.
func (t *TokenSource) Interceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			if err := t.SetToken(ctx); err != nil {
				return nil, connect.NewError(connect.CodeUnauthenticated, err)
			}
			req.Header().Set("Authorization", "Bearer "+t.token)
			req.Header().Set(meta.HeaderAgentID, t.claims.AgentID)

			resp, err := next(ctx, req)
			if connect.CodeOf(err) == connect.CodeUnauthenticated {
				t.mu.Lock()
				t.token = ""
				t.mu.Unlock()
			}
			return resp, err
		}
	}
}

func (t *TokenSource) GetToken() string {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.token
}

func (t *TokenSource) GetClient() mantraev1connect.AgentServiceClient {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.client
}

func (t *TokenSource) PrintConnection() {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.client != nil {
		slog.Info("Connected", "server", t.claims.ServerURL)
	}
}
