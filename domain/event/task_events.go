package event

// TaskCreatedEvent is fired when a new task is created
type TaskCreatedEvent struct {
	BaseDomainEvent
	ProjectID   string
	Title       string
	Description string
	AssigneeID  string
	Priority    string
}

// NewTaskCreatedEvent creates a new TaskCreatedEvent
func NewTaskCreatedEvent(
	taskID, projectID, title, description, assigneeID, priority string,
) TaskCreatedEvent {
	return TaskCreatedEvent{
		BaseDomainEvent: NewBaseDomainEvent("TaskCreated", taskID, "Task"),
		ProjectID:       projectID,
		Title:           title,
		Description:     description,
		AssigneeID:      assigneeID,
		Priority:        priority,
	}
}

// TaskAssignedEvent is fired when a task is assigned to a user
type TaskAssignedEvent struct {
	BaseDomainEvent
	AssigneeID string
	PreviousAssigneeID string
}

// NewTaskAssignedEvent creates a new TaskAssignedEvent
func NewTaskAssignedEvent(taskID, assigneeID, previousAssigneeID string) TaskAssignedEvent {
	return TaskAssignedEvent{
		BaseDomainEvent: NewBaseDomainEvent("TaskAssigned", taskID, "Task"),
		AssigneeID:      assigneeID,
		PreviousAssigneeID: previousAssigneeID,
	}
}

// TaskStatusChangedEvent is fired when a task status changes
type TaskStatusChangedEvent struct {
	BaseDomainEvent
	OldStatus string
	NewStatus string
}

// NewTaskStatusChangedEvent creates a new TaskStatusChangedEvent
func NewTaskStatusChangedEvent(taskID, oldStatus, newStatus string) TaskStatusChangedEvent {
	return TaskStatusChangedEvent{
		BaseDomainEvent: NewBaseDomainEvent("TaskStatusChanged", taskID, "Task"),
		OldStatus:       oldStatus,
		NewStatus:       newStatus,
	}
}

// TaskDeadlineSetEvent is fired when a deadline is set on a task
type TaskDeadlineSetEvent struct {
	BaseDomainEvent
	DueDate string // ISO 8601 format
}

// NewTaskDeadlineSetEvent creates a new TaskDeadlineSetEvent
func NewTaskDeadlineSetEvent(taskID, dueDate string) TaskDeadlineSetEvent {
	return TaskDeadlineSetEvent{
		BaseDomainEvent: NewBaseDomainEvent("TaskDeadlineSet", taskID, "Task"),
		DueDate:         dueDate,
	}
}

// TaskOverdueEvent is fired when a task becomes overdue
type TaskOverdueEvent struct {
	BaseDomainEvent
	DaysOverdue int
}

// NewTaskOverdueEvent creates a new TaskOverdueEvent
func NewTaskOverdueEvent(taskID string, daysOverdue int) TaskOverdueEvent {
	return TaskOverdueEvent{
		BaseDomainEvent: NewBaseDomainEvent("TaskOverdue", taskID, "Task"),
		DaysOverdue:     daysOverdue,
	}
}

// TaskCompletedEvent is fired when a task is completed
type TaskCompletedEvent struct {
	BaseDomainEvent
	CompletedBy string
	CompletionTime string // ISO 8601 format
}

// NewTaskCompletedEvent creates a new TaskCompletedEvent
func NewTaskCompletedEvent(taskID, completedBy, completionTime string) TaskCompletedEvent {
	return TaskCompletedEvent{
		BaseDomainEvent: NewBaseDomainEvent("TaskCompleted", taskID, "Task"),
		CompletedBy:     completedBy,
		CompletionTime:  completionTime,
	}
}

// TaskDeletedEvent is fired when a task is deleted
type TaskDeletedEvent struct {
	BaseDomainEvent
	ProjectID string
}

// NewTaskDeletedEvent creates a new TaskDeletedEvent
func NewTaskDeletedEvent(taskID, projectID string) TaskDeletedEvent {
	return TaskDeletedEvent{
		BaseDomainEvent: NewBaseDomainEvent("TaskDeleted", taskID, "Task"),
		ProjectID:       projectID,
	}
}

// EventPublisher defines the interface for publishing domain events
type EventPublisher interface {
	Publish(event DomainEvent) error
	PublishAll(events []DomainEvent) error
}

// EventStore defines the interface for storing domain events
type EventStore interface {
	Store(event DomainEvent) error
	GetEvents(aggregateID string) ([]DomainEvent, error)
	GetEventsSince(aggregateID string, since string) ([]DomainEvent, error)
}

// EventSubscriber defines the interface for subscribing to domain events
type EventSubscriber interface {
	Subscribe(eventType string, handler func(event DomainEvent) error) error
	Unsubscribe(eventType string) error
}