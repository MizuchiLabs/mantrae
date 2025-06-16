package middlewares

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"connectrpc.com/connect"
	"github.com/mizuchilabs/mantrae/internal/config"
	"github.com/mizuchilabs/mantrae/internal/db"
	"github.com/mizuchilabs/mantrae/internal/util"
	"github.com/mizuchilabs/mantrae/pkg/meta"
	"github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1/mantraev1connect"
	"golang.org/x/crypto/bcrypt"
)

type ctxKey string

const (
	AuthUserIDKey ctxKey = "user_id"
	AuthUserKey   ctxKey = "user"
	AuthAgentKey  ctxKey = "agent"
)

// BasicAuth middleware for simple authentication
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

		ctx := context.WithValue(r.Context(), AuthUserKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// JWT authentication middleware
func (h *MiddlewareHandler) JWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(util.CookieName)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token := cookie.Value
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims, err := util.DecodeUserJWT(token, h.app.Config.Secret)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if claims.Username == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Verify user exists in database
		q := h.app.Conn.GetQuery()
		user, err := q.GetUserByUsername(r.Context(), claims.Username)
		if err != nil {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}

		// Add user to context
		ctx := context.WithValue(r.Context(), AuthUserKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

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
				ctx = context.WithValue(ctx, AuthAgentKey, &agent)
				return next(ctx, req)
			}

			// Parse and validate the token
			claims, err := util.DecodeUserJWT(tokenString, app.Config.Secret)
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

func (h *MiddlewareHandler) AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the username from the request context
		user, ok := r.Context().Value(AuthUserKey).(db.GetUserByUsernameRow)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if !user.IsAdmin {
			http.Error(w, "Admin privileges required", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func GetAgentContext(ctx context.Context) *db.Agent {
	agent, ok := ctx.Value(AuthAgentKey).(*db.Agent)
	if !ok {
		return nil
	}
	return agent
}

func AgentAuth(app *config.App) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			auth := req.Header().Get("Authorization")
			if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
				return nil, connect.NewError(
					connect.CodeUnauthenticated,
					errors.New("invalid or missing authorization header"),
				)
			}

			agentID := req.Header().Get(meta.HeaderAgentID)
			if agentID == "" {
				return nil, connect.NewError(
					connect.CodeUnauthenticated,
					errors.New("missing mantrae-agent-id header"),
				)
			}

			q := app.Conn.GetQuery()
			dbAgent, err := q.GetAgent(ctx, agentID)
			if err != nil {
				return nil, connect.NewError(connect.CodeNotFound, errors.New("agent not found"))
			}

			token := strings.TrimPrefix(auth, "Bearer ")
			if dbAgent.Token != token {
				return nil, connect.NewError(
					connect.CodeUnauthenticated,
					errors.New("token mismatch"),
				)
			}

			// Store agent in context
			ctx = context.WithValue(ctx, AuthAgentKey, &dbAgent)

			return next(ctx, req)
		}
	}
}

func isPublicEndpoint(procedure string) bool {
	publicEndpoints := map[string]bool{
		mantraev1connect.UserServiceLoginUserProcedure: true,
		mantraev1connect.UserServiceVerifyOTPProcedure: true,
		mantraev1connect.UserServiceSendOTPProcedure:   true,
	}
	return publicEndpoints[procedure]
}

// Helper functions to get auth info from context
func GetUserIDFromContext(ctx context.Context) (int64, bool) {
	id, ok := ctx.Value(AuthUserIDKey).(int64)
	return id, ok
}
