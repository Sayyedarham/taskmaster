# TaskMaster Architecture

## Overview

TaskMaster follows **Clean Architecture** (Hexagonal Architecture) principles, ensuring that business logic is independent of frameworks, UI, and external services. This makes the codebase highly testable, maintainable, and adaptable to changing requirements.

## Architecture Layers

### 1. Domain Layer (Core)
**Location:** `backend/internal/domain/`

The innermost layer containing pure business entities. It has zero external dependencies.

```go
type Task struct {
    ID          uuid.UUID
    Title       string
    Status      TaskStatus
    Priority    Priority
    // ... pure Go structs
}
```

**Rules:**
- No imports of Gin, PostgreSQL, Redis, or AWS SDK
- Contains only structs, enums, and value objects
- Defines business rules through methods

### 2. Ports Layer (Interfaces)
**Location:** `backend/internal/ports/`

Defines contracts (interfaces) that the application layer expects from infrastructure.

```go
type UserRepository interface {
    Create(ctx context.Context, user *domain.User) error
    FindByEmail(ctx context.Context, email string) (*domain.User, error)
}
```

**Benefits:**
- Domain layer doesn't care if PostgreSQL or MongoDB is used
- Easy to swap implementations for testing
- Clear boundaries between layers

### 3. Application Layer (Use Cases)
**Location:** `backend/internal/service/`

Contains business logic orchestration. Services coordinate between domain entities and repositories.

```
HTTP Request → Handler → Service → Repository Interface → PostgreSQL
                                    ↑
                              Business Rules
                                    ↓
HTTP Response ← Handler ← Service ← Repository Interface ← PostgreSQL
```

**Key Services:**
- `AuthService`: Registration, login, JWT generation
- `TaskService`: Task CRUD with SLA validation, team limits, assignment

### 4. Infrastructure Layer (Adapters)
**Location:** `backend/internal/repository/`

Implements the ports interfaces using real technologies.

- **PostgreSQL Repository**: SQL queries, connection pooling
- **Redis Repository**: Caching, rate limiting, session storage
- **WebSocket Hub**: Real-time communication

### 5. External Layer (Frameworks)
**Location:** `backend/internal/handler/`, `backend/internal/middleware/`

The "dirty" outer layer that connects to the outside world.

- **HTTP Handlers**: Gin route handlers
- **Middleware**: JWT validation, RBAC, rate limiting, CORS, logging
- **WebSocket Handler**: Gorilla WebSocket upgrades

## Data Flow

### Scenario: Manager Creates a Task

```
1. Browser sends POST /api/v1/tasks
   ↓
2. ALB (AWS) routes to ECS Fargate
   ↓
3. Gin Router matches route
   ↓
4. Middleware Chain:
   - Logger: "POST /tasks started"
   - CORS: Check origin
   - Rate Limiter: Redis INCR ip:tasks:count
   - JWT Auth: Validate token, extract claims
   - RBAC: Verify manager has task:create permission
   ↓
5. TaskHandler parses JSON → TaskCreateDTO
   ↓
6. TaskService executes business logic:
   - Is assignee in same team? → TeamRepo.IsMember()
   - Does team have task limits? → CacheRepo.Get() + TaskRepo.CountByTeam()
   - Is due_date within SLA? → TeamRepo.GetSLA()
   - Create Task entity → TaskRepo.Create()
   - Invalidate cache → CacheRepo.Delete()
   - Audit log → AuditRepo.Create()
   - Broadcast real-time → wsHub.Broadcast()
   ↓
7. HTTP 201 response with task JSON
   ↓
8. WebSocket pushes to all team members' browsers
```

## Dependency Rule

**Inner layers know nothing about outer layers.**

```
Domain ← Application ← Infrastructure ← External
   ↑         ↑            ↑              ↑
  Pure    Uses ports   Implements    HTTP, WS,
  Go      only         interfaces    AWS SDK
```

## Multi-Container WebSocket Sync

When ECS scales to 2+ tasks, each container has its own in-memory WebSocket Hub. Redis Pub/Sub synchronizes messages across containers.

```
Container 1 ──► Local Hub ──► Redis PUBLISH ──► Container 2
                    │                              │
                    └──► Client A                  └──► Client B
```

## AWS Architecture

```
Internet
   │
Route53 (DNS)
   │
┌──┴──┐
│     │
▼     ▼
CloudFront    ALB
(S3 SPA)      (ECS Tasks)
   │            │
   ▼            ▼
  S3         ECS Fargate
             (Go API)
                │
    ┌───────────┼───────────┐
    ▼           ▼           ▼
   RDS      ElastiCache   Secrets
(PostgreSQL)  (Redis)     Manager
```

## Security

- **Authentication**: JWT tokens with HS256 signing
- **Authorization**: RBAC with role-based permissions
- **Rate Limiting**: Redis-based sliding window (100 req/min)
- **Secrets**: AWS Secrets Manager in production
- **CORS**: Configured for specific origins
- **Input Validation**: Gin binding validators
- **SQL Injection**: Parameterized queries via pgx
- **Passwords**: bcrypt hashing with default cost

## Testing Strategy

| Type | Scope | Tools |
|------|-------|-------|
| Unit | Individual functions | testify/mock |
| Integration | Full HTTP flows | testcontainers-go |
| E2E | Browser automation | Playwright |

## Technology Choices

| Decision | Rationale |
|----------|-----------|
| **Go** | Performance, concurrency, HENNGE requirement |
| **Gin** | Fast, middleware-friendly, widely used |
| **pgx** | Modern PostgreSQL driver with connection pooling |
| **go-redis** | Official Redis client, cluster support |
| **Clean Architecture** | Testability, maintainability, clear boundaries |
| **Docker** | Consistent environments, easy deployment |
| **Terraform** | Infrastructure as Code, reproducible |
| **GitHub Actions** | Native integration, free for public repos |
