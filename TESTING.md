# Testing Guide

## Overview

This project uses a multi-layered testing approach:

1. **Unit Tests**: Test individual domain components
2. **Integration Tests**: Test complete command/query flows
3. **Property-Based Tests** (optional): Test properties hold across inputs

## Running Tests

```bash
# All tests
make test

# Only unit tests
make test-unit

# Only integration tests
make test-int

# With verbose output and coverage
go test -v -cover ./...

# Run specific test
go test -v ./tests/unit -run TestTaskCreation

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Unit Tests

Located in: `/tests/unit/domain_test.go`

### Test Categories

#### 1. Domain Model Tests
Test aggregate creation and validation:

```go
func TestTaskCreation(t *testing.T) {
    // Setup
    taskID := value.GenerateTaskID()
    projectID := value.GenerateProjectID()
    priority, _ := value.NewPriority("HIGH")
    userID := value.GenerateUserID()

    // Execute
    task, err := aggregate.NewTask(
        taskID,
        projectID,
        "Test Task",
        "Test Description",
        priority,
        userID,
    )

    // Assert
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    
    if task.Title() != "Test Task" {
        t.Errorf("Expected title 'Test Task', got '%s'", task.Title())
    }
}
```

#### 2. Value Object Tests
Test value object validation:

```go
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
```

#### 3. Business Rules Tests
Test domain-specific rules:

```go
func TestTaskInvalidStatusTransition(t *testing.T) {
    // Create task
    task, _ := aggregate.NewTask(...)

    // Try invalid transition
    err := task.ChangeStatus(value.TaskStatusCompleted)
    
    if err == nil {
        t.Fatal("Expected error for invalid transition")
    }
}
```

### Creating New Unit Tests

1. Add test function to appropriate file
2. Use `Test` prefix and descriptive name
3. Follow Arrange-Act-Assert pattern
4. Use table-driven tests for multiple scenarios

Example:
```go
func TestPriorityComparison(t *testing.T) {
    tests := []struct {
        name     string
        p1       Priority
        p2       Priority
        expected bool
    }{
        {"Low vs High", PriorityLow, PriorityHigh, true},
        {"Same priority", PriorityMedium, PriorityMedium, false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := tt.p1.Numeric() < tt.p2.Numeric()
            if result != tt.expected {
                t.Errorf("got %v, want %v", result, tt.expected)
            }
        })
    }
}
```

## Integration Tests

Located in: `/tests/integration/command_test.go`

### Test Pattern

Integration tests verify complete command/query flows:

```go
func TestCreateTaskCommandFlow(t *testing.T) {
    // SETUP: Create dependencies
    container := di.NewContainer()
    
    // Create necessary test data
    user, _ := aggregate.NewUser(...)
    container.UserRepository.Save(user)
    
    project, _ := aggregate.NewProject(...)
    container.ProjectRepository.Save(project)

    // EXECUTE: Run command
    cmd := command.CreateTaskCommand{
        ProjectID: projectID.Value(),
        Title:     "Integration Test Task",
        Priority:  "HIGH",
        CreatedBy: userID.Value(),
    }
    
    result, err := container.CreateTaskCommandHandler.Handle(cmd)

    // ASSERT: Verify results
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }

    // Verify persistence
    taskID, _ := value.NewTaskID(result.TaskID)
    savedTask, _ := container.TaskRepository.GetByID(taskID)
    
    if savedTask.Title() != "Integration Test Task" {
        t.Errorf("Task not persisted correctly")
    }
}
```

### Creating New Integration Tests

1. Add test function to appropriate command test file
2. Setup DI container
3. Create test data (users, projects, workflows)
4. Execute command/query
5. Assert results and side effects
6. Verify repository persistence

Example structure:
```go
func TestNewCommandFlow(t *testing.T) {
    // 1. Setup
    container := di.NewContainer()
    testData := setupTestData(container)

    // 2. Execute
    cmd := NewCommand{...}
    result, err := handler.Handle(cmd)

    // 3. Assert
    if err != nil {
        t.Fatal(err)
    }

    // 4. Verify persistence
    retrieved, _ := container.Repository.GetByID(result.ID)
    if !match(retrieved, expected) {
        t.Error("Persistence failed")
    }
}
```

## Mocking and Test Doubles

### Mock Repository Example

```go
type MockTaskRepository struct {
    SavedTasks map[string]*aggregate.Task
}

func (m *MockTaskRepository) Save(task *aggregate.Task) error {
    m.SavedTasks[task.ID().Value()] = task
    return nil
}

func (m *MockTaskRepository) GetByID(id value.TaskID) (*aggregate.Task, error) {
    task, exists := m.SavedTasks[id.Value()]
    if !exists {
        return nil, fmt.Errorf("task not found")
    }
    return task, nil
}
```

Usage:
```go
func TestWithMockRepository(t *testing.T) {
    mockRepo := &MockTaskRepository{
        SavedTasks: make(map[string]*aggregate.Task),
    }

    // Use mock instead of real repository
    handler := command.NewCreateTaskCommandHandler(mockRepo, ...)
    result, err := handler.Handle(cmd)
}
```

## Test Coverage

### Checking Coverage

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...

# View in terminal
go tool cover -func=coverage.out

# Generate HTML report
go tool cover -html=coverage.out -o coverage.html
```

### Coverage Goals

- **Domain layer**: 80%+ (critical business logic)
- **Application layer**: 70%+ (orchestration logic)
- **Infrastructure layer**: 60%+ (harder to test external dependencies)

### Improving Coverage

Focus on testing:
1. Happy paths (normal operation)
2. Error conditions (validation failures)
3. Edge cases (boundaries, nulls)
4. Integration between layers

## Performance Testing

### Simple Benchmark

```go
func BenchmarkTaskCreation(b *testing.B) {
    projectID := value.GenerateProjectID()
    priority, _ := value.NewPriority("HIGH")
    userID := value.GenerateUserID()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        taskID := value.GenerateTaskID()
        aggregate.NewTask(
            taskID,
            projectID,
            "Task",
            "Description",
            priority,
            userID,
        )
    }
}

// Run: go test -bench=. ./tests/unit
```

## Test Organization

### Test Naming Convention

```go
// Format: Test<Component><Scenario><Expected>
func TestTaskCreationWithValidData(t *testing.T)
func TestTaskAssignmentToNonexistentUser(t *testing.T)
func TestStatusTransitionFromInProgressToReview(t *testing.T)
```

### Grouping Tests

Use subtests for related scenarios:

```go
func TestTaskStatusTransitions(t *testing.T) {
    t.Run("valid transitions", func(t *testing.T) {
        // Test valid transitions
    })

    t.Run("invalid transitions", func(t *testing.T) {
        // Test invalid transitions
    })

    t.Run("edge cases", func(t *testing.T) {
        // Test edge cases
    })
}
```

## CI/CD Testing

### GitHub Actions Example

```yaml
name: Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: 1.21
      - name: Run tests
        run: make test
      - name: Generate coverage
        run: go test -coverprofile=coverage.out ./...
```

## Testing Best Practices

1. **Keep tests focused**: One assertion per test when possible
2. **Use meaningful names**: Test name should describe what's being tested
3. **Setup and cleanup**: Use test helpers for common setup
4. **Don't test implementation**: Test behavior and contracts
5. **Test edge cases**: Nulls, empty values, boundaries
6. **Avoid test interdependence**: Each test should be independent
7. **Use table-driven tests**: For multiple similar scenarios
8. **Mock external dependencies**: Don't test 3rd-party code
9. **Test error conditions**: Not just happy paths
10. **Keep tests fast**: Unit tests should run quickly

## Troubleshooting Tests

### Test Fails Intermittently

- Check for race conditions: `go test -race ./...`
- Look for timing issues in async code
- Ensure test isolation (no shared state)

### Test Hangs

- Check for deadlocks in goroutines
- Verify goroutine cleanup
- Use timeouts in tests

### Coverage Missing

- Check if code is actually tested
- Look for error paths not tested
- Add negative test cases

## Further Reading

- Go Testing: https://golang.org/pkg/testing/
- Table-driven tests: https://golang.org/wiki/TableDrivenTests
- Testable code patterns: https://golang.org/doc/effective_go#interfaces