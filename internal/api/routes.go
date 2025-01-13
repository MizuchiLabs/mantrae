package api

import (
	"io/fs"
	"log"
	"net/http"

	"github.com/MizuchiLabs/mantrae/agent/proto/gen/agent/v1/agentv1connect"
	"github.com/MizuchiLabs/mantrae/pkg/util"
	"github.com/MizuchiLabs/mantrae/web"
)

func Server() http.Handler {
	agent := http.NewServeMux()
	agent.Handle(agentv1connect.NewAgentServiceHandler(&AgentServer{}))

	mux := http.NewServeMux()
	mux.Handle("/", staticHandler())

	// API routing
	mux.Handle("/api/", http.StripPrefix("/api", apiRoutes()))

	// gRPC routing
	mux.Handle("/grpc/", http.StripPrefix("/grpc", agent))

	return mux
}

func apiRoutes() http.Handler {
	api := http.NewServeMux()

	// Helper for middleware registration
	register := func(method, path string, chain Middleware, handler http.HandlerFunc) {
		api.Handle(method+" "+path, chain(handler))
	}

	// Middlewares
	logChain := Chain(Log)
	jwtChain := Chain(Log, JWT)
	basicChain := Chain(Log, BasicAuth)

	// Auth
	register("POST", "/login", logChain, Login)
	register("POST", "/verify", logChain, VerifyToken)
	register("POST", "/reset", logChain, ResetPassword)
	register("POST", "/reset/{name}", logChain, SendResetEmail)

	// Events
	register("GET", "/events", logChain, GetEvents)
	register("GET", "/version", logChain, GetVersion)

	// Profiles
	register("GET", "/profile", jwtChain, GetProfiles)
	register("GET", "/profile/{id}", jwtChain, GetProfile)
	register("PUT", "/profile", jwtChain, UpsertProfile)
	register("DELETE", "/profile/{id}", jwtChain, DeleteProfile)

	// External Routers
	register("GET", "/router/external/http/{id}", logChain, GetExternalHTTPRouters)
	register("GET", "/router/external/tcp/{id}", logChain, GetExternalTCPRouters)
	register("GET", "/router/external/udp/{id}", logChain, GetExternalUDPRouters)
	register("GET", "/router/external/http/{id}/{name}", logChain, GetExternalHTTPRouter)
	register("GET", "/router/external/tcp/{id}/{name}", logChain, GetExternalTCPRouter)
	register("GET", "/router/external/udp/{id}/{name}", logChain, GetExternalUDPRouter)

	// Internal Routers
	register("GET", "/router/http/{id}", logChain, GetInternalHTTPRouters)
	register("GET", "/router/tcp/{id}", logChain, GetInternalTCPRouters)
	register("GET", "/router/udp/{id}", logChain, GetInternalUDPRouters)
	register("GET", "/router/http/{id}/{name}", logChain, GetInternalHTTPRouter)
	register("GET", "/router/tcp/{id}/{name}", logChain, GetInternalTCPRouter)
	register("GET", "/router/udp/{id}/{name}", logChain, GetInternalUDPRouter)
	register("POST", "/router/http/{id}/{name}", logChain, UpsertInternalHTTPRouter)
	register("POST", "/router/tcp/{id}/{name}", logChain, UpsertInternalTCPRouter)
	register("POST", "/router/udp/{id}/{name}", logChain, UpsertInternalUDPRouter)
	register("DELETE", "/router/http/{id}/{name}", jwtChain, DeleteInternalHTTPRouter)
	register("DELETE", "/router/tcp/{id}/{name}", jwtChain, DeleteInternalTCPRouter)
	register("DELETE", "/router/udp/{id}/{name}", jwtChain, DeleteInternalUDPRouter)

	// External Services
	register("GET", "/service/external/http/{id}", jwtChain, GetExternalHTTPServices)
	register("GET", "/service/external/tcp/{id}", jwtChain, GetExternalTCPServices)
	register("GET", "/service/external/udp/{id}", jwtChain, GetExternalUDPServices)
	register("GET", "/service/external/http/{id}/{name}", jwtChain, GetExternalHTTPService)
	register("GET", "/service/external/tcp/{id}/{name}", jwtChain, GetExternalTCPService)
	register("GET", "/service/external/udp/{id}/{name}", jwtChain, GetExternalUDPService)

	// Internal Services
	register("GET", "/service/http/{id}", logChain, GetInternalHTTPServices)
	register("GET", "/service/tcp/{id}", logChain, GetInternalTCPServices)
	register("GET", "/service/udp/{id}", logChain, GetInternalUDPServices)
	register("GET", "/service/http/{id}/{name}", logChain, GetInternalHTTPService)
	register("GET", "/service/tcp/{id}/{name}", logChain, GetInternalTCPService)
	register("GET", "/service/udp/{id}/{name}", logChain, GetInternalUDPService)
	register("POST", "/service/http/{id}/{name}", logChain, UpsertInternalHTTPService)
	register("POST", "/service/tcp/{id}/{name}", logChain, UpsertInternalTCPService)
	register("POST", "/service/udp/{id}/{name}", logChain, UpsertInternalUDPService)
	register("DELETE", "/service/http/{id}/{name}", jwtChain, DeleteInternalHTTPService)
	register("DELETE", "/service/tcp/{id}/{name}", jwtChain, DeleteInternalTCPService)
	register("DELETE", "/service/udp/{id}/{name}", jwtChain, DeleteInternalUDPService)

	// External Middlewares
	register("GET", "/middleware/external/http/{id}", logChain, GetExternalHTTPMiddlewares)
	register("GET", "/middleware/external/tcp/{id}", logChain, GetExternalTCPMiddlewares)
	register("GET", "/middleware/external/http/{id}/{name}", logChain, GetExternalHTTPMiddleware)
	register("GET", "/middleware/external/tcp/{id}/{name}", logChain, GetExternalTCPMiddleware)

	// Internal Middlewares
	register("GET", "/middleware/http/{id}", logChain, GetInternalHTTPMiddlewares)
	register("GET", "/middleware/tcp/{id}", logChain, GetInternalTCPMiddlewares)
	register("GET", "/middleware/http/{id}/{name}", logChain, GetInternalHTTPMiddleware)
	register("GET", "/middleware/tcp/{id}/{name}", logChain, GetInternalTCPMiddleware)
	register("POST", "/middleware/http/{id}/{name}", logChain, UpsertInternalHTTPMiddleware)
	register("POST", "/middleware/tcp/{id}/{name}", logChain, UpsertInternalTCPMiddleware)
	register("DELETE", "/middleware/http/{id}/{name}", jwtChain, DeleteInternalHTTPMiddleware)
	register("DELETE", "/middleware/tcp/{id}/{name}", jwtChain, DeleteInternalTCPMiddleware)
	register("GET", "/middleware/plugins", logChain, GetMiddlewarePlugins)

	// Users
	register("GET", "/user", jwtChain, GetUsers)
	register("GET", "/user/{id}", jwtChain, GetUser)
	register("POST", "/user", jwtChain, UpsertUser)
	register("DELETE", "/user/{id}", jwtChain, DeleteUser)

	// DNS Provider
	register("GET", "/provider", jwtChain, GetProviders)
	register("POST", "/provider", jwtChain, UpsertProvider)
	register("DELETE", "/provider/{id}", jwtChain, DeleteProvider)

	// Settings
	register("GET", "/settings", jwtChain, GetSettings)
	register("GET", "/settings/{key}", jwtChain, GetSetting)
	register("PUT", "/settings", jwtChain, UpdateSetting)

	// Agent
	register("GET", "/agent/{id}", jwtChain, GetAgents)
	register("PUT", "/agent", jwtChain, UpsertAgent)
	register("DELETE", "/agent/{id}", jwtChain, DeleteAgent)
	register("POST", "/agent/token/{id}", jwtChain, RegenerateAgentToken)

	// Backup
	register("GET", "/backup", jwtChain, DownloadBackup)
	register("POST", "/backup", jwtChain, UploadBackup)

	// IP
	register("GET", "/ip/{id}", jwtChain, GetPublicIP)

	// Traefik
	register("GET", "/traefik/{id}", jwtChain, GetTraefikOverview)
	register("GET", "/traefik/entrypoints/{id}", logChain, GetTraefikEntrypoints)

	// Dynamic config
	if util.App.EnableBasicAuth {
		register("GET", "/{name}", basicChain, PublishTraefikConfig)
	} else {
		register("GET", "/{name}", logChain, PublishTraefikConfig)
	}

	return Cors(api)
}

func staticHandler() http.Handler {
	mux := http.NewServeMux()
	staticContent, err := fs.Sub(web.StaticFS, "build")
	if err != nil {
		log.Fatal(err)
	}

	mux.Handle("/", http.FileServer(http.FS(staticContent)))
	return mux
}
