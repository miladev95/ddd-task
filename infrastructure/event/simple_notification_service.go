package event

import (
	"fmt"

	"github.com/example/task-management/domain/aggregate"
	"github.com/example/task-management/domain/service"
)

// SimpleNotificationService is a basic implementation of NotificationService
type SimpleNotificationService struct {
	// In a real application, this would connect to a messaging service
	// or email service
}

// NewSimpleNotificationService creates a new SimpleNotificationService
func NewSimpleNotificationService() *SimpleNotificationService {
	return &SimpleNotificationService{}
}

// NotifyTaskOverdue sends a notification for an overdue task
func (s *SimpleNotificationService) NotifyTaskOverdue(task *aggregate.Task) error {
	if task.Assignee() == nil {
		return fmt.Errorf("task has no assignee")
	}

	// In real implementation, send notification via email/SMS/push notification
	fmt.Printf("NOTIFICATION: Task '%s' is overdue for user %s\n",
		task.Title(),
		task.Assignee().AssigneeID().Value(),
	)

	return nil
}

// NotifyTaskAssigned sends a notification for a task assignment
func (s *SimpleNotificationService) NotifyTaskAssigned(task *aggregate.Task, assigneeID string) error {
	// In real implementation, send notification
	fmt.Printf("NOTIFICATION: Task '%s' has been assigned to user %s\n",
		task.Title(),
		assigneeID,
	)

	return nil
}

// NotifyTaskStatusChanged sends a notification for a status change
func (s *SimpleNotificationService) NotifyTaskStatusChanged(
	task *aggregate.Task,
	oldStatus, newStatus string,
) error {
	if task.Assignee() == nil {
		return nil // No one assigned, no need to notify
	}

	// In real implementation, send notification
	fmt.Printf("NOTIFICATION: Task '%s' status changed from %s to %s\n",
		task.Title(),
		oldStatus,
		newStatus,
	)

	return nil
}

// Ensure SimpleNotificationService implements service.NotificationService
var _ service.NotificationService = (*SimpleNotificationService)(nil)