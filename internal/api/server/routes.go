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

	// Helper for middleware registration
	register := func(method, path string, chain middlewares.Middleware, handler http.HandlerFunc) {
		s.mux.Handle(method+" "+path, chain(handler))
	}

	// Auth
	register("POST", "/login", logChain, handler.Login(DB, s.app.Config.Secret))
	register("POST", "/verify", logChain, handler.VerifyToken(s.app.Config.Secret))
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
	register("POST", "/router/{id}", logChain, handler.UpsertRouter(DB))
	register("DELETE", "/router/{id}/{name}/{type}", jwtChain, handler.DeleteRouter(DB))

	// Middlewares
	register("POST", "/middleware/{id}", logChain, handler.UpsertMiddleware(DB))
	register("DELETE", "/middleware/{id}/{name}/{type}", jwtChain, handler.DeleteMiddleware(DB))
	register("GET", "/middleware/plugins", logChain, handler.GetMiddlewarePlugins)

	// Users
	register("GET", "/user", jwtChain, handler.ListUsers(DB))
	register("GET", "/user/{id}", jwtChain, handler.GetUser(DB))
	register("POST", "/user", jwtChain, handler.CreateUser(DB))
	register("PUT", "/user", jwtChain, handler.UpdateUser(DB))
	register("DELETE", "/user/{id}", jwtChain, handler.DeleteUser(DB))

	// DNS Provider
	register("GET", "/provider", jwtChain, handler.ListDNSProviders(DB))
	register("GET", "/provider/{id}", jwtChain, handler.GetDNSProvider(DB))
	register("POST", "/provider", jwtChain, handler.CreateDNSProvider(DB))
	register("PUT", "/provider", jwtChain, handler.UpdateDNSProvider(DB))
	register("DELETE", "/provider/{id}", jwtChain, handler.DeleteDNSProvider(DB))

	// Settings
	register("GET", "/settings", jwtChain, handler.ListSettings(DB))
	register("GET", "/settings/{key}", jwtChain, handler.GetSetting(DB))
	register("PUT", "/settings", jwtChain, handler.UpsertSetting(DB))

	// Agent
	register("GET", "/agent", jwtChain, handler.ListAgents(DB))
	register("GET", "/agent/{id}", jwtChain, handler.GetAgent(DB))
	register("POST", "/agent", jwtChain, handler.CreateAgent(DB))
	register("PUT", "/agent", jwtChain, handler.UpdateAgent(DB))
	register("DELETE", "/agent/{id}", jwtChain, handler.DeleteAgent(DB))
	// register("POST", "/agent/token/{id}", jwtChain, handler.RotateAgentToken(DB))

	// Backup
	register("GET", "/backups", jwtChain, handler.ListBackups(s.app.BM))
	register("GET", "/backup", jwtChain, handler.DownloadBackup(s.app.BM))
	register("POST", "/backup", jwtChain, handler.RestoreBackup(s.app.BM))

	// IP
	// register("GET", "/ip/{id}", jwtChain, GetPublicIP)

	// Traefik
	register("GET", "/traefik/{id}", jwtChain, handler.GetTraefikConfig(DB))

	// Dynamic config
	if s.app.Config.Server.BasicAuth {
		register("GET", "/{name}", basicChain, handler.PublishTraefikConfig(DB))
	} else {
		register("GET", "/{name}", logChain, handler.PublishTraefikConfig(DB))
	}
}
