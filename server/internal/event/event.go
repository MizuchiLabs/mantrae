package event

import (
	"context"
	"log/slog"
	"sync"

	mantraev1 "github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1"
)

type Broadcaster struct {
	mu      sync.RWMutex
	clients map[int]*subscriber
	nextID  int
	ctx     context.Context
	cancel  context.CancelFunc
}

type subscriber struct {
	profileID int64
	ch        chan *mantraev1.EventStreamResponse
}

func NewBroadcaster(parent context.Context) *Broadcaster {
	ctx, cancel := context.WithCancel(parent)
	b := &Broadcaster{
		clients: make(map[int]*subscriber),
		ctx:     ctx,
		cancel:  cancel,
	}

	go b.cleanup()
	return b
}

func (b *Broadcaster) Subscribe(
	profileID int64,
) (id int, ch <-chan *mantraev1.EventStreamResponse) {
	b.mu.Lock()
	defer b.mu.Unlock()
	id = b.nextID
	b.nextID++
	c := make(chan *mantraev1.EventStreamResponse, 32)
	b.clients[id] = &subscriber{
		profileID: profileID,
		ch:        c,
	}
	return id, c
}

func (b *Broadcaster) Unsubscribe(id int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if sub, ok := b.clients[id]; ok {
		close(sub.ch)
		delete(b.clients, id)
	}
}

func (b *Broadcaster) Broadcast(event *mantraev1.EventStreamResponse) {
	profileID, isGlobal := getProfileIDFromEvent(event)

	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, sub := range b.clients {
		if !isGlobal && sub.profileID != profileID {
			continue
		}
		select {
		case sub.ch <- event:
		default:
			close(sub.ch)
		}
	}
}

func (b *Broadcaster) cleanup() {
	<-b.ctx.Done()
	slog.Info("Broadcaster is shutting down")
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, client := range b.clients {
		close(client.ch)
	}
	b.clients = nil
}

func getProfileIDFromEvent(event *mantraev1.EventStreamResponse) (int64, bool) {
	switch d := event.Data.(type) {
	case *mantraev1.EventStreamResponse_Profile:
		return d.Profile.Id, false
	case *mantraev1.EventStreamResponse_Agent:
		return d.Agent.ProfileId, false
	case *mantraev1.EventStreamResponse_EntryPoint:
		return d.EntryPoint.ProfileId, false
	case *mantraev1.EventStreamResponse_Router:
		return d.Router.ProfileId, false
	case *mantraev1.EventStreamResponse_Service:
		return d.Service.ProfileId, false
	case *mantraev1.EventStreamResponse_Middleware:
		return d.Middleware.ProfileId, false
	case *mantraev1.EventStreamResponse_ServersTransport:
		return d.ServersTransport.ProfileId, false
	case *mantraev1.EventStreamResponse_TraefikInstance:
		return d.TraefikInstance.ProfileId, false
	case *mantraev1.EventStreamResponse_DnsProvider:
		return 0, true
	case *mantraev1.EventStreamResponse_User:
		return 0, true
	default:
		return 0, false
	}
}
