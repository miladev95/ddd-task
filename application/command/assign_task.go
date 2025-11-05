package command

import (
	"fmt"

	"github.com/miladev95/ddd-task/domain"
	"github.com/miladev95/ddd-task/domain/event"
	"github.com/miladev95/ddd-task/domain/service"
	"github.com/miladev95/ddd-task/domain/value"
)

// AssignTaskCommand represents a command to assign a task to a user
type AssignTaskCommand struct {
	TaskID     string
	AssigneeID string
	AssignedBy string
}

// AssignTaskCommandHandler handles AssignTaskCommand
type AssignTaskCommandHandler struct {
	taskRepository    domain.TaskRepository
	eventPublisher    event.EventPublisher
	assignmentService *service.TaskAssignmentService
}

// NewAssignTaskCommandHandler creates a new AssignTaskCommandHandler
func NewAssignTaskCommandHandler(
	taskRepository domain.TaskRepository,
	eventPublisher event.EventPublisher,
	assignmentService *service.TaskAssignmentService,
) *AssignTaskCommandHandler {
	return &AssignTaskCommandHandler{
		taskRepository:    taskRepository,
		eventPublisher:    eventPublisher,
		assignmentService: assignmentService,
	}
}

// AssignTaskResult represents the result of assigning a task
type AssignTaskResult struct {
	Error error
}

// Handle handles the AssignTaskCommand
func (h *AssignTaskCommandHandler) Handle(cmd AssignTaskCommand) (*AssignTaskResult, error) {
	// Parse IDs
	taskID, err := value.NewTaskID(cmd.TaskID)
	if err != nil {
		return nil, fmt.Errorf("invalid task id: %w", err)
	}

	assigneeID, err := value.NewUserID(cmd.AssigneeID)
	if err != nil {
		return nil, fmt.Errorf("invalid assignee id: %w", err)
	}

	assignedByID, err := value.NewUserID(cmd.AssignedBy)
	if err != nil {
		return nil, fmt.Errorf("invalid assigner id: %w", err)
	}

	// Get task
	task, err := h.taskRepository.GetByID(taskID)
	if err != nil {
		return nil, fmt.Errorf("task not found: %w", err)
	}

	// Assign task
	err = h.assignmentService.AssignTask(task, assigneeID, assignedByID)
	if err != nil {
		return nil, fmt.Errorf("failed to assign task: %w", err)
	}

	// Save task
	err = h.taskRepository.Update(task)
	if err != nil {
		return nil, fmt.Errorf("failed to save task: %w", err)
	}

	// Publish domain events
	for _, domainEvent := range task.DomainEvents() {
		err = h.eventPublisher.Publish(domainEvent)
		if err != nil {
			return nil, fmt.Errorf("failed to publish event: %w", err)
		}
	}

	task.ClearDomainEvents()

	return &AssignTaskResult{}, nil
}