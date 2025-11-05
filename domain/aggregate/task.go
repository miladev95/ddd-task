package aggregate

import (
	"fmt"
	"time"

	"github.com/example/task-management/domain/entity"
	"github.com/example/task-management/domain/event"
	"github.com/example/task-management/domain/value"
)

// Task is the aggregate root for the Task aggregate
type Task struct {
	id          value.TaskID
	projectID   value.ProjectID
	title       string
	description string
	status      value.TaskStatus
	priority    value.Priority
	assignee    *entity.Assignment
	deadline    *value.Deadline
	comments    []*entity.Comment
	createdAt   time.Time
	updatedAt   time.Time
	createdBy   value.UserID
	domainEvents []event.DomainEvent
}

// NewTask creates a new Task
func NewTask(
	id value.TaskID,
	projectID value.ProjectID,
	title, description string,
	priority value.Priority,
	createdBy value.UserID,
) (*Task, error) {
	if title == "" {
		return nil, fmt.Errorf("task title cannot be empty")
	}

	if !priority.IsValid() {
		return nil, fmt.Errorf("invalid priority")
	}

	task := &Task{
		id:           id,
		projectID:    projectID,
		title:        title,
		description:  description,
		status:       value.TaskStatusToDo,
		priority:     priority,
		comments:     make([]*entity.Comment, 0),
		createdAt:    time.Now(),
		updatedAt:    time.Now(),
		createdBy:    createdBy,
		domainEvents: make([]event.DomainEvent, 0),
	}

	// Raise domain event
	createdEvent := event.NewTaskCreatedEvent(
		id.Value(),
		projectID.Value(),
		title,
		description,
		"", // no assignee yet
		priority.Value(),
	)
	task.domainEvents = append(task.domainEvents, createdEvent)

	return task, nil
}

// ID returns the task ID
func (t *Task) ID() value.TaskID {
	return t.id
}

// ProjectID returns the project ID
func (t *Task) ProjectID() value.ProjectID {
	return t.projectID
}

// Title returns the task title
func (t *Task) Title() string {
	return t.title
}

// Description returns the task description
func (t *Task) Description() string {
	return t.description
}

// Status returns the task status
func (t *Task) Status() value.TaskStatus {
	return t.status
}

// Priority returns the task priority
func (t *Task) Priority() value.Priority {
	return t.priority
}

// Assignee returns the assignment if any
func (t *Task) Assignee() *entity.Assignment {
	return t.assignee
}

// Deadline returns the deadline if any
func (t *Task) Deadline() *value.Deadline {
	return t.deadline
}

// Comments returns all comments
func (t *Task) Comments() []*entity.Comment {
	return append([]*entity.Comment{}, t.comments...)
}

// CreatedAt returns when the task was created
func (t *Task) CreatedAt() time.Time {
	return t.createdAt
}

// UpdatedAt returns when the task was last updated
func (t *Task) UpdatedAt() time.Time {
	return t.updatedAt
}

// CreatedBy returns who created the task
func (t *Task) CreatedBy() value.UserID {
	return t.createdBy
}

// DomainEvents returns all uncommitted domain events
func (t *Task) DomainEvents() []event.DomainEvent {
	return append([]event.DomainEvent{}, t.domainEvents...)
}

// ClearDomainEvents clears all domain events after they have been published
func (t *Task) ClearDomainEvents() {
	t.domainEvents = make([]event.DomainEvent, 0)
}

// Assign assigns the task to a user
func (t *Task) Assign(assigneeID value.UserID, assignedBy value.UserID) error {
	previousAssigneeID := ""
	if t.assignee != nil {
		previousAssigneeID = t.assignee.AssigneeID().Value()
	}

	assignment, err := entity.NewAssignment(t.id, assigneeID, assignedBy)
	if err != nil {
		return err
	}

	t.assignee = assignment
	t.updatedAt = time.Now()

	// Raise domain event
	assignedEvent := event.NewTaskAssignedEvent(
		t.id.Value(),
		assigneeID.Value(),
		previousAssigneeID,
	)
	t.domainEvents = append(t.domainEvents, assignedEvent)

	return nil
}

// ChangeStatus changes the task status with validation
func (t *Task) ChangeStatus(newStatus value.TaskStatus) error {
	if !newStatus.IsValid() {
		return fmt.Errorf("invalid status: %s", newStatus.Value())
	}

	if !t.status.CanTransitionTo(newStatus) {
		return fmt.Errorf("cannot transition from %s to %s", t.status.Value(), newStatus.Value())
	}

	oldStatus := t.status
	t.status = newStatus
	t.updatedAt = time.Now()

	// Raise domain event
	statusChangedEvent := event.NewTaskStatusChangedEvent(
		t.id.Value(),
		oldStatus.Value(),
		newStatus.Value(),
	)
	t.domainEvents = append(t.domainEvents, statusChangedEvent)

	// If completed, raise completion event
	if newStatus == value.TaskStatusCompleted {
		completedEvent := event.NewTaskCompletedEvent(
			t.id.Value(),
			t.assignee.AssigneeID().Value(),
			time.Now().Format(time.RFC3339),
		)
		t.domainEvents = append(t.domainEvents, completedEvent)
	}

	return nil
}

// SetDeadline sets the deadline for the task
func (t *Task) SetDeadline(deadline value.Deadline) error {
	t.deadline = &deadline
	t.updatedAt = time.Now()

	// Raise domain event
	deadlineEvent := event.NewTaskDeadlineSetEvent(
		t.id.Value(),
		deadline.Value().Format(time.RFC3339),
	)
	t.domainEvents = append(t.domainEvents, deadlineEvent)

	return nil
}

// AddComment adds a comment to the task
func (t *Task) AddComment(comment *entity.Comment) error {
	if comment == nil {
		return fmt.Errorf("comment cannot be nil")
	}

	t.comments = append(t.comments, comment)
	t.updatedAt = time.Now()

	return nil
}

// UpdateTitle updates the task title
func (t *Task) UpdateTitle(newTitle string) error {
	if newTitle == "" {
		return fmt.Errorf("title cannot be empty")
	}

	t.title = newTitle
	t.updatedAt = time.Now()

	return nil
}

// UpdateDescription updates the task description
func (t *Task) UpdateDescription(newDescription string) error {
	t.description = newDescription
	t.updatedAt = time.Now()

	return nil
}

// UpdatePriority updates the task priority
func (t *Task) UpdatePriority(newPriority value.Priority) error {
	if !newPriority.IsValid() {
		return fmt.Errorf("invalid priority")
	}

	t.priority = newPriority
	t.updatedAt = time.Now()

	return nil
}

// CheckDeadlineStatus checks and updates the deadline status
func (t *Task) CheckDeadlineStatus() {
	if t.deadline == nil {
		return
	}

	if t.deadline.IsOverdue() && t.status != value.TaskStatusCompleted && t.status != value.TaskStatusCancelled {
		daysOverdue := -t.deadline.DaysUntilDue()
		overdueEvent := event.NewTaskOverdueEvent(t.id.Value(), daysOverdue)
		t.domainEvents = append(t.domainEvents, overdueEvent)
	}
}

// UpdateStatus is a convenience method for status update (without validation)
func (t *Task) UpdateStatus(newStatus value.TaskStatus) {
	t.status = newStatus
	t.updatedAt = time.Now()
}