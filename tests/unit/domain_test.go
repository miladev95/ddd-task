package unit

import (
	"testing"
	"time"

	"github.com/example/task-management/domain/aggregate"
	"github.com/example/task-management/domain/entity"
	"github.com/example/task-management/domain/value"
)

// TestTaskCreation tests creating a new task
func TestTaskCreation(t *testing.T) {
	taskID := value.GenerateTaskID()
	projectID := value.GenerateProjectID()
	priority, _ := value.NewPriority("HIGH")
	userID := value.GenerateUserID()

	task, err := aggregate.NewTask(
		taskID,
		projectID,
		"Test Task",
		"Test Description",
		priority,
		userID,
	)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if task == nil {
		t.Fatal("Expected task to be created")
	}

	if task.Title() != "Test Task" {
		t.Errorf("Expected title 'Test Task', got '%s'", task.Title())
	}

	if task.Status() != value.TaskStatusToDo {
		t.Errorf("Expected status TO_DO, got %s", task.Status().Value())
	}
}

// TestTaskAssignment tests assigning a task to a user
func TestTaskAssignment(t *testing.T) {
	taskID := value.GenerateTaskID()
	projectID := value.GenerateProjectID()
	priority, _ := value.NewPriority("MEDIUM")
	creatorID := value.GenerateUserID()
	assigneeID := value.GenerateUserID()

	task, _ := aggregate.NewTask(taskID, projectID, "Task", "Description", priority, creatorID)

	err := task.Assign(assigneeID, creatorID)
	if err != nil {
		t.Fatalf("Expected no error assigning task, got %v", err)
	}

	if task.Assignee() == nil {
		t.Fatal("Expected task to be assigned")
	}

	if !task.Assignee().IsAssignedTo(assigneeID) {
		t.Error("Expected assignee to match")
	}
}

// TestTaskStatusTransition tests valid status transitions
func TestTaskStatusTransition(t *testing.T) {
	taskID := value.GenerateTaskID()
	projectID := value.GenerateProjectID()
	priority, _ := value.NewPriority("LOW")
	userID := value.GenerateUserID()

	task, _ := aggregate.NewTask(taskID, projectID, "Task", "Description", priority, userID)

	// Valid transition: TO_DO -> IN_PROGRESS
	err := task.ChangeStatus(value.TaskStatusInProgress)
	if err != nil {
		t.Fatalf("Expected valid transition, got error: %v", err)
	}

	if task.Status() != value.TaskStatusInProgress {
		t.Errorf("Expected status IN_PROGRESS, got %s", task.Status().Value())
	}
}

// TestTaskInvalidStatusTransition tests invalid status transitions
func TestTaskInvalidStatusTransition(t *testing.T) {
	taskID := value.GenerateTaskID()
	projectID := value.GenerateProjectID()
	priority, _ := value.NewPriority("CRITICAL")
	userID := value.GenerateUserID()

	task, _ := aggregate.NewTask(taskID, projectID, "Task", "Description", priority, userID)

	// Invalid transition: TO_DO -> COMPLETED (must go through IN_PROGRESS and IN_REVIEW)
	err := task.ChangeStatus(value.TaskStatusCompleted)
	if err == nil {
		t.Fatal("Expected error for invalid transition")
	}
}

// TestDeadlineValidation tests deadline validation
func TestDeadlineValidation(t *testing.T) {
	futureDate := time.Now().AddDate(0, 0, 7)
	deadline, err := value.NewDeadline(futureDate)

	if err != nil {
		t.Fatalf("Expected no error for future deadline, got %v", err)
	}

	if deadline.IsOverdue() {
		t.Error("Expected deadline not to be overdue")
	}
}

// TestDeadlineOverdue tests overdue deadline detection
func TestDeadlineOverdue(t *testing.T) {
	pastDate := time.Now().AddDate(0, 0, -1)
	_, err := value.NewDeadline(pastDate)

	if err == nil {
		t.Fatal("Expected error for past deadline")
	}
}

// TestCommentCreation tests creating a comment
func TestCommentCreation(t *testing.T) {
	taskID := value.GenerateTaskID()
	authorID := value.GenerateUserID()

	comment, err := entity.NewComment(taskID, authorID, "This is a test comment")

	if err != nil {
		t.Fatalf("Expected no error creating comment, got %v", err)
	}

	if comment.Content() != "This is a test comment" {
		t.Errorf("Expected comment content to match")
	}

	if !comment.TaskID().Equals(taskID) {
		t.Error("Expected task ID to match")
	}
}

// TestProjectCreation tests creating a project
func TestProjectCreation(t *testing.T) {
	projectID := value.GenerateProjectID()
	ownerID := value.GenerateUserID()
	workflowID := value.GenerateWorkflowID()

	project, err := aggregate.NewProject(
		projectID,
		"Test Project",
		"Test project description",
		ownerID,
		workflowID,
	)

	if err != nil {
		t.Fatalf("Expected no error creating project, got %v", err)
	}

	if project.Name() != "Test Project" {
		t.Errorf("Expected project name 'Test Project', got '%s'", project.Name())
	}

	if project.IsArchived() {
		t.Error("Expected project not to be archived by default")
	}
}

// TestProjectAddTask tests adding a task to a project
func TestProjectAddTask(t *testing.T) {
	projectID := value.GenerateProjectID()
	ownerID := value.GenerateUserID()
	workflowID := value.GenerateWorkflowID()
	taskID := value.GenerateTaskID()

	project, _ := aggregate.NewProject(projectID, "Project", "Description", ownerID, workflowID)

	err := project.AddTask(taskID)
	if err != nil {
		t.Fatalf("Expected no error adding task, got %v", err)
	}

	if project.TaskCount() != 1 {
		t.Errorf("Expected 1 task, got %d", project.TaskCount())
	}
}

// TestUserCreation tests creating a user
func TestUserCreation(t *testing.T) {
	userID := value.GenerateUserID()

	user, err := aggregate.NewUser(
		userID,
		"test@example.com",
		"John",
		"Doe",
	)

	if err != nil {
		t.Fatalf("Expected no error creating user, got %v", err)
	}

	if user.FullName() != "John Doe" {
		t.Errorf("Expected full name 'John Doe', got '%s'", user.FullName())
	}

	if !user.IsActive() {
		t.Error("Expected user to be active by default")
	}
}

// TestValueObjectEquality tests value object equality
func TestValueObjectEquality(t *testing.T) {
	id1 := value.GenerateTaskID()
	id2 := value.GenerateTaskID()

	if id1.Equals(id2) {
		t.Error("Expected different IDs to not be equal")
	}

	if !id1.Equals(id1) {
		t.Error("Expected same ID to be equal")
	}
}