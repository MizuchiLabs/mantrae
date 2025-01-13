package server

import (
	"net/http"

	"github.com/MizuchiLabs/mantrae/internal/api/handler"
	"github.com/MizuchiLabs/mantrae/internal/api/middlewares"
	"github.com/MizuchiLabs/mantrae/pkg/util"
)

func (s *Server) routes() {
	// Helper for middleware registration
	register := func(method, path string, chain middlewares.Middleware, handler http.HandlerFunc) {
		s.mux.Handle(method+" "+path, chain(handler))
	}
	// Middlewares
	logChain := middlewares.Chain(middlewares.Log)
	jwtChain := middlewares.Chain(middlewares.Log, middlewares.JWT)
	basicChain := middlewares.Chain(middlewares.Log, middlewares.BasicAuth)

	// Auth
	register("POST", "/login", logChain, handler.Login)
	register("POST", "/verify", logChain, handler.VerifyToken)
	register("POST", "/reset", logChain, handler.ResetPassword)
	register("POST", "/reset/{name}", logChain, handler.SendResetEmail)

	// Events
	register("GET", "/events", logChain, handler.GetEvents)
	register("GET", "/version", logChain, handler.GetVersion)

	// Profiles
	register("GET", "/profile", jwtChain, handler.ListProfiles(s.db))
	register("GET", "/profile/{id}", jwtChain, handler.GetProfile(s.db))
	register("POST", "/profile", jwtChain, handler.CreateProfile(s.db))
	register("PUT", "/profile", jwtChain, handler.UpdateProfile(s.db))
	register("DELETE", "/profile/{id}", jwtChain, handler.DeleteProfile(s.db))

	// Routers
	register("GET", "/routers/http", logChain, handler.GetHTTPRoutersBySource(s.db))
	register("GET", "/routers/tcp", logChain, handler.GetTCPRoutersBySource(s.db))
	register("GET", "/routers/udp", logChain, handler.GetUDPRoutersBySource(s.db))
	register("GET", "/router/http", logChain, handler.GetHTTPRouterByName(s.db))
	register("GET", "/router/tcp", logChain, handler.GetTCPRouterByName(s.db))
	register("GET", "/router/udp", logChain, handler.GetUDPRouterByName(s.db))
	register("POST", "/router/http", logChain, handler.UpsertHTTPRouter(s.db))
	register("POST", "/router/tcp", logChain, handler.UpsertTCPRouter(s.db))
	register("POST", "/router/udp", logChain, handler.UpsertUDPRouter(s.db))
	register("DELETE", "/router/http", jwtChain, handler.DeleteHTTPRouter(s.db))
	register("DELETE", "/router/tcp", jwtChain, handler.DeleteTCPRouter(s.db))
	register("DELETE", "/router/udp", jwtChain, handler.DeleteUDPRouter(s.db))

	// Services
	register("GET", "/services/http", logChain, handler.GetHTTPServicesBySource(s.db))
	register("GET", "/services/tcp", logChain, handler.GetTCPServicesBySource(s.db))
	register("GET", "/services/udp", logChain, handler.GetUDPServicesBySource(s.db))
	register("GET", "/service/http", logChain, handler.GetHTTPServiceByName(s.db))
	register("GET", "/service/tcp", logChain, handler.GetTCPServiceByName(s.db))
	register("GET", "/service/udp", logChain, handler.GetUDPServiceByName(s.db))
	register("POST", "/service/http", logChain, handler.UpsertHTTPService(s.db))
	register("POST", "/service/tcp", logChain, handler.UpsertTCPService(s.db))
	register("POST", "/service/udp", logChain, handler.UpsertUDPService(s.db))
	register("DELETE", "/service/http", jwtChain, handler.DeleteHTTPService(s.db))
	register("DELETE", "/service/tcp", jwtChain, handler.DeleteTCPService(s.db))
	register("DELETE", "/service/udp", jwtChain, handler.DeleteUDPService(s.db))

	// Middlewares
	register("GET", "/middleware/http", logChain, handler.GetHTTPMiddlewaresBySource(s.db))
	register("GET", "/middleware/tcp", logChain, handler.GetTCPMiddlewaresBySource(s.db))
	register("GET", "/middleware/http", logChain, handler.GetHTTPMiddlewareByName(s.db))
	register("GET", "/middleware/tcp", logChain, handler.GetTCPMiddlewareByName(s.db))
	register("POST", "/middleware/http", logChain, handler.UpsertHTTPMiddleware(s.db))
	register("POST", "/middleware/tcp", logChain, handler.UpsertTCPMiddleware(s.db))
	register("DELETE", "/middleware/http", jwtChain, handler.DeleteHTTPMiddleware(s.db))
	register("DELETE", "/middleware/tcp", jwtChain, handler.DeleteTCPMiddleware(s.db))
	register("GET", "/middleware/plugins", logChain, handler.GetMiddlewarePlugins)

	// Users
	register("GET", "/user", jwtChain, handler.ListUsers(s.db))
	register("GET", "/user/{id}", jwtChain, handler.GetUser(s.db))
	register("POST", "/user", jwtChain, handler.CreateUser(s.db))
	register("PUT", "/user", jwtChain, handler.UpdateUser(s.db))
	register("DELETE", "/user/{id}", jwtChain, handler.DeleteUser(s.db))

	// DNS Provider
	register("GET", "/provider", jwtChain, handler.ListDNSProviders(s.db))
	register("GET", "/provider/{id}", jwtChain, handler.GetDNSProvider(s.db))
	register("POST", "/provider", jwtChain, handler.CreateDNSProvider(s.db))
	register("PUT", "/provider", jwtChain, handler.UpdateDNSProvider(s.db))
	register("DELETE", "/provider/{id}", jwtChain, handler.DeleteDNSProvider(s.db))

	// Settings
	register("GET", "/settings", jwtChain, handler.ListSettings(s.db))
	register("GET", "/settings/{key}", jwtChain, handler.GetSetting(s.db))
	register("PUT", "/settings", jwtChain, handler.UpsertSetting(s.db))

	// Agent
	register("GET", "/agent", jwtChain, handler.ListAgents(s.db))
	register("GET", "/agent/{id}", jwtChain, handler.GetAgent(s.db))
	register("POST", "/agent", jwtChain, handler.CreateAgent(s.db))
	register("PUT", "/agent", jwtChain, handler.UpdateAgent(s.db))
	register("DELETE", "/agent/{id}", jwtChain, handler.DeleteAgent(s.db))
	// register("POST", "/agent/token/{id}", jwtChain, handler.RotateAgentToken(s.db))

	// Backup
	register("GET", "/backup", jwtChain, handler.DownloadBackup)
	register("POST", "/backup", jwtChain, handler.UploadBackup)

	// IP
	// register("GET", "/ip/{id}", jwtChain, GetPublicIP)

	// Traefik
	register("GET", "/traefik/{id}", jwtChain, handler.GetTraefikConfig(s.db))

	// Dynamic config
	if util.App.EnableBasicAuth {
		register("GET", "/{name}", basicChain, handler.PublishTraefikConfig(s.db))
	} else {
		register("GET", "/{name}", logChain, handler.PublishTraefikConfig(s.db))
	}
}
