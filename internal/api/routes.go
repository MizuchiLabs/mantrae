package api

import "net/http"

func Routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/login", Login)
	mux.HandleFunc("POST /api/verify", VerifyToken)

	mux.HandleFunc("POST /api/profiles", JWT(CreateProfile))
	mux.HandleFunc("GET /api/profiles", JWT(GetProfiles))
	mux.HandleFunc("PUT /api/profiles/{name}", JWT(UpdateProfile))
	mux.HandleFunc("DELETE /api/profiles/{name}", JWT(DeleteProfile))

	mux.HandleFunc("PUT /api/routers/{profile}/{router}", JWT(UpdateRouter))
	mux.HandleFunc("DELETE /api/routers/{profile}/{router}", JWT(DeleteRouter))

	mux.HandleFunc("PUT /api/services/{profile}/{service}", JWT(UpdateService))
	mux.HandleFunc("DELETE /api/services/{profile}/{service}", JWT(DeleteService))

	mux.HandleFunc("PUT /api/middlewares/{profile}/{middleware}", JWT(UpdateMiddleware))
	mux.HandleFunc("DELETE /api/middlewares/{profile}/{middleware}", JWT(DeleteMiddleware))

	mux.HandleFunc("GET /api/{name}", GetConfig)

	return mux
}
