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
	"time"

	"connectrpc.com/connect"
	"github.com/mizuchilabs/mantrae/agent/internal/collector"
	"github.com/mizuchilabs/mantrae/pkg/meta"
	"github.com/mizuchilabs/mantrae/pkg/util"
	mantraev1 "github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1"
	"github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1/mantraev1connect"
)

const (
	tokenFile          = "data/.mantrae-token"
	connectionTimeout  = 10 * time.Second
	healthCheckTimeout = 5 * time.Second
)

type ConnectionState int

const (
	StateDisconnected ConnectionState = iota
	StateConnecting
	StateConnected
	StateError
)

func (s ConnectionState) String() string {
	switch s {
	case StateDisconnected:
		return "disconnected"
	case StateConnecting:
		return "connecting"
	case StateConnected:
		return "connected"
	case StateError:
		return "error"
	default:
		return "unknown"
	}
}

type TokenSource struct {
	mu       sync.RWMutex
	client   mantraev1connect.AgentServiceClient
	token    string
	claims   *meta.AgentClaims
	activeIP string
	state    ConnectionState
}

func NewTokenSource() *TokenSource {
	t := &TokenSource{state: StateDisconnected}
	if err := t.prepare(); err != nil {
		slog.Error("failed to prepare token source", "error", err)
	}
	return t
}

func (t *TokenSource) prepare() error {
	t.state = StateConnecting

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

	// already initialized?
	if t.client != nil && t.token != "" {
		return nil
	}

	t.state = StateConnecting

	// load token if needed
	if t.token == "" {
		data, _ := os.ReadFile(tokenFile)
		t.token = strings.TrimSpace(string(data))
		if t.token == "" {
			t.state = StateDisconnected
			return errors.New("no token found")
		}
		if err := os.WriteFile(tokenFile, []byte(t.token), 0o600); err != nil {
			slog.Warn("could not write token file", "err", err)
		}
	}

	var err error
	t.claims, err = meta.DecodeAgentToken(t.token, "")
	if err != nil {
		return err
	}

	slog.Info("Starting agent...", "server", func() string {
		if t.claims != nil {
			return t.claims.ServerURL
		}
		return "unknown"
	}())

	httpClient := &http.Client{
		Timeout: connectionTimeout,
	}

	t.client = mantraev1connect.NewAgentServiceClient(
		httpClient,
		util.CleanURL(t.claims.ServerURL),
		connect.WithInterceptors(t.Interceptor()),
	)
	return nil
}

// Refresh does a health‐check, rotates token or falls back on unauthenticated.
func (t *TokenSource) Refresh(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, healthCheckTimeout)
	defer cancel()

	if err := t.ensure(); err != nil {
		slog.Error("Failed to connect to server", "error", err)
		return
	}

	t.state = StateConnecting

	info := collector.GetMachineInfo()
	req := connect.NewRequest(&mantraev1.HealthCheckRequest{
		MachineId: info.MachineID,
		Hostname:  info.Hostname,
		PrivateIp: info.PrivateIP,
		PublicIp:  info.PublicIPs.IPv4,
	})

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
			slog.Error("Failed to connect to server after token refresh", "error", err)
			t.state = StateError
			return
		}

		t.Refresh(ctx)
		return
	} else if err != nil {
		slog.Error("Health check failed", "error", err)
		t.state = StateError
		return
	}

	if resp.Msg.Agent == nil {
		slog.Error("Agent deleted on server")
		t.state = StateError
		return
	}

	// Successfully connected
	t.state = StateConnected

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
				t.state = StateError
			}
			return resp, err
		}
	}
}

func (t *TokenSource) IsConnected() bool {
	return t.state == StateConnected
}

func (t *TokenSource) GetToken() string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.token
}

func (t *TokenSource) GetClient() mantraev1connect.AgentServiceClient {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.client
}

func (t *TokenSource) PrintConnection() {
	t.mu.RLock()
	defer t.mu.RUnlock()

	switch t.state {
	case StateConnected:
		slog.Info("Agent connected", "services_ip", t.activeIP)
	case StateConnecting:
		slog.Info("Agent connecting", "server", func() string {
			if t.claims != nil {
				return t.claims.ServerURL
			}
			return "unknown"
		}())
	case StateDisconnected:
		slog.Info("Agent disconnected")
	case StateError:
		slog.Error("Agent connection error")
	}
}
