package server

import (
	"net/http"

	"github.com/MizuchiLabs/mantrae/internal/api/handler"
	"github.com/MizuchiLabs/mantrae/internal/api/middlewares"
)

func (s *Server) routes() {
	DB := s.app.DB
	// Create middleware handler with database access
	mw := middlewares.NewMiddlewareHandler(DB, *s.app.Config)

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

	// adminChain := middlewares.Chain(
	// 	mw.Logger,
	// 	mw.JWT,
	// 	mw.AdminOnly,
	// )

	// Helper for middleware registration
	register := func(method, path string, chain middlewares.Middleware, handler http.HandlerFunc) {
		s.mux.Handle(method+" /api"+path, chain(handler))
	}

	// Auth
	register("POST", "/login", logChain, handler.Login(DB, s.app.Config.Secret))
	register("POST", "/verify", logChain, handler.VerifyToken(DB, s.app.Config.Secret))
	register("POST", "/reset", logChain, handler.ResetPassword(DB, s.app.Config.Secret))
	register("POST", "/reset/{name}", logChain, handler.SendResetEmail(DB, s.app.Config.Secret))

	// Events
	register("GET", "/events", logChain, handler.GetEvents)
	register("GET", "/version", logChain, handler.GetVersion)

	// Profiles
	register("GET", "/profile", jwtChain, handler.ListProfiles(DB))
	register("GET", "/profile/{id}", jwtChain, handler.GetProfile(DB))
	register("POST", "/profile", jwtChain, handler.CreateProfile(DB))
	register("PUT", "/profile", jwtChain, handler.UpdateProfile(DB))
	register("DELETE", "/profile/{id}", jwtChain, handler.DeleteProfile(DB))

	// Routers/Services
	register("POST", "/router/{id}", jwtChain, handler.UpsertRouter(DB))
	register("DELETE", "/router/{id}/{name}/{protocol}", jwtChain, handler.DeleteRouter(DB))

	// Middlewares
	register("POST", "/middleware/{id}", jwtChain, handler.UpsertMiddleware(DB))
	register("DELETE", "/middleware/{id}/{name}/{protocol}", jwtChain, handler.DeleteMiddleware(DB))
	register("GET", "/middleware/plugins", jwtChain, handler.GetMiddlewarePlugins)

	// Users
	register("GET", "/user", jwtChain, handler.ListUsers(DB))
	register("GET", "/user/{id}", jwtChain, handler.GetUser(DB))
	register("POST", "/user", jwtChain, handler.CreateUser(DB))
	register("PUT", "/user", jwtChain, handler.UpdateUser(DB))
	register("DELETE", "/user/{id}", jwtChain, handler.DeleteUser(DB))

	// DNS Provider
	register("GET", "/dns", jwtChain, handler.ListDNSProviders(DB))
	register("GET", "/dns/{id}", jwtChain, handler.GetDNSProvider(DB))
	register("POST", "/dns", jwtChain, handler.CreateDNSProvider(DB))
	register("PUT", "/dns", jwtChain, handler.UpdateDNSProvider(DB))
	register("DELETE", "/dns/{id}", jwtChain, handler.DeleteDNSProvider(DB))

	// DNS To Router
	register("GET", "/dns/router", jwtChain, handler.GetRouterDNSProvider(DB))
	register("GET", "/dns/router/{id}", jwtChain, handler.ListRouterDNSProviders(DB))
	register("POST", "/dns/router", jwtChain, handler.SetRouterDNSProvider(DB))
	register("DELETE", "/dns/router", jwtChain, handler.DeleteRouterDNSProvider(DB))

	// Settings
	register("GET", "/settings", jwtChain, handler.ListSettings(s.app.SM))
	register("GET", "/settings/{key}", jwtChain, handler.GetSetting(s.app.SM))
	register("POST", "/settings", jwtChain, handler.UpsertSetting(s.app.SM))

	// Agent
	register("GET", "/agent", jwtChain, handler.ListAgents(DB))
	register("GET", "/agent/list/{id}", jwtChain, handler.ListAgentsByProfile(DB))
	register("GET", "/agent/{id}", jwtChain, handler.GetAgent(DB))
	register("POST", "/agent/{id}", jwtChain, handler.CreateAgent(s.app))
	register("PUT", "/agent", jwtChain, handler.UpdateAgent(DB))
	register("DELETE", "/agent/{id}", jwtChain, handler.DeleteAgent(DB))
	register("POST", "/agent/token/{id}", jwtChain, handler.RotateAgentToken(s.app))

	// Backup
	register("GET", "/backups", jwtChain, handler.ListBackups(s.app.BM))
	register("GET", "/backups/download", jwtChain, handler.DownloadBackup(s.app.BM))
	register(
		"GET",
		"/backups/download/{filename}",
		jwtChain,
		handler.DownloadBackupByName(s.app.BM),
	)
	register("POST", "/backups", jwtChain, handler.CreateBackup(s.app.BM))
	register("POST", "/backups/restore", jwtChain, handler.RestoreBackup(s.app.BM))
	register("DELETE", "/backups/{filename}", jwtChain, handler.DeleteBackup(s.app.BM))

	// IP
	// register("GET", "/ip/{id}", jwtChain, GetPublicIP)

	// Traefik
	register("GET", "/traefik/{id}/{source}", jwtChain, handler.GetTraefikConfig(DB))

	// Dynamic config
	if s.app.Config.Server.BasicAuth {
		register("GET", "/{name}", basicChain, handler.PublishTraefikConfig(DB))
	} else {
		register("GET", "/{name}", logChain, handler.PublishTraefikConfig(DB))
	}
}
