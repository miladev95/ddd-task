package event

import "time"

// DomainEvent is the interface that all domain events must implement
type DomainEvent interface {
	EventType() string
	OccurredAt() time.Time
	AggregateID() string
	AggregateType() string
}

// BaseDomainEvent provides common functionality for domain events
type BaseDomainEvent struct {
	eventType     string
	occurredAt    time.Time
	aggregateID   string
	aggregateType string
}

// NewBaseDomainEvent creates a new base domain event
func NewBaseDomainEvent(eventType, aggregateID, aggregateType string) BaseDomainEvent {
	return BaseDomainEvent{
		eventType:     eventType,
		occurredAt:    time.Now(),
		aggregateID:   aggregateID,
		aggregateType: aggregateType,
	}
}

// EventType returns the event type
func (b BaseDomainEvent) EventType() string {
	return b.eventType
}

// OccurredAt returns when the event occurred
func (b BaseDomainEvent) OccurredAt() time.Time {
	return b.occurredAt
}

// AggregateID returns the aggregate ID
func (b BaseDomainEvent) AggregateID() string {
	return b.aggregateID
}

// AggregateType returns the aggregate type
func (b BaseDomainEvent) AggregateType() string {
	return b.aggregateType
}