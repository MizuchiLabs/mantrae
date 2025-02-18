package handler

import (
	"net/http"

	"github.com/MizuchiLabs/mantrae/internal/util"
)

// GetEvents streams server-sent events (SSE) for real-time updates.
func GetEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Register the client to receive updates
	util.ClientsMutex.Lock()
	util.Clients[w] = true
	util.ClientsMutex.Unlock()

	clientDone := make(chan struct{})
	defer func() {
		// Unregister the client when the connection is closed
		util.ClientsMutex.Lock()
		delete(util.Clients, w)
		util.ClientsMutex.Unlock()
	}()

	select {
	case <-r.Context().Done():
		return
	case <-util.SSEDone:
		return
	case <-clientDone:
		return
	}
}
