package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MizuchiLabs/mantrae/pkg/util"
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

	defer func() {
		// Unregister the client when the connection is closed
		util.ClientsMutex.Lock()
		delete(util.Clients, w)
		util.ClientsMutex.Unlock()
	}()

	for {
		select {
		case message := <-util.Broadcast:
			// Serialize the EventMessage to JSON
			data, err := json.Marshal(message)
			if err != nil {
				fmt.Printf("Error marshalling message: %v\n", err)
				continue
			}
			// Send the data to the client
			fmt.Fprintf(w, "data: %s\n\n", data)
			w.(http.Flusher).Flush()
		case <-r.Context().Done():
			return
		}
	}
}
