package dto

import "time"

// TaskDTO is the data transfer object for Task
type TaskDTO struct {
	ID          string            `json:"id"`
	ProjectID   string            `json:"project_id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      string            `json:"status"`
	Priority    string            `json:"priority"`
	Assignee    *AssignmentDTO    `json:"assignee,omitempty"`
	Deadline    *DeadlineDTO      `json:"deadline,omitempty"`
	Comments    []CommentDTO      `json:"comments,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	CreatedBy   string            `json:"created_by"`
}

// CommentDTO is the data transfer object for Comment
type CommentDTO struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	AuthorID  string    `json:"author_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AssignmentDTO is the data transfer object for Assignment
type AssignmentDTO struct {
	AssigneeID string    `json:"assignee_id"`
	AssignedAt time.Time `json:"assigned_at"`
	AssignedBy string    `json:"assigned_by"`
}

// DeadlineDTO is the data transfer object for Deadline
type DeadlineDTO struct {
	DueDate     time.Time `json:"due_date"`
	IsOverdue   bool      `json:"is_overdue"`
	DaysUntil   int       `json:"days_until"`
}

// CreateTaskRequest represents the request to create a task
type CreateTaskRequest struct {
	ProjectID   string `json:"project_id" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Priority    string `json:"priority" binding:"required"`
	AssigneeID  string `json:"assignee_id"`
	Deadline    string `json:"deadline"`
}

// UpdateTaskRequest represents the request to update a task
type UpdateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	Status      string `json:"status"`
	Deadline    string `json:"deadline"`
}

// AssignTaskRequest represents the request to assign a task
type AssignTaskRequest struct {
	AssigneeID string `json:"assignee_id" binding:"required"`
}

// UpdateTaskStatusRequest represents the request to update task status
type UpdateTaskStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// AddCommentRequest represents the request to add a comment
type AddCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

// SetDeadlineRequest represents the request to set a deadline
type SetDeadlineRequest struct {
	DueDate string `json:"due_date" binding:"required"`
}