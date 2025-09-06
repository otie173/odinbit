package events

import "sync"

type Event struct{}

type EventBuffer interface {
	Push(event Event)
	PopAll() []Event
}

type buffer struct {
	events []Event
	mu     sync.Mutex
}

func NewBuffer(capacity int) *buffer {
	events := make([]Event, capacity)

	return &buffer{
		events: events,
	}
}

func (b *buffer) Push(event Event) {
	b.mu.Lock()
	b.events = append(b.events, event)
	b.mu.Unlock()
}

func (b *buffer) PopAll() []Event {
	b.mu.Lock()
	events := b.events
	b.events = b.events[:0]
	b.mu.Unlock()
	return events
}
