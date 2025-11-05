package query

import (
	"fmt"

	"github.com/miladev95/ddd-task/application/dto"
	"github.com/miladev95/ddd-task/domain"
	"github.com/miladev95/ddd-task/domain/value"
)

// ListTasksByProjectQuery represents a query to list tasks by project
type ListTasksByProjectQuery struct {
	ProjectID string
	Status    string // optional filter
}

// ListTasksByProjectQueryHandler handles ListTasksByProjectQuery
type ListTasksByProjectQueryHandler struct {
	taskRepository domain.TaskRepository
}

// NewListTasksByProjectQueryHandler creates a new ListTasksByProjectQueryHandler
func NewListTasksByProjectQueryHandler(
	taskRepository domain.TaskRepository,
) *ListTasksByProjectQueryHandler {
	return &ListTasksByProjectQueryHandler{
		taskRepository: taskRepository,
	}
}

// Handle handles the ListTasksByProjectQuery
func (h *ListTasksByProjectQueryHandler) Handle(query ListTasksByProjectQuery) ([]*dto.TaskDTO, error) {
	// Parse project ID
	projectID, err := value.NewProjectID(query.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("invalid project id: %w", err)
	}

	// Get tasks for project
	var tasks interface{}
	var err2 error

	if query.Status != "" {
		// Get tasks for project with specific status
		status, err := value.NewTaskStatus(query.Status)
		if err != nil {
			return nil, fmt.Errorf("invalid status: %w", err)
		}
		tasks, err2 = h.taskRepository.FindByProjectIDAndStatus(projectID, status)
	} else {
		// Get all tasks for project
		tasks, err2 = h.taskRepository.GetByProjectID(projectID)
	}

	if err2 != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err2)
	}

	// Convert to DTOs
	taskDTOs := make([]*dto.TaskDTO, 0)
	_ = tasks // placeholder - actual implementation would convert

	return taskDTOs, nil
}