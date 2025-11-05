package integration

import (
	"testing"
	"time"

	"github.com/miladev95/ddd-task/application/command"
	"github.com/miladev95/ddd-task/domain/aggregate"
	"github.com/miladev95/ddd-task/domain/value"
	"github.com/miladev95/ddd-task/shared/di"
)

// TestCreateTaskCommandFlow tests the complete create task command flow
func TestCreateTaskCommandFlow(t *testing.T) {
	// Setup
	container := di.NewContainer()

	// Create users
	userID := value.GenerateUserID()
	user, _ := aggregate.NewUser(userID, "user@example.com", "Test", "User")
	container.UserRepository.Save(user)

	// Create project
	projectID := value.GenerateProjectID()
	workflowID := value.GenerateWorkflowID()
	workflow, _ := aggregate.NewWorkflow(
		workflowID,
		"Test Workflow",
		"Test",
		[]aggregate.WorkflowStatus{
			aggregate.NewWorkflowStatus("TO_DO", "To Do", 1, false),
			aggregate.NewWorkflowStatus("DONE", "Done", 2, true),
		},
	)
	container.WorkflowRepository.Save(workflow)

	project, _ := aggregate.NewProject(projectID, "Test Project", "Test", userID, workflowID)
	container.ProjectRepository.Save(project)

	// Create command
	cmd := command.CreateTaskCommand{
		ProjectID:   projectID.Value(),
		Title:       "Integration Test Task",
		Description: "Testing command flow",
		Priority:    "HIGH",
		AssigneeID:  userID.Value(),
		Deadline:    time.Now().AddDate(0, 0, 1).Format(time.RFC3339),
		CreatedBy:   userID.Value(),
	}

	// Execute
	result, err := container.CreateTaskCommandHandler.Handle(cmd)

	// Verify
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.TaskID == "" {
		t.Fatal("Expected task ID to be generated")
	}

	// Verify task was saved
	taskID, _ := value.NewTaskID(result.TaskID)
	savedTask, err := container.TaskRepository.GetByID(taskID)
	if err != nil {
		t.Fatalf("Expected saved task, got error: %v", err)
	}

	if savedTask.Title() != "Integration Test Task" {
		t.Errorf("Expected title match")
	}

	if savedTask.Assignee() == nil {
		t.Error("Expected task to be assigned")
	}

	if savedTask.Deadline() == nil {
		t.Error("Expected task to have deadline")
	}
}

// TestAssignTaskCommandFlow tests the complete assign task command flow
func TestAssignTaskCommandFlow(t *testing.T) {
	// Setup
	container := di.NewContainer()

	// Create users
	creatorID := value.GenerateUserID()
	creator, _ := aggregate.NewUser(creatorID, "creator@example.com", "Creator", "User")
	container.UserRepository.Save(creator)

	assigneeID := value.GenerateUserID()
	assignee, _ := aggregate.NewUser(assigneeID, "assignee@example.com", "Assignee", "User")
	container.UserRepository.Save(assignee)

	// Create task
	taskID := value.GenerateTaskID()
	projectID := value.GenerateProjectID()
	priority, _ := value.NewPriority("MEDIUM")
	task, _ := aggregate.NewTask(taskID, projectID, "Test Task", "Description", priority, creatorID)
	container.TaskRepository.Save(task)

	// Create command
	cmd := command.AssignTaskCommand{
		TaskID:     taskID.Value(),
		AssigneeID: assigneeID.Value(),
		AssignedBy: creatorID.Value(),
	}

	// Execute
	result, err := container.AssignTaskCommandHandler.Handle(cmd)

	// Verify
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.Error != nil {
		t.Fatalf("Expected no error in result, got %v", result.Error)
	}

	// Verify task was updated
	updatedTask, _ := container.TaskRepository.GetByID(taskID)
	if updatedTask.Assignee() == nil {
		t.Error("Expected task to be assigned")
	}
}

// TestUpdateTaskStatusCommandFlow tests the complete update status command flow
func TestUpdateTaskStatusCommandFlow(t *testing.T) {
	// Setup
	container := di.NewContainer()

	// Create task
	userID := value.GenerateUserID()
	taskID := value.GenerateTaskID()
	projectID := value.GenerateProjectID()
	priority, _ := value.NewPriority("LOW")

	user, _ := aggregate.NewUser(userID, "test@example.com", "Test", "User")
	container.UserRepository.Save(user)

	task, _ := aggregate.NewTask(taskID, projectID, "Test Task", "Description", priority, userID)
	// Assign task first (business rule: must be assigned before IN_PROGRESS)
	task.Assign(userID, userID)
	// Set deadline (business rule: must have deadline before completion)
	futureDate := time.Now().AddDate(0, 0, 7) // 7 days from now
	deadline, _ := value.NewDeadline(futureDate)
	task.SetDeadline(deadline)
	container.TaskRepository.Save(task)

	// Create command
	cmd := command.UpdateTaskStatusCommand{
		TaskID:    taskID.Value(),
		NewStatus: "IN_PROGRESS",
	}

	// Execute
	result, err := container.UpdateTaskStatusCommandHandler.Handle(cmd)

	// Verify
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.Error != nil {
		t.Fatalf("Expected no error in result, got %v", result.Error)
	}

	// Verify task was updated
	updatedTask, _ := container.TaskRepository.GetByID(taskID)
	if updatedTask.Status() != value.TaskStatusInProgress {
		t.Errorf("Expected status IN_PROGRESS, got %s", updatedTask.Status().Value())
	}
}