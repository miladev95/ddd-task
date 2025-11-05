package entity

import (
	"fmt"
	"time"

	"github.com/example/task-management/domain/value"
)

// Assignment represents the assignment of a task to a user
type Assignment struct {
	taskID    value.TaskID
	assigneeID value.UserID
	assignedAt time.Time
	assignedBy value.UserID
}

// NewAssignment creates a new Assignment
func NewAssignment(taskID value.TaskID, assigneeID value.UserID, assignedBy value.UserID) (*Assignment, error) {
	if assigneeID.Equals(value.UserID{}) {
		return nil, fmt.Errorf("assignee cannot be empty")
	}

	return &Assignment{
		taskID:     taskID,
		assigneeID: assigneeID,
		assignedAt: time.Now(),
		assignedBy: assignedBy,
	}, nil
}

// TaskID returns the task ID
func (a *Assignment) TaskID() value.TaskID {
	return a.taskID
}

// AssigneeID returns the assignee ID
func (a *Assignment) AssigneeID() value.UserID {
	return a.assigneeID
}

// AssignedAt returns when the task was assigned
func (a *Assignment) AssignedAt() time.Time {
	return a.assignedAt
}

// AssignedBy returns who assigned the task
func (a *Assignment) AssignedBy() value.UserID {
	return a.assignedBy
}

// IsAssignedTo checks if the assignment is for a specific user
func (a *Assignment) IsAssignedTo(userID value.UserID) bool {
	return a.assigneeID.Equals(userID)
}