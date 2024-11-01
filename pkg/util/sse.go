package util

import (
	"net/http"
	"sync"
)

type EventMessage struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

// A channel to broadcast updates
var Broadcast = make(chan EventMessage)

// A list of clients connected to SSE
var (
	Clients      = make(map[http.ResponseWriter]bool)
	ClientsMutex = &sync.Mutex{}
)
