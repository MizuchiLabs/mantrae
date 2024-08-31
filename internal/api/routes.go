package api

import (
	"io/fs"
	"log"
	"net/http"

	"github.com/MizuchiLabs/mantrae/web"
)

func Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/login", Login)
	mux.HandleFunc("POST /api/verify", VerifyToken)

	mux.HandleFunc("GET /api/profiles", JWT(GetProfiles))
	mux.HandleFunc("GET /api/profile/{name}", JWT(GetProfile))
	mux.HandleFunc("POST /api/profiles", JWT(CreateProfile))
	mux.HandleFunc("PUT /api/profiles/{name}", JWT(UpdateProfile))
	mux.HandleFunc("DELETE /api/profiles/{name}", JWT(DeleteProfile))

	mux.HandleFunc("GET /api/providers", JWT(GetProviders))
	mux.HandleFunc("PUT /api/providers/{name}", JWT(UpdateProvider))
	mux.HandleFunc("DELETE /api/providers/{name}", JWT(DeleteProvider))

	mux.HandleFunc("PUT /api/routers/{profile}/{router}", JWT(UpdateRouter))
	mux.HandleFunc("DELETE /api/routers/{profile}/{router}", JWT(DeleteRouter))

	mux.HandleFunc("PUT /api/services/{profile}/{service}", JWT(UpdateService))
	mux.HandleFunc("DELETE /api/services/{profile}/{service}", JWT(DeleteService))

	mux.HandleFunc("PUT /api/middlewares/{profile}/{middleware}", JWT(UpdateMiddleware))
	mux.HandleFunc("DELETE /api/middlewares/{profile}/{middleware}", JWT(DeleteMiddleware))

	mux.HandleFunc("GET /api/{name}", GetConfig)

	staticContent, err := fs.Sub(web.StaticFS, "build")
	if err != nil {
		log.Fatal(err)
	}

	mux.Handle("/", http.FileServer(http.FS(staticContent)))

	middle := Chain(Log, Cors)

	return middle(mux)
}
