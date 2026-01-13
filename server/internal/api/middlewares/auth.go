// Package middlewares provides authentication and logging middleware.
package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"connectrpc.com/connect"
	"github.com/mizuchilabs/mantrae/pkg/meta"
	"github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1/mantraev1connect"
	"github.com/mizuchilabs/mantrae/server/internal/config"
)

type ctxKey string

const (
	AuthUserIDKey  ctxKey = "user_id"
	AuthAgentIDKey ctxKey = "agent_id"
)

type AuthInterceptor struct {
	app *config.App
}

func NewAuthInterceptor(app *config.App) *AuthInterceptor {
	return &AuthInterceptor{app: app}
}

func (i *AuthInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return connect.UnaryFunc(
		func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			// Skip authentication for public endpoints
			if isPublicEndpoint(req.Spec().Procedure) {
				return next(ctx, req)
			}

			authedCtx, err := i.authenticateRequest(ctx, req.Header())
			if err != nil {
				return nil, err
			}

			return next(authedCtx, req)
		},
	)
}

func (i *AuthInterceptor) WrapStreamingClient(
	next connect.StreamingClientFunc,
) connect.StreamingClientFunc {
	return connect.StreamingClientFunc(
		func(ctx context.Context, spec connect.Spec) connect.StreamingClientConn {
			return next(ctx, spec)
		},
	)
}

func (i *AuthInterceptor) WrapStreamingHandler(
	next connect.StreamingHandlerFunc,
) connect.StreamingHandlerFunc {
	return connect.StreamingHandlerFunc(
		func(ctx context.Context, conn connect.StreamingHandlerConn) error {
			// Skip authentication for public endpoints (if any streaming endpoints are public)
			if isPublicEndpoint(conn.Spec().Procedure) {
				return next(ctx, conn)
			}

			authedCtx, err := i.authenticateRequest(ctx, conn.RequestHeader())
			if err != nil {
				return err
			}

			return next(authedCtx, conn)
		},
	)
}

// HTTP wrappers --------------------------------------------------------------

func (a *AuthInterceptor) WithAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Authenticate user using the same logic as Connect-RPC
		authedCtx, err := a.authenticateRequest(r.Context(), r.Header)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Continue with authenticated context
		next.ServeHTTP(w, r.WithContext(authedCtx))
	})
}

// Authentication logic -------------------------------------------------------

func (i *AuthInterceptor) authenticateRequest(
	ctx context.Context,
	header http.Header,
) (context.Context, error) {
	// Agent request (Bearer) -------------------------------------------------
	if agentID := header.Get(meta.HeaderAgentID); agentID != "" {
		agent, err := i.app.Conn.Query.GetAgent(ctx, agentID)
		if err != nil {
			return nil, connect.NewError(
				connect.CodeNotFound,
				errors.New("unauthorized"),
			)
		}
		if agent.Token != getBearerToken(header) {
			return nil, connect.NewError(
				connect.CodeUnauthenticated,
				errors.New("unauthorized"),
			)
		}

		return context.WithValue(ctx, AuthAgentIDKey, agent.ID), nil
	}

	// User request (Cookie/Bearer) -------------------------------------------
	if token := getCookieToken(header); token != "" {
		claims, err := meta.DecodeUserToken(token, i.app.Secret)
		if err != nil {
			return nil, connect.NewError(
				connect.CodeUnauthenticated,
				errors.New("unauthorized"),
			)
		}
		if claims.IsExpired() {
			return nil, connect.NewError(
				connect.CodeUnauthenticated,
				errors.New("unauthorized"),
			)
		}
		user, err := i.app.Conn.Query.GetUserByID(ctx, claims.UserID)
		if err != nil {
			return nil, connect.NewError(
				connect.CodeUnauthenticated,
				errors.New("unauthorized"),
			)
		}
		return context.WithValue(ctx, AuthUserIDKey, user.ID), nil
	}
	if token := getBearerToken(header); token != "" {
		claims, err := meta.DecodeUserToken(token, i.app.Secret)
		if err != nil {
			return nil, connect.NewError(
				connect.CodeUnauthenticated,
				errors.New("unauthorized"),
			)
		}
		if claims.IsExpired() {
			return nil, connect.NewError(
				connect.CodeUnauthenticated,
				errors.New("unauthorized"),
			)
		}
		user, err := i.app.Conn.Query.GetUserByID(ctx, claims.UserID)
		if err != nil {
			return nil, connect.NewError(
				connect.CodeUnauthenticated,
				errors.New("unauthorized"),
			)
		}
		return context.WithValue(ctx, AuthUserIDKey, user.ID), nil
	}

	// Unauthorized -----------------------------------------------------------
	return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("unauthorized"))
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
	const prefix = "Bearer "
	auth := header.Get("Authorization")
	// Case insensitive prefix match. See RFC 9110 Section 11.1.
	if len(auth) < len(prefix) || !strings.EqualFold(auth[:len(prefix)], prefix) {
		return ""
	}
	return auth[len(prefix):]
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

func GetUserIDFromContext(ctx context.Context) *string {
	if id := ctx.Value(AuthUserIDKey); id != nil {
		if userID, ok := id.(string); ok && userID != "" {
			return &userID
		}
	}
	return nil
}

func GetAgentIDFromContext(ctx context.Context) *string {
	if agent := ctx.Value(AuthAgentIDKey); agent != nil {
		if agentID, ok := agent.(string); ok && agentID != "" {
			return &agentID
		}
	}
	return nil
}
