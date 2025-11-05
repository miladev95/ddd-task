package main

import (
	"fmt"
	"time"

	"github.com/miladev95/ddd-task/application/command"
	"github.com/miladev95/ddd-task/domain/aggregate"
	"github.com/miladev95/ddd-task/domain/value"
	"github.com/miladev95/ddd-task/shared/di"
)

// This file demonstrates how to use the DDD architecture

func main() {
	// Initialize DI container
	container := di.NewContainer()

	// Example 1: Create users
	fmt.Println("=== Creating Users ===")
	user1ID := value.GenerateUserID()
	user1, err := aggregate.NewUser(user1ID, "alice@example.com", "Alice", "Johnson")
	if err != nil {
		fmt.Printf("Error creating user: %v\n", err)
		return
	}
	container.UserRepository.Save(user1)
	fmt.Printf("Created user: %s (%s)\n", user1.FullName(), user1ID.Value())

	user2ID := value.GenerateUserID()
	user2, err := aggregate.NewUser(user2ID, "bob@example.com", "Bob", "Smith")
	if err != nil {
		fmt.Printf("Error creating user: %v\n", err)
		return
	}
	container.UserRepository.Save(user2)
	fmt.Printf("Created user: %s (%s)\n", user2.FullName(), user2ID.Value())

	// Example 2: Create a workflow
	fmt.Println("\n=== Creating Workflow ===")
	workflowID := value.GenerateWorkflowID()
	statuses := []aggregate.WorkflowStatus{
		aggregate.NewWorkflowStatus("Backlog", "Tasks in backlog", 1, false),
		aggregate.NewWorkflowStatus("To Do", "Tasks to do", 2, false),
		aggregate.NewWorkflowStatus("In Progress", "Tasks in progress", 3, false),
		aggregate.NewWorkflowStatus("In Review", "Tasks in review", 4, false),
		aggregate.NewWorkflowStatus("Completed", "Tasks completed", 5, true),
	}
	workflow, err := aggregate.NewWorkflow(workflowID, "Default Workflow", "Default project workflow", statuses)
	if err != nil {
		fmt.Printf("Error creating workflow: %v\n", err)
		return
	}
	container.WorkflowRepository.Save(workflow)
	fmt.Printf("Created workflow: %s\n", workflow.Name())

	// Example 3: Create a project
	fmt.Println("\n=== Creating Project ===")
	projectID := value.GenerateProjectID()
	project, err := aggregate.NewProject(
		projectID,
		"Web Application",
		"Build a modern web application",
		user1ID,
		workflowID,
	)
	if err != nil {
		fmt.Printf("Error creating project: %v\n", err)
		return
	}
	container.ProjectRepository.Save(project)
	fmt.Printf("Created project: %s (Owner: %s)\n", project.Name(), user1.FullName())

	// Example 4: Create a task using command handler
	fmt.Println("\n=== Creating Task (via Command) ===")
	createTaskCmd := command.CreateTaskCommand{
		ProjectID:   projectID.Value(),
		Title:       "Implement user authentication",
		Description: "Add JWT-based authentication to the application",
		Priority:    "HIGH",
		AssigneeID:  user2ID.Value(),
		Deadline:    time.Now().AddDate(0, 0, 7).Format(time.RFC3339),
		CreatedBy:   user1ID.Value(),
	}

	result, err := container.CreateTaskCommandHandler.Handle(createTaskCmd)
	if err != nil {
		fmt.Printf("Error creating task: %v\n", err)
		return
	}
	fmt.Printf("Created task: %s (ID: %s)\n", createTaskCmd.Title, result.TaskID)

	// Example 5: Retrieve and display the task
	fmt.Println("\n=== Retrieving Task ===")
	taskID, _ := value.NewTaskID(result.TaskID)
	task, err := container.TaskRepository.GetByID(taskID)
	if err != nil {
		fmt.Printf("Error retrieving task: %v\n", err)
		return
	}

	fmt.Printf("Task ID: %s\n", task.ID().Value())
	fmt.Printf("Title: %s\n", task.Title())
	fmt.Printf("Status: %s\n", task.Status().Value())
	fmt.Printf("Priority: %s\n", task.Priority().Value())
	if task.Assignee() != nil {
		fmt.Printf("Assigned to: %s\n", task.Assignee().AssigneeID().Value())
	}
	if task.Deadline() != nil {
		fmt.Printf("Deadline: %s\n", task.Deadline().String())
	}

	// Example 6: Update task status
	fmt.Println("\n=== Updating Task Status ===")
	statusCmd := command.UpdateTaskStatusCommand{
		TaskID:    result.TaskID,
		NewStatus: "IN_PROGRESS",
	}

	_, err = container.UpdateTaskStatusCommandHandler.Handle(statusCmd)
	if err != nil {
		fmt.Printf("Error updating task status: %v\n", err)
		return
	}
	fmt.Printf("Task status updated to: IN_PROGRESS\n")

	// Example 7: List tasks by project
	fmt.Println("\n=== Listing Tasks by Project ===")
	tasks, err := container.TaskRepository.GetByProjectID(projectID)
	if err != nil {
		fmt.Printf("Error listing tasks: %v\n", err)
		return
	}

	fmt.Printf("Found %d task(s) in project\n", len(tasks))
	for _, t := range tasks {
		fmt.Printf("- %s (%s)\n", t.Title(), t.Status().Value())
	}

	// Example 8: Get tasks assigned to user
	fmt.Println("\n=== Getting Tasks Assigned to User ===")
	assignedTasks, err := container.TaskRepository.GetByAssigneeID(user2ID)
	if err != nil {
		fmt.Printf("Error getting assigned tasks: %v\n", err)
		return
	}

	fmt.Printf("%s has %d task(s) assigned\n", user2.FullName(), len(assignedTasks))
	for _, t := range assignedTasks {
		fmt.Printf("- %s (%s)\n", t.Title(), t.Status().Value())
	}

	// Example 9: Domain event demonstration
	fmt.Println("\n=== Domain Events ===")
	events := task.DomainEvents()
	fmt.Printf("Task generated %d domain events:\n", len(events))
	for _, evt := range events {
		fmt.Printf("- Event Type: %s, Aggregate: %s\n", evt.EventType(), evt.AggregateType())
	}

	// Example 10: Verify business rules - invalid status transition
	fmt.Println("\n=== Testing Business Rules ===")
	invalidStatusCmd := command.UpdateTaskStatusCommand{
		TaskID:    result.TaskID,
		NewStatus: "BACKLOG", // Invalid transition from IN_PROGRESS
	}

	_, err = container.UpdateTaskStatusCommandHandler.Handle(invalidStatusCmd)
	if err != nil {
		fmt.Printf("Expected error (invalid transition): %v\n", err)
	}

	fmt.Println("\n=== Example Complete ===")
}