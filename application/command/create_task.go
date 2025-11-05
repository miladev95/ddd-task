package command

import (
	"fmt"
	"time"

	"github.com/miladev95/ddd-task/domain"
	"github.com/miladev95/ddd-task/domain/aggregate"
	"github.com/miladev95/ddd-task/domain/event"
	"github.com/miladev95/ddd-task/domain/service"
	"github.com/miladev95/ddd-task/domain/value"
)

// CreateTaskCommand represents a command to create a task
type CreateTaskCommand struct {
	ProjectID   string
	Title       string
	Description string
	Priority    string
	AssigneeID  string
	Deadline    string
	CreatedBy   string
}

// CreateTaskCommandHandler handles CreateTaskCommand
type CreateTaskCommandHandler struct {
	taskRepository       domain.TaskRepository
	projectRepository    domain.ProjectRepository
	userRepository       domain.UserRepository
	workflowRepository   domain.WorkflowRepository
	eventPublisher       event.EventPublisher
	assignmentService    *service.TaskAssignmentService
	deadlineService      *service.DeadlineEnforcementService
}

// NewCreateTaskCommandHandler creates a new CreateTaskCommandHandler
func NewCreateTaskCommandHandler(
	taskRepository domain.TaskRepository,
	projectRepository domain.ProjectRepository,
	userRepository domain.UserRepository,
	workflowRepository domain.WorkflowRepository,
	eventPublisher event.EventPublisher,
	assignmentService *service.TaskAssignmentService,
	deadlineService *service.DeadlineEnforcementService,
) *CreateTaskCommandHandler {
	return &CreateTaskCommandHandler{
		taskRepository:       taskRepository,
		projectRepository:    projectRepository,
		userRepository:       userRepository,
		workflowRepository:   workflowRepository,
		eventPublisher:       eventPublisher,
		assignmentService:    assignmentService,
		deadlineService:      deadlineService,
	}
}

// CreateTaskResult represents the result of creating a task
type CreateTaskResult struct {
	TaskID string
	Error  error
}

// Handle handles the CreateTaskCommand
func (h *CreateTaskCommandHandler) Handle(cmd CreateTaskCommand) (*CreateTaskResult, error) {
	// Validate project exists
	projectID, err := value.NewProjectID(cmd.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("invalid project id: %w", err)
	}

	project, err := h.projectRepository.GetByID(projectID)
	if err != nil {
		return nil, fmt.Errorf("project not found: %w", err)
	}

	// Validate priority
	priority, err := value.NewPriority(cmd.Priority)
	if err != nil {
		return nil, fmt.Errorf("invalid priority: %w", err)
	}

	// Validate created by user
	createdByID, err := value.NewUserID(cmd.CreatedBy)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}

	_, err = h.userRepository.GetByID(createdByID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Generate new task ID
	taskID := value.GenerateTaskID()

	// Create task aggregate
	task, err := aggregate.NewTask(
		taskID,
		projectID,
		cmd.Title,
		cmd.Description,
		priority,
		createdByID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	// Assign task if assignee provided
	if cmd.AssigneeID != "" {
		assigneeID, err := value.NewUserID(cmd.AssigneeID)
		if err != nil {
			return nil, fmt.Errorf("invalid assignee id: %w", err)
		}

		err = h.assignmentService.AssignTask(task, assigneeID, createdByID)
		if err != nil {
			return nil, fmt.Errorf("failed to assign task: %w", err)
		}
	}

	// Set deadline if provided
	if cmd.Deadline != "" {
		dueDate, err := time.Parse(time.RFC3339, cmd.Deadline)
		if err != nil {
			return nil, fmt.Errorf("invalid deadline format: %w", err)
		}

		deadline, err := value.NewDeadline(dueDate)
		if err != nil {
			return nil, fmt.Errorf("invalid deadline: %w", err)
		}

		err = h.deadlineService.SetDeadline(task, deadline)
		if err != nil {
			return nil, fmt.Errorf("failed to set deadline: %w", err)
		}
	}

	// Add task to project
	err = project.AddTask(taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to add task to project: %w", err)
	}

	// Save task
	err = h.taskRepository.Save(task)
	if err != nil {
		return nil, fmt.Errorf("failed to save task: %w", err)
	}

	// Update project
	err = h.projectRepository.Update(project)
	if err != nil {
		return nil, fmt.Errorf("failed to update project: %w", err)
	}

	// Publish domain events
	for _, domainEvent := range task.DomainEvents() {
		err = h.eventPublisher.Publish(domainEvent)
		if err != nil {
			return nil, fmt.Errorf("failed to publish event: %w", err)
		}
	}

	task.ClearDomainEvents()

	return &CreateTaskResult{
		TaskID: taskID.Value(),
	}, nil
}