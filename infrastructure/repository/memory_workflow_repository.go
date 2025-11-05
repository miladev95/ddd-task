package repository

import (
	"fmt"
	"sync"

	"github.com/miladev95/ddd-task/domain"
	"github.com/miladev95/ddd-task/domain/aggregate"
	"github.com/miladev95/ddd-task/domain/value"
)

// InMemoryWorkflowRepository is an in-memory implementation of WorkflowRepository
type InMemoryWorkflowRepository struct {
	workflows map[string]*aggregate.Workflow
	mu        sync.RWMutex
}

// NewInMemoryWorkflowRepository creates a new InMemoryWorkflowRepository
func NewInMemoryWorkflowRepository() *InMemoryWorkflowRepository {
	return &InMemoryWorkflowRepository{
		workflows: make(map[string]*aggregate.Workflow),
	}
}

// Save persists a workflow to the repository
func (r *InMemoryWorkflowRepository) Save(workflow *aggregate.Workflow) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if workflow == nil {
		return fmt.Errorf("workflow cannot be nil")
	}

	r.workflows[workflow.ID().Value()] = workflow
	return nil
}

// GetByID retrieves a workflow by ID
func (r *InMemoryWorkflowRepository) GetByID(id value.WorkflowID) (*aggregate.Workflow, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	workflow, exists := r.workflows[id.Value()]
	if !exists {
		return nil, fmt.Errorf("workflow not found")
	}

	return workflow, nil
}

// GetByName retrieves a workflow by name
func (r *InMemoryWorkflowRepository) GetByName(name string) (*aggregate.Workflow, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, workflow := range r.workflows {
		if workflow.Name() == name {
			return workflow, nil
		}
	}

	return nil, fmt.Errorf("workflow not found")
}

// GetAll retrieves all workflows
func (r *InMemoryWorkflowRepository) GetAll() ([]*aggregate.Workflow, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	workflows := make([]*aggregate.Workflow, 0, len(r.workflows))
	for _, workflow := range r.workflows {
		workflows = append(workflows, workflow)
	}

	return workflows, nil
}

// Delete removes a workflow from the repository
func (r *InMemoryWorkflowRepository) Delete(id value.WorkflowID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.workflows[id.Value()]; !exists {
		return fmt.Errorf("workflow not found")
	}

	delete(r.workflows, id.Value())
	return nil
}

// Update updates an existing workflow
func (r *InMemoryWorkflowRepository) Update(workflow *aggregate.Workflow) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if workflow == nil {
		return fmt.Errorf("workflow cannot be nil")
	}

	if _, exists := r.workflows[workflow.ID().Value()]; !exists {
		return fmt.Errorf("workflow not found")
	}

	r.workflows[workflow.ID().Value()] = workflow
	return nil
}

// GetActive retrieves all active workflows
func (r *InMemoryWorkflowRepository) GetActive() ([]*aggregate.Workflow, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	workflows := make([]*aggregate.Workflow, 0)
	for _, workflow := range r.workflows {
		if workflow.IsActive() {
			workflows = append(workflows, workflow)
		}
	}

	return workflows, nil
}

// Ensure InMemoryWorkflowRepository implements domain.WorkflowRepository
var _ domain.WorkflowRepository = (*InMemoryWorkflowRepository)(nil)