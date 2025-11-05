package value

import (
	"fmt"
	"time"
)

// Deadline represents a task deadline
type Deadline struct {
	dueDate time.Time
}

// NewDeadline creates a new Deadline
func NewDeadline(dueDate time.Time) (Deadline, error) {
	if dueDate.Before(time.Now()) {
		return Deadline{}, fmt.Errorf("deadline cannot be in the past")
	}
	return Deadline{dueDate: dueDate}, nil
}

// Value returns the time.Time representation
func (d Deadline) Value() time.Time {
	return d.dueDate
}

// IsOverdue checks if the deadline is overdue
func (d Deadline) IsOverdue() bool {
	return d.dueDate.Before(time.Now())
}

// IsDueSoon checks if the deadline is due within the specified duration
func (d Deadline) IsDueSoon(duration time.Duration) bool {
	now := time.Now()
	return d.dueDate.After(now) && d.dueDate.Before(now.Add(duration))
}

// DaysUntilDue returns the number of days until the deadline
func (d Deadline) DaysUntilDue() int {
	now := time.Now()
	return int(d.dueDate.Sub(now).Hours() / 24)
}

// String returns the string representation
func (d Deadline) String() string {
	return d.dueDate.Format("2006-01-02 15:04:05")
}