package di

import (
	"github.com/example/task-management/application/command"
	"github.com/example/task-management/application/query"
	"github.com/example/task-management/domain"
	"github.com/example/task-management/domain/event"
	"github.com/example/task-management/domain/service"
	infraEvent "github.com/example/task-management/infrastructure/event"
	"github.com/example/task-management/infrastructure/repository"
)

// Container holds all application dependencies
type Container struct {
	// Repositories
	TaskRepository      domain.TaskRepository
	ProjectRepository   domain.ProjectRepository
	UserRepository      domain.UserRepository
	WorkflowRepository  domain.WorkflowRepository

	// Event
	EventPublisher      event.EventPublisher
	NotificationService service.NotificationService

	// Domain Services
	TaskAssignmentService    *service.TaskAssignmentService
	StatusTransitionService  *service.StatusTransitionService
	DeadlineEnforcementService *service.DeadlineEnforcementService

	// Command Handlers
	CreateTaskCommandHandler       *command.CreateTaskCommandHandler
	AssignTaskCommandHandler       *command.AssignTaskCommandHandler
	UpdateTaskStatusCommandHandler *command.UpdateTaskStatusCommandHandler

	// Query Handlers
	GetTaskQueryHandler               *query.GetTaskQueryHandler
	ListTasksByProjectQueryHandler    *query.ListTasksByProjectQueryHandler
}

// NewContainer creates and initializes a new dependency injection container
func NewContainer() *Container {
	c := &Container{}

	// Initialize repositories (using in-memory implementations for demo)
	c.TaskRepository = repository.NewInMemoryTaskRepository()
	c.ProjectRepository = repository.NewInMemoryProjectRepository()
	c.UserRepository = repository.NewInMemoryUserRepository()
	c.WorkflowRepository = repository.NewInMemoryWorkflowRepository()

	// Initialize event publisher
	c.EventPublisher = infraEvent.NewSimpleEventPublisher()

	// Initialize notification service
	c.NotificationService = infraEvent.NewSimpleNotificationService()

	// Initialize domain services
	c.TaskAssignmentService = service.NewTaskAssignmentService(
		c.UserRepository.(service.UserRepository),
		c.TaskRepository.(service.TaskRepository),
	)

	c.StatusTransitionService = service.NewStatusTransitionService(
		c.WorkflowRepository,
	)

	c.DeadlineEnforcementService = service.NewDeadlineEnforcementService(
		c.NotificationService,
	)

	// Initialize command handlers
	c.CreateTaskCommandHandler = command.NewCreateTaskCommandHandler(
		c.TaskRepository,
		c.ProjectRepository,
		c.UserRepository,
		c.WorkflowRepository,
		c.EventPublisher,
		c.TaskAssignmentService,
		c.DeadlineEnforcementService,
	)

	c.AssignTaskCommandHandler = command.NewAssignTaskCommandHandler(
		c.TaskRepository,
		c.EventPublisher,
		c.TaskAssignmentService,
	)

	c.UpdateTaskStatusCommandHandler = command.NewUpdateTaskStatusCommandHandler(
		c.TaskRepository,
		c.EventPublisher,
		c.StatusTransitionService,
	)

	// Initialize query handlers
	c.GetTaskQueryHandler = query.NewGetTaskQueryHandler(
		c.TaskRepository,
		c.UserRepository,
	)

	c.ListTasksByProjectQueryHandler = query.NewListTasksByProjectQueryHandler(
		c.TaskRepository,
	)

	return c
}