package api

import (
	"io/fs"
	"log"
	"net/http"

	"github.com/MizuchiLabs/mantrae/pkg/util"
	"github.com/MizuchiLabs/mantrae/web"
)

func Routes() http.Handler {
	mux := http.NewServeMux()

	logChain := Chain(Log)
	jwtChain := Chain(Log, JWT)
	// adminChain := Chain(Log, AdminOnly, JWT) // Handle perms later
	basicChain := Chain(Log, BasicAuth)

	mux.Handle("POST /api/login", logChain(Login))
	mux.Handle("POST /api/verify", logChain(VerifyToken))

	mux.Handle("GET /api/version", logChain(GetVersion))
	mux.Handle("GET /api/events", logChain(GetEvents))

	mux.Handle("GET /api/profile", jwtChain(GetProfiles))
	mux.Handle("GET /api/profile/{id}", jwtChain(GetProfile))
	mux.Handle("POST /api/profile", jwtChain(CreateProfile))
	mux.Handle("PUT /api/profile", jwtChain(UpdateProfile))
	mux.Handle("DELETE /api/profile/{id}", jwtChain(DeleteProfile))

	mux.Handle("GET /api/router/{id}", jwtChain(GetRouters))
	mux.Handle("POST /api/router", jwtChain(UpsertRouter))
	mux.Handle("DELETE /api/router/{id}", jwtChain(DeleteRouter))

	mux.Handle("GET /api/service/{id}", jwtChain(GetServices))
	mux.Handle("POST /api/service", jwtChain(UpsertService))
	mux.Handle("DELETE /api/service/{id}", jwtChain(DeleteService))

	mux.Handle("GET /api/middleware/{id}", jwtChain(GetMiddlewares))
	mux.Handle("POST /api/middleware", jwtChain(UpsertMiddleware))
	mux.Handle("DELETE /api/middleware/{id}", jwtChain(DeleteMiddleware))

	mux.Handle("GET /api/entrypoint/{id}", jwtChain(GetEntryPoints))
	mux.Handle("GET /api/middleware/plugins", logChain(GetMiddlewarePlugins))

	mux.Handle("GET /api/user", jwtChain(GetUsers))
	mux.Handle("GET /api/user/{id}", jwtChain(GetUser))
	mux.Handle("POST /api/user", jwtChain(UpsertUser))
	mux.Handle("DELETE /api/user/{id}", jwtChain(DeleteUser))

	mux.Handle("GET /api/provider", jwtChain(GetProviders))
	mux.Handle("GET /api/provider/{id}", jwtChain(GetProvider))
	mux.Handle("POST /api/provider", jwtChain(CreateProvider))
	mux.Handle("PUT /api/provider", jwtChain(UpdateProvider))
	mux.Handle("DELETE /api/provider/{id}", jwtChain(DeleteProvider))
	mux.Handle("POST /api/dns", jwtChain(DeleteRouterDNS)) // Extra route for deleting DNS records

	mux.Handle("GET /api/settings", jwtChain(GetSettings))
	mux.Handle("GET /api/settings/{key}", jwtChain(GetSetting))
	mux.Handle("PUT /api/settings", jwtChain(UpdateSetting))

	mux.Handle("GET /api/agent/{id}", jwtChain(GetAgents))
	mux.Handle("GET /api/agent/token/{id}", jwtChain(GetAgentToken))
	mux.Handle("PUT /api/agent/{id}", jwtChain(UpsertAgent))
	mux.Handle("DELETE /api/agent/{id}/{type}", jwtChain(DeleteAgent))

	mux.Handle("GET /api/ip/{id}", jwtChain(GetPublicIP))

	mux.Handle("GET /api/backup", jwtChain(DownloadBackup))
	mux.Handle("POST /api/restore", jwtChain(UploadBackup))

	mux.Handle("GET /api/traefik/{id}", jwtChain(GetTraefikOverview))

	if util.App.EnableBasicAuth {
		mux.Handle("GET /api/{name}", basicChain(GetTraefikConfig))
	} else {
		mux.Handle("GET /api/{name}", logChain(GetTraefikConfig))
	}

	staticContent, err := fs.Sub(web.StaticFS, "build")
	if err != nil {
		log.Fatal(err)
	}

	mux.Handle("/", http.FileServer(http.FS(staticContent)))

	return Cors(mux)
}
