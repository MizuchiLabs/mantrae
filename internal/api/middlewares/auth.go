package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"connectrpc.com/connect"
	"github.com/mizuchilabs/mantrae/internal/config"
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

// JWT authentication middleware for http endpoints
func (h *MiddlewareHandler) JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		token = strings.TrimPrefix(token, "Bearer ")

		claims, err := meta.DecodeUserToken(token, h.app.Secret)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if claims.UserID == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Verify user exists in database
		q := h.app.Conn.GetQuery()
		user, err := q.GetUserByID(r.Context(), claims.UserID)
		if err != nil {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}

		// Add user to context
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

			// Agent request (Bearer) -----------------------------------------
			if agentID := req.Header().Get(meta.HeaderAgentID); agentID != "" {
				agent, err := app.Conn.GetQuery().GetAgent(ctx, agentID)
				if err != nil {
					return nil, connect.NewError(
						connect.CodeNotFound,
						errors.New("agent not found"),
					)
				}
				if agent.Token != getBearerToken(req.Header()) {
					return nil, connect.NewError(
						connect.CodeUnauthenticated,
						errors.New("token mismatch"),
					)
				}
				ctx = context.WithValue(ctx, AuthAgentIDKey, agent.ID)
				return next(ctx, req)
			}

			// User request (Cookie/Bearer) -----------------------------------
			if token := getCookieToken(req.Header()); token != "" {
				claims, err := meta.DecodeUserToken(token, app.Secret)
				if err != nil {
					return nil, connect.NewError(
						connect.CodeUnauthenticated,
						errors.New("invalid token"),
					)
				}
				ctx = context.WithValue(ctx, AuthUserIDKey, claims.UserID)
				return next(ctx, req)
			}
			if token := getBearerToken(req.Header()); token != "" {
				claims, err := meta.DecodeUserToken(token, app.Secret)
				if err != nil {
					return nil, connect.NewError(
						connect.CodeUnauthenticated,
						errors.New("invalid token"),
					)
				}
				ctx = context.WithValue(ctx, AuthUserIDKey, claims.UserID)
				return next(ctx, req)
			}

			// Unauthorized ---------------------------------------------------
			return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("unauthorized"))
		}
	})
}

// Helper
func isPublicEndpoint(procedure string) bool {
	publicEndpoints := map[string]bool{
		mantraev1connect.UserServiceLoginUserProcedure:     true,
		mantraev1connect.UserServiceVerifyOTPProcedure:     true,
		mantraev1connect.UserServiceSendOTPProcedure:       true,
		mantraev1connect.UserServiceGetOIDCStatusProcedure: true,
	}
	return publicEndpoints[procedure]
}

func getBearerToken(header http.Header) string {
	authHeader := header.Get("Authorization")
	if authHeader == "" {
		return ""
	}
	return strings.TrimPrefix(authHeader, "Bearer ")
}

func getCookieToken(header http.Header) string {
	cookieHeader := header.Get("Cookie")
	if cookieHeader == "" {
		return ""
	}
	cookies, err := http.ParseCookie(cookieHeader)
	if err != nil {
		return ""
	}
	var token string
	for _, cookie := range cookies {
		if cookie.Name == meta.CookieName {
			token = cookie.Value
		}
	}
	return token
}

func GetUserIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(AuthUserIDKey).(string)
	return id, ok
}

func GetAgentIDFromContext(ctx context.Context) (string, bool) {
	agent, ok := ctx.Value(AuthAgentIDKey).(string)
	return agent, ok
}
