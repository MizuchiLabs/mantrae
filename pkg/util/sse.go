package util

import "sync"

// A channel to broadcast updates
var Broadcast chan string

// A list of clients connected to SSE
var (
	Clients      = make(map[chan string]bool)
	ClientsMutex sync.Mutex
)

func init() {
	Broadcast = make(chan string)
}
