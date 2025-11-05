package value

import (
	"fmt"

	"github.com/google/uuid"
)

// TaskID represents a unique identifier for a Task
type TaskID struct {
	value string
}

// NewTaskID creates a new TaskID
func NewTaskID(id string) (TaskID, error) {
	if id == "" {
		return TaskID{}, fmt.Errorf("task id cannot be empty")
	}
	return TaskID{value: id}, nil
}

// GenerateTaskID generates a new random TaskID
func GenerateTaskID() TaskID {
	return TaskID{value: uuid.New().String()}
}

// Value returns the string representation of TaskID
func (t TaskID) Value() string {
	return t.value
}

// Equals compares two TaskIDs for equality
func (t TaskID) Equals(other TaskID) bool {
	return t.value == other.value
}

// ProjectID represents a unique identifier for a Project
type ProjectID struct {
	value string
}

// NewProjectID creates a new ProjectID
func NewProjectID(id string) (ProjectID, error) {
	if id == "" {
		return ProjectID{}, fmt.Errorf("project id cannot be empty")
	}
	return ProjectID{value: id}, nil
}

// GenerateProjectID generates a new random ProjectID
func GenerateProjectID() ProjectID {
	return ProjectID{value: uuid.New().String()}
}

// Value returns the string representation of ProjectID
func (p ProjectID) Value() string {
	return p.value
}

// Equals compares two ProjectIDs for equality
func (p ProjectID) Equals(other ProjectID) bool {
	return p.value == other.value
}

// UserID represents a unique identifier for a User
type UserID struct {
	value string
}

// NewUserID creates a new UserID
func NewUserID(id string) (UserID, error) {
	if id == "" {
		return UserID{}, fmt.Errorf("user id cannot be empty")
	}
	return UserID{value: id}, nil
}

// GenerateUserID generates a new random UserID
func GenerateUserID() UserID {
	return UserID{value: uuid.New().String()}
}

// Value returns the string representation of UserID
func (u UserID) Value() string {
	return u.value
}

// Equals compares two UserIDs for equality
func (u UserID) Equals(other UserID) bool {
	return u.value == other.value
}

// WorkflowID represents a unique identifier for a Workflow
type WorkflowID struct {
	value string
}

// NewWorkflowID creates a new WorkflowID
func NewWorkflowID(id string) (WorkflowID, error) {
	if id == "" {
		return WorkflowID{}, fmt.Errorf("workflow id cannot be empty")
	}
	return WorkflowID{value: id}, nil
}

// GenerateWorkflowID generates a new random WorkflowID
func GenerateWorkflowID() WorkflowID {
	return WorkflowID{value: uuid.New().String()}
}

// Value returns the string representation of WorkflowID
func (w WorkflowID) Value() string {
	return w.value
}

// Equals compares two WorkflowIDs for equality
func (w WorkflowID) Equals(other WorkflowID) bool {
	return w.value == other.value
}