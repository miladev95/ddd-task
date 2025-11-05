# Deployment Guide

## Overview

This guide covers deploying the Task Management System to various environments and platforms.

## Development Environment

### Local Development

```bash
# 1. Install dependencies
go mod download

# 2. Run the application
go run main.go

# 3. Run tests
go test ./...

# 4. Run example
go run examples/usage_example.go
```

The API will be available at `http://localhost:8080`

### Docker Development

```dockerfile
# Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/app .

EXPOSE 8080

CMD ["./app"]
```

Build and run:
```bash
docker build -t task-management:latest .
docker run -p 8080:8080 task-management:latest
```

## Staging Environment

### Configuration Management

Create `.env.staging`:
```
DB_HOST=postgres-staging
DB_PORT=5432
DB_USER=app_user
DB_PASSWORD=secure_password
DB_NAME=task_management
LOG_LEVEL=info
API_PORT=8080
```

### Docker Compose for Staging

```yaml
version: '3.8'

services:
  api:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=app_user
      - DB_PASSWORD=app_password
      - DB_NAME=task_management
    depends_on:
      - postgres
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3

  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: app_user
      POSTGRES_PASSWORD: app_password
      POSTGRES_DB: task_management
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U app_user"]
      interval: 10s
      timeout: 5s
      retries: 5

volumes:
  postgres_data:
```

Deploy to staging:
```bash
docker-compose -f docker-compose.staging.yml up -d
```

## Production Environment

### Binary Build

```bash
# Build optimized binary
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/task-management main.go

# Create release artifact
tar -czf task-management-linux-amd64.tar.gz bin/task-management
```

### Production Configuration

`.env.production`:
```
# Database
DB_HOST=prod-postgres.example.com
DB_PORT=5432
DB_USER=prod_app
DB_PASSWORD=${DB_PASSWORD}  # From secrets management
DB_NAME=task_management_prod
DB_POOL_SIZE=25

# API
API_PORT=8080
API_TIMEOUT=30s

# Logging
LOG_LEVEL=warn
LOG_FORMAT=json

# Performance
CACHE_ENABLED=true
CACHE_TTL=3600

# Security
ENABLE_HTTPS=true
TLS_CERT_PATH=/etc/certs/server.crt
TLS_KEY_PATH=/etc/certs/server.key
```

### Kubernetes Deployment

#### namespace.yaml
```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: task-management
```

#### configmap.yaml
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
  namespace: task-management
data:
  LOG_LEVEL: "info"
  LOG_FORMAT: "json"
  CACHE_ENABLED: "true"
```

#### secret.yaml
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: app-secrets
  namespace: task-management
type: Opaque
stringData:
  DB_PASSWORD: "your-secure-password"
  API_KEY: "your-api-key"
```

#### deployment.yaml
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: task-management-api
  namespace: task-management
spec:
  replicas: 3
  selector:
    matchLabels:
      app: task-management-api
  template:
    metadata:
      labels:
        app: task-management-api
    spec:
      containers:
      - name: api
        image: your-registry/task-management:1.0.0
        ports:
        - containerPort: 8080
        env:
        - name: DB_HOST
          value: "postgres.task-management.svc.cluster.local"
        - name: DB_PORT
          value: "5432"
        - name: DB_USER
          value: "app_user"
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: app-secrets
              key: DB_PASSWORD
        - name: DB_NAME
          value: "task_management"
        - name: LOG_LEVEL
          valueFrom:
            configMapKeyRef:
              name: app-config
              key: LOG_LEVEL
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "512Mi"
            cpu: "500m"
```

#### service.yaml
```yaml
apiVersion: v1
kind: Service
metadata:
  name: task-management-api
  namespace: task-management
spec:
  selector:
    app: task-management-api
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
  type: LoadBalancer
```

Deploy to Kubernetes:
```bash
kubectl apply -f namespace.yaml
kubectl apply -f configmap.yaml
kubectl apply -f secret.yaml
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml
```

### AWS EC2 Deployment

#### User Data Script
```bash
#!/bin/bash
set -e

# Update system
sudo apt-get update
sudo apt-get install -y postgresql-client

# Download and install Go (if needed)
wget https://go.dev/dl/go1.21.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.linux-amd64.tar.gz

# Clone application
cd /opt
sudo git clone https://github.com/example/task-management.git
cd task-management

# Build application
go build -o bin/task-management main.go

# Create systemd service
sudo tee /etc/systemd/system/task-management.service > /dev/null <<EOF
[Unit]
Description=Task Management API
After=network.target

[Service]
Type=simple
User=app
WorkingDirectory=/opt/task-management
EnvironmentFile=/opt/task-management/.env.production
ExecStart=/opt/task-management/bin/task-management
Restart=on-failure
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF

# Enable and start service
sudo systemctl daemon-reload
sudo systemctl enable task-management
sudo systemctl start task-management
```

### GCP App Engine Deployment

#### app.yaml
```yaml
runtime: go121

env: standard

handlers:
- url: /.*
  script: auto

env_variables:
  DB_HOST: "cloudsql-instance"
  DB_PORT: "5432"
  LOG_LEVEL: "info"
```

Deploy:
```bash
gcloud app deploy
```

## Monitoring & Logging

### Application Metrics

Add Prometheus support:
```go
import "github.com/prometheus/client_golang/prometheus"

var (
    taskCreatedCounter = prometheus.NewCounter(prometheus.CounterOpts{
        Name: "tasks_created_total",
        Help: "Total number of tasks created",
    })
)

// Use in handlers
taskCreatedCounter.Inc()
```

### Logging Configuration

```go
import "log/slog"

logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

logger.Info("Task created", 
    "task_id", taskID,
    "project_id", projectID,
)
```

### Health Checks

Kubernetes health check endpoint `/health`:
```bash
curl http://localhost:8080/health
# {"status":"healthy"}
```

## CI/CD Pipeline

### GitHub Actions

`.github/workflows/deploy.yml`:
```yaml
name: Deploy

on:
  push:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_PASSWORD: postgres
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5

    steps:
    - uses: actions/checkout@v2
    
    - name: Set up Go
      uses: actions/setup-go@v2
      with:
        go-version: 1.21
    
    - name: Test
      run: go test ./...
    
    - name: Build
      run: go build -o bin/task-management main.go
    
    - name: Push to registry
      run: |
        docker build -t your-registry/task-management:${{ github.sha }} .
        docker push your-registry/task-management:${{ github.sha }}
    
    - name: Deploy to Kubernetes
      run: |
        kubectl set image deployment/task-management-api \
          api=your-registry/task-management:${{ github.sha }} \
          --record
```

## Backup Strategy

### Database Backups

```bash
# Daily backup to S3
0 2 * * * pg_dump -h $DB_HOST -U $DB_USER $DB_NAME | \
    gzip | \
    aws s3 cp - s3://backups/task-management/$(date +\%Y-\%m-\%d).sql.gz
```

### Restore from Backup

```bash
aws s3 cp s3://backups/task-management/2024-01-01.sql.gz - | \
    gunzip | \
    psql -h $DB_HOST -U $DB_USER $DB_NAME
```

## Rolling Updates

### Zero-Downtime Deployment

```bash
# Update deployment
kubectl set image deployment/task-management-api \
  api=your-registry/task-management:new-version

# Monitor rollout
kubectl rollout status deployment/task-management-api

# Rollback if needed
kubectl rollout undo deployment/task-management-api
```

## Performance Optimization

### Caching Layer

Add Redis:
```yaml
services:
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
```

### Database Connection Pooling

Update connection pool settings:
```go
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(5)
db.SetConnMaxLifetime(time.Hour)
```

## Security Checklist

- [ ] Use HTTPS with valid certificates
- [ ] Implement authentication/authorization
- [ ] Validate all inputs
- [ ] Use secrets management (HashiCorp Vault, AWS Secrets Manager)
- [ ] Enable CORS if needed
- [ ] Rate limiting
- [ ] SQL injection prevention (parameterized queries)
- [ ] CSRF protection
- [ ] Security headers
- [ ] Logging and monitoring
- [ ] Regular security updates

## Troubleshooting

### Application won't start

```bash
# Check logs
docker logs <container_id>

# Verify environment variables
printenv | grep DB_

# Test database connection
psql -h $DB_HOST -U $DB_USER -d $DB_NAME
```

### High memory usage

- Check for goroutine leaks
- Adjust connection pool size
- Enable profiling

### Slow queries

- Check database indexes
- Use query analyzer: `EXPLAIN ANALYZE`
- Consider caching

## Rollback Procedure

```bash
# Kubernetes rollback
kubectl rollout undo deployment/task-management-api

# Docker Compose rollback
docker-compose down
docker-compose up -d  # With previous image
```

## Maintenance

### Regular Tasks

- Weekly: Monitor logs and metrics
- Monthly: Review and optimize slow queries
- Quarterly: Security updates and dependencies
- Annually: Disaster recovery drill

### Update Procedures

```bash
# Test in staging first
# Review changelog
# Plan maintenance window
# Take backup
# Deploy update
# Verify health
# Monitor for issues
```