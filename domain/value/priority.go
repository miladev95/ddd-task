package value

import "fmt"

// Priority represents the priority level of a task
type Priority string

const (
	PriorityLow    Priority = "LOW"
	PriorityMedium Priority = "MEDIUM"
	PriorityHigh   Priority = "HIGH"
	PriorityCritical Priority = "CRITICAL"
)

// NewPriority creates a new Priority from string
func NewPriority(priority string) (Priority, error) {
	p := Priority(priority)
	switch p {
	case PriorityLow, PriorityMedium, PriorityHigh, PriorityCritical:
		return p, nil
	default:
		return "", fmt.Errorf("invalid priority: %s", priority)
	}
}

// Value returns the string representation
func (p Priority) Value() string {
	return string(p)
}

// IsValid checks if the priority is valid
func (p Priority) IsValid() bool {
	switch p {
	case PriorityLow, PriorityMedium, PriorityHigh, PriorityCritical:
		return true
	default:
		return false
	}
}

// Numeric returns a numeric representation for comparison
func (p Priority) Numeric() int {
	switch p {
	case PriorityLow:
		return 1
	case PriorityMedium:
		return 2
	case PriorityHigh:
		return 3
	case PriorityCritical:
		return 4
	default:
		return 0
	}
}