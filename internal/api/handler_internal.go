package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/MizuchiLabs/mantrae/internal/db"
	"github.com/traefik/traefik/v3/pkg/config/dynamic"
)

// General ---------------------------------------------------------------------

func GetInternalTraefikConfig(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}
	config, err := db.Query.GetInternalTraefikConfigByProfileID(r.Context(), id)
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

// Internal Router -------------------------------------------------------------

func GetInternalHTTPRouters(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}
	router, err := db.Query.GetInternalHTTPRoutersByProfileID(r.Context(), id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Query error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	writeJSON(w, router)
}

func GetInternalTCPRouters(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}
	router, err := db.Query.GetInternalTCPRoutersByProfileID(r.Context(), id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Query error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	writeJSON(w, router)
}

func GetInternalUDPRouters(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}
	router, err := db.Query.GetInternalUDPRoutersByProfileID(r.Context(), id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Query error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	writeJSON(w, router)
}

func GetInternalHTTPRouter(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}
	name := r.PathValue("name")
	router, err := db.Query.GetInternalHTTPRouterByName(r.Context(), db.GetInternalHTTPRouterByNameParams{
		ProfileID: id,
		Name:      &name,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Query error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	writeJSON(w, router)
}

func GetInternalTCPRouter(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}
	name := r.PathValue("name")
	router, err := db.Query.GetInternalTCPRouterByName(r.Context(), db.GetInternalTCPRouterByNameParams{
		ProfileID: id,
		Name:      &name,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Query error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	writeJSON(w, router)
}

func GetInternalUDPRouter(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}
	name := r.PathValue("name")
	router, err := db.Query.GetInternalUDPRouterByName(r.Context(), db.GetInternalUDPRouterByNameParams{
		ProfileID: id,
		Name:      &name,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Query error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	writeJSON(w, router)
}

func UpsertInternalHTTPRouter(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}
	name := r.PathValue("name")

	var router dynamic.Router
	if err := json.NewDecoder(r.Body).Decode(&router); err != nil {
		http.Error(w, fmt.Sprintf("Decode error: %s", err.Error()), http.StatusBadRequest)
		return
	}
	routerJSON, err := json.Marshal(router)
	if err != nil {
		http.Error(w, fmt.Sprintf("Marshal error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	err = db.Query.UpsertInternalHTTPRouter(r.Context(), db.UpsertInternalHTTPRouterParams{
		ProfileID: id,
		Name:      &name,
		Body:      string(routerJSON),
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Upsert error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func UpsertInternalTCPRouter(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}

	var router dynamic.TCPRouter
	if err := json.NewDecoder(r.Body).Decode(&router); err != nil {
		http.Error(w, fmt.Sprintf("Decode error: %s", err.Error()), http.StatusBadRequest)
		return
	}

	routerJSON, err := json.Marshal(router)
	if err != nil {
		http.Error(w, fmt.Sprintf("Marshal error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	if err := db.Query.UpsertInternalTCPRouter(r.Context(), db.UpsertInternalTCPRouterParams{
		ProfileID: id,
		Name:      &name,
		Body:      string(routerJSON),
	}); err != nil {
		http.Error(w, fmt.Sprintf("Upsert error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func UpsertInternalUDPRouter(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}

	var router dynamic.UDPRouter
	if err := json.NewDecoder(r.Body).Decode(&router); err != nil {
		http.Error(w, fmt.Sprintf("Decode error: %s", err.Error()), http.StatusBadRequest)
		return
	}

	routerJSON, err := json.Marshal(router)
	if err != nil {
		http.Error(w, fmt.Sprintf("Marshal error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	if err := db.Query.UpsertInternalUDPRouter(r.Context(), db.UpsertInternalUDPRouterParams{
		ProfileID: id,
		Name:      &name,
		Body:      string(routerJSON),
	}); err != nil {
		http.Error(w, fmt.Sprintf("Upsert error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteInternalHTTPRouter(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}

	if err := db.Query.DeleteInternalHTTPRouter(r.Context(), db.DeleteInternalHTTPRouterParams{
		ProfileID: id,
		Name:      &name,
	}); err != nil {
		http.Error(w, fmt.Sprintf("Delete error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteInternalTCPRouter(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}

	if err := db.Query.DeleteInternalTCPRouter(r.Context(), db.DeleteInternalTCPRouterParams{
		ProfileID: id,
		Name:      &name,
	}); err != nil {
		http.Error(w, fmt.Sprintf("Delete error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteInternalUDPRouter(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}

	if err := db.Query.DeleteInternalUDPRouter(r.Context(), db.DeleteInternalUDPRouterParams{
		ProfileID: id,
		Name:      &name,
	}); err != nil {
		http.Error(w, fmt.Sprintf("Delete error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetInternalHTTPServices(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}
	services, err := db.Query.GetInternalHTTPServicesByProfileID(r.Context(), id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Query error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	writeJSON(w, services)
}

func GetInternalTCPServices(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}
	services, err := db.Query.GetInternalTCPServicesByProfileID(r.Context(), id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Query error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	writeJSON(w, services)
}

func GetInternalUDPServices(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}
	services, err := db.Query.GetInternalUDPServicesByProfileID(r.Context(), id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Query error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	writeJSON(w, services)
}

func GetInternalHTTPService(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}
	name := r.PathValue("name")
	service, err := db.Query.GetInternalHTTPServiceByName(r.Context(), db.GetInternalHTTPServiceByNameParams{
		ProfileID: id,
		Name:      &name,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Query error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	writeJSON(w, service)
}

func GetInternalTCPService(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}
	name := r.PathValue("name")
	service, err := db.Query.GetInternalTCPServiceByName(r.Context(), db.GetInternalTCPServiceByNameParams{
		ProfileID: id,
		Name:      &name,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Query error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	writeJSON(w, service)
}

func GetInternalUDPService(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}
	name := r.PathValue("name")
	service, err := db.Query.GetInternalUDPServiceByName(r.Context(), db.GetInternalUDPServiceByNameParams{
		ProfileID: id,
		Name:      &name,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Query error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	writeJSON(w, service)
}

func UpsertInternalHTTPService(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}

	var service dynamic.Service
	if err := json.NewDecoder(r.Body).Decode(&service); err != nil {
		http.Error(w, fmt.Sprintf("Decode error: %s", err.Error()), http.StatusBadRequest)
		return
	}

	serviceJSON, err := json.Marshal(service)
	if err != nil {
		http.Error(w, fmt.Sprintf("Marshal error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	if err := db.Query.UpsertInternalHTTPService(r.Context(), db.UpsertInternalHTTPServiceParams{
		ProfileID: id,
		Name:      &name,
		Body:      string(serviceJSON),
	}); err != nil {
		http.Error(w, fmt.Sprintf("Upsert error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func UpsertInternalTCPService(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}

	var service dynamic.TCPService
	if err := json.NewDecoder(r.Body).Decode(&service); err != nil {
		http.Error(w, fmt.Sprintf("Decode error: %s", err.Error()), http.StatusBadRequest)
		return
	}

	serviceJSON, err := json.Marshal(service)
	if err != nil {
		http.Error(w, fmt.Sprintf("Marshal error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	if err := db.Query.UpsertInternalTCPService(r.Context(), db.UpsertInternalTCPServiceParams{
		ProfileID: id,
		Name:      &name,
		Body:      string(serviceJSON),
	}); err != nil {
		http.Error(w, fmt.Sprintf("Upsert error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func UpsertInternalUDPService(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}

	var service dynamic.UDPService
	if err := json.NewDecoder(r.Body).Decode(&service); err != nil {
		http.Error(w, fmt.Sprintf("Decode error: %s", err.Error()), http.StatusBadRequest)
		return
	}

	serviceJSON, err := json.Marshal(service)
	if err != nil {
		http.Error(w, fmt.Sprintf("Marshal error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	if err := db.Query.UpsertInternalUDPService(r.Context(), db.UpsertInternalUDPServiceParams{
		ProfileID: id,
		Name:      &name,
		Body:      string(serviceJSON),
	}); err != nil {
		http.Error(w, fmt.Sprintf("Upsert error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteInternalHTTPService(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}

	if err := db.Query.DeleteInternalHTTPService(r.Context(), db.DeleteInternalHTTPServiceParams{
		ProfileID: id,
		Name:      &name,
	}); err != nil {
		http.Error(w, fmt.Sprintf("Delete error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteInternalTCPService(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}

	if err := db.Query.DeleteInternalTCPService(r.Context(), db.DeleteInternalTCPServiceParams{
		ProfileID: id,
		Name:      &name,
	}); err != nil {
		http.Error(w, fmt.Sprintf("Delete error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteInternalUDPService(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}

	if err := db.Query.DeleteInternalUDPService(r.Context(), db.DeleteInternalUDPServiceParams{
		ProfileID: id,
		Name:      &name,
	}); err != nil {
		http.Error(w, fmt.Sprintf("Delete error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetInternalHTTPMiddlewares(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}
	middlewares, err := db.Query.GetInternalHTTPMiddlewaresByProfileID(r.Context(), id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Query error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	writeJSON(w, middlewares)
}

func GetInternalTCPMiddlewares(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}
	middlewares, err := db.Query.GetInternalTCPMiddlewaresByProfileID(r.Context(), id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Query error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	writeJSON(w, middlewares)
}

func GetInternalHTTPMiddleware(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}
	name := r.PathValue("name")
	middleware, err := db.Query.GetInternalHTTPMiddlewareByName(r.Context(), db.GetInternalHTTPMiddlewareByNameParams{
		ProfileID: id,
		Name:      &name,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Query error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	writeJSON(w, middleware)
}

func GetInternalTCPMiddleware(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}
	name := r.PathValue("name")
	middleware, err := db.Query.GetInternalTCPMiddlewareByName(r.Context(), db.GetInternalTCPMiddlewareByNameParams{
		ProfileID: id,
		Name:      &name,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Query error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	writeJSON(w, middleware)
}

func UpsertInternalHTTPMiddleware(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}

	var middleware dynamic.Middleware
	if err := json.NewDecoder(r.Body).Decode(&middleware); err != nil {
		http.Error(w, fmt.Sprintf("Decode error: %s", err.Error()), http.StatusBadRequest)
		return
	}

	middlewareJSON, err := json.Marshal(middleware)
	if err != nil {
		http.Error(w, fmt.Sprintf("Marshal error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	if err := db.Query.UpsertInternalHTTPMiddleware(r.Context(), db.UpsertInternalHTTPMiddlewareParams{
		ProfileID: id,
		Name:      &name,
		Body:      string(middlewareJSON),
	}); err != nil {
		http.Error(w, fmt.Sprintf("Upsert error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func UpsertInternalTCPMiddleware(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}

	var middleware dynamic.TCPMiddleware
	if err := json.NewDecoder(r.Body).Decode(&middleware); err != nil {
		http.Error(w, fmt.Sprintf("Decode error: %s", err.Error()), http.StatusBadRequest)
		return
	}

	middlewareJSON, err := json.Marshal(middleware)
	if err != nil {
		http.Error(w, fmt.Sprintf("Marshal error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	if err := db.Query.UpsertInternalTCPMiddleware(r.Context(), db.UpsertInternalTCPMiddlewareParams{
		ProfileID: id,
		Name:      &name,
		Body:      string(middlewareJSON),
	}); err != nil {
		http.Error(w, fmt.Sprintf("Upsert error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteInternalHTTPMiddleware(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}

	if err := db.Query.DeleteInternalHTTPMiddleware(r.Context(), db.DeleteInternalHTTPMiddlewareParams{
		ProfileID: id,
		Name:      &name,
	}); err != nil {
		http.Error(w, fmt.Sprintf("Delete error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteInternalTCPMiddleware(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse error: %s", err.Error()), http.StatusNotFound)
		return
	}

	if err := db.Query.DeleteInternalTCPMiddleware(r.Context(), db.DeleteInternalTCPMiddlewareParams{
		ProfileID: id,
		Name:      &name,
	}); err != nil {
		http.Error(w, fmt.Sprintf("Delete error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
