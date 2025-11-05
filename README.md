# Task Management System - DDD Backend Architecture

A comprehensive Domain-Driven Design (DDD) implementation in Go for a task management system.

## Architecture Overview

```
task-management/
├── domain/                          # Domain Layer (Business Rules)
│   ├── aggregate/
│   │   ├── task.go                 # Task Aggregate Root
│   │   ├── project.go              # Project Aggregate Root
│   │   ├── user.go                 # User Aggregate Root
│   │   └── workflow.go             # Workflow Aggregate Root
│   ├── entity/
│   │   ├── comment.go              # Comment Entity
│   │   └── assignment.go           # Assignment Entity
│   ├── value/
│   │   ├── task_status.go          # Task Status Value Object
│   │   ├── priority.go             # Priority Value Object
│   │   ├── deadline.go             # Deadline Value Object
│   │   └── identifier.go           # Domain Identifiers
│   ├── event/
│   │   ├── domain_event.go         # Domain Event Interface
│   │   └── task_events.go          # Task-related Events
│   ├── service/
│   │   ├── task_assignment.go      # Task Assignment Service
│   │   ├── status_transition.go    # Status Transition Service
│   │   ├── deadline_enforcement.go # Deadline Enforcement Service
│   │   └── notification_service.go # Notification Service
│   └── repository.go               # Repository Interfaces
│
├── application/                     # Application Layer (Use Cases)
│   ├── command/
│   │   ├── create_task.go          # Command Handlers
│   │   ├── assign_task.go
│   │   ├── update_task_status.go
│   │   └── create_project.go
│   ├── query/
│   │   ├── get_task.go             # Query Handlers
│   │   ├── list_tasks_by_project.go
│   │   └── get_user_tasks.go
│   ├── dto/
│   │   ├── task_dto.go             # Data Transfer Objects
│   │   └── project_dto.go
│   └── service/
│       ├── task_application_service.go
│       └── project_application_service.go
│
├── infrastructure/                  # Infrastructure Layer
│   ├── repository/
│   │   ├── task_repository.go      # Repository Implementations
│   │   ├── project_repository.go
│   │   └── user_repository.go
│   ├── persistence/
│   │   ├── migrations.go           # Database Migrations
│   │   └── connection.go           # DB Connection Setup
│   ├── event/
│   │   └── event_publisher.go      # Event Publishing
│   └── config/
│       └── config.go               # Configuration Management
│
├── interface/                       # Interface/Presentation Layer
│   ├── http/
│   │   ├── handler/
│   │   │   ├── task_handler.go     # HTTP Handlers
│   │   │   └── project_handler.go
│   │   ├── middleware/
│   │   │   ├── auth.go             # Authentication
│   │   │   └── error_handler.go    # Error Handling
│   │   └── router.go               # Route Configuration
│   └── cli/
│       └── commands.go             # CLI Commands
│
├── shared/                          # Shared/Cross-cutting Concerns
│   ├── di/
│   │   └── container.go            # Dependency Injection
│   ├── errors/
│   │   └── errors.go               # Domain Errors
│   ├── logger/
│   │   └── logger.go               # Logging
│   └── utils/
│       └── validator.go            # Validation Utilities
│
├── tests/                           # Tests
│   ├── integration/
│   │   ├── task_service_test.go
│   │   └── repository_test.go
│   └── unit/
│       ├── domain_test.go
│       └── service_test.go
│
├── main.go                          # Application Entry Point
└── Makefile                         # Build Automation
```

## Key Design Patterns

### 1. Aggregates
- **Task Aggregate**: Task (root) with Comments as part entities
- **Project Aggregate**: Project (root) with Tasks as part entities
- **User Aggregate**: User (root) with Profile as part entity
- **Workflow Aggregate**: Workflow (root) with Status definitions

### 2. Domain Services
- **TaskAssignmentService**: Handles task assignment logic
- **StatusTransitionService**: Manages valid status transitions
- **DeadlineEnforcementService**: Validates and enforces deadlines
- **NotificationService**: Triggers notifications on domain events

### 3. Value Objects
- TaskStatus, Priority, Deadline
- TaskID, ProjectID, UserID (Identifiers)
- Email, Password (for User)

### 4. Repository Pattern
- Abstractions in domain layer
- Implementations in infrastructure layer
- Supports dependency injection

### 5. Dependency Injection
- Container-based DI in shared/di
- Constructor injection throughout
- Loose coupling between layers

## Running the Application

```bash
# Initialize database
make migrate

# Run the application
make run

# Run tests
make test

# Build
make build
```

## Example Flows

### Command Flow: Create Task
1. HTTP Handler receives request
2. Handler creates CreateTaskCommand
3. Application Service processes command
4. Domain Service validates business rules
5. Task Aggregate is created
6. Repository persists Task
7. Domain Event is published
8. Response is returned to client

### Query Flow: List Tasks by Project
1. HTTP Handler receives request
2. Handler creates ListTasksByProjectQuery
3. Query Handler fetches from repository
4. DTOs are returned
5. Response is serialized and sent to client

## Testing Strategy

- **Unit Tests**: Domain models and services
- **Integration Tests**: Repository and database interactions
- **Acceptance Tests**: End-to-end command/query flows
- **Testability**: All dependencies are injected, enabling easy mocking

## Scalability Features

- Repository interfaces allow switching persistence mechanisms
- Domain events support event-driven architecture
- Service layer enables async processing
- Event publishing allows subscribing services
- Modular structure supports microservices decomposition