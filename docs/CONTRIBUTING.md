# Contributing to TaskMaster

## Development Workflow

### Branch Naming
- `feature/description` — New features
- `bugfix/description` — Bug fixes
- `hotfix/description` — Production fixes
- `docs/description` — Documentation updates

### Commit Messages
Follow conventional commits:
```
feat: add task assignment endpoint
fix: resolve race condition in WebSocket hub
docs: update API documentation
refactor: extract validation logic
test: add integration tests for auth flow
```

### Code Review Checklist
- [ ] Tests pass (`make test`)
- [ ] Linting passes (`make lint`)
- [ ] No hardcoded secrets
- [ ] Error handling is explicit
- [ ] Context is passed through all layers
- [ ] Database queries use parameterized statements
- [ ] No memory leaks (check goroutine usage)

---

## Architecture Decisions

### Why Clean Architecture?
- **Testability**: Business logic has no external dependencies
- **Flexibility**: Swap PostgreSQL for MongoDB without touching services
- **Clarity**: Clear boundaries between layers

### Why Go?
- Performance for high-concurrency WebSocket connections
- Static typing catches errors at compile time
- Excellent standard library for networking
- HENNGE requires Go proficiency

### Why PostgreSQL?
- ACID compliance for financial/audit data
- JSONB support for flexible metadata
- Excellent Go driver (pgx)
- RDS managed service on AWS

### Why Redis?
- Session storage with TTL
- Rate limiting counters
- Cache invalidation
- WebSocket Pub/Sub for multi-container sync

---

## Testing Guidelines

### Unit Tests
```go
func TestService_Method(t *testing.T) {
    // Arrange
    mockRepo := new(MockRepository)
    svc := NewService(mockRepo)

    // Act
    result, err := svc.Method(ctx, input)

    // Assert
    assert.NoError(t, err)
    assert.Equal(t, expected, result)
    mockRepo.AssertExpectations(t)
}
```

### Integration Tests
```go
func TestAPI_CreateTask(t *testing.T) {
    // Setup test database
    db := setupTestDB(t)

    // Create handler with real repositories
    router := setupRouter(db)

    // Execute request
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/api/v1/tasks", body)
    router.ServeHTTP(w, req)

    // Assert
    assert.Equal(t, 201, w.Code)
}
```

### Test Coverage Targets
- Domain layer: 100%
- Service layer: 90%
- Handlers: 80%
- Integration: Critical paths

---

## Performance Guidelines

### Database
- Use connection pooling (pgxpool)
- Index foreign keys and query columns
- Use EXPLAIN ANALYZE for slow queries
- Batch inserts when possible

### Caching
- Cache read-heavy data (task lists, user profiles)
- Invalidate on write
- Use Redis TTL for automatic expiration

### WebSocket
- Limit connections per client
- Use ping/pong for connection health
- Batch broadcasts when possible

### Memory
- Avoid goroutine leaks (always close channels)
- Use sync.Pool for object reuse
- Profile with pprof in production

---

## Security Guidelines

### Authentication
- Always validate JWT before processing
- Use HTTPS in production
- Rotate secrets regularly
- Implement refresh token flow

### Authorization
- Check permissions at handler level
- Verify resource ownership
- Log access denials

### Input Validation
- Validate all user input
- Sanitize before database queries
- Limit request body size

### Secrets
- Never commit secrets to git
- Use AWS Secrets Manager in production
- Rotate credentials quarterly
