package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/MizuchiLabs/mantrae/internal/db"
)

// General ---------------------------------------------------------------------

func GetTraefikEntrypoints(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}
	entrypoints, err := db.Query.GetTraefikEntrypointsByProfileID(r.Context(), id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	writeJSON(w, entrypoints)
}

func GetTraefikOverview(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}
	overview, err := db.Query.GetTraefikOverviewByProfileID(r.Context(), id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	writeJSON(w, overview)
}

func GetExternalTraefikConfig(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}
	config, err := db.Query.GetExternalTraefikConfigByProfileID(r.Context(), id)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Failed to get config: %s", err.Error()),
			http.StatusInternalServerError,
		)
		return
	}
	writeJSON(w, config)
}

// External Router -------------------------------------------------------------

func GetExternalHTTPRouters(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}
	router, err := db.Query.GetExternalHTTPRoutersByProfileID(r.Context(), id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Query error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	writeJSON(w, router)
}

func GetExternalTCPRouters(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}
	router, err := db.Query.GetExternalTCPRoutersByProfileID(r.Context(), id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Query error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	writeJSON(w, router)
}

func GetExternalUDPRouters(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}
	router, err := db.Query.GetExternalUDPRoutersByProfileID(r.Context(), id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Query error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	writeJSON(w, router)
}

func GetExternalHTTPRouter(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}
	name := r.PathValue("name")
	router, err := db.Query.GetExternalHTTPRouterByName(r.Context(), db.GetExternalHTTPRouterByNameParams{
		ProfileID: id,
		Name:      &name,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Query error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	writeJSON(w, router)
}

func GetExternalTCPRouter(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}
	name := r.PathValue("name")
	router, err := db.Query.GetExternalTCPRouterByName(r.Context(), db.GetExternalTCPRouterByNameParams{
		ProfileID: id,
		Name:      &name,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Query error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	writeJSON(w, router)
}

func GetExternalUDPRouter(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}
	name := r.PathValue("name")
	router, err := db.Query.GetExternalUDPRouterByName(r.Context(), db.GetExternalUDPRouterByNameParams{
		ProfileID: id,
		Name:      &name,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Query error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	writeJSON(w, router)
}

func GetExternalHTTPServices(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}
	services, err := db.Query.GetExternalHTTPServicesByProfileID(r.Context(), id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Query error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	writeJSON(w, services)
}

func GetExternalTCPServices(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}
	services, err := db.Query.GetExternalTCPServicesByProfileID(r.Context(), id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Query error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	writeJSON(w, services)
}

func GetExternalUDPServices(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}
	services, err := db.Query.GetExternalUDPServicesByProfileID(r.Context(), id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Query error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	writeJSON(w, services)
}

func GetExternalHTTPService(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}
	name := r.PathValue("name")
	service, err := db.Query.GetExternalHTTPServiceByName(r.Context(), db.GetExternalHTTPServiceByNameParams{
		ProfileID: id,
		Name:      &name,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Query error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	writeJSON(w, service)
}

func GetExternalTCPService(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}
	name := r.PathValue("name")
	service, err := db.Query.GetExternalTCPServiceByName(r.Context(), db.GetExternalTCPServiceByNameParams{
		ProfileID: id,
		Name:      &name,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Query error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	writeJSON(w, service)
}

func GetExternalUDPService(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}
	name := r.PathValue("name")
	service, err := db.Query.GetExternalUDPServiceByName(r.Context(), db.GetExternalUDPServiceByNameParams{
		ProfileID: id,
		Name:      &name,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Query error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	writeJSON(w, service)
}

func GetExternalHTTPMiddlewares(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}
	middlewares, err := db.Query.GetExternalHTTPMiddlewaresByProfileID(r.Context(), id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Query error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	writeJSON(w, middlewares)
}

func GetExternalTCPMiddlewares(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}
	middlewares, err := db.Query.GetExternalTCPMiddlewaresByProfileID(r.Context(), id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Query error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	writeJSON(w, middlewares)
}

func GetExternalHTTPMiddleware(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}
	name := r.PathValue("name")
	middleware, err := db.Query.GetExternalHTTPMiddlewareByName(r.Context(), db.GetExternalHTTPMiddlewareByNameParams{
		ProfileID: id,
		Name:      &name,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Query error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	writeJSON(w, middleware)
}

func GetExternalTCPMiddleware(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}
	name := r.PathValue("name")
	middleware, err := db.Query.GetExternalTCPMiddlewareByName(r.Context(), db.GetExternalTCPMiddlewareByNameParams{
		ProfileID: id,
		Name:      &name,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Query error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	writeJSON(w, middleware)
}
