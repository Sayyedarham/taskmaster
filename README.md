# TaskMaster — HENNGE-Aligned Full-Stack Project

> **Production-grade task management system** built with Go, AWS, Terraform, Docker, and React.

---

## Tech Stack

| Layer | Technology |
|-------|-----------|
| **Backend** | Go 1.22 + Gin + Clean Architecture |
| **Frontend** | React 18 + TypeScript + Vite (Phase 2) |
| **Database** | PostgreSQL 16 (RDS in prod) |
| **Cache** | Redis 7 (ElastiCache in prod) |
| **Cloud** | AWS (ECS Fargate, ALB, RDS, ElastiCache, S3, CloudFront, Route53) |
| **IaC** | Terraform |
| **CI/CD** | GitHub Actions → ECR → ECS |
| **Auth** | JWT + RBAC + Rate Limiting |
| **Real-time** | WebSocket + Redis Pub/Sub |

---

## Quick Start (Local Development)

### Windows Users

On Windows, use the batch scripts instead of shell scripts:

```powershell
# Setup everything
.\scripts\setup.bat

# Or manually:
copy .env.example .env
docker compose up -d postgres redis

# Wait for PostgreSQL, then run migrations
docker compose up migrate

# Seed data
.\scripts\seed.bat

# Start API and Frontend
docker compose up -d api frontend
```

**Note:** If you see `the attribute 'version' is obsolete`, this is just a warning and can be safely ignored. We've removed the `version` attribute from `docker-compose.yml` to fix this.

### macOS/Linux Users

### Prerequisites
- Docker & Docker Compose
- Go 1.22+ (for local Go dev)
- Make (optional)

### 1. Clone & Start Infrastructure
```bash
git clone <your-repo>
cd taskmaster

# Start PostgreSQL + Redis + API
docker-compose up -d

# Run migrations automatically
docker-compose logs migrate
```

### 2. Test the API
```bash
# Health check
curl http://localhost:8080/health

# Register
curl -X POST http://localhost:8080/api/v1/auth/register   -H "Content-Type: application/json"   -d '{"email":"dev@example.com","password":"password123","name":"Dev User"}'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login   -H "Content-Type: application/json"   -d '{"email":"dev@example.com","password":"password123"}'

# Create Task (use token from login)
curl -X POST http://localhost:8080/api/v1/tasks   -H "Content-Type: application/json"   -H "Authorization: Bearer <TOKEN>"   -d '{"title":"Build HENNGE project","priority":"high"}'
```

### 3. WebSocket Test
```bash
# Connect to WebSocket (use wscat or browser console)
wscat -c ws://localhost:8080/ws
# You'll receive real-time broadcasts when tasks are created
```

---

## Project Structure

```
taskmaster/
├── backend/
│   ├── cmd/server/           # Entry point
│   ├── internal/
│   │   ├── config/           # DB, Redis, env loader
│   │   ├── domain/           # Entities (User, Task, Team)
│   │   ├── ports/            # Interfaces (Repository contracts)
│   │   ├── service/          # Business logic (Auth, Task)
│   │   ├── handler/          # HTTP handlers (Gin)
│   │   ├── middleware/       # JWT, RBAC, Rate Limit, CORS, Logger
│   │   ├── repository/
│   │   │   ├── postgres/     # SQL implementations
│   │   │   └── redis/        # Cache implementations
│   │   └── websocket/        # WS Hub + Client management
│   ├── migrations/           # golang-migrate SQL files
│   ├── tests/                # Unit + Integration tests
│   ├── Dockerfile            # Multi-stage build
│   └── go.mod
├── frontend/                 # React SPA (Phase 2)
├── infrastructure/
│   ├── modules/              # VPC, ECS, RDS, ElastiCache, ALB, S3/CloudFront
│   └── environments/         # dev / prod
├── .github/workflows/         # CI/CD pipeline
├── docker-compose.yml         # Local dev stack
└── README.md
```

---

## Architecture: Clean / Hexagonal

```
HTTP Request
    ↓
[Gin Handler] ──→ [Middleware: JWT, RBAC, RateLimit]
    ↓
[Service Layer] ──→ Business Rules (SLA, Limits, Assignment)
    ↓
[Repository Interface] ──→ [PostgreSQL / Redis / WebSocket]
    ↓
[Domain Entities] (pure Go structs, no framework deps)
```

**Dependency Rule:** Inner layers know NOTHING about outer layers.
- **Domain:** Pure structs + interfaces
- **Application:** Services use domain interfaces
- **Infrastructure:** Implements interfaces with real tech
- **External:** HTTP, WS, AWS SDK

---

## Development Roadmap

### ✅ Phase 1: Backend Core (COMPLETE IN STARTER)
- [x] Go project scaffold with Clean Architecture
- [x] PostgreSQL + Redis connection management
- [x] Domain entities (User, Task, Team)
- [x] Repository interfaces + PostgreSQL implementations
- [x] Redis cache repository
- [x] Auth service (Register/Login with bcrypt + JWT)
- [x] Task service (Create with SLA/Limit validation)
- [x] Gin handlers + middleware (JWT, RBAC, RateLimit, CORS, Logger)
- [x] WebSocket Hub for real-time updates
- [x] Docker + docker-compose local setup
- [x] Database migrations (golang-migrate)
- [x] Terraform skeleton (VPC, ECS, RDS modules)
- [x] GitHub Actions CI/CD pipeline

### 🔨 Phase 2: Frontend & Polish
- [ ] React 18 + TypeScript + Vite scaffold
- [ ] Auth pages (Login/Register)
- [ ] Dashboard with task list
- [ ] Task creation modal
- [ ] Real-time WebSocket integration
- [ ] Role-based UI rendering
- [ ] S3 + CloudFront hosting via Terraform

### 🔨 Phase 3: AWS Production
- [ ] Complete Terraform modules (RDS, ElastiCache, ALB, ECS, IAM)
- [ ] AWS Secrets Manager integration
- [ ] ECR repository + image push
- [ ] ECS Fargate service with auto-scaling
- [ ] Blue-green deployment strategy
- [ ] CloudWatch logging + metrics + alarms
- [ ] Route53 + SSL certificate setup

### 🔨 Phase 4: Advanced Features
- [ ] Multi-container WebSocket sync via Redis Pub/Sub
- [ ] Audit logging system
- [ ] Task assignment + notifications
- [ ] Team management endpoints
- [ ] Integration tests with testcontainers
- [ ] E2E tests with Playwright
- [ ] API documentation (Swagger/OpenAPI)

---

## Checklist Alignment

| Expectation | How TaskMaster Demonstrates It |
|-------------------|-------------------------------|
| **Go proficiency** | Clean Architecture, interfaces, context usage, goroutines |
| **TypeScript + React** | Phase 2: Type-safe frontend with modern hooks |
| **AWS Cloud** | ECS, RDS, ElastiCache, S3, CloudFront, ALB, Route53 |
| **Docker** | Multi-stage builds, docker-compose, ECS containers |
| **Terraform** | Modular IaC for entire AWS infrastructure |
| **CI/CD** | GitHub Actions: lint → test → build → push → deploy |
| **Testing** | Unit (mocked), Integration (testcontainers), E2E |
| **Auth/Security** | JWT, RBAC, bcrypt, rate limiting, secrets manager |
| **DevOps** | Zero-downtime deployment, auto-scaling, monitoring |
| **Linux/Unix** | Alpine containers, signal handling, file permissions |

---

## Testing Strategy

```bash
cd backend

# Unit tests (fast, mocked)
go test ./... -short

# Integration tests (requires Docker)
go test ./tests/integration/... -v

# Coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

---

## Deployment Flow

```
Developer pushes to main
    ↓
GitHub Actions
    ├── golangci-lint
    ├── go test (unit + integration)
    ├── docker build
    ├── push to ECR
    ├── terraform apply
    └── ecs update-service --force-new-deployment
    ↓
AWS ECS
    ├── Start new container (Blue)
    ├── ALB health check
    ├── Register to target group
    ├── Drain old connections
    └── Stop old container (Green)
```

---

## Environment Variables

| Variable | Local Default | Production Source |
|----------|--------------|-------------------|
| `DATABASE_URL` | `postgres://...` | AWS Secrets Manager |
| `REDIS_ADDR` | `localhost:6379` | ElastiCache endpoint |
| `JWT_SECRET` | `dev-secret...` | AWS Secrets Manager |
| `PORT` | `8080` | `8080` |

Copy `.env.example` to `.env` for local development.

---

## Contributing

1. Fork the repo
2. Create a feature branch: `git checkout -b feature/amazing-thing`
3. Commit changes: `git commit -m 'Add amazing thing'`
4. Push: `git push origin feature/amazing-thing`
5. Open a Pull Request

---

## License

MIT — Built for HENNGE Global Internship Program preparation.

---

**Ready to build?** Start with `docker-compose up -d` and hit `http://localhost:8080/health` 🚀


## 🚀 Complete Setup Guide

For a **step-by-step walkthrough** from zero to deployed on AWS, see:

**[📖 COMPLETE_SETUP.md](docs/COMPLETE_SETUP.md)**

This covers:
- AWS account creation (free tier)
- IAM user setup with correct permissions
- GitHub repository + secrets configuration
- Local Docker development
- First AWS deployment with Terraform
- CI/CD pipeline activation
- What's built vs. what remains to develop
