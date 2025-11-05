# Project Index & Navigation Guide

## 📖 Documentation Files

Start here to understand the project:

### Core Documentation
1. **[README.md](README.md)** - Project overview and architecture diagram
   - High-level overview
   - Architecture layers
   - Key design patterns
   - Running the application

2. **[QUICKSTART.md](QUICKSTART.md)** - Get started quickly
   - Installation steps
   - Running the application
   - API quick reference
   - Common tasks

3. **[ARCHITECTURE.md](ARCHITECTURE.md)** - Comprehensive architecture guide (50+ pages)
   - Detailed layer-by-layer breakdown
   - Design patterns explained
   - Business rules and transitions
   - Data flow diagrams
   - Testing strategy
   - Extensibility guide

### Implementation Guides
4. **[DATABASE.md](DATABASE.md)** - Database implementation
   - PostgreSQL schema and repository implementation
   - MongoDB example
   - Migration strategies
   - Performance optimization
   - Production checklist

5. **[DEPLOYMENT.md](DEPLOYMENT.md)** - Deployment guide
   - Development environment
   - Docker deployment
   - Kubernetes setup
   - AWS/GCP deployment
   - CI/CD pipelines
   - Monitoring and logging

6. **[TESTING.md](TESTING.md)** - Testing guide
   - Unit testing
   - Integration testing
   - Mocking strategies
   - Coverage tracking
   - Best practices

### Project Information
7. **[PROJECT_SUMMARY.md](PROJECT_SUMMARY.md)** - Project summary
   - What's included
   - Architecture highlights
   - Real-world scenarios
   - Production migration steps

## 🏗️ Source Code Structure

### Domain Layer (`/domain`)
Pure business logic - no external dependencies

#### Value Objects (`/domain/value/`)
- **[identifier.go](domain/value/identifier.go)** - TaskID, ProjectID, UserID, WorkflowID
- **[task_status.go](domain/value/task_status.go)** - Task status with valid transitions
- **[priority.go](domain/value/priority.go)** - Priority levels (LOW, MEDIUM, HIGH, CRITICAL)
- **[deadline.go](domain/value/deadline.go)** - Deadline with overdue detection

#### Entities (`/domain/entity/`)
- **[comment.go](domain/entity/comment.go)** - Comment entity (part of Task aggregate)
- **[assignment.go](domain/entity/assignment.go)** - Task assignment entity

#### Aggregates (`/domain/aggregate/`)
- **[task.go](domain/aggregate/task.go)** - Task aggregate root
  - Status management with transitions
  - Assignment management
  - Deadline management
  - Comment management
  - Domain event raising
  
- **[project.go](domain/aggregate/project.go)** - Project aggregate root
  - Task collection management
  - Archive functionality
  
- **[user.go](domain/aggregate/user.go)** - User aggregate root
  - User profile management
  - Preference management
  
- **[workflow.go](domain/aggregate/workflow.go)** - Workflow aggregate root
  - Workflow status definition
  - Workflow management

#### Domain Services (`/domain/service/`)
Business logic that spans multiple aggregates

- **[task_assignment.go](domain/service/task_assignment.go)** - Task assignment service
  - Validates assignees
  - Manages reassignment
  - Checks assignment capacity
  
- **[status_transition.go](domain/service/status_transition.go)** - Status transition service
  - Validates state transitions
  - Returns valid next states
  - Enforces prerequisites
  
- **[deadline_enforcement.go](domain/service/deadline_enforcement.go)** - Deadline enforcement
  - Validates deadlines
  - Detects overdue tasks
  - Notifies about deadlines

#### Domain Events (`/domain/event/`)
- **[domain_event.go](domain/event/domain_event.go)** - Base event interface and class
- **[task_events.go](domain/event/task_events.go)** - Task-related events
  - TaskCreatedEvent
  - TaskAssignedEvent
  - TaskStatusChangedEvent
  - TaskDeadlineSetEvent
  - TaskOverdueEvent
  - TaskCompletedEvent
  - TaskDeletedEvent

#### Repository Interfaces (`/domain/repository.go`)
- TaskRepository interface
- ProjectRepository interface
- UserRepository interface
- WorkflowRepository interface
- UnitOfWork interface

### Application Layer (`/application`)
Use cases and orchestration - coordinates domain objects

#### Commands (`/application/command/`)
Modify state operations

- **[create_task.go](application/command/create_task.go)** - CreateTaskCommand and handler
- **[assign_task.go](application/command/assign_task.go)** - AssignTaskCommand and handler
- **[update_task_status.go](application/command/update_task_status.go)** - UpdateTaskStatusCommand and handler

#### Queries (`/application/query/`)
Read-only state operations

- **[get_task.go](application/query/get_task.go)** - GetTaskQuery and handler
- **[list_tasks_by_project.go](application/query/list_tasks_by_project.go)** - ListTasksByProjectQuery and handler

#### DTOs (`/application/dto/`)
Data transfer objects for HTTP communication

- **[task_dto.go](application/dto/task_dto.go)** - Task DTOs and request/response objects
- **[project_dto.go](application/dto/project_dto.go)** - Project DTOs and request/response objects

### Infrastructure Layer (`/infrastructure`)
Technical implementations and external services

#### Repositories (`/infrastructure/repository/`)
In-memory implementations (switch out for database implementations)

- **[memory_task_repository.go](infrastructure/repository/memory_task_repository.go)** - In-memory task repository
- **[memory_project_repository.go](infrastructure/repository/memory_project_repository.go)** - In-memory project repository
- **[memory_user_repository.go](infrastructure/repository/memory_user_repository.go)** - In-memory user repository
- **[memory_workflow_repository.go](infrastructure/repository/memory_workflow_repository.go)** - In-memory workflow repository

#### Event Publishing (`/infrastructure/event/`)
- **[simple_event_publisher.go](infrastructure/event/simple_event_publisher.go)** - In-memory event publisher
- **[simple_notification_service.go](infrastructure/event/simple_notification_service.go)** - Basic notification service

### Interface Layer (`/interface`)
HTTP API exposure

#### HTTP Handlers (`/interface/http/handler/`)
- **[task_handler.go](interface/http/handler/task_handler.go)** - Task HTTP handler
  - CreateTask: POST /api/tasks
  - GetTask: GET /api/tasks/get
  - ListTasksByProject: GET /api/tasks
  - AssignTask: POST /api/tasks/assign
  - UpdateTaskStatus: PUT /api/tasks/status

#### Middleware (`/interface/http/middleware/`)
- **[error_handler.go](interface/http/middleware/error_handler.go)** - HTTP error handling
  - Converts domain errors to HTTP responses
  - Consistent error format

#### Router (`/interface/http/router.go`)
- HTTP route setup and configuration

### Shared Layer (`/shared`)
Cross-cutting concerns

#### Dependency Injection (`/shared/di/`)
- **[container.go](shared/di/container.go)** - DI container
  - Initializes all repositories
  - Creates domain services
  - Wires command/query handlers
  - Manages component lifetime

## 🧪 Tests

### Unit Tests (`/tests/unit/`)
- **[domain_test.go](tests/unit/domain_test.go)** - Domain model tests
  - Task creation tests
  - Status transition tests
  - Deadline validation tests
  - Value object tests
  - Aggregate tests

### Integration Tests (`/tests/integration/`)
- **[command_test.go](tests/integration/command_test.go)** - Command flow tests
  - CreateTaskCommandFlow test
  - AssignTaskCommandFlow test
  - UpdateTaskStatusCommandFlow test

## 🚀 Executable Files

### Application Entry Point
- **[main.go](main.go)** - Application entry point
  - Initializes DI container
  - Sets up HTTP router
  - Starts HTTP server

### Examples
- **[examples/usage_example.go](examples/usage_example.go)** - Complete usage example
  - Creating users
  - Creating workflow
  - Creating project
  - Creating task via command
  - Retrieving task
  - Updating task status
  - Listing tasks
  - Domain event demonstration

## 📋 Configuration Files

- **[go.mod](go.mod)** - Go module definition
- **[Makefile](Makefile)** - Build automation
  - `make build` - Build executable
  - `make run` - Run application
  - `make test` - Run tests
  - `make test-unit` - Run unit tests only
  - `make test-int` - Run integration tests only
  - `make lint` - Run linter
  - `make fmt` - Format code
  - `make example` - Run example
  - `make clean` - Clean build artifacts

- **[.gitignore](.gitignore)** - Git ignore rules

## 📊 Project Statistics

- **Total Files**: 40+
- **Go Source Files**: 37
- **Documentation Files**: 7
- **Configuration Files**: 3
- **Total Lines of Code**: 5000+
- **Documentation Lines**: 3000+

## 🎯 Quick Navigation by Task

### I want to...

#### **Learn DDD principles**
→ Read [ARCHITECTURE.md](ARCHITECTURE.md) - Detailed explanation of all DDD patterns

#### **Get started quickly**
→ Read [QUICKSTART.md](QUICKSTART.md) then run `make run`

#### **Understand the code structure**
→ Explore [PROJECT_SUMMARY.md](PROJECT_SUMMARY.md) and this INDEX file

#### **Run the application**
→ Follow [QUICKSTART.md](QUICKSTART.md) - "Running the Application"

#### **Run tests**
→ Execute `make test` or see [TESTING.md](TESTING.md)

#### **Add a new feature**
→ Read [ARCHITECTURE.md](ARCHITECTURE.md) - "Extensibility" section

#### **Deploy to production**
→ See [DEPLOYMENT.md](DEPLOYMENT.md) and [DATABASE.md](DATABASE.md)

#### **Switch to a real database**
→ See [DATABASE.md](DATABASE.md) - "PostgreSQL Implementation Example"

#### **Understand command flow**
→ Look at [CreateTaskCommandFlow](application/command/create_task.go) and [ARCHITECTURE.md](ARCHITECTURE.md) - "Data Flow"

#### **Write tests**
→ See [TESTING.md](TESTING.md) and [tests/](tests/) directory

#### **Understand domain models**
→ Study [domain/aggregate/](domain/aggregate/) and [domain/value/](domain/value/)

#### **Learn HTTP API design**
→ See [interface/http/handler/task_handler.go](interface/http/handler/task_handler.go)

## 📚 Reading Order (Recommended)

For a complete understanding, read in this order:

1. Start: [README.md](README.md) - Overview
2. Quick Start: [QUICKSTART.md](QUICKSTART.md) - Get it running
3. Learn: [PROJECT_SUMMARY.md](PROJECT_SUMMARY.md) - What's included
4. Deep Dive: [ARCHITECTURE.md](ARCHITECTURE.md) - All the details
5. Production: [DATABASE.md](DATABASE.md) + [DEPLOYMENT.md](DEPLOYMENT.md)
6. Testing: [TESTING.md](TESTING.md) - Quality assurance
7. Code Review: Explore files in this order:
   - Domain: `/domain/value/` → `/domain/entity/` → `/domain/aggregate/`
   - Services: `/domain/service/`
   - Events: `/domain/event/`
   - Application: `/application/command/` → `/application/query/`
   - Infrastructure: `/infrastructure/repository/` → `/infrastructure/event/`
   - Interface: `/interface/http/`
   - Shared: `/shared/di/`

## 🔑 Key Concepts Quick Reference

**Aggregate**: Task, Project, User, Workflow (root objects)

**Entity**: Comment, Assignment (objects with identity within aggregates)

**Value Object**: TaskStatus, Priority, Deadline, IDs (immutable objects)

**Domain Service**: TaskAssignmentService, StatusTransitionService, DeadlineEnforcementService

**Domain Event**: TaskCreatedEvent, TaskAssignedEvent, TaskStatusChangedEvent, etc.

**Command**: CreateTaskCommand, AssignTaskCommand, UpdateTaskStatusCommand

**Query**: GetTaskQuery, ListTasksByProjectQuery

**Repository**: Persists and retrieves aggregates

**DI Container**: Central location for component creation and wiring

## 💡 Tips for Navigation

- Use IDE search (Ctrl+Shift+F / Cmd+Shift+F) to find specific classes/methods
- Start with `main.go` to understand application startup
- Follow the DI container to see how components are wired
- Domain layer files (`/domain/`) are independent - start there
- Application layer (`/application/`) coordinates domain objects
- Infrastructure (`/infrastructure/`) is replaceable
- Interface layer (`/interface/`) exposes to HTTP

## 🆘 Need Help?

1. **Confused about architecture?** → [ARCHITECTURE.md](ARCHITECTURE.md)
2. **Don't know how to run it?** → [QUICKSTART.md](QUICKSTART.md)
3. **Want to understand a specific file?** → Look at inline code comments
4. **Need production setup?** → [DEPLOYMENT.md](DEPLOYMENT.md) and [DATABASE.md](DATABASE.md)
5. **Want to add features?** → [ARCHITECTURE.md](ARCHITECTURE.md) "Extensibility"
6. **Testing questions?** → [TESTING.md](TESTING.md)

---

**Happy learning! This index should help you navigate the entire project.** 🎉