// Package client for authentication and connection to the mantrae server.
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
}

func NewTokenSource() *TokenSource {
	t := &TokenSource{}
	if err := t.prepare(); err != nil {
		slog.Error("failed to prepare token source", "error", err)
	}
	return t
}

func (t *TokenSource) prepare() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Create token file if it doesn't exist
	if err := os.MkdirAll("data", 0o755); err != nil {
		return err
	}

	token, ok := os.LookupEnv("TOKEN")
	if !ok {
		return nil
	}

	if err := os.WriteFile(tokenFile, []byte(token), 0o600); err != nil {
		slog.Warn("could not write token file", "err", err)
	}
	return nil
}

// ensure loads token from disk, writes it back, decodes claims and initializes client.
func (t *TokenSource) ensure() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	// already ready?
	if t.client != nil && t.token != "" {
		return nil
	}

	// load token if needed
	if t.token == "" {
		data, _ := os.ReadFile(tokenFile)
		t.token = strings.TrimSpace(string(data))
		if t.token == "" {
			return errors.New("no token found")
		}
		if err := os.WriteFile(tokenFile, []byte(t.token), 0o600); err != nil {
			slog.Warn("could not write token file", "err", err)
		}
	}

	// parse & client
	claims, err := meta.DecodeAgentToken(t.token, "")
	if err != nil {
		return err
	}
	t.claims = claims
	t.client = mantraev1connect.NewAgentServiceClient(
		http.DefaultClient,
		claims.ServerURL,
		connect.WithInterceptors(t.Interceptor()),
	)
	return nil
}

// Refresh does a health‐check, rotates token or falls back on unauthenticated.
func (t *TokenSource) Refresh(ctx context.Context) {
	if err := t.ensure(); err != nil {
		slog.Error("Failed to connect to server", "error", err)
		return
	}

	info := collector.GetMachineInfo()
	req := connect.NewRequest(&mantraev1.HealthCheckRequest{
		MachineId: info.MachineID,
		Hostname:  info.Hostname,
		PublicIp:  info.PublicIPs.IPv4,
		PrivateIp: info.PrivateIPs.IPv4,
	})
	req.Header().Set("Authorization", "Bearer "+t.token)
	req.Header().Set(meta.HeaderAgentID, t.claims.AgentID)

	resp, err := t.client.HealthCheck(ctx, req)
	if connect.CodeOf(err) == connect.CodeUnauthenticated {
		// remove stored token and retry once from env
		if err = os.Remove(tokenFile); err != nil {
			slog.Warn("Could not remove token file", "err", err)
		}
		t.mu.Lock()
		t.token, t.client = "", nil
		t.mu.Unlock()
		if err = t.ensure(); err != nil {
			slog.Error("Failed to connect to server", "error", err)
			return
		}
		t.Refresh(ctx)
	} else if err != nil {
		slog.Error("Failed to connect to server", "error", err)
		return
	}

	if resp.Msg.Agent == nil {
		slog.Error("Agent deleted")
		return
	}

	// rotate token if changed
	if nt := resp.Msg.Agent.Token; nt != "" && nt != t.token {
		t.mu.Lock()
		t.token = nt
		if err = os.WriteFile(tokenFile, []byte(nt), 0o600); err != nil {
			slog.Warn("Could not write token file", "err", err)
		}
		t.mu.Unlock()
	}
	// rotate active IP
	if ip := resp.Msg.Agent.ActiveIp; ip != "" && ip != t.activeIP {
		t.mu.Lock()
		t.activeIP = ip
		t.mu.Unlock()
	}
}

// Interceptor injects Authorization header, auto-refreshing on 401.
func (t *TokenSource) Interceptor() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			if err := t.ensure(); err != nil {
				return nil, connect.NewError(connect.CodeUnauthenticated, err)
			}
			req.Header().Set("Authorization", "Bearer "+t.token)
			req.Header().Set(meta.HeaderAgentID, t.claims.AgentID)

			resp, err := next(ctx, req)
			if connect.CodeOf(err) == connect.CodeUnauthenticated {
				if err = os.Remove(tokenFile); err != nil {
					slog.Warn("Could not remove token file", "err", err)
				}
				t.mu.Lock()
				t.token, t.client = "", nil
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
