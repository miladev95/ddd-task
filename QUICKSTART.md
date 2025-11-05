# Quick Start Guide

## Prerequisites

- Go 1.21 or later
- Git

## Installation

```bash
# Clone the repository
git clone <repository-url>
cd task-management

# Download dependencies
go mod download
```

## Running the Application

### Option 1: Direct Execution
```bash
go run main.go
```

The API will be available at `http://localhost:8080`

### Option 2: Build and Run
```bash
make build
./bin/task-management
```

### Option 3: Using Make
```bash
make run
```

## Running Tests

```bash
# All tests
make test

# Unit tests only
make test-unit

# Integration tests only
make test-int
```

## Running Examples

```bash
make example
```

This demonstrates creating users, projects, tasks, and performing various operations.

## API Quick Reference

### Health Check
```bash
curl http://localhost:8080/health
```

### Create a Task
```bash
curl -X POST http://localhost:8080/api/tasks \
  -H "Content-Type: application/json" \
  -H "X-User-ID: user-123" \
  -d '{
    "project_id": "proj-123",
    "title": "My Task",
    "description": "Task description",
    "priority": "HIGH",
    "assignee_id": "user-456",
    "deadline": "2024-12-31T23:59:59Z"
  }'
```

### Get a Task
```bash
curl "http://localhost:8080/api/tasks/get?id=task-uuid"
```

### List Tasks by Project
```bash
curl "http://localhost:8080/api/tasks?project_id=proj-123"
```

### Assign a Task
```bash
curl -X POST "http://localhost:8080/api/tasks/assign?id=task-uuid" \
  -H "Content-Type: application/json" \
  -H "X-User-ID: user-123" \
  -d '{
    "assignee_id": "user-456"
  }'
```

### Update Task Status
```bash
curl -X PUT "http://localhost:8080/api/tasks/status?id=task-uuid" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "IN_PROGRESS"
  }'
```

## Project Structure Overview

```
task-management/
├── domain/              # Business logic and rules
├── application/         # Use cases and commands/queries
├── infrastructure/      # Repositories and external services
├── interface/          # HTTP API handlers
├── shared/             # Cross-cutting concerns (DI, etc)
├── tests/              # Unit and integration tests
├── examples/           # Usage examples
├── main.go             # Entry point
└── Makefile            # Build automation
```

## Key Concepts

### Aggregates
Root objects that encapsulate business logic:
- **Task**: Represents a work item with status, priority, assignee, deadline
- **Project**: Container for tasks
- **User**: Represents system users
- **Workflow**: Defines task status workflow

### Domain Services
Business logic that spans multiple aggregates:
- **TaskAssignmentService**: Handles task assignments
- **StatusTransitionService**: Manages status changes
- **DeadlineEnforcementService**: Validates and enforces deadlines
- **NotificationService**: Sends notifications

### Commands
Operations that modify state:
- `CreateTaskCommand`: Create a new task
- `AssignTaskCommand`: Assign task to user
- `UpdateTaskStatusCommand`: Change task status

### Queries
Operations that read state:
- `GetTaskQuery`: Retrieve a task
- `ListTasksByProjectQuery`: List tasks in a project

## Common Tasks

### Create a New User
In Go code:
```go
userID := value.GenerateUserID()
user, err := aggregate.NewUser(
    userID,
    "john@example.com",
    "John",
    "Doe",
)
if err != nil {
    log.Fatal(err)
}
container.UserRepository.Save(user)
```

### Create a Project
```go
projectID := value.GenerateProjectID()
project, err := aggregate.NewProject(
    projectID,
    "My Project",
    "Project description",
    userID,
    workflowID,
)
```

### Create a Task with Status Transition
```go
taskID := value.GenerateTaskID()
task, err := aggregate.NewTask(
    taskID,
    projectID,
    "Task Title",
    "Task Description",
    priority,
    createdByID,
)

// Assign task
err = task.Assign(assigneeID, createdByID)

// Set deadline
deadline, _ := value.NewDeadline(time.Now().AddDate(0, 0, 7))
err = task.SetDeadline(deadline)

// Change status
err = task.ChangeStatus(value.TaskStatusInProgress)

// Save
container.TaskRepository.Save(task)
```

## Troubleshooting

### Port Already in Use
If port 8080 is in use, modify `main.go` to use a different port:
```go
port := ":9000" // Change to desired port
```

### Tests Failing
Ensure Go 1.21+ is installed:
```bash
go version
```

### Module Not Found
Ensure go.mod is present and download dependencies:
```bash
go mod download
go mod tidy
```

## Next Steps

1. Read the [ARCHITECTURE.md](ARCHITECTURE.md) for detailed architecture information
2. Explore the `examples/` directory for more usage patterns
3. Check the `tests/` directory for testing patterns
4. Implement database repositories for production use
5. Add authentication and authorization
6. Set up event persistence and replay
7. Add API documentation with Swagger

## Support

For questions or issues, refer to:
- Architecture documentation: `ARCHITECTURE.md`
- Examples: `examples/usage_example.go`
- Tests: `tests/` directory
- Code comments in respective files

## License

[Add your license here]