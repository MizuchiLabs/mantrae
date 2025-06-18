package middlewares

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"connectrpc.com/connect"
	"github.com/mizuchilabs/mantrae/internal/config"
	"github.com/mizuchilabs/mantrae/internal/util"
	"github.com/mizuchilabs/mantrae/pkg/meta"
	"github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1/mantraev1connect"
	"golang.org/x/crypto/bcrypt"
)

type ctxKey string

const (
	AuthUserIDKey  ctxKey = "user_id"
	AuthAgentIDKey ctxKey = "agent_id"
)

// BasicAuth middleware for http endpoints
func (h *MiddlewareHandler) BasicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		q := h.app.Conn.GetQuery()
		user, err := q.GetUserByUsername(r.Context(), username)
		if err != nil {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), AuthUserIDKey, user.ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Authentication middleware for gRPC endpoints
func Authentication(app *config.App) connect.UnaryInterceptorFunc {
	return connect.UnaryInterceptorFunc(func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			// Skip authentication for certain endpoints (like login)
			if isPublicEndpoint(req.Spec().Procedure) {
				return next(ctx, req)
			}

			authHeader := req.Header().Get("Authorization")
			if authHeader == "" {
				return nil, connect.NewError(
					connect.CodeUnauthenticated,
					fmt.Errorf("missing authorization header"),
				)
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				return nil, connect.NewError(
					connect.CodeUnauthenticated,
					fmt.Errorf("invalid authorization header"),
				)
			}

			// Check if it's an agent request
			agentID := req.Header().Get(meta.HeaderAgentID)
			if agentID != "" {
				agent, err := app.Conn.GetQuery().GetAgent(ctx, agentID)
				if err != nil {
					return nil, connect.NewError(
						connect.CodeNotFound,
						errors.New("agent not found"),
					)
				}
				if agent.Token != tokenString {
					return nil, connect.NewError(
						connect.CodeUnauthenticated,
						errors.New("token mismatch"),
					)
				}
				ctx = context.WithValue(ctx, AuthAgentIDKey, agent.ID)
				return next(ctx, req)
			}

			// Parse and validate the token
			claims, err := util.DecodeUserJWT(tokenString, app.Secret)
			if err != nil {
				return nil, connect.NewError(
					connect.CodeUnauthenticated,
					fmt.Errorf("invalid token: %w", err),
				)
			}

			// Add claims to context
			ctx = context.WithValue(ctx, AuthUserIDKey, claims.ID)
			return next(ctx, req)
		}
	})
}

// Helper
func isPublicEndpoint(procedure string) bool {
	publicEndpoints := map[string]bool{
		mantraev1connect.UserServiceLoginUserProcedure: true,
		mantraev1connect.UserServiceVerifyOTPProcedure: true,
		mantraev1connect.UserServiceSendOTPProcedure:   true,
	}
	return publicEndpoints[procedure]
}

func GetUserIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(AuthUserIDKey).(string)
	return id, ok
}

func GetAgentIDFromContext(ctx context.Context) (string, bool) {
	agent, ok := ctx.Value(AuthAgentIDKey).(string)
	return agent, ok
}
