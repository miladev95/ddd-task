# Task Management API - Setup Guide

## Quick Start

This guide explains how to use the Task Management API with the Postman collection.

## Prerequisites

1. **Go 1.19+** - Build and run the application
2. **Postman** - Test the API endpoints
3. **Git** - Clone the repository

## Step 1: Start the Server

```bash
# Build the application
go build -o bin/task-management ./main.go

# Run the server
./bin/task-management
# or
go run main.go
```

The server will start on `http://localhost:8080`

## Step 2: Import Postman Collection

1. Open **Postman**
2. Click **Import** (top-left corner)
3. Select the file: `Task_Management_API.postman_collection.json`
4. Click **Import**

All API endpoints will be organized in folders.

## Step 3: API Setup Workflow

> ⚠️ **Important**: Follow these steps in order. Each step depends on the previous one.

### Step 3.1: Create Users

Users must exist before creating projects and assigning tasks.

**Endpoint**: `POST /api/users`

1. In Postman, go to **Setup - Step 1: Users** → **Create User (Alice)**
2. Click **Send**
3. Copy the `user_id` from the response
4. Repeat for **Create User (Bob)**
5. Save both user IDs for later use

**Example Response**:
```json
{
  "user_id": "550e8400-e29b-41d4-a716-446655440001",
  "email": "alice@example.com",
  "first_name": "Alice",
  "last_name": "Johnson",
  "full_name": "Alice Johnson",
  "message": "User created successfully"
}
```

### Step 3.2: Create Workflow

Workflows define the available task statuses (states).

**Endpoint**: `POST /api/workflows`

1. Go to **Setup - Step 2: Workflows** → **Create Workflow (Default)**
2. Click **Send**
3. Copy the `workflow_id` from the response

**Default Workflow Statuses**:
- Backlog (order 1)
- To Do (order 2)
- In Progress (order 3)
- In Review (order 4)
- Completed (order 5, final)

**Example Response**:
```json
{
  "workflow_id": "550e8400-e29b-41d4-a716-446655440002",
  "name": "Default Workflow",
  "description": "Standard project workflow with 5 statuses",
  "message": "Workflow created successfully"
}
```

### Step 3.3: Create Project

Projects contain tasks.

**Endpoint**: `POST /api/projects`

1. Go to **Setup - Step 3: Projects** → **Create Project**
2. Update the request body with:
   - `owner_id`: Use Alice's user ID from Step 3.1
   - `workflow_id`: Use the workflow ID from Step 3.2
3. Click **Send**
4. Copy the `project_id` from the response

**Example Request**:
```json
{
  "name": "Web Application Development",
  "description": "Build a modern web application with authentication and API",
  "owner_id": "550e8400-e29b-41d4-a716-446655440001",
  "workflow_id": "550e8400-e29b-41d4-a716-446655440002"
}
```

**Example Response**:
```json
{
  "project_id": "550e8400-e29b-41d4-a716-446655440003",
  "name": "Web Application Development",
  "message": "Project created successfully"
}
```

### Step 3.4: Create Tasks

Now you can create tasks in the project!

**Endpoint**: `POST /api/tasks`

1. Go to **Task Management** → **Create Task**
2. Update the request body with:
   - `project_id`: Use the project ID from Step 3.3
   - `assignee_id`: Use Bob's user ID from Step 3.1
   - Set `X-User-ID` header to Alice's user ID
3. Click **Send**
4. Copy the `task_id` from the response

**Example Request**:
```json
{
  "project_id": "550e8400-e29b-41d4-a716-446655440003",
  "title": "Implement User Authentication",
  "description": "Add JWT-based authentication to the API",
  "priority": "HIGH",
  "assignee_id": "550e8400-e29b-41d4-a716-446655440002",
  "deadline": "2025-12-31T23:59:59Z"
}
```

**Priority Values**: `LOW`, `MEDIUM`, `HIGH`, `CRITICAL`

**Example Response**:
```json
{
  "task_id": "550e8400-e29b-41d4-a716-446655440004",
  "message": "Task created successfully"
}
```

## Step 4: Manage Tasks

### Get Task Details

**Endpoint**: `GET /api/tasks/get?id={task_id}`

Retrieve all details of a task including status, assignee, deadline, and comments.

### List Tasks by Project

**Endpoint**: `GET /api/tasks?project_id={project_id}&status={status}`

- `project_id` (required): The project ID
- `status` (optional): Filter by status (BACKLOG, TO_DO, IN_PROGRESS, IN_REVIEW, COMPLETED, CANCELLED)

### Assign Task to User

**Endpoint**: `POST /api/tasks/assign?id={task_id}`

Assign a task to a user. This is **required** before the task can move to `IN_PROGRESS`.

**Request Body**:
```json
{
  "assignee_id": "550e8400-e29b-41d4-a716-446655440002"
}
```

### Update Task Status

**Endpoint**: `PUT /api/tasks/status?id={task_id}`

Change the task status following the workflow state machine.

**Valid Status Transitions**:
```
BACKLOG → TO_DO → IN_PROGRESS → IN_REVIEW → COMPLETED
              ↓
            (or skip to)
                ↓
            CANCELLED
```

**Business Rules**:
- Task must be **assigned** before moving to `IN_PROGRESS`
- Task must have a **deadline** before moving to `COMPLETED`
- Once in a final state (COMPLETED, CANCELLED), no further transitions allowed

**Request Body**:
```json
{
  "status": "IN_PROGRESS"
}
```

## API Endpoints Summary

### Users
| Method | Endpoint | Purpose |
|--------|----------|---------|
| POST | `/api/users` | Create a new user |
| GET | `/api/users/get?id={user_id}` | Get user details |

### Workflows
| Method | Endpoint | Purpose |
|--------|----------|---------|
| POST | `/api/workflows` | Create a new workflow |
| GET | `/api/workflows/get?id={workflow_id}` | Get workflow details |

### Projects
| Method | Endpoint | Purpose |
|--------|----------|---------|
| POST | `/api/projects` | Create a new project |
| GET | `/api/projects/get?id={project_id}` | Get project details |

### Tasks
| Method | Endpoint | Purpose |
|--------|----------|---------|
| POST | `/api/tasks` | Create a new task |
| GET | `/api/tasks/get?id={task_id}` | Get task details |
| GET | `/api/tasks?project_id={project_id}&status={status}` | List project tasks |
| POST | `/api/tasks/assign?id={task_id}` | Assign task to user |
| PUT | `/api/tasks/status?id={task_id}` | Update task status |

### Health
| Method | Endpoint | Purpose |
|--------|----------|---------|
| GET | `/health` | Health check |

## Error Handling

All API errors follow this format:

```json
{
  "code": 400,
  "message": "Error description"
}
```

**Common Error Codes**:
- `400 Bad Request` - Invalid request data
- `404 Not Found` - Resource not found (project, user, task, etc.)
- `500 Internal Server Error` - Server error

## Example: Complete Workflow

```bash
# 1. Create User Alice
POST /api/users
{
  "email": "alice@example.com",
  "first_name": "Alice",
  "last_name": "Johnson"
}
# Response: user_id = user-1

# 2. Create User Bob
POST /api/users
{
  "email": "bob@example.com",
  "first_name": "Bob",
  "last_name": "Smith"
}
# Response: user_id = user-2

# 3. Create Workflow
POST /api/workflows
{
  "name": "Default Workflow",
  "statuses": [...]
}
# Response: workflow_id = workflow-1

# 4. Create Project
POST /api/projects
{
  "name": "Web App",
  "owner_id": "user-1",
  "workflow_id": "workflow-1"
}
# Response: project_id = proj-1

# 5. Create Task
POST /api/tasks
{
  "project_id": "proj-1",
  "title": "Implement Auth",
  "priority": "HIGH",
  "assignee_id": "user-2"
}
# Response: task_id = task-1

# 6. Update Task Status
PUT /api/tasks/status?id=task-1
{
  "status": "IN_PROGRESS"
}
# Success response

# 7. Get Task Details
GET /api/tasks/get?id=task-1
# Response: Full task details with status, assignee, etc.
```

## Troubleshooting

### "project not found" Error
- Make sure you created a project first (Step 3.3)
- Verify the `project_id` in the task creation request matches the created project

### "user not found" Error
- Create users before creating projects or assigning tasks
- Verify the user IDs are correct

### "Task must be assigned" Error
- Assign the task using `/api/tasks/assign` before moving to `IN_PROGRESS`

### "Invalid status transition" Error
- Follow the workflow state machine: BACKLOG → TO_DO → IN_PROGRESS → IN_REVIEW → COMPLETED
- You cannot skip states or go backwards

## Advanced: Modifying the Collection

To save your actual IDs for easy access:

1. Create a new **Environment** in Postman
2. Add variables:
   - `alice_user_id` = (copy from Create User response)
   - `bob_user_id` = (copy from Create User response)
   - `workflow_id` = (copy from Create Workflow response)
   - `project_id` = (copy from Create Project response)
   - `task_id` = (copy from Create Task response)

3. Update request bodies to use `{{alice_user_id}}` instead of hardcoded IDs

## Next Steps

1. **Explore the domain model**: Check `domain/aggregate/task.go` to understand task business logic
2. **Add more endpoints**: Follow the same pattern to add update/delete operations
3. **Database integration**: Replace in-memory repositories with PostgreSQL (see DATABASE.md)
4. **Authentication**: Add JWT middleware to protect endpoints
5. **Event processing**: Implement async event handlers

## Architecture Notes

This API follows **Domain-Driven Design (DDD)** principles:

- **Domain Layer**: Core business logic (validation, state machine)
- **Application Layer**: Use case orchestration (command handlers)
- **Infrastructure Layer**: Technical implementations (repositories, event publishing)
- **Interface Layer**: HTTP API exposure

All business rules are enforced in the domain layer and cannot be bypassed.

## Need Help?

- Read the **ARCHITECTURE.md** for design details
- Check the **examples/usage_example.go** for programmatic usage
- Review **TESTING.md** for test structure