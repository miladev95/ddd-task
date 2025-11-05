package event

import (
	"fmt"
	"sync"

	"github.com/miladev95/ddd-task/domain/event"
)

// SimpleEventPublisher is a basic in-memory event publisher implementation
type SimpleEventPublisher struct {
	subscribers map[string][]func(event.DomainEvent) error
	mu          sync.RWMutex
}

// NewSimpleEventPublisher creates a new SimpleEventPublisher
func NewSimpleEventPublisher() *SimpleEventPublisher {
	return &SimpleEventPublisher{
		subscribers: make(map[string][]func(event.DomainEvent) error),
	}
}

// Publish publishes a domain event to all subscribers
func (p *SimpleEventPublisher) Publish(evt event.DomainEvent) error {
	p.mu.RLock()
	subscribers, exists := p.subscribers[evt.EventType()]
	p.mu.RUnlock()

	if !exists {
		return nil // No subscribers, but not an error
	}

	var errs []error
	for _, handler := range subscribers {
		if err := handler(evt); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to publish event: %v", errs)
	}

	return nil
}

// PublishAll publishes multiple domain events
func (p *SimpleEventPublisher) PublishAll(events []event.DomainEvent) error {
	for _, evt := range events {
		if err := p.Publish(evt); err != nil {
			return err
		}
	}
	return nil
}

// Subscribe subscribes a handler to an event type
func (p *SimpleEventPublisher) Subscribe(
	eventType string,
	handler func(event.DomainEvent) error,
) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, exists := p.subscribers[eventType]; !exists {
		p.subscribers[eventType] = make([]func(event.DomainEvent) error, 0)
	}

	p.subscribers[eventType] = append(p.subscribers[eventType], handler)
	return nil
}

// Unsubscribe unsubscribes all handlers for an event type
func (p *SimpleEventPublisher) Unsubscribe(eventType string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	delete(p.subscribers, eventType)
	return nil
}

// Ensure SimpleEventPublisher implements event.EventPublisher
var _ event.EventPublisher = (*SimpleEventPublisher)(nil)