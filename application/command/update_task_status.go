package command

import (
	"fmt"

	"github.com/example/task-management/domain"
	"github.com/example/task-management/domain/event"
	"github.com/example/task-management/domain/service"
	"github.com/example/task-management/domain/value"
)

// UpdateTaskStatusCommand represents a command to update task status
type UpdateTaskStatusCommand struct {
	TaskID    string
	NewStatus string
}

// UpdateTaskStatusCommandHandler handles UpdateTaskStatusCommand
type UpdateTaskStatusCommandHandler struct {
	taskRepository        domain.TaskRepository
	eventPublisher        event.EventPublisher
	statusTransitionService *service.StatusTransitionService
}

// NewUpdateTaskStatusCommandHandler creates a new UpdateTaskStatusCommandHandler
func NewUpdateTaskStatusCommandHandler(
	taskRepository domain.TaskRepository,
	eventPublisher event.EventPublisher,
	statusTransitionService *service.StatusTransitionService,
) *UpdateTaskStatusCommandHandler {
	return &UpdateTaskStatusCommandHandler{
		taskRepository:        taskRepository,
		eventPublisher:        eventPublisher,
		statusTransitionService: statusTransitionService,
	}
}

// UpdateTaskStatusResult represents the result of updating task status
type UpdateTaskStatusResult struct {
	Error error
}

// Handle handles the UpdateTaskStatusCommand
func (h *UpdateTaskStatusCommandHandler) Handle(cmd UpdateTaskStatusCommand) (*UpdateTaskStatusResult, error) {
	// Parse IDs
	taskID, err := value.NewTaskID(cmd.TaskID)
	if err != nil {
		return nil, fmt.Errorf("invalid task id: %w", err)
	}

	// Parse new status
	newStatus, err := value.NewTaskStatus(cmd.NewStatus)
	if err != nil {
		return nil, fmt.Errorf("invalid status: %w", err)
	}

	// Get task
	task, err := h.taskRepository.GetByID(taskID)
	if err != nil {
		return nil, fmt.Errorf("task not found: %w", err)
	}

	// Transition task status
	err = h.statusTransitionService.TransitionTask(task, newStatus)
	if err != nil {
		return nil, fmt.Errorf("failed to update status: %w", err)
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

	return &UpdateTaskStatusResult{}, nil
}