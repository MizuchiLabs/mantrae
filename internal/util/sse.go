package util

import (
	"net/http"
	"sync"
)

type EventMessage struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

const (
	EventTypeError  = "error"
	EventTypeInfo   = "info"
	EventTypeCreate = "create"
	EventTypeUpdate = "update"
	EventTypeDelete = "delete"
)

// A channel to broadcast updates
var Broadcast = make(chan EventMessage)

// A list of clients connected to SSE
var (
	Clients      = make(map[http.ResponseWriter]bool)
	ClientsMutex = &sync.Mutex{}
)
