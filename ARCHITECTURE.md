# QMS Engine Architecture

## Overview

The QMS Engine follows a **Clean Architecture** pattern with clear separation of concerns across multiple layers.

## Architecture Layers

```
┌─────────────────────────────────────────────────────────┐
│                    HTTP Layer (Gin)                      │
│                  handlers/ (Controllers)                 │
└───────────────────────────┬─────────────────────────────┘
                            │
                            ↓
┌─────────────────────────────────────────────────────────┐
│                  Business Logic Layer                    │
│              internal/service/ (Use Cases)               │
└───────────────────────────┬─────────────────────────────┘
                            │
                            ↓
┌─────────────────────────────────────────────────────────┐
│                   Repository Layer                       │
│         internal/repository/ (Data Access)               │
└───────────────────────────┬─────────────────────────────┘
                            │
                            ↓
┌─────────────────────────────────────────────────────────┐
│                    Database Layer                        │
│              MySQL (Master/Slave Setup)                  │
└─────────────────────────────────────────────────────────┘
```

## Project Structure

```
qms-engine/
├── cmd/
│   └── qms-engine/
│       └── main.go                     # Application entry point
│
├── server/
│   ├── serve.go                        # Server initialization & DI
│   └── config/                         # Configuration management
│       ├── config.go                   # Config structures
│       ├── db.go                       # Database initialization
│       ├── gin.go                      # HTTP server setup
│       └── ...
│
├── handlers/                           # HTTP Handlers (Controllers)
│   ├── qmsengine.service.go           # Main service with DI
│   ├── qmsengine.routes.go            # Route definitions
│   └── qmsengine_create_project_handler.go  # CRUD handlers
│
├── internal/
│   ├── entity/                         # Domain entities
│   │   └── project.go                 # Project entity
│   │
│   ├── model/                          # Request/Response models
│   │   └── project_model.go          # Project DTOs
│   │
│   ├── service/                        # Business logic
│   │   └── project_usecase.go        # Project use cases
│   │
│   └── repository/                     # Data access layer
│       ├── project_repository.go      # Repository interface
│       └── mysql/                     # MySQL implementation
│           ├── mysql_project_repository.go
│           └── repository.go          # Generic repository (optional)
│
└── db/
    └── mysql/
        └── deploy/
            └── 0000-init.sql          # Database schema
```

## Component Details

### 1. Handler Layer (`handlers/`)

**Purpose:** HTTP request/response handling

**Components:**
- `QMSEngineService`: Main service struct with dependencies
- Route registration and grouping
- Request validation
- Response formatting
- Error handling

**Example:**
```go
func (s *QMSEngineService) CreateProject(c *gin.Context) {
    var req model.CreateProjectRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    project, err := s.projectUsecase.CreateProject(c.Request.Context(), req)
    // ... handle response
}
```

### 2. Use Case Layer (`internal/service/`)

**Purpose:** Business logic and orchestration

**Responsibilities:**
- Validate business rules
- Orchestrate repository calls
- Handle transactions (if needed)
- Log business events
- Transform between DTOs and entities

**Example:**
```go
func (u *ProjectUsecase) CreateProject(ctx context.Context, req model.CreateProjectRequest) (*entity.Project, error) {
    // Business logic here
    project := entity.Project{
        Name:        req.Name,
        Description: req.Description,
    }
    
    if err := u.repo.Save(ctx, project); err != nil {
        return nil, err
    }
    
    return &project, nil
}
```

### 3. Repository Layer (`internal/repository/`)

**Purpose:** Data access abstraction

**Features:**
- Interface-based design for testability
- Master/Slave database pattern
- CRUD operations
- Soft delete support
- Query methods

**Example:**
```go
type IProjectRepository interface {
    Save(ctx context.Context, project entity.Project) error
    GetByID(ctx context.Context, id int64) (*entity.Project, error)
    ListProjects(ctx context.Context, limit, offset int) ([]*entity.Project, error)
    Update(ctx context.Context, id int64, project entity.Project) error
    Delete(ctx context.Context, id int64) error
}
```

### 4. Entity Layer (`internal/entity/`)

**Purpose:** Domain models

**Features:**
- Pure domain objects
- Database mapping tags
- Table name method

**Example:**
```go
type Project struct {
    ID          int        `json:"id" db:"id"`
    Name        string     `json:"name" db:"name"`
    Description string     `json:"description" db:"description"`
    CreatedAt   time.Time  `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
    DeletedAt   *time.Time `json:"deleted_at" db:"deleted_at"`
}
```

## Dependency Injection Flow

```
serve.go (main)
    ↓
1. Load Config
    ↓
2. Initialize Logger
    ↓
3. Initialize Database (Master/Slave)
    ↓
4. Initialize Repository (with DB)
    ↓
5. Initialize UseCase (with Repository + Logger)
    ↓
6. Initialize Handler (with UseCase + Logger)
    ↓
7. Initialize Gin Engine
    ↓
8. Register Routes
    ↓
9. Start Server
```

## Database Pattern

### Master/Slave Configuration

- **Master DB**: Used for all write operations (INSERT, UPDATE, DELETE)
- **Slave DB**: Used for all read operations (SELECT)

This pattern provides:
- ✅ Read scalability
- ✅ Load distribution
- ✅ High availability
- ✅ Data redundancy

### Soft Delete

All DELETE operations set `deleted_at` timestamp instead of removing records:
- Maintains data history
- Allows data recovery
- Supports audit trails

## Key Features

### ✅ Implemented

1. **Clean Architecture** - Separation of concerns
2. **Dependency Injection** - Testable and maintainable
3. **Repository Pattern** - Database abstraction
4. **Context Propagation** - Cancellation and timeout support
5. **Structured Logging** - JSON format with fields
6. **Master/Slave DB** - Read/write splitting
7. **Soft Delete** - Data preservation
8. **RESTful API** - Standard HTTP methods
9. **Request Validation** - Input sanitization
10. **Error Handling** - Comprehensive error responses
11. **CORS Support** - Cross-origin requests
12. **Health Checks** - Service monitoring

## API Routes

```
GET    /health                    # Health check
GET    /ping                      # Ping endpoint
POST   /api/v1/projects          # Create project
GET    /api/v1/projects          # List projects (paginated)
GET    /api/v1/projects/:id      # Get project by ID
PUT    /api/v1/projects/:id      # Update project
DELETE /api/v1/projects/:id      # Delete project (soft)
```

## Configuration

Configuration is loaded from `config_files/service-config.json`:

```json
{
  "env": "development",
  "host": "0.0.0.0",
  "port": 8080,
  "data": {
    "mysql": {
      "master": {
        "dsn": "user:password@tcp(localhost:3306)/qms_db",
        "maxIdle": 10,
        "maxOpen": 100,
        "connMaxLifetime": "1h"
      },
      "slave": {
        "dsn": "user:password@tcp(localhost:3306)/qms_db",
        "maxIdle": 10,
        "maxOpen": 100,
        "connMaxLifetime": "1h"
      }
    }
  }
}
```

## Testing Strategy

### Unit Tests
- Test use cases with mock repositories
- Test handlers with mock use cases
- Test repository methods with test database

### Integration Tests
- Test full request/response flow
- Test database operations
- Test error scenarios

### Example Mock
```go
type MockProjectRepository struct {
    mock.Mock
}

func (m *MockProjectRepository) Save(ctx context.Context, project entity.Project) error {
    args := m.Called(ctx, project)
    return args.Error(0)
}
```

## Future Enhancements

1. **Authentication & Authorization**
   - JWT tokens
   - Role-based access control

2. **Caching Layer**
   - Redis integration
   - Cache invalidation

3. **Event System**
   - Kafka integration
   - Event-driven architecture

4. **API Versioning**
   - Multiple API versions support

5. **Rate Limiting**
   - Request throttling
   - DDoS protection

6. **Metrics & Monitoring**
   - Prometheus integration
   - Grafana dashboards

7. **API Documentation**
   - Swagger/OpenAPI
   - Auto-generated docs

## Running the Application

```bash
# Build
make build

# Run
make run

# Test
make test

# Database migration
make db-migrate
```

