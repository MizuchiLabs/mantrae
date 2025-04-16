package util

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
)

type EventMessage struct {
	Type     string `json:"type"`
	Category string `json:"category"`
	Message  string `json:"message"`
}

const (
	// Types
	EventTypeError  = "error"
	EventTypeInfo   = "info"
	EventTypeCreate = "create"
	EventTypeUpdate = "update"
	EventTypeDelete = "delete"

	// Categories
	EventCategoryProfile = "profile"
	EventCategoryTraefik = "traefik"
	EventCategoryDNS     = "dns"
	EventCategoryUser    = "user"
	EventCategoryAgent   = "agent"
	EventCategorySetting = "setting"
)

var (
	Broadcast    = make(chan EventMessage, 100)
	SSEDone      = make(chan struct{})
	Clients      = make(map[http.ResponseWriter]bool)
	ClientsMutex = &sync.Mutex{}
)

func StartEventProcessor(ctx context.Context) {
	go func() {
		defer close(SSEDone)

		for {
			select {
			case msg := <-Broadcast:
				ClientsMutex.Lock()
				for client := range Clients {
					// Non-blocking send to each client
					go func(w http.ResponseWriter, message EventMessage) {
						if err := SendEventToClient(w, message); err != nil {
							slog.Error("Failed to send event", "error", err)
							// Remove failed clients
							ClientsMutex.Lock()
							delete(Clients, w)
							ClientsMutex.Unlock()
						}
					}(client, msg)
				}
				ClientsMutex.Unlock()

				// If no clients, log the dropped event
				if len(Clients) == 0 {
					slog.Debug("Event dropped - no clients connected",
						"type", msg.Type,
						"message", msg.Message)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

func SendEventToClient(w http.ResponseWriter, msg EventMessage) error {
	// Implementation of sending event to a single client
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "data: %s\n\n", data)
	if err != nil {
		return err
	}

	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
	return nil
}
