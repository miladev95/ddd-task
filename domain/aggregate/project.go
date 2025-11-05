package aggregate

import (
	"fmt"
	"time"

	"github.com/example/task-management/domain/event"
	"github.com/example/task-management/domain/value"
)

// Project is the aggregate root for the Project aggregate
type Project struct {
	id          value.ProjectID
	name        string
	description string
	ownerID     value.UserID
	taskIDs     []value.TaskID
	workflowID  value.WorkflowID
	createdAt   time.Time
	updatedAt   time.Time
	archived    bool
	domainEvents []event.DomainEvent
}

// NewProject creates a new Project
func NewProject(
	id value.ProjectID,
	name, description string,
	ownerID value.UserID,
	workflowID value.WorkflowID,
) (*Project, error) {
	if name == "" {
		return nil, fmt.Errorf("project name cannot be empty")
	}

	return &Project{
		id:           id,
		name:         name,
		description:  description,
		ownerID:      ownerID,
		workflowID:   workflowID,
		taskIDs:      make([]value.TaskID, 0),
		createdAt:    time.Now(),
		updatedAt:    time.Now(),
		archived:     false,
		domainEvents: make([]event.DomainEvent, 0),
	}, nil
}

// ID returns the project ID
func (p *Project) ID() value.ProjectID {
	return p.id
}

// Name returns the project name
func (p *Project) Name() string {
	return p.name
}

// Description returns the project description
func (p *Project) Description() string {
	return p.description
}

// OwnerID returns the owner ID
func (p *Project) OwnerID() value.UserID {
	return p.ownerID
}

// WorkflowID returns the workflow ID
func (p *Project) WorkflowID() value.WorkflowID {
	return p.workflowID
}

// TaskIDs returns all task IDs in the project
func (p *Project) TaskIDs() []value.TaskID {
	return append([]value.TaskID{}, p.taskIDs...)
}

// CreatedAt returns when the project was created
func (p *Project) CreatedAt() time.Time {
	return p.createdAt
}

// UpdatedAt returns when the project was last updated
func (p *Project) UpdatedAt() time.Time {
	return p.updatedAt
}

// IsArchived returns whether the project is archived
func (p *Project) IsArchived() bool {
	return p.archived
}

// DomainEvents returns all uncommitted domain events
func (p *Project) DomainEvents() []event.DomainEvent {
	return append([]event.DomainEvent{}, p.domainEvents...)
}

// ClearDomainEvents clears all domain events after they have been published
func (p *Project) ClearDomainEvents() {
	p.domainEvents = make([]event.DomainEvent, 0)
}

// AddTask adds a task to the project
func (p *Project) AddTask(taskID value.TaskID) error {
	if taskID.Equals(value.TaskID{}) {
		return fmt.Errorf("task id cannot be empty")
	}

	// Check if task already exists
	for _, existingID := range p.taskIDs {
		if existingID.Equals(taskID) {
			return fmt.Errorf("task already exists in project")
		}
	}

	p.taskIDs = append(p.taskIDs, taskID)
	p.updatedAt = time.Now()

	return nil
}

// RemoveTask removes a task from the project
func (p *Project) RemoveTask(taskID value.TaskID) error {
	for i, id := range p.taskIDs {
		if id.Equals(taskID) {
			p.taskIDs = append(p.taskIDs[:i], p.taskIDs[i+1:]...)
			p.updatedAt = time.Now()
			return nil
		}
	}

	return fmt.Errorf("task not found in project")
}

// UpdateName updates the project name
func (p *Project) UpdateName(newName string) error {
	if newName == "" {
		return fmt.Errorf("project name cannot be empty")
	}

	p.name = newName
	p.updatedAt = time.Now()

	return nil
}

// UpdateDescription updates the project description
func (p *Project) UpdateDescription(newDescription string) error {
	p.description = newDescription
	p.updatedAt = time.Now()

	return nil
}

// Archive archives the project
func (p *Project) Archive() error {
	p.archived = true
	p.updatedAt = time.Now()

	return nil
}

// Unarchive unarchives the project
func (p *Project) Unarchive() error {
	p.archived = false
	p.updatedAt = time.Now()

	return nil
}

// TaskCount returns the number of tasks in the project
func (p *Project) TaskCount() int {
	return len(p.taskIDs)
}