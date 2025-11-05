package aggregate

import (
	"fmt"
	"time"

	"github.com/example/task-management/domain/event"
	"github.com/example/task-management/domain/value"
)

// User is the aggregate root for the User aggregate
type User struct {
	id           value.UserID
	email        string
	firstName    string
	lastName     string
	active       bool
	createdAt    time.Time
	updatedAt    time.Time
	lastLogin    *time.Time
	preferences  map[string]string
	domainEvents []event.DomainEvent
}

// NewUser creates a new User
func NewUser(
	id value.UserID,
	email, firstName, lastName string,
) (*User, error) {
	if email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}

	if firstName == "" || lastName == "" {
		return nil, fmt.Errorf("first and last name cannot be empty")
	}

	return &User{
		id:           id,
		email:        email,
		firstName:    firstName,
		lastName:     lastName,
		active:       true,
		createdAt:    time.Now(),
		updatedAt:    time.Now(),
		preferences:  make(map[string]string),
		domainEvents: make([]event.DomainEvent, 0),
	}, nil
}

// ID returns the user ID
func (u *User) ID() value.UserID {
	return u.id
}

// Email returns the user email
func (u *User) Email() string {
	return u.email
}

// FirstName returns the user first name
func (u *User) FirstName() string {
	return u.firstName
}

// LastName returns the user last name
func (u *User) LastName() string {
	return u.lastName
}

// FullName returns the user full name
func (u *User) FullName() string {
	return u.firstName + " " + u.lastName
}

// IsActive returns whether the user is active
func (u *User) IsActive() bool {
	return u.active
}

// CreatedAt returns when the user was created
func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

// UpdatedAt returns when the user was last updated
func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}

// LastLogin returns the last login time
func (u *User) LastLogin() *time.Time {
	return u.lastLogin
}

// DomainEvents returns all uncommitted domain events
func (u *User) DomainEvents() []event.DomainEvent {
	return append([]event.DomainEvent{}, u.domainEvents...)
}

// ClearDomainEvents clears all domain events after they have been published
func (u *User) ClearDomainEvents() {
	u.domainEvents = make([]event.DomainEvent, 0)
}

// Activate activates the user
func (u *User) Activate() error {
	if u.active {
		return fmt.Errorf("user is already active")
	}

	u.active = true
	u.updatedAt = time.Now()

	return nil
}

// Deactivate deactivates the user
func (u *User) Deactivate() error {
	if !u.active {
		return fmt.Errorf("user is already inactive")
	}

	u.active = false
	u.updatedAt = time.Now()

	return nil
}

// UpdateLastLogin updates the last login time
func (u *User) UpdateLastLogin() {
	now := time.Now()
	u.lastLogin = &now
	u.updatedAt = now
}

// UpdateEmail updates the user email
func (u *User) UpdateEmail(newEmail string) error {
	if newEmail == "" {
		return fmt.Errorf("email cannot be empty")
	}

	u.email = newEmail
	u.updatedAt = time.Now()

	return nil
}

// UpdateName updates the user name
func (u *User) UpdateName(firstName, lastName string) error {
	if firstName == "" || lastName == "" {
		return fmt.Errorf("first and last name cannot be empty")
	}

	u.firstName = firstName
	u.lastName = lastName
	u.updatedAt = time.Now()

	return nil
}

// SetPreference sets a user preference
func (u *User) SetPreference(key, value string) {
	u.preferences[key] = value
	u.updatedAt = time.Now()
}

// GetPreference gets a user preference
func (u *User) GetPreference(key string) (string, bool) {
	value, exists := u.preferences[key]
	return value, exists
}

// GetPreferences returns all preferences
func (u *User) GetPreferences() map[string]string {
	// Return a copy to prevent external modification
	prefs := make(map[string]string)
	for k, v := range u.preferences {
		prefs[k] = v
	}
	return prefs
}