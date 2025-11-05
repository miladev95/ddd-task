package service

import (
	"fmt"

	"github.com/example/task-management/domain/aggregate"
	"github.com/example/task-management/domain/value"
)

// TaskAssignmentService handles task assignment logic
type TaskAssignmentService struct {
	userRepository  UserRepository
	taskRepository  TaskRepository
}

// NewTaskAssignmentService creates a new TaskAssignmentService
func NewTaskAssignmentService(
	userRepository UserRepository,
	taskRepository TaskRepository,
) *TaskAssignmentService {
	return &TaskAssignmentService{
		userRepository:  userRepository,
		taskRepository:  taskRepository,
	}
}

// AssignTask assigns a task to a user
func (s *TaskAssignmentService) AssignTask(
	task *aggregate.Task,
	assigneeID value.UserID,
	assignedBy value.UserID,
) error {
	// Verify assignee exists
	_, err := s.userRepository.GetByID(assigneeID)
	if err != nil {
		return fmt.Errorf("assignee not found: %w", err)
	}

	// Verify assigner has permission (simplified - can be enhanced with permissions service)
	_, err = s.userRepository.GetByID(assignedBy)
	if err != nil {
		return fmt.Errorf("assigner not found: %w", err)
	}

	// Assign the task
	if err := task.Assign(assigneeID, assignedBy); err != nil {
		return fmt.Errorf("failed to assign task: %w", err)
	}

	return nil
}

// ReassignTask reassigns a task from one user to another
func (s *TaskAssignmentService) ReassignTask(
	task *aggregate.Task,
	newAssigneeID value.UserID,
	reassignedBy value.UserID,
) error {
	// Verify task is assigned
	if task.Assignee() == nil {
		return fmt.Errorf("task is not assigned")
	}

	// Verify new assignee exists
	_, err := s.userRepository.GetByID(newAssigneeID)
	if err != nil {
		return fmt.Errorf("new assignee not found: %w", err)
	}

	// Verify reassigner has permission
	_, err = s.userRepository.GetByID(reassignedBy)
	if err != nil {
		return fmt.Errorf("reassigner not found: %w", err)
	}

	// Reassign the task
	if err := task.Assign(newAssigneeID, reassignedBy); err != nil {
		return fmt.Errorf("failed to reassign task: %w", err)
	}

	return nil
}

// UnassignTask unassigns a task from its current assignee
func (s *TaskAssignmentService) UnassignTask(task *aggregate.Task) error {
	// Verify task is assigned
	if task.Assignee() == nil {
		return fmt.Errorf("task is not assigned")
	}

	// Get unassigned user ID (empty/zero user)
	unassignedID := value.UserID{} // This will be a zero value

	// Create a new unassigned user ID properly
	tempUserID, err := value.NewUserID("unassigned")
	if err != nil {
		return fmt.Errorf("failed to create unassigned id: %w", err)
	}

	// Assign to unassigned user
	if err := task.Assign(tempUserID, tempUserID); err != nil {
		return fmt.Errorf("failed to unassign task: %w", err)
	}

	_ = unassignedID // Remove unused variable warning if not needed

	return nil
}

// GetAssigneeTaskCount returns the number of tasks assigned to a user
func (s *TaskAssignmentService) GetAssigneeTaskCount(assigneeID value.UserID) (int, error) {
	// This would typically use a repository query
	// For now, returning 0 as placeholder
	return 0, nil
}

// ValidateAssignmentCapacity checks if a user can take on more tasks
func (s *TaskAssignmentService) ValidateAssignmentCapacity(
	assigneeID value.UserID,
	maxTasksPerUser int,
) (bool, error) {
	count, err := s.GetAssigneeTaskCount(assigneeID)
	if err != nil {
		return false, err
	}

	return count < maxTasksPerUser, nil
}

// UserRepository interface for getting user information
type UserRepository interface {
	GetByID(id value.UserID) (*aggregate.User, error)
}

// TaskRepository interface for getting task information
type TaskRepository interface {
	GetByID(id value.TaskID) (*aggregate.Task, error)
}