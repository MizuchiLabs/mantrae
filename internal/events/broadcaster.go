// Package events provides a centralized event broadcasting system.
package events

import (
	"slices"
	"sync"

	mantraev1 "github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1"
)

type EventBroadcaster struct {
	mu             sync.RWMutex
	profileClients map[int64]map[chan *mantraev1.ProfileEvent]*EventFilter
	globalClients  map[chan *mantraev1.GlobalEvent]*GlobalEventFilter
}

type EventFilter struct {
	ProfileID     int64
	ResourceTypes []mantraev1.ResourceType
}

type GlobalEventFilter struct {
	ResourceTypes []mantraev1.ResourceType
}

func NewEventBroadcaster() *EventBroadcaster {
	return &EventBroadcaster{
		profileClients: make(map[int64]map[chan *mantraev1.ProfileEvent]*EventFilter),
		globalClients:  make(map[chan *mantraev1.GlobalEvent]*GlobalEventFilter),
	}
}

func (b *EventBroadcaster) RegisterProfileClient(
	filter *EventFilter,
	ch chan *mantraev1.ProfileEvent,
) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.profileClients[filter.ProfileID] == nil {
		b.profileClients[filter.ProfileID] = make(map[chan *mantraev1.ProfileEvent]*EventFilter)
	}
	b.profileClients[filter.ProfileID][ch] = filter
}

func (b *EventBroadcaster) UnregisterProfileClient(ch chan *mantraev1.ProfileEvent) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for profileID, clients := range b.profileClients {
		if _, exists := clients[ch]; exists {
			delete(clients, ch)
			if len(clients) == 0 {
				delete(b.profileClients, profileID)
			}
			close(ch)
			return
		}
	}
}

func (b *EventBroadcaster) RegisterGlobalClient(
	filter *GlobalEventFilter,
	ch chan *mantraev1.GlobalEvent,
) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.globalClients[ch] == nil {
		b.globalClients[ch] = filter
	}
}

func (b *EventBroadcaster) UnregisterGlobalClient(ch chan *mantraev1.GlobalEvent) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if _, exists := b.globalClients[ch]; exists {
		delete(b.globalClients, ch)
		close(ch)
	}
}

func (b *EventBroadcaster) BroadcastProfileEvent(event *mantraev1.ProfileEvent) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if clients, exists := b.profileClients[event.ProfileId]; exists {
		for ch, filter := range clients {
			// Apply filtering
			if b.matchesFilter(event, filter) {
				select {
				case ch <- event:
				default:
					// Channel is full, skip
				}
			}
		}
	}
}

func (b *EventBroadcaster) BroadcastGlobalEvent(event *mantraev1.GlobalEvent) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for ch, filter := range b.globalClients {
		// Apply filtering for global events
		if b.matchesGlobalFilter(event, filter) {
			select {
			case ch <- event:
			default:
				// Channel is full, skip
			}
		}
	}
}

func (b *EventBroadcaster) matchesFilter(event *mantraev1.ProfileEvent, filter *EventFilter) bool {
	// Filter by resource types if specified
	if len(filter.ResourceTypes) > 0 {
		return slices.Contains(filter.ResourceTypes, event.ResourceType)
	}

	return true
}

func (b *EventBroadcaster) matchesGlobalFilter(
	event *mantraev1.GlobalEvent,
	filter *GlobalEventFilter,
) bool {
	if len(filter.ResourceTypes) > 0 {
		return slices.Contains(filter.ResourceTypes, event.ResourceType)
	}
	return true
}
