package domain

import (
	"github.com/miladev95/ddd-task/domain/aggregate"
	"github.com/miladev95/ddd-task/domain/value"
)

// TaskRepository defines the interface for task persistence
type TaskRepository interface {
	// Save persists a task to the repository
	Save(task *aggregate.Task) error

	// GetByID retrieves a task by ID
	GetByID(id value.TaskID) (*aggregate.Task, error)

	// GetByProjectID retrieves all tasks for a project
	GetByProjectID(projectID value.ProjectID) ([]*aggregate.Task, error)

	// GetByAssigneeID retrieves all tasks assigned to a user
	GetByAssigneeID(userID value.UserID) ([]*aggregate.Task, error)

	// GetByStatus retrieves all tasks with a specific status
	GetByStatus(status value.TaskStatus) ([]*aggregate.Task, error)

	// GetAll retrieves all tasks
	GetAll() ([]*aggregate.Task, error)

	// Delete removes a task from the repository
	Delete(id value.TaskID) error

	// Update updates an existing task
	Update(task *aggregate.Task) error

	// FindByProjectIDAndStatus retrieves tasks for a project with specific status
	FindByProjectIDAndStatus(projectID value.ProjectID, status value.TaskStatus) ([]*aggregate.Task, error)
}

// ProjectRepository defines the interface for project persistence
type ProjectRepository interface {
	// Save persists a project to the repository
	Save(project *aggregate.Project) error

	// GetByID retrieves a project by ID
	GetByID(id value.ProjectID) (*aggregate.Project, error)

	// GetByOwnerID retrieves all projects owned by a user
	GetByOwnerID(userID value.UserID) ([]*aggregate.Project, error)

	// GetAll retrieves all projects
	GetAll() ([]*aggregate.Project, error)

	// Delete removes a project from the repository
	Delete(id value.ProjectID) error

	// Update updates an existing project
	Update(project *aggregate.Project) error

	// GetActive retrieves all active projects
	GetActive() ([]*aggregate.Project, error)
}

// UserRepository defines the interface for user persistence
type UserRepository interface {
	// Save persists a user to the repository
	Save(user *aggregate.User) error

	// GetByID retrieves a user by ID
	GetByID(id value.UserID) (*aggregate.User, error)

	// GetByEmail retrieves a user by email
	GetByEmail(email string) (*aggregate.User, error)

	// GetAll retrieves all users
	GetAll() ([]*aggregate.User, error)

	// Delete removes a user from the repository
	Delete(id value.UserID) error

	// Update updates an existing user
	Update(user *aggregate.User) error

	// GetActive retrieves all active users
	GetActive() ([]*aggregate.User, error)
}

// WorkflowRepository defines the interface for workflow persistence
type WorkflowRepository interface {
	// Save persists a workflow to the repository
	Save(workflow *aggregate.Workflow) error

	// GetByID retrieves a workflow by ID
	GetByID(id value.WorkflowID) (*aggregate.Workflow, error)

	// GetByName retrieves a workflow by name
	GetByName(name string) (*aggregate.Workflow, error)

	// GetAll retrieves all workflows
	GetAll() ([]*aggregate.Workflow, error)

	// Delete removes a workflow from the repository
	Delete(id value.WorkflowID) error

	// Update updates an existing workflow
	Update(workflow *aggregate.Workflow) error

	// GetActive retrieves all active workflows
	GetActive() ([]*aggregate.Workflow, error)
}

// UnitOfWork defines the interface for transaction management
type UnitOfWork interface {
	// BeginTransaction starts a new transaction
	BeginTransaction() error

	// Commit commits the current transaction
	Commit() error

	// Rollback rolls back the current transaction
	Rollback() error

	// GetTaskRepository returns the task repository
	GetTaskRepository() TaskRepository

	// GetProjectRepository returns the project repository
	GetProjectRepository() ProjectRepository

	// GetUserRepository returns the user repository
	GetUserRepository() UserRepository

	// GetWorkflowRepository returns the workflow repository
	GetWorkflowRepository() WorkflowRepository
}