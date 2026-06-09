# Troubleshooting Guide

## Common Issues & Solutions

### 1. Docker Compose Fails to Start

**Symptom:** `docker-compose up` fails with "Cannot start service"

**Diagnosis:**
```bash
# Check Docker daemon
docker info

# Check port conflicts
lsof -ti:8080  # API port
lsof -ti:3000  # Frontend port
lsof -ti:5432  # PostgreSQL port
lsof -ti:6379  # Redis port
```

**Solutions:**
```bash
# Kill conflicting processes
lsof -ti:8080 | xargs kill -9

# Or use different ports in .env
PORT=8081
DATABASE_URL=postgres://...:5433/...  # Note: change docker-compose.yml too
```

---

### 2. PostgreSQL Connection Refused

**Symptom:** API logs show "connection refused" to postgres:5432

**Diagnosis:**
```bash
docker-compose logs postgres
docker-compose exec postgres pg_isready -U taskmaster
```

**Solutions:**
```bash
# Wait for PostgreSQL to be fully ready
docker-compose up -d postgres
sleep 5  # Give it time to initialize

# Check if database exists
docker-compose exec postgres psql -U taskmaster -l

# Reset if corrupted
docker-compose down -v
docker-compose up -d postgres
```

---

### 3. Migrations Fail

**Symptom:** `migrate` container exits with error

**Common Errors:**
- `Dirty database version X`: Migration was partially applied
- `Table already exists`: Running migration twice
- `Column does not exist`: Migration order wrong

**Solutions:**
```bash
# Fix dirty database
cd backend
migrate -path migrations -database "postgres://taskmaster:taskmaster@localhost:5432/taskmaster?sslmode=disable" force <LAST_SUCCESSFUL_VERSION>

# Reset everything
docker-compose down -v
docker-compose up -d postgres
migrate -path migrations -database "postgres://taskmaster:taskmaster@localhost:5432/taskmaster?sslmode=disable" up
```

---

### 4. JWT Authentication Fails

**Symptom:** API returns 401 "invalid token"

**Diagnosis:**
```bash
# Check token format
echo $TOKEN | cut -d'.' -f2 | base64 -d  # Decode payload

# Check JWT secret matches
# backend/.env JWT_SECRET must match the secret used to sign the token
```

**Solutions:**
```bash
# Re-login to get fresh token
curl -X POST http://localhost:8080/api/v1/auth/login   -H "Content-Type: application/json"   -d '{"email":"admin@taskmaster.com","password":"password123"}'
```

---

### 5. Frontend Cannot Reach API

**Symptom:** Browser console shows "Failed to fetch" or CORS errors

**Diagnosis:**
```bash
# Check API is running
curl http://localhost:8080/health

# Check CORS headers
curl -I -X OPTIONS http://localhost:8080/api/v1/tasks   -H "Origin: http://localhost:3000"   -H "Access-Control-Request-Method: POST"
```

**Solutions:**
```bash
# If using Docker frontend, check nginx.conf proxy settings
# If using local frontend, check vite.config.ts proxy

# Ensure API CORS allows frontend origin
# In backend/internal/middleware/cors.go:
# c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
```

---

### 6. WebSocket Not Working

**Symptom:** No real-time updates in dashboard

**Diagnosis:**
```bash
# Test WebSocket manually
wscat -c ws://localhost:8080/ws

# Check Redis is running
docker-compose exec redis redis-cli ping

# Check API logs for WS errors
docker-compose logs -f api | grep -i websocket
```

**Solutions:**
```bash
# Restart API to reconnect WS hub
docker-compose restart api

# Check firewall/proxy settings
# WebSocket needs Upgrade header support
```

---

### 7. Go Build Fails

**Symptom:** `go build` fails with module errors

**Solutions:**
```bash
cd backend

# Download dependencies
go mod download

# Tidy modules
go mod tidy

# Verify modules
go mod verify

# If using private repos:
export GOPRIVATE=github.com/yourorg
```

---

### 8. Frontend Build Fails

**Symptom:** `npm run build` fails

**Solutions:**
```bash
cd frontend

# Clean install
rm -rf node_modules package-lock.json
npm install

# Check TypeScript errors
npx tsc --noEmit

# Check for missing types
npm install --save-dev @types/missing-package
```

---

### 9. Rate Limiting Blocks Requests

**Symptom:** API returns 429 "Too Many Requests"

**Diagnosis:**
```bash
# Check Redis keys
docker-compose exec redis redis-cli KEYS "ratelimit:*"

# Check current count
docker-compose exec redis redis-cli GET "ratelimit:127.0.0.1:/api/v1/tasks"
```

**Solutions:**
```bash
# Wait 60 seconds for rate limit to reset
# Or clear Redis:
docker-compose exec redis redis-cli FLUSHALL

# Adjust limit in middleware/ratelimit.go if needed for development
```

---

### 10. Memory Issues

**Symptom:** Docker containers use too much RAM

**Solutions:**
```bash
# Limit container resources in docker-compose.yml:
services:
  api:
    deploy:
      resources:
        limits:
          memory: 512M

# Or use Docker Desktop settings to limit total Docker memory
```

---

## Debug Mode

Enable debug logging:

```bash
# Backend
docker-compose exec api sh
export GIN_MODE=debug
# Or set in .env: ENV=development

# Frontend
# React DevTools browser extension
# Vite: npm run dev (already in debug mode)
```

## Getting Help

1. Check logs: `docker-compose logs -f <service>`
2. Check health: `curl http://localhost:8080/health`
3. Review `.env` configuration
4. Check GitHub Issues (when published)
5. Run diagnostics: `./scripts/test.sh`
