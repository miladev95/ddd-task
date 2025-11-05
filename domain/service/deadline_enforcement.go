package service

import (
	"fmt"
	"time"

	"github.com/example/task-management/domain/aggregate"
	"github.com/example/task-management/domain/value"
)

// DeadlineEnforcementService handles deadline validation and enforcement
type DeadlineEnforcementService struct {
	notificationService NotificationService
}

// NewDeadlineEnforcementService creates a new DeadlineEnforcementService
func NewDeadlineEnforcementService(
	notificationService NotificationService,
) *DeadlineEnforcementService {
	return &DeadlineEnforcementService{
		notificationService: notificationService,
	}
}

// ValidateDeadline validates that a deadline is reasonable
func (s *DeadlineEnforcementService) ValidateDeadline(deadline value.Deadline) error {
	if deadline.IsOverdue() {
		return fmt.Errorf("deadline cannot be in the past")
	}

	// Check if deadline is more than 5 years in the future (arbitrary validation)
	futureThreshold := time.Now().AddDate(5, 0, 0)
	if deadline.Value().After(futureThreshold) {
		return fmt.Errorf("deadline too far in the future")
	}

	return nil
}

// SetDeadline sets a deadline on a task with validation
func (s *DeadlineEnforcementService) SetDeadline(
	task *aggregate.Task,
	deadline value.Deadline,
) error {
	// Validate deadline
	if err := s.ValidateDeadline(deadline); err != nil {
		return fmt.Errorf("invalid deadline: %w", err)
	}

	// Task cannot have a deadline if already completed or cancelled
	if task.Status() == value.TaskStatusCompleted || task.Status() == value.TaskStatusCancelled {
		return fmt.Errorf("cannot set deadline for completed or cancelled tasks")
	}

	// Set the deadline
	if err := task.SetDeadline(deadline); err != nil {
		return fmt.Errorf("failed to set deadline: %w", err)
	}

	return nil
}

// CheckOverdueStatus checks and notifies about overdue tasks
func (s *DeadlineEnforcementService) CheckOverdueStatus(task *aggregate.Task) error {
	if task.Deadline() == nil {
		return nil
	}

	if task.Deadline().IsOverdue() && task.Status() != value.TaskStatusCompleted && task.Status() != value.TaskStatusCancelled {
		task.CheckDeadlineStatus()

		// Notify assignee if assigned
		if task.Assignee() != nil {
			err := s.notificationService.NotifyTaskOverdue(task)
			if err != nil {
				// Log but don't fail - notification error shouldn't block deadline check
				return nil
			}
		}
	}

	return nil
}

// GetTasksDueWithin returns tasks that are due within a duration
func (s *DeadlineEnforcementService) GetTasksDueWithin(
	tasks []*aggregate.Task,
	duration time.Duration,
) []*aggregate.Task {
	dueTasks := make([]*aggregate.Task, 0)

	for _, task := range tasks {
		if task.Deadline() != nil && task.Deadline().IsDueSoon(duration) {
			if task.Status() != value.TaskStatusCompleted && task.Status() != value.TaskStatusCancelled {
				dueTasks = append(dueTasks, task)
			}
		}
	}

	return dueTasks
}

// GetOverdueTasks returns all overdue tasks
func (s *DeadlineEnforcementService) GetOverdueTasks(tasks []*aggregate.Task) []*aggregate.Task {
	overdueTasks := make([]*aggregate.Task, 0)

	for _, task := range tasks {
		if task.Deadline() != nil && task.Deadline().IsOverdue() {
			if task.Status() != value.TaskStatusCompleted && task.Status() != value.TaskStatusCancelled {
				overdueTasks = append(overdueTasks, task)
			}
		}
	}

	return overdueTasks
}

// ExtendDeadline extends a task's deadline
func (s *DeadlineEnforcementService) ExtendDeadline(
	task *aggregate.Task,
	newDeadline value.Deadline,
) error {
	// Validate new deadline is later than current
	if task.Deadline() != nil {
		if newDeadline.Value().Before(task.Deadline().Value()) {
			return fmt.Errorf("new deadline must be after current deadline")
		}
	}

	return s.SetDeadline(task, newDeadline)
}

// NotificationService interface for sending notifications
type NotificationService interface {
	NotifyTaskOverdue(task *aggregate.Task) error
	NotifyTaskAssigned(task *aggregate.Task, assigneeID string) error
	NotifyTaskStatusChanged(task *aggregate.Task, oldStatus, newStatus string) error
}