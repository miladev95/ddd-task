package value

import "fmt"

// TaskStatus represents the status of a task
type TaskStatus string

const (
	TaskStatusBacklog    TaskStatus = "BACKLOG"
	TaskStatusToDo       TaskStatus = "TO_DO"
	TaskStatusInProgress TaskStatus = "IN_PROGRESS"
	TaskStatusInReview   TaskStatus = "IN_REVIEW"
	TaskStatusCompleted  TaskStatus = "COMPLETED"
	TaskStatusCancelled  TaskStatus = "CANCELLED"
)

// NewTaskStatus creates a new TaskStatus from string
func NewTaskStatus(status string) (TaskStatus, error) {
	ts := TaskStatus(status)
	switch ts {
	case TaskStatusBacklog, TaskStatusToDo, TaskStatusInProgress, TaskStatusInReview, TaskStatusCompleted, TaskStatusCancelled:
		return ts, nil
	default:
		return "", fmt.Errorf("invalid task status: %s", status)
	}
}

// Value returns the string representation
func (t TaskStatus) Value() string {
	return string(t)
}

// IsValid checks if the status is valid
func (t TaskStatus) IsValid() bool {
	switch t {
	case TaskStatusBacklog, TaskStatusToDo, TaskStatusInProgress, TaskStatusInReview, TaskStatusCompleted, TaskStatusCancelled:
		return true
	default:
		return false
	}
}

// CanTransitionTo checks if transition from current status to target is valid
func (t TaskStatus) CanTransitionTo(target TaskStatus) bool {
	validTransitions := map[TaskStatus][]TaskStatus{
		TaskStatusBacklog: {TaskStatusToDo, TaskStatusCancelled},
		TaskStatusToDo: {TaskStatusInProgress, TaskStatusBacklog, TaskStatusCancelled},
		TaskStatusInProgress: {TaskStatusInReview, TaskStatusToDo, TaskStatusCancelled},
		TaskStatusInReview: {TaskStatusCompleted, TaskStatusInProgress, TaskStatusCancelled},
		TaskStatusCompleted: {}, // Cannot transition from completed
		TaskStatusCancelled: {}, // Cannot transition from cancelled
	}

	allowed, exists := validTransitions[t]
	if !exists {
		return false
	}

	for _, s := range allowed {
		if s == target {
			return true
		}
	}
	return false
}