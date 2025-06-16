package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"connectrpc.com/connect"
	"github.com/mizuchilabs/mantrae/internal/config"
	"github.com/mizuchilabs/mantrae/internal/db"
	"github.com/mizuchilabs/mantrae/internal/util"
	"github.com/mizuchilabs/mantrae/pkg/meta"
	"golang.org/x/crypto/bcrypt"
)

type ctxKey string

const (
	AuthUserKey  ctxKey = "user"
	AuthAgentKey ctxKey = "agent"
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

		userPassword, err := q.GetUserPassword(r.Context(), user.ID)
		if err != nil {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(password)); err != nil {
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
