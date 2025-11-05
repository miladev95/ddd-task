package query

import (
	"fmt"

	"github.com/example/task-management/application/dto"
	"github.com/example/task-management/domain"
	"github.com/example/task-management/domain/value"
)

// GetTaskQuery represents a query to get a task by ID
type GetTaskQuery struct {
	TaskID string
}

// GetTaskQueryHandler handles GetTaskQuery
type GetTaskQueryHandler struct {
	taskRepository domain.TaskRepository
	userRepository domain.UserRepository
}

// NewGetTaskQueryHandler creates a new GetTaskQueryHandler
func NewGetTaskQueryHandler(
	taskRepository domain.TaskRepository,
	userRepository domain.UserRepository,
) *GetTaskQueryHandler {
	return &GetTaskQueryHandler{
		taskRepository: taskRepository,
		userRepository: userRepository,
	}
}

// Handle handles the GetTaskQuery
func (h *GetTaskQueryHandler) Handle(query GetTaskQuery) (*dto.TaskDTO, error) {
	// Parse task ID
	taskID, err := value.NewTaskID(query.TaskID)
	if err != nil {
		return nil, fmt.Errorf("invalid task id: %w", err)
	}

	// Get task
	task, err := h.taskRepository.GetByID(taskID)
	if err != nil {
		return nil, fmt.Errorf("task not found: %w", err)
	}

	// Convert to DTO
	return convertTaskToDTO(task), nil
}

// Helper function to convert task aggregate to DTO
func convertTaskToDTO(task interface{}) *dto.TaskDTO {
	// This is a placeholder - actual implementation would convert task aggregate to DTO
	// For now returning a basic structure
	return &dto.TaskDTO{
		Status: "TO_DO",
	}
}