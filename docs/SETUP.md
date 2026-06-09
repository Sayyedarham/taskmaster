# Complete Local Setup Guide — TaskMaster

> **Goal**: Get the entire stack running locally in under 5 minutes with one command.

---

## Prerequisites

| Tool | Version | Check Command |
|------|---------|---------------|
| Docker | 24.0+ | `docker --version` |
| Docker Compose | 2.20+ | `docker compose version` |
| Git | 2.40+ | `git --version` |
| Make | 3.81+ | `make --version` |
| Go | 1.22+ (optional, for local dev) | `go version` |
| Node.js | 20+ (optional, for frontend dev) | `node --version` |

**Install Docker Desktop** (includes Compose):
- **macOS**: `brew install --cask docker` or download from [docker.com](https://www.docker.com/products/docker-desktop)
- **Windows**: Download from [docker.com](https://www.docker.com/products/docker-desktop)
- **Linux**: `sudo apt install docker.io docker-compose-plugin`

---

## Step 1: Clone & Extract

```bash
# If you have the zip file
unzip taskmaster-starter.zip
cd taskmaster-starter

# Or clone from GitHub (when published)
git clone https://github.com/yourusername/taskmaster.git
cd taskmaster
```

---

## Step 2: One-Command Setup (Recommended)

```bash
./scripts/setup.sh
```

This script automatically:
1. Checks Docker is running
2. Creates `.env` from `.env.example`
3. Starts PostgreSQL and Redis containers
4. Waits for PostgreSQL to be healthy
5. Runs database migrations
6. Seeds demo data (3 users, 5 tasks, 1 team)
7. Starts the Go API and React frontend

**Output you should see:**
```
🚀 TaskMaster Setup Script
===========================
Checking prerequisites...
✅ Docker and Docker Compose found
Creating .env from example...
✅ .env created.
Starting PostgreSQL, Redis, and API...
Waiting for PostgreSQL to be ready...
✅ PostgreSQL is ready
Running database migrations...
Seeding demo data...
Starting API and Frontend...

✅ Setup complete!

📍 Access points:
   API:      http://localhost:8080
   Frontend: http://localhost:3000
   Health:   http://localhost:8080/health

🔑 Demo credentials:
   Admin:   admin@taskmaster.com / password123
   Manager: manager@taskmaster.com / password123
   Member:  member@taskmaster.com / password123
```

---

## Step 3: Verify Everything Works

### 3.1 Health Check
```bash
curl http://localhost:8080/health
```
**Expected:** `{"status":"ok","time":"2026-05-23T..."}`

### 3.2 API Test — Register a User
```bash
curl -X POST http://localhost:8080/api/v1/auth/register   -H "Content-Type: application/json"   -d '{"email":"you@example.com","password":"password123","name":"Your Name"}'
```
**Expected:** `{"user":{"id":"...","email":"you@example.com","name":"Your Name","role":"member"}}`

### 3.3 API Test — Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login   -H "Content-Type: application/json"   -d '{"email":"admin@taskmaster.com","password":"password123"}'
```
**Expected:** `{"token":"eyJ...","user":{"id":"...","email":"admin@taskmaster.com","role":"admin"}}`

Save the token for the next step.

### 3.4 API Test — Create a Task
```bash
curl -X POST http://localhost:8080/api/v1/tasks   -H "Content-Type: application/json"   -H "Authorization: Bearer <YOUR_TOKEN>"   -d '{"title":"My first task","priority":"high","description":"Learning TaskMaster"}'
```

### 3.5 API Test — List Tasks
```bash
curl http://localhost:8080/api/v1/tasks   -H "Authorization: Bearer <YOUR_TOKEN>"
```

### 3.6 Frontend Test
Open **http://localhost:3000** in your browser:
- Login with demo credentials
- See the dashboard with pre-seeded tasks
- Create a new task (if you're admin/manager)
- Open browser DevTools → Network tab to see WebSocket traffic

### 3.7 WebSocket Test
```bash
# Install wscat globally
npm install -g wscat

# Connect to WebSocket
wscat -c ws://localhost:8080/ws
```
Then create a task via API — you should see a real-time message in the wscat terminal.

---

## Step 4: Manual Setup (If Script Fails)

If the automated script doesn't work, here's the manual process:

### 4.1 Create Environment File
```bash
cp .env.example .env
```

### 4.2 Start Infrastructure Only
```bash
docker-compose up -d postgres redis
```

### 4.3 Wait for PostgreSQL
```bash
# Check health
docker-compose exec postgres pg_isready -U taskmaster

# Or wait manually
until docker-compose exec -T postgres pg_isready -U taskmaster > /dev/null 2>&1; do
    echo "Waiting for PostgreSQL..."
    sleep 1
done
```

### 4.4 Run Migrations
```bash
# Using Docker (recommended)
docker-compose up migrate

# Or using local golang-migrate (if installed)
cd backend
migrate -path migrations -database "postgres://taskmaster:taskmaster@localhost:5432/taskmaster?sslmode=disable" up
```

### 4.5 Seed Data (Optional)
```bash
# Using Docker
docker-compose exec -T postgres psql -U taskmaster -d taskmaster -f /migrations/002_seed.up.sql

# Or using local psql
psql postgres://taskmaster:taskmaster@localhost:5432/taskmaster -f backend/migrations/002_seed.up.sql
```

### 4.6 Start API
```bash
# Using Docker (production-like)
docker-compose up -d api

# Or run Go locally (for hot reload development)
cd backend
go mod download
go run cmd/server/main.go
```

### 4.7 Start Frontend
```bash
# Using Docker
docker-compose up -d frontend

# Or run React locally (for hot reload)
cd frontend
npm install
npm run dev
```

---

## Step 5: Development Workflow

### 5.1 Hot Reload Development (Recommended)

For the best development experience, run services outside Docker:

**Terminal 1 — Infrastructure:**
```bash
docker-compose up -d postgres redis
```

**Terminal 2 — Go API (with auto-reload):**
```bash
cd backend
# Install air for hot reload: go install github.com/cosmtrek/air@latest
air
# Or manually:
go run cmd/server/main.go
```

**Terminal 3 — React Frontend (with HMR):**
```bash
cd frontend
npm install
npm run dev
```

**Benefits:**
- Instant code changes without rebuilding Docker images
- Better debugging experience
- Faster test execution

### 5.2 Docker-Only Development

If you prefer everything in containers:

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f api
docker-compose logs -f frontend

# Rebuild after code changes
docker-compose up -d --build api
```

---

## Step 6: Running Tests

### 6.1 Unit Tests (Fast, Mocked)
```bash
# Using Make
make test

# Or manually
cd backend
go test ./tests/unit/... -v -race -cover
```

**Expected output:**
```
=== RUN   TestAuthService_Register_Success
--- PASS: TestAuthService_Register_Success (0.00s)
=== RUN   TestAuthService_Register_DuplicateEmail
--- PASS: TestAuthService_Register_DuplicateEmail (0.00s)
=== RUN   TestTaskService_Create_Success
--- PASS: TestTaskService_Create_Success (0.00s)
=== RUN   TestTaskService_Create_ExceedsLimit
--- PASS: TestTaskService_Create_ExceedsLimit (0.00s)
PASS
coverage: 78.5% of statements
```

### 6.2 Integration Tests (Requires Docker)
```bash
# Using Make
make test-integration

# Or manually
cd backend
go test ./tests/integration/... -v
```

**Note:** Integration tests use testcontainers-go to spin up real PostgreSQL and Redis containers automatically.

### 6.3 Linting
```bash
make lint
# Or: cd backend && golangci-lint run
```

---

## Step 7: Database Operations

### 7.1 Create a New Migration
```bash
make migrate-create name=add_notifications

# Or manually
cd backend
migrate create -ext sql -dir migrations add_notifications
```

This creates:
- `migrations/003_add_notifications.up.sql`
- `migrations/003_add_notifications.down.sql`

### 7.2 Apply Migrations
```bash
make migrate-up
```

### 7.3 Rollback One Migration
```bash
make migrate-down
```

### 7.4 Connect to Database Directly
```bash
# Using Docker
docker-compose exec postgres psql -U taskmaster -d taskmaster

# Using local psql
psql postgres://taskmaster:taskmaster@localhost:5432/taskmaster

# Common queries
\dt                    -- List tables
\d users              -- Describe users table
SELECT * FROM tasks;    -- View all tasks
```

### 7.5 Reset Database (⚠️ Destroys all data)
```bash
docker-compose down -v
docker-compose up -d postgres redis
make migrate-up
```

---

## Step 8: Troubleshooting

### Issue: "Cannot connect to Docker daemon"
**Fix:**
```bash
# macOS/Windows: Start Docker Desktop
# Linux:
sudo systemctl start docker
sudo usermod -aG docker $USER  # Log out and back in
```

### Issue: "Port 8080 already in use"
**Fix:**
```bash
# Find and kill process
lsof -ti:8080 | xargs kill -9

# Or change port in .env
echo "PORT=8081" >> .env
```

### Issue: "Migration failed: dirty database"
**Fix:**
```bash
cd backend
migrate -path migrations -database "postgres://taskmaster:taskmaster@localhost:5432/taskmaster?sslmode=disable" force <VERSION>
```

### Issue: "Frontend shows 'Cannot connect to API'"
**Fix:**
```bash
# Check API is running
curl http://localhost:8080/health

# Check CORS settings in backend/internal/middleware/cors.go
# Ensure vite.config.ts proxy points to correct API port
```

### Issue: "WebSocket not receiving updates"
**Fix:**
```bash
# Check WebSocket endpoint
curl -i -N -H "Connection: Upgrade" -H "Upgrade: websocket" http://localhost:8080/ws

# Verify Redis is running
docker-compose exec redis redis-cli ping
```

### Issue: "Permission denied" on scripts
**Fix:**
```bash
chmod +x scripts/*.sh
```

---

## Step 9: Useful Commands Cheat Sheet

```bash
# === LIFECYCLE ===
make dev          # Start all services
make dev-down     # Stop all services
make dev-logs     # Tail all logs
make clean        # Nuke everything including volumes

# === TESTING ===
make test              # Unit tests
make test-integration  # Integration tests
make lint              # Code linting

# === DATABASE ===
make migrate-up        # Apply migrations
make migrate-down      # Rollback one
make migrate-create    # Create new migration
make seed              # Seed demo data

# === INFRASTRUCTURE ===
make infra-init        # Terraform init
make infra-plan        # Terraform plan
make infra-apply       # Terraform apply
make infra-destroy     # Destroy AWS resources

# === DOCKER ===
docker-compose ps              # List running containers
docker-compose exec api sh     # Shell into API container
docker-compose logs -f api     # Follow API logs
docker-compose up -d --build   # Rebuild and start
```

---

## Step 10: Next Steps After Setup

1. **Explore the codebase** — Start with `backend/internal/domain/entities.go` to understand the data model
2. **Add a feature** — Try adding "task categories" or "due date notifications"
3. **Write tests** — Add unit tests for your new feature using the mock patterns in `tests/unit/`
4. **Deploy to AWS** — Follow `docs/DEPLOYMENT.md` for Terraform setup
5. **Customize** — Modify `frontend/src/index.css` for your own branding

---

## Architecture Quick Reference

```
Browser → Frontend (React + Vite) @ :3000
              ↓ Proxy (vite.config.ts)
         API (Go + Gin) @ :8080
              ↓
    ┌─────────┼─────────┐
    ▼         ▼         ▼
 PostgreSQL   Redis   WebSocket
   :5432     :6379     :8080/ws
```

**Data Flow:**
```
HTTP Request → Gin Router → Middleware (JWT, RBAC, RateLimit)
                              ↓
                         Handler → Service (Business Logic)
                              ↓
                         Repository Interface → PostgreSQL/Redis
                              ↓
                         WebSocket Broadcast → All connected clients
```

---

**You're all set!** If anything breaks, check `docker-compose logs` and refer to `docs/TROUBLESHOOTING.md` (create it as you learn).
