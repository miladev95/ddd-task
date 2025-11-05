package repository

import (
	"fmt"
	"sync"

	"github.com/miladev95/ddd-task/domain"
	"github.com/miladev95/ddd-task/domain/aggregate"
	"github.com/miladev95/ddd-task/domain/value"
)

// InMemoryProjectRepository is an in-memory implementation of ProjectRepository
type InMemoryProjectRepository struct {
	projects map[string]*aggregate.Project
	mu       sync.RWMutex
}

// NewInMemoryProjectRepository creates a new InMemoryProjectRepository
func NewInMemoryProjectRepository() *InMemoryProjectRepository {
	return &InMemoryProjectRepository{
		projects: make(map[string]*aggregate.Project),
	}
}

// Save persists a project to the repository
func (r *InMemoryProjectRepository) Save(project *aggregate.Project) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if project == nil {
		return fmt.Errorf("project cannot be nil")
	}

	r.projects[project.ID().Value()] = project
	return nil
}

// GetByID retrieves a project by ID
func (r *InMemoryProjectRepository) GetByID(id value.ProjectID) (*aggregate.Project, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	project, exists := r.projects[id.Value()]
	if !exists {
		return nil, fmt.Errorf("project not found")
	}

	return project, nil
}

// GetByOwnerID retrieves all projects owned by a user
func (r *InMemoryProjectRepository) GetByOwnerID(userID value.UserID) ([]*aggregate.Project, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	projects := make([]*aggregate.Project, 0)
	for _, project := range r.projects {
		if project.OwnerID().Equals(userID) {
			projects = append(projects, project)
		}
	}

	return projects, nil
}

// GetAll retrieves all projects
func (r *InMemoryProjectRepository) GetAll() ([]*aggregate.Project, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	projects := make([]*aggregate.Project, 0, len(r.projects))
	for _, project := range r.projects {
		projects = append(projects, project)
	}

	return projects, nil
}

// Delete removes a project from the repository
func (r *InMemoryProjectRepository) Delete(id value.ProjectID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.projects[id.Value()]; !exists {
		return fmt.Errorf("project not found")
	}

	delete(r.projects, id.Value())
	return nil
}

// Update updates an existing project
func (r *InMemoryProjectRepository) Update(project *aggregate.Project) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if project == nil {
		return fmt.Errorf("project cannot be nil")
	}

	if _, exists := r.projects[project.ID().Value()]; !exists {
		return fmt.Errorf("project not found")
	}

	r.projects[project.ID().Value()] = project
	return nil
}

// GetActive retrieves all active projects
func (r *InMemoryProjectRepository) GetActive() ([]*aggregate.Project, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	projects := make([]*aggregate.Project, 0)
	for _, project := range r.projects {
		if !project.IsArchived() {
			projects = append(projects, project)
		}
	}

	return projects, nil
}

// Ensure InMemoryProjectRepository implements domain.ProjectRepository
var _ domain.ProjectRepository = (*InMemoryProjectRepository)(nil)