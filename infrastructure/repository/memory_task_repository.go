package repository

import (
	"fmt"
	"sync"

	"github.com/example/task-management/domain"
	"github.com/example/task-management/domain/aggregate"
	"github.com/example/task-management/domain/value"
)

// InMemoryTaskRepository is an in-memory implementation of TaskRepository for testing and demo
type InMemoryTaskRepository struct {
	tasks map[string]*aggregate.Task
	mu    sync.RWMutex
}

// NewInMemoryTaskRepository creates a new InMemoryTaskRepository
func NewInMemoryTaskRepository() *InMemoryTaskRepository {
	return &InMemoryTaskRepository{
		tasks: make(map[string]*aggregate.Task),
	}
}

// Save persists a task to the repository
func (r *InMemoryTaskRepository) Save(task *aggregate.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if task == nil {
		return fmt.Errorf("task cannot be nil")
	}

	r.tasks[task.ID().Value()] = task
	return nil
}

// GetByID retrieves a task by ID
func (r *InMemoryTaskRepository) GetByID(id value.TaskID) (*aggregate.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	task, exists := r.tasks[id.Value()]
	if !exists {
		return nil, fmt.Errorf("task not found")
	}

	return task, nil
}

// GetByProjectID retrieves all tasks for a project
func (r *InMemoryTaskRepository) GetByProjectID(projectID value.ProjectID) ([]*aggregate.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tasks := make([]*aggregate.Task, 0)
	for _, task := range r.tasks {
		if task.ProjectID().Equals(projectID) {
			tasks = append(tasks, task)
		}
	}

	return tasks, nil
}

// GetByAssigneeID retrieves all tasks assigned to a user
func (r *InMemoryTaskRepository) GetByAssigneeID(userID value.UserID) ([]*aggregate.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tasks := make([]*aggregate.Task, 0)
	for _, task := range r.tasks {
		if task.Assignee() != nil && task.Assignee().IsAssignedTo(userID) {
			tasks = append(tasks, task)
		}
	}

	return tasks, nil
}

// GetByStatus retrieves all tasks with a specific status
func (r *InMemoryTaskRepository) GetByStatus(status value.TaskStatus) ([]*aggregate.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tasks := make([]*aggregate.Task, 0)
	for _, task := range r.tasks {
		if task.Status() == status {
			tasks = append(tasks, task)
		}
	}

	return tasks, nil
}

// GetAll retrieves all tasks
func (r *InMemoryTaskRepository) GetAll() ([]*aggregate.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tasks := make([]*aggregate.Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// Delete removes a task from the repository
func (r *InMemoryTaskRepository) Delete(id value.TaskID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tasks[id.Value()]; !exists {
		return fmt.Errorf("task not found")
	}

	delete(r.tasks, id.Value())
	return nil
}

// Update updates an existing task
func (r *InMemoryTaskRepository) Update(task *aggregate.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if task == nil {
		return fmt.Errorf("task cannot be nil")
	}

	if _, exists := r.tasks[task.ID().Value()]; !exists {
		return fmt.Errorf("task not found")
	}

	r.tasks[task.ID().Value()] = task
	return nil
}

// FindByProjectIDAndStatus retrieves tasks for a project with specific status
func (r *InMemoryTaskRepository) FindByProjectIDAndStatus(
	projectID value.ProjectID,
	status value.TaskStatus,
) ([]*aggregate.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tasks := make([]*aggregate.Task, 0)
	for _, task := range r.tasks {
		if task.ProjectID().Equals(projectID) && task.Status() == status {
			tasks = append(tasks, task)
		}
	}

	return tasks, nil
}

// Ensure InMemoryTaskRepository implements domain.TaskRepository
var _ domain.TaskRepository = (*InMemoryTaskRepository)(nil)