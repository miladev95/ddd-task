package repository

import (
	"fmt"
	"sync"

	"github.com/example/task-management/domain"
	"github.com/example/task-management/domain/aggregate"
	"github.com/example/task-management/domain/value"
)

// InMemoryUserRepository is an in-memory implementation of UserRepository
type InMemoryUserRepository struct {
	users map[string]*aggregate.User
	mu    sync.RWMutex
}

// NewInMemoryUserRepository creates a new InMemoryUserRepository
func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[string]*aggregate.User),
	}
}

// Save persists a user to the repository
func (r *InMemoryUserRepository) Save(user *aggregate.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if user == nil {
		return fmt.Errorf("user cannot be nil")
	}

	r.users[user.ID().Value()] = user
	return nil
}

// GetByID retrieves a user by ID
func (r *InMemoryUserRepository) GetByID(id value.UserID) (*aggregate.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id.Value()]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

// GetByEmail retrieves a user by email
func (r *InMemoryUserRepository) GetByEmail(email string) (*aggregate.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Email() == email {
			return user, nil
		}
	}

	return nil, fmt.Errorf("user not found")
}

// GetAll retrieves all users
func (r *InMemoryUserRepository) GetAll() ([]*aggregate.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]*aggregate.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}

	return users, nil
}

// Delete removes a user from the repository
func (r *InMemoryUserRepository) Delete(id value.UserID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[id.Value()]; !exists {
		return fmt.Errorf("user not found")
	}

	delete(r.users, id.Value())
	return nil
}

// Update updates an existing user
func (r *InMemoryUserRepository) Update(user *aggregate.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if user == nil {
		return fmt.Errorf("user cannot be nil")
	}

	if _, exists := r.users[user.ID().Value()]; !exists {
		return fmt.Errorf("user not found")
	}

	r.users[user.ID().Value()] = user
	return nil
}

// GetActive retrieves all active users
func (r *InMemoryUserRepository) GetActive() ([]*aggregate.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]*aggregate.User, 0)
	for _, user := range r.users {
		if user.IsActive() {
			users = append(users, user)
		}
	}

	return users, nil
}

// Ensure InMemoryUserRepository implements domain.UserRepository
var _ domain.UserRepository = (*InMemoryUserRepository)(nil)