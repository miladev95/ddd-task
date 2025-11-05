# Domain-Driven Design Architecture - Task Management System

## Overview

This is a comprehensive Domain-Driven Design (DDD) implementation in Go for a task management system. The architecture is organized into four distinct layers, each with clear responsibilities and well-defined interfaces.

## Architecture Layers

### 1. **Domain Layer** (`/domain`)

The core of the application containing all business logic and rules.

#### Components:

**Value Objects** (`/domain/value/`)
- **Identifiers**: `TaskID`, `ProjectID`, `UserID`, `WorkflowID`
  - Immutable, unique identifiers for aggregates
  - Implemented with validation and equality methods
  
- **Task Status**: Represents task states with valid transitions
  - States: BACKLOG, TO_DO, IN_PROGRESS, IN_REVIEW, COMPLETED, CANCELLED
  - Enforces state machine transitions at value level
  
- **Priority**: Represents task priority levels
  - Levels: LOW, MEDIUM, HIGH, CRITICAL
  - Numeric representation for comparison
  
- **Deadline**: Represents task deadlines
  - Enforces future dates only
  - Provides overdue detection and duration calculations

**Entities** (`/domain/entity/`)
- **Comment**: Belongs to Task, mutable within aggregate
  - Tracks author, content, creation/update times
  
- **Assignment**: Represents task-to-user assignment
  - Tracks assignee, assignment time, and assigner

**Aggregates** (`/domain/aggregate/`)
- **Task**: Root aggregate
  - Owns Comments as part entities
  - Maintains consistency of task state
  - Raises domain events for important state changes
  
- **Project**: Root aggregate
  - Manages collection of task references
  - Enforces project-level business rules
  - Supports archiving
  
- **User**: Root aggregate
  - Manages user information and preferences
  - Tracks activity (login times)
  - Supports activation/deactivation
  
- **Workflow**: Root aggregate
  - Defines available statuses for tasks
  - Supports custom workflow definitions
  - Can be activated/deactivated

**Domain Services** (`/domain/service/`)
- **TaskAssignmentService**: Handles task assignment logic
  - Validates assignees exist
  - Manages reassignment
  - Checks assignment capacity
  
- **StatusTransitionService**: Manages valid status transitions
  - Enforces state machine rules
  - Validates transition prerequisites
  - Returns valid next states
  
- **DeadlineEnforcementService**: Validates and enforces deadlines
  - Validates deadline reasonableness
  - Detects overdue tasks
  - Notifies about upcoming deadlines
  
- **NotificationService**: Interface for notifications
  - Abstracted to allow different implementations
  - Triggered by domain events

**Domain Events** (`/domain/event/`)
- **TaskCreatedEvent**: Raised when a task is created
- **TaskAssignedEvent**: Raised when a task is assigned
- **TaskStatusChangedEvent**: Raised when status changes
- **TaskDeadlineSetEvent**: Raised when deadline is set
- **TaskOverdueEvent**: Raised when task becomes overdue
- **TaskCompletedEvent**: Raised when task is completed
- **TaskDeletedEvent**: Raised when task is deleted

**Repository Interfaces** (`/domain/repository.go`)
- Define contracts for persistence without implementation details
- Separate from infrastructure concerns
- Support dependency injection and testing

### 2. **Application Layer** (`/application`)

Coordinates domain objects and implements use cases through commands and queries.

#### Components:

**Commands** (`/application/command/`)
- **CreateTaskCommand**: Creates a new task
  - Validates project exists
  - Assigns task if assignee provided
  - Sets deadline if provided
  - Publishes events
  
- **AssignTaskCommand**: Assigns task to user
  - Validates users exist
  - Updates assignment entity
  - Publishes assignment event
  
- **UpdateTaskStatusCommand**: Changes task status
  - Validates status transition
  - Checks prerequisites
  - Publishes status change event

**Queries** (`/application/query/`)
- **GetTaskQuery**: Retrieves a single task
  - Converts aggregate to DTO
  
- **ListTasksByProjectQuery**: Lists tasks in a project
  - Optional status filter
  - Returns multiple DTOs

**DTOs** (`/application/dto/`)
- **TaskDTO**: Read-only representation for responses
- **ProjectDTO**: Project data for responses
- **Request/Response Objects**: For HTTP layer

### 3. **Infrastructure Layer** (`/infrastructure`)

Implements technical concerns and persistence mechanisms.

#### Components:

**Repositories** (`/infrastructure/repository/`)
- **InMemoryTaskRepository**: In-memory Task persistence
- **InMemoryProjectRepository**: In-memory Project persistence
- **InMemoryUserRepository**: In-memory User persistence
- **InMemoryWorkflowRepository**: In-memory Workflow persistence

All repositories implement domain interfaces and support:
- Create, Read, Update, Delete operations
- Complex queries (by project, by assignee, by status)
- Full in-memory storage for demo/testing

**Event Publishing** (`/infrastructure/event/`)
- **SimpleEventPublisher**: In-memory event publishing
  - Supports subscriber registration
  - Publishes to all interested subscribers
  
- **SimpleNotificationService**: Basic notification implementation
  - Prints notifications to console
  - Can be replaced with email/SMS/push service

### 4. **Interface Layer** (`/interface`)

Exposes application through HTTP API.

#### Components:

**HTTP Handlers** (`/interface/http/handler/`)
- **TaskHandler**: Handles task-related HTTP requests
  - POST /api/tasks: Create task
  - GET /api/tasks: List tasks
  - POST /api/tasks/assign: Assign task
  - PUT /api/tasks/status: Update status

**Middleware** (`/interface/http/middleware/`)
- **ErrorHandler**: Converts domain errors to HTTP responses
- Provides consistent error format

**Router** (`/interface/http/router.go`)
- Sets up HTTP routes
- Integrates handlers with dependency injection

### 5. **Shared/Cross-cutting** (`/shared`)

Common concerns used across layers.

#### Components:

**Dependency Injection** (`/shared/di/`)
- **Container**: Centralized component creation
  - Initializes all repositories
  - Creates domain services
  - Wires command/query handlers
  - Manages lifetime of dependencies

## Data Flow

### Command Flow: Creating a Task

```
HTTP Request
    ↓
TaskHandler.CreateTask()
    ↓
CreateTaskCommand
    ↓
CreateTaskCommandHandler.Handle()
    ↓
Domain Services:
  - Validate project exists
  - Validate priority
  - Validate assignee (if provided)
  - Validate deadline (if provided)
    ↓
Task Aggregate.NewTask()
    ↓
Task.Assign() (if assignee)
    ↓
Task.SetDeadline() (if deadline)
    ↓
Domain Events raised:
  - TaskCreatedEvent
  - TaskAssignedEvent (if assigned)
  - TaskDeadlineSetEvent (if deadline)
    ↓
Repository.Save(task)
    ↓
EventPublisher.Publish(events)
    ↓
HTTP Response (201 Created)
```

### Query Flow: Listing Tasks by Project

```
HTTP Request
    ↓
TaskHandler.ListTasksByProject()
    ↓
ListTasksByProjectQuery
    ↓
ListTasksByProjectQueryHandler.Handle()
    ↓
TaskRepository.GetByProjectID()
    ↓
Convert aggregates to DTOs
    ↓
HTTP Response (200 OK)
```

## Key Design Patterns

### 1. **Aggregate Pattern**
- Task, Project, User, Workflow are aggregate roots
- Each owns their entities and value objects
- Encapsulate all business rules
- Raise domain events for external communication

### 2. **Repository Pattern**
- Abstract persistence behind interfaces
- Enable dependency injection
- Support multiple implementations
- Simplify testing

### 3. **Value Object Pattern**
- TaskStatus, Priority, Deadline, Identifiers are immutable
- Encapsulate validation logic
- Provide domain-specific methods

### 4. **Domain Events**
- Communicate important business events
- Decouple aggregates
- Support event-driven architecture
- Enable audit trails and event sourcing

### 5. **Command/Query Separation**
- Commands modify state (CreateTaskCommand)
- Queries read state (GetTaskQuery, ListTasksByProjectQuery)
- Clear separation of concerns
- Different optimization strategies

### 6. **Dependency Injection**
- All dependencies passed via constructors
- Enables loose coupling
- Facilitates testing with mock implementations
- Centralized configuration

## Business Rules

### Task Status Transitions
```
BACKLOG → TO_DO, CANCELLED
TO_DO → IN_PROGRESS, BACKLOG, CANCELLED
IN_PROGRESS → IN_REVIEW, TO_DO, CANCELLED
IN_REVIEW → COMPLETED, IN_PROGRESS, CANCELLED
COMPLETED → (no transitions)
CANCELLED → (no transitions)
```

### Task Assignment Rules
- Task must have valid assignee
- Assignee must be an active user
- Assignment tracked with timestamp and assigner

### Deadline Rules
- Deadline must be in future (for new tasks)
- Deadline cannot be more than 5 years away
- Task must be assigned before completion with deadline
- Overdue detection triggers notifications

## Testing Strategy

### Unit Tests (`/tests/unit/`)
- Test individual aggregates
- Test value object validation
- Test domain services logic
- Mock external dependencies

### Integration Tests (`/tests/integration/`)
- Test complete command flows
- Test multiple aggregates interacting
- Test repository operations
- Verify event publishing

### Test Patterns
- Setup test data
- Execute command/query
- Assert results and side effects

## Extensibility

### Adding New Aggregate
1. Create aggregate in `/domain/aggregate/`
2. Define value objects in `/domain/value/`
3. Define domain events in `/domain/event/`
4. Create repository interface in `/domain/repository.go`
5. Implement repository in `/infrastructure/repository/`
6. Create commands in `/application/command/`
7. Create queries in `/application/query/`
8. Add handlers to DI container
9. Create HTTP handlers in `/interface/http/handler/`
10. Register routes in `/interface/http/router.go`

### Switching Persistence
1. Implement domain repository interface with new database (PostgreSQL, MongoDB, etc.)
2. Update DI container to use new repository
3. No changes needed to domain or application layers

### Adding Event Subscribers
1. Create subscriber implementation
2. Register with EventPublisher in DI container
3. Subscriber receives events when published

## Running the Application

### Build
```bash
make build
```

### Run
```bash
make run
```

### Test
```bash
make test          # All tests
make test-unit     # Unit tests only
make test-int      # Integration tests only
```

### Example
```bash
make example
```

## HTTP API Endpoints

### Create Task
```
POST /api/tasks
Content-Type: application/json

{
  "project_id": "project-uuid",
  "title": "Implement feature X",
  "description": "Detailed description",
  "priority": "HIGH",
  "assignee_id": "user-uuid",
  "deadline": "2024-12-31T23:59:59Z"
}
```

### Get Task
```
GET /api/tasks/get?id=task-uuid
```

### List Tasks by Project
```
GET /api/tasks?project_id=project-uuid&status=IN_PROGRESS
```

### Assign Task
```
POST /api/tasks/assign?id=task-uuid
Content-Type: application/json

{
  "assignee_id": "user-uuid"
}
```

### Update Task Status
```
PUT /api/tasks/status?id=task-uuid
Content-Type: application/json

{
  "status": "IN_PROGRESS"
}
```

## Future Enhancements

1. **Event Sourcing**: Store events instead of current state
2. **CQRS**: Separate read and write models
3. **Workflow Engine**: Dynamic workflow definitions
4. **Notifications**: Real email/SMS/push notifications
5. **Permissions**: Fine-grained access control
6. **Audit Trail**: Complete audit logging
7. **Search**: Full-text search capabilities
8. **Webhooks**: External system integration
9. **Real Database**: PostgreSQL/MongoDB persistence
10. **Caching**: Redis for performance optimization