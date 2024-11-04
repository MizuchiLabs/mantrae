package api

import (
	"io/fs"
	"log"
	"net/http"

	"github.com/MizuchiLabs/mantrae/web"
)

func Routes(useAuth bool) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/login", Login)
	mux.HandleFunc("POST /api/verify", VerifyToken)

	mux.HandleFunc("GET /api/version", GetVersion)
	mux.HandleFunc("GET /api/events", GetEvents)

	mux.HandleFunc("GET /api/profile", JWT(GetProfiles))
	mux.HandleFunc("GET /api/profile/{id}", JWT(GetProfile))
	mux.HandleFunc("POST /api/profile", JWT(CreateProfile))
	mux.HandleFunc("PUT /api/profile", JWT(UpdateProfile))
	mux.HandleFunc("DELETE /api/profile/{id}", JWT(DeleteProfile))

	mux.HandleFunc("GET /api/router/{id}", JWT(GetRouters))
	mux.HandleFunc("POST /api/router", JWT(UpsertRouter))
	mux.HandleFunc("DELETE /api/router/{id}", JWT(DeleteRouter))

	mux.HandleFunc("GET /api/service/{id}", JWT(GetServices))
	mux.HandleFunc("POST /api/service", JWT(UpsertService))
	mux.HandleFunc("DELETE /api/service/{id}", JWT(DeleteService))

	mux.HandleFunc("GET /api/middleware/{id}", JWT(GetMiddlewares))
	mux.HandleFunc("POST /api/middleware", JWT(UpsertMiddleware))
	mux.HandleFunc("DELETE /api/middleware/{id}", JWT(DeleteMiddleware))

	mux.HandleFunc("GET /api/entrypoint/{id}", JWT(GetEntryPoints))
	mux.HandleFunc("GET /api/middleware/plugins", GetMiddlewarePlugins)

	mux.HandleFunc("GET /api/user", JWT(GetUsers))
	mux.HandleFunc("GET /api/user/{id}", JWT(GetUser))
	mux.HandleFunc("POST /api/user", JWT(CreateUser))
	mux.HandleFunc("PUT /api/user", JWT(UpdateUser))
	mux.HandleFunc("DELETE /api/user/{id}", JWT(DeleteUser))

	mux.HandleFunc("GET /api/provider", JWT(GetProviders))
	mux.HandleFunc("GET /api/provider/{id}", JWT(GetProvider))
	mux.HandleFunc("POST /api/provider", JWT(CreateProvider))
	mux.HandleFunc("PUT /api/provider", JWT(UpdateProvider))
	mux.HandleFunc("DELETE /api/provider/{id}", JWT(DeleteProvider))
	mux.HandleFunc("POST /api/dns", JWT(DeleteRouterDNS)) // Extra route for deleting DNS records

	mux.HandleFunc("GET /api/settings", JWT(GetSettings))
	mux.HandleFunc("GET /api/settings/{key}", JWT(GetSetting))
	mux.HandleFunc("PUT /api/settings", JWT(UpdateSetting))

	mux.HandleFunc("GET /api/agent/{id}", JWT(GetAgents))
	mux.HandleFunc("PUT /api/agent/{id}", JWT(UpsertAgent))
	mux.HandleFunc("DELETE /api/agent/{id}/{type}", JWT(DeleteAgent))
	mux.HandleFunc("POST /api/agent/token", JWT(GetAgentToken))

	mux.HandleFunc("GET /api/ip/{id}", JWT(GetPublicIP))

	mux.HandleFunc("GET /api/backup", JWT(DownloadBackup))
	mux.HandleFunc("POST /api/restore", JWT(UploadBackup))

	mux.HandleFunc("GET /api/traefik/{id}", JWT(GetTraefikOverview))

	if useAuth {
		mux.HandleFunc("GET /api/{name}", BasicAuth(GetTraefikConfig))
	} else {
		mux.HandleFunc("GET /api/{name}", GetTraefikConfig)
	}

	staticContent, err := fs.Sub(web.StaticFS, "build")
	if err != nil {
		log.Fatal(err)
	}

	mux.Handle("/", http.FileServer(http.FS(staticContent)))

	middle := Chain(Log, Cors)

	return middle(mux)
}
