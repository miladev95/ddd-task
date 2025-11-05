package aggregate

import (
	"fmt"
	"time"

	"github.com/miladev95/ddd-task/domain/event"
	"github.com/miladev95/ddd-task/domain/value"
)

// WorkflowStatus represents a status in a workflow
type WorkflowStatus struct {
	name        string
	description string
	order       int
	isFinal     bool
}

// Workflow is the aggregate root for the Workflow aggregate
type Workflow struct {
	id           value.WorkflowID
	name         string
	description  string
	statuses     []WorkflowStatus
	createdAt    time.Time
	updatedAt    time.Time
	active       bool
	domainEvents []event.DomainEvent
}

// NewWorkflow creates a new Workflow
func NewWorkflow(
	id value.WorkflowID,
	name, description string,
	statuses []WorkflowStatus,
) (*Workflow, error) {
	if name == "" {
		return nil, fmt.Errorf("workflow name cannot be empty")
	}

	if len(statuses) == 0 {
		return nil, fmt.Errorf("workflow must have at least one status")
	}

	// Validate statuses
	seenNames := make(map[string]bool)
	for _, status := range statuses {
		if status.name == "" {
			return nil, fmt.Errorf("status name cannot be empty")
		}
		if seenNames[status.name] {
			return nil, fmt.Errorf("duplicate status name: %s", status.name)
		}
		seenNames[status.name] = true
	}

	return &Workflow{
		id:           id,
		name:         name,
		description:  description,
		statuses:     statuses,
		createdAt:    time.Now(),
		updatedAt:    time.Now(),
		active:       true,
		domainEvents: make([]event.DomainEvent, 0),
	}, nil
}

// ID returns the workflow ID
func (w *Workflow) ID() value.WorkflowID {
	return w.id
}

// Name returns the workflow name
func (w *Workflow) Name() string {
	return w.name
}

// Description returns the workflow description
func (w *Workflow) Description() string {
	return w.description
}

// Statuses returns the workflow statuses
func (w *Workflow) Statuses() []WorkflowStatus {
	return append([]WorkflowStatus{}, w.statuses...)
}

// CreatedAt returns when the workflow was created
func (w *Workflow) CreatedAt() time.Time {
	return w.createdAt
}

// UpdatedAt returns when the workflow was last updated
func (w *Workflow) UpdatedAt() time.Time {
	return w.updatedAt
}

// IsActive returns whether the workflow is active
func (w *Workflow) IsActive() bool {
	return w.active
}

// DomainEvents returns all uncommitted domain events
func (w *Workflow) DomainEvents() []event.DomainEvent {
	return append([]event.DomainEvent{}, w.domainEvents...)
}

// ClearDomainEvents clears all domain events after they have been published
func (w *Workflow) ClearDomainEvents() {
	w.domainEvents = make([]event.DomainEvent, 0)
}

// GetStatusByName gets a status by name
func (w *Workflow) GetStatusByName(name string) (*WorkflowStatus, error) {
	for i, status := range w.statuses {
		if status.name == name {
			return &w.statuses[i], nil
		}
	}
	return nil, fmt.Errorf("status not found: %s", name)
}

// IsValidStatus checks if a status exists in the workflow
func (w *Workflow) IsValidStatus(statusName string) bool {
	for _, status := range w.statuses {
		if status.name == statusName {
			return true
		}
	}
	return false
}

// Activate activates the workflow
func (w *Workflow) Activate() error {
	if w.active {
		return fmt.Errorf("workflow is already active")
	}

	w.active = true
	w.updatedAt = time.Now()

	return nil
}

// Deactivate deactivates the workflow
func (w *Workflow) Deactivate() error {
	if !w.active {
		return fmt.Errorf("workflow is already inactive")
	}

	w.active = false
	w.updatedAt = time.Now()

	return nil
}

// UpdateName updates the workflow name
func (w *Workflow) UpdateName(newName string) error {
	if newName == "" {
		return fmt.Errorf("workflow name cannot be empty")
	}

	w.name = newName
	w.updatedAt = time.Now()

	return nil
}

// NewWorkflowStatus creates a new workflow status
func NewWorkflowStatus(name, description string, order int, isFinal bool) WorkflowStatus {
	return WorkflowStatus{
		name:        name,
		description: description,
		order:       order,
		isFinal:     isFinal,
	}
}

// GetName returns the status name
func (ws *WorkflowStatus) GetName() string {
	return ws.name
}

// GetDescription returns the status description
func (ws *WorkflowStatus) GetDescription() string {
	return ws.description
}

// GetOrder returns the status order
func (ws *WorkflowStatus) GetOrder() int {
	return ws.order
}

// IsFinal returns whether this is a final status
func (ws *WorkflowStatus) IsFinal() bool {
	return ws.isFinal
}