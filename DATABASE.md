# Database Implementation Guide

## Current State

The current implementation uses **in-memory repositories** for demonstration and testing purposes. These repositories:
- Store data in Go maps
- Use RWMutex for thread-safe access
- Are suitable for testing and development
- Do not persist data between restarts

## Production Database Setup

### PostgreSQL Implementation Example

Here's how to implement PostgreSQL repositories following the same pattern:

#### 1. Create Database Schema

```sql
-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    preferences JSONB DEFAULT '{}'::jsonb
);

-- Workflows table
CREATE TABLE workflows (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    statuses JSONB NOT NULL,
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Projects table
CREATE TABLE projects (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    owner_id UUID NOT NULL REFERENCES users(id),
    workflow_id UUID NOT NULL REFERENCES workflows(id),
    archived BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tasks table
CREATE TABLE tasks (
    id UUID PRIMARY KEY,
    project_id UUID NOT NULL REFERENCES projects(id),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) NOT NULL,
    priority VARCHAR(50) NOT NULL,
    deadline TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by UUID NOT NULL REFERENCES users(id),
    INDEX idx_project_id (project_id),
    INDEX idx_status (status),
    INDEX idx_created_by (created_by)
);

-- Task assignments table
CREATE TABLE assignments (
    task_id UUID NOT NULL REFERENCES tasks(id),
    assignee_id UUID NOT NULL REFERENCES users(id),
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    assigned_by UUID NOT NULL REFERENCES users(id),
    PRIMARY KEY (task_id, assignee_id)
);

-- Comments table
CREATE TABLE comments (
    id UUID PRIMARY KEY,
    task_id UUID NOT NULL REFERENCES tasks(id),
    author_id UUID NOT NULL REFERENCES users(id),
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_task_id (task_id)
);

-- Domain events table (for event sourcing)
CREATE TABLE domain_events (
    id UUID PRIMARY KEY,
    event_type VARCHAR(100) NOT NULL,
    aggregate_id UUID NOT NULL,
    aggregate_type VARCHAR(100) NOT NULL,
    payload JSONB NOT NULL,
    occurred_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_aggregate (aggregate_id, aggregate_type),
    INDEX idx_event_type (event_type)
);
```

#### 2. PostgreSQL Repository Implementation

```go
package repository

import (
	"database/sql"
	"fmt"

	"github.com/example/task-management/domain"
	"github.com/example/task-management/domain/aggregate"
	"github.com/example/task-management/domain/value"
	_ "github.com/lib/pq"
)

// PostgresTaskRepository implements TaskRepository using PostgreSQL
type PostgresTaskRepository struct {
	db *sql.DB
}

// NewPostgresTaskRepository creates a new PostgresTaskRepository
func NewPostgresTaskRepository(db *sql.DB) *PostgresTaskRepository {
	return &PostgresTaskRepository{db: db}
}

// Save persists a task to PostgreSQL
func (r *PostgresTaskRepository) Save(task *aggregate.Task) error {
	if task == nil {
		return fmt.Errorf("task cannot be nil")
	}

	query := `
		INSERT INTO tasks (
			id, project_id, title, description, status, priority,
			deadline, created_at, updated_at, created_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	deadline := task.Deadline()
	var deadlineTime interface{}
	if deadline != nil {
		deadlineTime = deadline.Value()
	}

	_, err := r.db.Exec(
		query,
		task.ID().Value(),
		task.ProjectID().Value(),
		task.Title(),
		task.Description(),
		task.Status().Value(),
		task.Priority().Value(),
		deadlineTime,
		task.CreatedAt(),
		task.UpdatedAt(),
		task.CreatedBy().Value(),
	)

	if err != nil {
		return fmt.Errorf("failed to save task: %w", err)
	}

	// Handle assignment if present
	if task.Assignee() != nil {
		assignQuery := `
			INSERT INTO assignments (task_id, assignee_id, assigned_at, assigned_by)
			VALUES ($1, $2, $3, $4)
		`

		_, err = r.db.Exec(
			assignQuery,
			task.ID().Value(),
			task.Assignee().AssigneeID().Value(),
			task.Assignee().AssignedAt(),
			task.Assignee().AssignedBy().Value(),
		)

		if err != nil {
			return fmt.Errorf("failed to save assignment: %w", err)
		}
	}

	return nil
}

// GetByID retrieves a task from PostgreSQL by ID
func (r *PostgresTaskRepository) GetByID(id value.TaskID) (*aggregate.Task, error) {
	query := `
		SELECT id, project_id, title, description, status, priority,
		       deadline, created_at, updated_at, created_by
		FROM tasks
		WHERE id = $1
	`

	row := r.db.QueryRow(query, id.Value())

	// Scan row and reconstruct task aggregate
	// (implementation omitted for brevity)
	
	return nil, nil // Placeholder
}

// GetByProjectID retrieves all tasks for a project
func (r *PostgresTaskRepository) GetByProjectID(
	projectID value.ProjectID,
) ([]*aggregate.Task, error) {
	query := `
		SELECT id, project_id, title, description, status, priority,
		       deadline, created_at, updated_at, created_by
		FROM tasks
		WHERE project_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, projectID.Value())
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	tasks := make([]*aggregate.Task, 0)
	
	for rows.Next() {
		// Scan and reconstruct task
		// (implementation omitted for brevity)
	}

	return tasks, nil
}

// Update updates an existing task
func (r *PostgresTaskRepository) Update(task *aggregate.Task) error {
	query := `
		UPDATE tasks
		SET title = $1, description = $2, status = $3, priority = $4,
		    deadline = $5, updated_at = $6
		WHERE id = $7
	`

	deadline := task.Deadline()
	var deadlineTime interface{}
	if deadline != nil {
		deadlineTime = deadline.Value()
	}

	result, err := r.db.Exec(
		query,
		task.Title(),
		task.Description(),
		task.Status().Value(),
		task.Priority().Value(),
		deadlineTime,
		task.UpdatedAt(),
		task.ID().Value(),
	)

	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

// Other methods (GetByAssigneeID, GetByStatus, GetAll, Delete, FindByProjectIDAndStatus)
// follow the same pattern...

// Ensure PostgresTaskRepository implements domain.TaskRepository
var _ domain.TaskRepository = (*PostgresTaskRepository)(nil)
```

#### 3. Database Connection Setup

```go
package persistence

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewConnection creates a new database connection
func NewConnection(config DatabaseConfig) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DBName,
		config.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	log.Println("Database connection established")

	return db, nil
}

// RunMigrations runs database migrations
func RunMigrations(db *sql.DB) error {
	// Execute migration SQL
	_, err := db.Exec(schemaSQL)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database migrations completed")
	return nil
}

// schemaSQL contains all table creation DDL
const schemaSQL = `
	-- Tables definitions here
`
```

#### 4. Update DI Container for PostgreSQL

```go
// In shared/di/container.go

// NewContainerWithPostgres creates container with PostgreSQL
func NewContainerWithPostgres(dbConfig persistence.DatabaseConfig) (*Container, error) {
	c := &Container{}

	// Create database connection
	db, err := persistence.NewConnection(dbConfig)
	if err != nil {
		return nil, err
	}

	// Run migrations
	if err := persistence.RunMigrations(db); err != nil {
		return nil, err
	}

	// Initialize PostgreSQL repositories
	c.TaskRepository = repository.NewPostgresTaskRepository(db)
	c.ProjectRepository = repository.NewPostgresProjectRepository(db)
	c.UserRepository = repository.NewPostgresUserRepository(db)
	c.WorkflowRepository = repository.NewPostgresWorkflowRepository(db)

	// Rest of container initialization...
	return c, nil
}
```

#### 5. Update Main to Use PostgreSQL

```go
// In main.go

func main() {
	// Database configuration (from environment variables)
	dbConfig := persistence.DatabaseConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     5432,
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  "disable",
	}

	// Initialize container with PostgreSQL
	container, err := di.NewContainerWithPostgres(dbConfig)
	if err != nil {
		log.Fatalf("Failed to initialize container: %v", err)
	}

	// Rest of application setup...
}
```

## Alternative Databases

### MongoDB Implementation

```go
type MongoTaskRepository struct {
	db *mongo.Client
	coll *mongo.Collection
}

func NewMongoTaskRepository(client *mongo.Client) *MongoTaskRepository {
	return &MongoTaskRepository{
		db: client,
		coll: client.Database("task_management").Collection("tasks"),
	}
}

func (r *MongoTaskRepository) Save(task *aggregate.Task) error {
	doc := bson.M{
		"_id": task.ID().Value(),
		"project_id": task.ProjectID().Value(),
		"title": task.Title(),
		// ... more fields
	}

	_, err := r.coll.InsertOne(context.TODO(), doc)
	return err
}
```

### MySQL Implementation

Similar to PostgreSQL with appropriate SQL syntax changes.

## Migration Strategy

### Using golang-migrate

```bash
# Install
go get -u github.com/golang-migrate/migrate/cmd/migrate

# Create migration
migrate create -ext sql -dir db/migrations -seq create_users_table

# Run migrations
migrate -path db/migrations -database "postgres://..." up

# Rollback
migrate -path db/migrations -database "postgres://..." down
```

## Performance Optimization

### Indexing

```sql
-- Create indexes for common queries
CREATE INDEX idx_tasks_project_id ON tasks(project_id);
CREATE INDEX idx_tasks_created_by ON tasks(created_by);
CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_assignments_assignee_id ON assignments(assignee_id);
```

### Query Optimization

- Use prepared statements to prevent SQL injection
- Implement result caching for read-heavy operations
- Use connection pooling
- Monitor slow queries

### Pagination

```go
func (r *PostgresTaskRepository) GetByProjectID(
	projectID value.ProjectID,
	limit int,
	offset int,
) ([]*aggregate.Task, error) {
	query := `
		SELECT ... FROM tasks
		WHERE project_id = $1
		LIMIT $2 OFFSET $3
	`
	// Execute query with pagination
}
```

## Data Integrity

### Constraints

- Foreign keys for referential integrity
- Unique constraints on emails, IDs
- Check constraints for valid statuses
- NOT NULL constraints on required fields

### Transactions

```go
func (r *PostgresTaskRepository) CreateWithAssignment(
	task *aggregate.Task,
) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// Save task
	// Save assignment
	// Commit or rollback

	return tx.Commit().Error
}
```

## Backup and Recovery

```bash
# PostgreSQL backup
pg_dump -U username -h localhost dbname > backup.sql

# Restore from backup
psql -U username -h localhost dbname < backup.sql
```

## Monitoring

### Health Check Queries

```go
func (r *PostgresTaskRepository) Health() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return r.db.PingContext(ctx)
}
```

## Migration from In-Memory to Database

1. Keep both repositories simultaneously
2. Read from in-memory, write to both
3. Verify data consistency
4. Switch to database-only
5. Archive in-memory data

## Production Checklist

- [ ] Database backups configured
- [ ] Connection pooling optimized
- [ ] Indexes created
- [ ] Query performance tested
- [ ] Transactions properly handled
- [ ] Error handling comprehensive
- [ ] Monitoring and logging in place
- [ ] Data validation at database level
- [ ] Migrations versioned
- [ ] Rollback procedures documented