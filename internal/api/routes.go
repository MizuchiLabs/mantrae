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

	mux.HandleFunc("GET /api/events", GetEvents)

	mux.HandleFunc("GET /api/profile", JWT(GetProfiles))
	mux.HandleFunc("GET /api/profile/{id}", JWT(GetProfile))
	mux.HandleFunc("POST /api/profile", JWT(CreateProfile))
	mux.HandleFunc("PUT /api/profile", JWT(UpdateProfile))
	mux.HandleFunc("DELETE /api/profile/{id}", JWT(DeleteProfile))

	mux.HandleFunc("GET /api/config/{id}", JWT(GetConfig))
	mux.HandleFunc("PUT /api/config/{id}", JWT(UpdateConfig))

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

	mux.HandleFunc("GET /api/backup", JWT(DownloadBackup))
	mux.HandleFunc("POST /api/restore", JWT(UploadBackup))

	mux.HandleFunc("GET /api/{name}", GetTraefikConfig)

	staticContent, err := fs.Sub(web.StaticFS, "build")
	if err != nil {
		log.Fatal(err)
	}

	mux.Handle("/", http.FileServer(http.FS(staticContent)))

	middle := Chain(Log, Cors)

	return middle(mux)
}
