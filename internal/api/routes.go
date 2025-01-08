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

	// Routers
	register("GET", "/router/{id}", jwtChain, GetRouters)
	register("POST", "/router", jwtChain, UpsertRouter)
	register("DELETE", "/router/{id}", jwtChain, DeleteRouter)

	// Services
	register("GET", "/service/{id}", jwtChain, GetServices)
	register("POST", "/service", jwtChain, UpsertService)
	register("DELETE", "/service/{id}", jwtChain, DeleteService)

	// Middlewares
	register("GET", "/middleware/{id}", jwtChain, GetMiddlewares)
	register("POST", "/middleware", jwtChain, UpsertMiddleware)
	register("DELETE", "/middleware/{id}", jwtChain, DeleteMiddleware)
	register("GET", "/middleware/plugins", logChain, GetMiddlewarePlugins)

	// EntryPoints
	register("GET", "/entrypoint/{id}", jwtChain, GetEntryPoints)

	// Users
	register("GET", "/user", jwtChain, GetUsers)
	register("GET", "/user/{id}", jwtChain, GetUser)
	register("POST", "/user", jwtChain, UpsertUser)
	register("DELETE", "/user/{id}", jwtChain, DeleteUser)

	// DNS Provider
	register("GET", "/provider", jwtChain, GetProviders)
	register("POST", "/provider", jwtChain, UpsertProvider)
	register("DELETE", "/provider/{id}", jwtChain, DeleteProvider)

	// Extra route for deleting DNS records
	register("POST", "/dns", jwtChain, DeleteRouterDNS)

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

	// Dynamic config
	if util.App.EnableBasicAuth {
		register("GET", "/{name}", basicChain, GetTraefikConfig)
	} else {
		register("GET", "/{name}", logChain, GetTraefikConfig)
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
