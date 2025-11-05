package service

import (
	"fmt"

	"github.com/example/task-management/domain/aggregate"
	"github.com/example/task-management/domain/value"
)

// StatusTransitionService handles task status transitions with business rule validation
type StatusTransitionService struct {
	workflowRepository WorkflowRepository
}

// NewStatusTransitionService creates a new StatusTransitionService
func NewStatusTransitionService(
	workflowRepository WorkflowRepository,
) *StatusTransitionService {
	return &StatusTransitionService{
		workflowRepository: workflowRepository,
	}
}

// CanTransition checks if a task can transition to a new status
func (s *StatusTransitionService) CanTransition(
	task *aggregate.Task,
	newStatus value.TaskStatus,
) bool {
	return task.Status().CanTransitionTo(newStatus)
}

// TransitionTask transitions a task to a new status with validation
func (s *StatusTransitionService) TransitionTask(
	task *aggregate.Task,
	newStatus value.TaskStatus,
) error {
	// Check if transition is allowed
	if !s.CanTransition(task, newStatus) {
		return fmt.Errorf(
			"invalid status transition from %s to %s",
			task.Status().Value(),
			newStatus.Value(),
		)
	}

	// Additional validation: task must be assigned before moving to in-progress
	if newStatus == value.TaskStatusInProgress && task.Assignee() == nil {
		return fmt.Errorf("task must be assigned before moving to in-progress")
	}

	// Additional validation: task must have a deadline before completing
	if newStatus == value.TaskStatusCompleted && task.Deadline() == nil {
		return fmt.Errorf("task must have a deadline before completion")
	}

	// Perform the transition
	if err := task.ChangeStatus(newStatus); err != nil {
		return fmt.Errorf("failed to change task status: %w", err)
	}

	return nil
}

// GetValidNextStatuses returns the valid next statuses for a task
func (s *StatusTransitionService) GetValidNextStatuses(
	currentStatus value.TaskStatus,
) []value.TaskStatus {
	validStatuses := make([]value.TaskStatus, 0)

	allStatuses := []value.TaskStatus{
		value.TaskStatusBacklog,
		value.TaskStatusToDo,
		value.TaskStatusInProgress,
		value.TaskStatusInReview,
		value.TaskStatusCompleted,
		value.TaskStatusCancelled,
	}

	for _, status := range allStatuses {
		if currentStatus.CanTransitionTo(status) {
			validStatuses = append(validStatuses, status)
		}
	}

	return validStatuses
}

// StartTask starts a task (transitions to in-progress)
func (s *StatusTransitionService) StartTask(task *aggregate.Task) error {
	return s.TransitionTask(task, value.TaskStatusInProgress)
}

// CompleteTask completes a task (transitions to completed)
func (s *StatusTransitionService) CompleteTask(task *aggregate.Task) error {
	return s.TransitionTask(task, value.TaskStatusCompleted)
}

// CancelTask cancels a task (transitions to cancelled)
func (s *StatusTransitionService) CancelTask(task *aggregate.Task) error {
	return s.TransitionTask(task, value.TaskStatusCancelled)
}

// MoveToReview moves a task to review status
func (s *StatusTransitionService) MoveToReview(task *aggregate.Task) error {
	return s.TransitionTask(task, value.TaskStatusInReview)
}

// WorkflowRepository interface for workflow operations
type WorkflowRepository interface {
	GetByID(id value.WorkflowID) (*aggregate.Workflow, error)
	GetByName(name string) (*aggregate.Workflow, error)
}