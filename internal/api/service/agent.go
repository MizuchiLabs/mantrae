package service

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"connectrpc.com/connect"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mizuchilabs/mantrae/internal/api/middlewares"
	"github.com/mizuchilabs/mantrae/internal/config"
	"github.com/mizuchilabs/mantrae/internal/settings"
	"github.com/mizuchilabs/mantrae/internal/store/db"
	"github.com/mizuchilabs/mantrae/internal/store/schema"
	"github.com/mizuchilabs/mantrae/internal/util"
	"github.com/mizuchilabs/mantrae/pkg/meta"
	mantraev1 "github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1"
	"github.com/traefik/paerser/parser"
	"github.com/traefik/traefik/v3/pkg/config/dynamic"
)

type AgentService struct {
	app *config.App
}

func NewAgentService(app *config.App) *AgentService {
	return &AgentService{app: app}
}

func (s *AgentService) HealthCheck(
	ctx context.Context,
	req *connect.Request[mantraev1.HealthCheckRequest],
) (*connect.Response[mantraev1.HealthCheckResponse], error) {
	// Rotate Token
	token, err := s.updateToken(ctx, req.Header().Get(meta.HeaderAgentID))
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	util.Broadcast <- util.EventMessage{
		Type:     util.EventTypeUpdate,
		Category: util.EventCategoryAgent,
	}
	return connect.NewResponse(&mantraev1.HealthCheckResponse{Ok: true, Token: *token}), nil
}

func (s *AgentService) GetContainer(
	ctx context.Context,
	req *connect.Request[mantraev1.GetContainerRequest],
) (*connect.Response[mantraev1.GetContainerResponse], error) {
	agentID, ok := middlewares.GetAgentIDFromContext(ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeInternal, errors.New("agent context missing"))
	}

	// Upsert agent
	params := db.UpdateAgentParams{
		ID:       agentID,
		Hostname: &req.Msg.Hostname,
		PublicIp: &req.Msg.PublicIp,
	}
	// if agent.ActiveIp == nil {
	// 	params.ActiveIp = &req.Msg.PublicIp
	// }

	privateIPs := schema.AgentPrivateIPs{
		IPs: make([]string, len(req.Msg.PrivateIps)),
	}
	privateIPs.IPs = req.Msg.PrivateIps
	params.PrivateIps = &privateIPs

	var containers schema.AgentContainers
	for _, container := range req.Msg.Containers {
		created := container.Created.AsTime()
		containers = append(containers, schema.AgentContainer{
			ID:      container.Id,
			Name:    container.Name,
			Labels:  container.Labels,
			Image:   container.Image,
			Portmap: container.Portmap,
			Status:  container.Status,
			Created: &created,
		})
	}
	params.Containers = &containers

	q := s.app.Conn.GetQuery()
	updatedAgent, err := q.UpdateAgent(ctx, params)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// Update dynamic config
	if err = s.DecodeAgentConfig(updatedAgent); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	util.Broadcast <- util.EventMessage{
		Type:     util.EventTypeUpdate,
		Category: util.EventCategoryAgent,
	}
	return connect.NewResponse(&mantraev1.GetContainerResponse{}), nil
}

func (s *AgentService) updateToken(ctx context.Context, id string) (*string, error) {
	q := s.app.Conn.GetQuery()
	agent, err := q.GetAgent(ctx, id)
	if err != nil {
		return nil, err
	}

	claims, err := DecodeJWT(agent.Token, s.app.Secret)
	if err != nil {
		return nil, err
	}

	// Only update the token if it's close to expiring (less than 25%)
	lifetime := claims.ExpiresAt.Sub(claims.IssuedAt.Time)
	remaining := time.Until(claims.ExpiresAt.Time)
	if remaining > lifetime/4 {
		return &agent.Token, nil // Still valid
	}

	agentInterval, ok := s.app.SM.Get(settings.KeyAgentCleanupInterval)
	if !ok {
		return nil, errors.New("failed to get agent cleanup interval setting")
	}

	token, err := claims.EncodeJWT(s.app.Secret, settings.AsDuration(agentInterval))
	if err != nil {
		return nil, err
	}

	err = q.UpdateAgentToken(ctx, db.UpdateAgentTokenParams{ID: agent.ID, Token: token})
	if err != nil {
		return nil, err
	}
	slog.Info("Rotating agent token", "agentID", agent.ID, "token", token)

	return &token, nil
}

// Helpers --------------------------------------------------------------------
type AgentClaims struct {
	AgentID   string `json:"agentId,omitempty"`
	ProfileID int64  `json:"profileId,omitempty"`
	ServerURL string `json:"serverUrl,omitempty"`
	jwt.RegisteredClaims
}

// EncodeJWT generates a JWT for agents
func (a *AgentClaims) EncodeJWT(secret string, expirationTime time.Duration) (string, error) {
	if a.ServerURL == "" || a.ProfileID == 0 {
		return "", errors.New("serverUrl and profileID cannot be empty")
	}

	if expirationTime == 0 {
		expirationTime = time.Hour * 24
	}

	claims := &AgentClaims{
		AgentID:   a.AgentID,
		ProfileID: a.ProfileID,
		ServerURL: a.ServerURL,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expirationTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// DecodeJWT decodes the agent token and returns claims if valid
func DecodeJWT(tokenString, secret string) (*AgentClaims, error) {
	claims := &AgentClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (any, error) {
			return []byte(secret), nil
		},
	)

	if err != nil || !token.Valid {
		return nil, err
	}
	return claims, nil
}

func (s *AgentService) DecodeAgentConfig(agent db.Agent) error {
	ctx := context.Background()

	q := s.app.Conn.GetQuery()
	for _, container := range *agent.Containers {
		dynConfig := &dynamic.Configuration{
			HTTP: &dynamic.HTTPConfiguration{},
			TCP:  &dynamic.TCPConfiguration{},
			UDP:  &dynamic.UDPConfiguration{},
			TLS:  &dynamic.TLSConfiguration{},
		}

		err := parser.Decode(
			container.Labels,
			dynConfig,
			parser.DefaultRootName,
			"traefik.http",
			"traefik.tcp",
			"traefik.udp",
			"traefik.tls.stores.default",
		)
		if err != nil {
			return err
		}

		for k, r := range dynConfig.HTTP.Routers {
			q.CreateHttpRouter(ctx, db.CreateHttpRouterParams{
				ProfileID: agent.ProfileID,
				AgentID:   &agent.ID,
				Name:      k,
				Config:    r,
			})
		}
		for k, r := range dynConfig.TCP.Routers {
			q.CreateTcpRouter(ctx, db.CreateTcpRouterParams{
				ProfileID: agent.ProfileID,
				AgentID:   &agent.ID,
				Name:      k,
				Config:    r,
			})
		}
		for k, r := range dynConfig.UDP.Routers {
			q.CreateUdpRouter(ctx, db.CreateUdpRouterParams{
				ProfileID: agent.ProfileID,
				AgentID:   &agent.ID,
				Name:      k,
				Config:    r,
			})
		}

		for k, r := range dynConfig.HTTP.Services {
			q.CreateHttpService(ctx, db.CreateHttpServiceParams{
				ProfileID: agent.ProfileID,
				AgentID:   &agent.ID,
				Name:      k,
				Config:    r,
			})
		}
		for k, r := range dynConfig.TCP.Services {
			q.CreateTcpService(ctx, db.CreateTcpServiceParams{
				ProfileID: agent.ProfileID,
				AgentID:   &agent.ID,
				Name:      k,
				Config:    r,
			})
		}
		for k, r := range dynConfig.UDP.Services {
			q.CreateUdpService(ctx, db.CreateUdpServiceParams{
				ProfileID: agent.ProfileID,
				AgentID:   &agent.ID,
				Name:      k,
				Config:    r,
			})
		}

		for k, r := range dynConfig.HTTP.Middlewares {
			q.CreateHttpMiddleware(ctx, db.CreateHttpMiddlewareParams{
				ProfileID: agent.ProfileID,
				AgentID:   &agent.ID,
				Name:      k,
				Config:    r,
			})
		}
		for k, r := range dynConfig.TCP.Middlewares {
			q.CreateTcpMiddleware(ctx, db.CreateTcpMiddlewareParams{
				ProfileID: agent.ProfileID,
				AgentID:   &agent.ID,
				Name:      k,
				Config:    r,
			})
		}
	}

	return nil
}
