package entity

import (
	"fmt"
	"time"

	"github.com/example/task-management/domain/value"
	"github.com/google/uuid"
)

// Comment represents a comment on a task
type Comment struct {
	id        string
	taskID    value.TaskID
	authorID  value.UserID
	content   string
	createdAt time.Time
	updatedAt time.Time
}

// NewComment creates a new Comment
func NewComment(taskID value.TaskID, authorID value.UserID, content string) (*Comment, error) {
	if content == "" {
		return nil, fmt.Errorf("comment content cannot be empty")
	}

	return &Comment{
		id:        uuid.New().String(),
		taskID:    taskID,
		authorID:  authorID,
		content:   content,
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}, nil
}

// ID returns the comment ID
func (c *Comment) ID() string {
	return c.id
}

// TaskID returns the task ID
func (c *Comment) TaskID() value.TaskID {
	return c.taskID
}

// AuthorID returns the author ID
func (c *Comment) AuthorID() value.UserID {
	return c.authorID
}

// Content returns the comment content
func (c *Comment) Content() string {
	return c.content
}

// CreatedAt returns the creation timestamp
func (c *Comment) CreatedAt() time.Time {
	return c.createdAt
}

// UpdatedAt returns the last update timestamp
func (c *Comment) UpdatedAt() time.Time {
	return c.updatedAt
}

// Update updates the comment content
func (c *Comment) Update(newContent string) error {
	if newContent == "" {
		return fmt.Errorf("comment content cannot be empty")
	}
	c.content = newContent
	c.updatedAt = time.Now()
	return nil
}