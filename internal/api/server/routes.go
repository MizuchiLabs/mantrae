package server

import (
	"net/http"

	"github.com/MizuchiLabs/mantrae/internal/api/handler"
	"github.com/MizuchiLabs/mantrae/internal/api/middlewares"
)

func (s *Server) routes() {
	// Create middleware handler with database access
	mw := middlewares.NewMiddlewareHandler(s.app)

	// Middleware chains
	logChain := middlewares.Chain(
		mw.Logger,
	)

	jwtChain := middlewares.Chain(
		mw.Logger,
		mw.JWT,
	)

	basicChain := middlewares.Chain(
		mw.Logger,
		mw.BasicAuth,
	)

	adminChain := middlewares.Chain(
		mw.Logger,
		mw.JWT,
		mw.AdminOnly,
	)

	// Helper for middleware registration
	register := func(method, path string, chain middlewares.Middleware, handler http.HandlerFunc) {
		s.mux.Handle(method+" /api"+path, chain(handler))
	}

	// Auth
	register("POST", "/login", logChain, handler.Login(s.app))
	register("POST", "/logout", logChain, handler.Logout)
	register("GET", "/verify", jwtChain, handler.Verify)
	register("POST", "/verify/otp", logChain, handler.VerifyOTP(s.app))
	register("POST", "/reset/{name}", logChain, handler.SendResetEmail(s.app))
	register("GET", "/oidc/callback", logChain, handler.OIDCCallback(s.app))
	register("GET", "/oidc/login", logChain, handler.OIDCLogin(s.app))
	register("GET", "/oidc/status", logChain, handler.OIDCStatus(s.app))

	// Events
	register("GET", "/events", logChain, handler.GetEvents)
	register("GET", "/version", logChain, handler.GetVersion)

	// Profiles
	register("GET", "/profile", jwtChain, handler.ListProfiles(s.app))
	register("GET", "/profile/{id}", jwtChain, handler.GetProfile(s.app))
	register("POST", "/profile", jwtChain, handler.CreateProfile(s.app))
	register("PUT", "/profile", jwtChain, handler.UpdateProfile(s.app))
	register("DELETE", "/profile/{id}", jwtChain, handler.DeleteProfile(s.app))

	// Routers/Services
	register("POST", "/router/{id}", jwtChain, handler.UpsertRouter(s.app))
	register("DELETE", "/router", jwtChain, handler.DeleteRouter(s.app))
	register("DELETE", "/router/bulk", jwtChain, handler.BulkDeleteRouter(s.app))

	// Middlewares
	register("POST", "/middleware/{id}", jwtChain, handler.UpsertMiddleware(s.app))
	register("DELETE", "/middleware", jwtChain, handler.DeleteMiddleware(s.app))
	register("DELETE", "/middleware/bulk", jwtChain, handler.BulkDeleteMiddleware(s.app))
	register("GET", "/middleware/plugins", jwtChain, handler.GetMiddlewarePlugins)

	// Users
	register("GET", "/user", jwtChain, handler.ListUsers(s.app))
	register("GET", "/user/{id}", jwtChain, handler.GetUser(s.app))
	register("POST", "/user", adminChain, handler.CreateUser(s.app))
	register("PUT", "/user", jwtChain, handler.UpdateUser(s.app))
	register("DELETE", "/user/{id}", adminChain, handler.DeleteUser(s.app))
	register("POST", "/user/password", adminChain, handler.UpdateUserPassword(s.app))

	// DNS Provider
	register("GET", "/dns", jwtChain, handler.ListDNSProviders(s.app))
	register("GET", "/dns/{id}", jwtChain, handler.GetDNSProvider(s.app))
	register("POST", "/dns", adminChain, handler.CreateDNSProvider(s.app))
	register("PUT", "/dns", adminChain, handler.UpdateDNSProvider(s.app))
	register("DELETE", "/dns/{id}", adminChain, handler.DeleteDNSProvider(s.app))

	// DNS To Router
	register("GET", "/dns/router", jwtChain, handler.GetRouterDNSProvider(s.app))
	register("GET", "/dns/router/{id}", jwtChain, handler.ListRouterDNSProviders(s.app))
	register("POST", "/dns/router", jwtChain, handler.SetRouterDNSProvider(s.app))
	register("DELETE", "/dns/router", jwtChain, handler.DeleteRouterDNSProvider(s.app))

	// Settings (admin only)
	register("GET", "/settings", adminChain, handler.ListSettings(s.app.SM))
	register("GET", "/settings/{key}", adminChain, handler.GetSetting(s.app.SM))
	register("POST", "/settings", adminChain, handler.UpsertSetting(s.app.SM))

	// Agent
	register("GET", "/agent", jwtChain, handler.ListAgents(s.app))
	register("GET", "/agent/list/{id}", jwtChain, handler.ListAgentsByProfile(s.app))
	register("GET", "/agent/{id}", jwtChain, handler.GetAgent(s.app))
	register("POST", "/agent/{id}", jwtChain, handler.CreateAgent(s.app))
	register("PUT", "/agent", jwtChain, handler.UpdateAgentIP(s.app))
	register("DELETE", "/agent/{id}", jwtChain, handler.DeleteAgent(s.app))
	register("POST", "/agent/token/{id}", jwtChain, handler.RotateAgentToken(s.app))

	// Backup & Restore (admin only)
	register("GET", "/backups", adminChain, handler.ListBackups(s.app.BM))
	register("GET", "/backups/download", adminChain, handler.DownloadBackup(s.app.BM))
	register("GET", "/backups/download/{name}", adminChain, handler.DownloadBackupByName(s.app.BM))
	register("POST", "/backups", adminChain, handler.CreateBackup(s.app.BM))
	register("POST", "/backups/restore", adminChain, handler.RestoreBackup(s.app.BM))
	register("POST", "/backups/restore/{name}", adminChain, handler.RestoreBackupByName(s.app.BM))
	register("DELETE", "/backups/{name}", adminChain, handler.DeleteBackup(s.app.BM))
	register("POST", "/dynamic/restore/{id}", adminChain, handler.RestoreDynamicConfig(s.app.BM))

	// Errors
	register("GET", "/errors", jwtChain, handler.ListErrors(s.app))
	register("GET", "/errors/{id}", jwtChain, handler.GetErrorsByProfile(s.app))
	register("DELETE", "/errors/{id}", jwtChain, handler.DeleteErrorByID(s.app))
	register("DELETE", "/errors/profile/{id}", jwtChain, handler.DeleteErrorsByProfile(s.app))

	// Current IP
	register("GET", "/ip", jwtChain, handler.GetPublicIP)

	// Traefik
	register("GET", "/traefik/{id}/{source}", jwtChain, handler.GetTraefikConfig(s.app))

	// Dynamic config
	if s.app.Config.Server.BasicAuth {
		register("GET", "/{name}", basicChain, handler.PublishTraefikConfig(s.app))
	} else {
		register("GET", "/{name}", logChain, handler.PublishTraefikConfig(s.app))
	}
}
