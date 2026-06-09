# Complete Setup Guide — TaskMaster

> **Prerequisites**: GitHub account, AWS free tier account, Docker installed locally.
> **Time**: ~45 minutes for first setup.

---

## PHASE 1: AWS ACCOUNT SETUP (10 min)

### Step 1.1: Create AWS Account
1. Go to https://aws.amazon.com/free
2. Click "Create a Free Account"
3. Enter email, password, account name
4. Choose "Personal" account type
5. Enter your name, phone number, address
6. Add credit/debit card (required for verification, won't charge if under free tier)
7. Verify identity via SMS/call
8. Select "Basic Support" (free)
9. Click "Complete Sign Up"
10. Sign in to AWS Management Console

### Step 1.2: Create IAM User (Never use root account for deployment)
1. In AWS Console, search "IAM" → click it
2. Left sidebar → Users → "Add users"
3. User name: `taskmaster-deploy`
4. Check "Provide user access to the AWS Management Console" (optional, but useful)
5. Click "Next"
6. Permissions options → "Attach policies directly"
7. Search and check these 8 policies:
   - `AmazonEC2FullAccess`
   - `AmazonRDSFullAccess`
   - `AmazonElastiCacheFullAccess`
   - `AmazonECS_FullAccess`
   - `AmazonS3FullAccess`
   - `AmazonRoute53FullAccess`
   - `CloudFrontFullAccess`
   - `IAMFullAccess`
8. Click "Next" → "Create user"
9. **IMPORTANT**: On the success screen, click "Download .csv file"
10. This file contains your Access Key ID and Secret Access Key
11. **Store this file securely** — you cannot download it again

### Step 1.3: Install AWS CLI Locally
**macOS:**
```bash
brew install awscli
```

**Windows:**
```bash
# Download installer from https://aws.amazon.com/cli/
# Or use winget:
winget install Amazon.AWSCLI
```

**Linux (Ubuntu/Debian):**
```bash
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
unzip awscliv2.zip
sudo ./aws/install
```

### Step 1.4: Configure AWS CLI
```bash
aws configure
```
Enter:
- AWS Access Key ID: (from the .csv file)
- AWS Secret Access Key: (from the .csv file)
- Default region name: `us-east-1`
- Default output format: `json`

Verify:
```bash
aws sts get-caller-identity
```
Should show your Account, UserId, and ARN.

### Step 1.5: Create Terraform State Bucket
```bash
aws s3 mb s3://taskmaster-terraform-state --region us-east-1
aws s3api put-bucket-versioning   --bucket taskmaster-terraform-state   --versioning-configuration Status=Enabled
```

---

## PHASE 2: GITHUB REPO SETUP (5 min)

### Step 2.1: Create GitHub Repository
1. Go to https://github.com/new
2. Repository name: `taskmaster`
3. Visibility: Public (free GitHub Actions minutes)
4. Check "Add a README file"
5. Click "Create repository"

### Step 2.2: Add AWS Secrets to GitHub
1. In your repo, click "Settings" tab
2. Left sidebar → "Secrets and variables" → "Actions"
3. Click "New repository secret"
4. Name: `AWS_ACCESS_KEY_ID`
   Value: (your Access Key ID from .csv)
5. Click "Add secret"
6. Click "New repository secret" again
7. Name: `AWS_SECRET_ACCESS_KEY`
   Value: (your Secret Access Key from .csv)
8. Click "Add secret"

### Step 2.3: Push Project Code
```bash
# Extract the taskmaster-starter.zip
cd taskmaster-starter

# Initialize git (if not already)
git init
git add .
git commit -m "Initial commit: TaskMaster starter"

# Connect to your GitHub repo
git remote add origin https://github.com/YOUR_USERNAME/taskmaster.git
git branch -M main
git push -u origin main
```

---

## PHASE 3: LOCAL DEVELOPMENT SETUP (10 min)

### Step 3.1: Verify Docker
```bash
docker --version
docker compose version
```
Both should show version numbers. If not, install Docker Desktop.

### Step 3.2: Run Local Stack
```bash
# In project root
./scripts/setup.sh
```

Wait for output:
```
✅ Setup complete!
📍 Access points:
   API:      http://localhost:8080
   Frontend: http://localhost:3000
```

### Step 3.3: Verify Everything Works
```bash
# Test health
curl http://localhost:8080/health

# Test login
curl -X POST http://localhost:8080/api/v1/auth/login   -H "Content-Type: application/json"   -d '{"email":"admin@taskmaster.com","password":"password123"}'

# Open browser → http://localhost:3000
# Login with admin@taskmaster.com / password123
```

### Step 3.4: Run Tests
```bash
make test
```
Should show all tests passing.

---

## PHASE 4: FIRST AWS DEPLOYMENT (15 min)

### Step 4.1: Build & Push First Docker Image
ECS needs an image in ECR before it can start. Do this once manually.

```bash
# Login to ECR
aws ecr get-login-password --region us-east-1 |   docker login --username AWS --password-stdin   $(aws sts get-caller-identity --query Account --output text).dkr.ecr.us-east-1.amazonaws.com

# Get your AWS account ID
ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)

# Build image
docker build -t taskmaster-api ./backend

# Tag image
docker tag taskmaster-api:latest   ${ACCOUNT_ID}.dkr.ecr.us-east-1.amazonaws.com/taskmaster/api:latest

# Push image
docker push   ${ACCOUNT_ID}.dkr.ecr.us-east-1.amazonaws.com/taskmaster/api:latest
```

### Step 4.2: Deploy Infrastructure with Terraform
```bash
cd infrastructure/environments/dev

# Initialize Terraform (downloads providers)
terraform init

# Preview changes
terraform plan

# Apply (creates everything: VPC, RDS, ElastiCache, ECS, ALB, S3, CloudFront)
terraform apply
```

Type `yes` when prompted.

Wait 10-15 minutes. Terraform will output:
- `alb_dns` — your API endpoint
- `db_endpoint` — PostgreSQL endpoint
- `redis_endpoint` — Redis endpoint
- `cloudfront_domain` — your frontend URL

### Step 4.3: Verify Deployment
```bash
# Test health on AWS
curl http://$(terraform output -raw alb_dns)/health

# Should return: {"status":"ok","time":"..."}
```

### Step 4.4: Configure Frontend for Production
Update `frontend/vite.config.ts` production build to point to your ALB DNS.

Then build and upload to S3:
```bash
cd frontend
npm install
npm run build

# Upload to S3
aws s3 sync dist/ s3://$(cd ../infrastructure && terraform output -raw s3_bucket)

# Invalidate CloudFront cache
aws cloudfront create-invalidation   --distribution-id $(cd ../infrastructure && terraform output -raw cloudfront_id)   --paths "/*"
```

---

## PHASE 5: ENABLE CI/CD (5 min)

### Step 5.1: Verify GitHub Actions Workflow
The `.github/workflows/ci.yml` is already in your repo. It triggers on every push to `main`.

### Step 5.2: Test CI/CD
Make a small change and push:
```bash
# Edit any file, e.g., backend/internal/handler/auth_handler.go
# Add a comment or change a log message

git add .
git commit -m "Test CI/CD pipeline"
git push origin main
```

### Step 5.3: Watch Pipeline Run
1. Go to your GitHub repo
2. Click "Actions" tab
3. You should see a workflow running
4. Click it to see real-time logs:
   - Lint → Test → Build → Push to ECR → Terraform Apply → ECS Deploy

### Step 5.4: Verify Auto-Deployment
```bash
# After pipeline completes (~5-10 minutes)
curl http://$(cd infrastructure/environments/dev && terraform output -raw alb_dns)/health
```

Should still return `{"status":"ok"}` with your new code changes.

---

## PHASE 6: DOMAIN & SSL (Optional, 10 min)

### Step 6.1: Buy/Use Domain
- Use Route 53 to buy a domain (e.g., `taskmaster.yourname.com`)
- Or use an existing domain

### Step 6.2: Update Terraform
In `infrastructure/modules/route53/main.tf`:
```hcl
resource "aws_route53_record" "api" {
  zone_id = var.zone_id
  name    = "api.taskmaster.com"
  type    = "A"
  alias {
    name                   = module.alb.dns_name
    zone_id                = module.alb.zone_id
    evaluate_target_health = true
  }
}
```

### Step 6.3: Add SSL Certificate
```bash
# Request certificate in ACM
aws acm request-certificate   --domain-name api.taskmaster.com   --validation-method DNS   --region us-east-1
```

Add certificate ARN to ALB listener in Terraform.

---

## COMPLETE FILE CHECKLIST

After setup, your project should have these files:

```
taskmaster/
├── .env                          # Local environment (gitignored)
├── .env.example                  # Template for .env
├── .gitignore                    # Git ignore rules
├── Makefile                      # Common commands
├── README.md                     # Project overview
├── docker-compose.yml            # Local infrastructure
├── backend/
│   ├── cmd/server/main.go        # Entry point
│   ├── go.mod                    # Go dependencies
│   ├── Dockerfile                # API container
│   ├── internal/                 # All Go code
│   ├── migrations/               # SQL migrations
│   └── tests/                    # Unit + integration tests
├── frontend/
│   ├── package.json              # Node dependencies
│   ├── vite.config.ts            # Build config
│   ├── Dockerfile                # Frontend container
│   ├── nginx.conf                # Production nginx
│   └── src/                      # React components
├── infrastructure/
│   ├── main.tf                   # Root Terraform
│   ├── variables.tf              # Input variables
│   ├── outputs.tf                # Output values
│   ├── modules/                  # AWS modules
│   └── environments/             # dev / prod
├── .github/
│   └── workflows/
│       └── ci.yml                # CI/CD pipeline
├── scripts/
│   ├── setup.sh                  # One-command local setup
│   ├── seed.sh                   # Database seeding
│   └── test.sh                   # Test runner
└── docs/
    ├── SETUP.md                  # This guide
    ├── DEPLOYMENT.md             # AWS deployment details
    ├── ARCHITECTURE.md           # Clean Architecture docs
    ├── API_REFERENCE.md          # API documentation
    ├── TROUBLESHOOTING.md        # Common issues
    ├── CONTRIBUTING.md           # Development workflow
    └── TaskMaster-API.postman_collection.json
```

---

## WHAT'S ALREADY BUILT (✅)

### Backend
- [x] Go project structure with Clean Architecture
- [x] PostgreSQL connection with connection pooling (pgxpool)
- [x] Redis connection (go-redis)
- [x] Domain entities: User, Task, Team
- [x] Repository interfaces (ports)
- [x] PostgreSQL implementations: UserRepo, TaskRepo, TeamRepo
- [x] Redis implementation: CacheRepo
- [x] Auth service: Register, Login, JWT generation (bcrypt + HS256)
- [x] Task service: Create, List, Get, Delete with business rules
- [x] Business rule validation: Team membership, task limits, SLA checks
- [x] Gin HTTP handlers: Auth, Task, WebSocket
- [x] Middleware chain: Logger, CORS, JWT, RBAC, Rate Limiting
- [x] WebSocket Hub with client management (ping/pong, broadcast)
- [x] Graceful shutdown (DB, Redis, WS connections)
- [x] AWS Secrets Manager integration (production)
- [x] Unit tests with testify/mock for Auth and Task services
- [x] Integration test skeleton
- [x] Database migrations (001_init, 002_seed)
- [x] Multi-stage Dockerfile (builder + Alpine runtime)

### Frontend
- [x] React 18 + TypeScript + Vite scaffold
- [x] AuthContext with JWT persistence (localStorage)
- [x] Login/Register page with toggle
- [x] Dashboard with task list grid
- [x] Task creation form (admin/manager only)
- [x] Task deletion
- [x] Real-time WebSocket integration
- [x] Role-based UI rendering
- [x] Responsive CSS with CSS variables
- [x] Vite proxy for API during development
- [x] Production Dockerfile with nginx
- [x] TypeScript type definitions

### Infrastructure
- [x] Terraform root configuration
- [x] VPC module (VPC, subnets, IGW, NAT GW)
- [x] RDS module (PostgreSQL, security groups)
- [x] ElastiCache module (Redis cluster)
- [x] ALB module (Application Load Balancer, target groups)
- [x] ECS module (Fargate cluster, task definitions, service)
- [x] S3 + CloudFront module (frontend hosting)
- [x] IAM module (ECS task roles, Secrets Manager access)
- [x] Secrets module (AWS Secrets Manager, random passwords)
- [x] Environment-specific configs (dev / prod)
- [x] Terraform outputs (ALB DNS, DB endpoint, etc.)

### DevOps
- [x] Docker Compose for local development
- [x] Multi-stage Dockerfiles (backend + frontend)
- [x] GitHub Actions CI/CD pipeline
- [x] Makefile with common commands
- [x] Setup script (./scripts/setup.sh)
- [x] Seed script for demo data
- [x] Postman API collection
- [x] Complete documentation (7 markdown files)

---

## WHAT REMAINS TO BE DEVELOPED (🔨)

### Phase 2: Frontend Polish (2-3 days)
- [ ] **Task detail page** — View full task with comments/history
- [ ] **Task editing** — Update title, description, status, priority
- [ ] **Task assignment UI** — Dropdown to assign to team members
- [ ] **Team management page** — Create teams, invite members
- [ ] **User profile page** — Update name, password, avatar
- [ ] **Notifications** — Toast messages for WS events
- [ ] **Loading states** — Skeleton screens, spinners
- [ ] **Error boundaries** — Graceful error handling
- [ ] **Responsive mobile design** — Current CSS is desktop-first
- [ ] **Dark mode toggle** — CSS variables make this easy
- [ ] **Search & filters** — Filter tasks by status, priority, assignee
- [ ] **Pagination** — Current list loads all tasks
- [ ] **Drag & drop** — Kanban board view for tasks

### Phase 3: AWS Production Hardening (2-3 days)
- [ ] **Route53 module** — DNS records, custom domain
- [ ] **ACM certificate** — SSL/TLS for ALB and CloudFront
- [ ] **CloudWatch alarms** — CPU, memory, error rate alerts
- [ ] **CloudWatch dashboards** — Visual metrics
- [ ] **X-Ray tracing** — Request tracing across services
- [ ] **ECS auto-scaling** — Target tracking policies
- [ ] **RDS read replicas** — For read-heavy workloads
- [ ] **RDS backups** — Automated snapshots
- [ ] **ElastiCache cluster mode** — For production Redis
- [ ] **S3 bucket policies** — Restrict public access
- [ ] **WAF (Web Application Firewall)** — DDoS protection
- [ ] **VPC Flow Logs** — Network traffic monitoring

### Phase 4: Advanced Features (3-5 days)
- [ ] **Multi-container WebSocket sync** — Redis Pub/Sub for ECS scaling
- [ ] **Audit logging system** — Track all actions with immutable logs
- [ ] **File attachments** — S3 upload for task attachments
- [ ] **Email notifications** — SES integration for task assignments
- [ ] **Slack/Discord integration** — Webhook notifications
- [ ] **Task templates** — Reusable task patterns
- [ ] **Recurring tasks** — Cron-like task creation
- [ ] **Time tracking** — Log hours spent on tasks
- [ ] **Reports & analytics** — Task completion rates, team velocity
- [ ] **API rate limiting per user** — Currently per IP only
- [ ] **API versioning** — v1, v2 strategy
- [ ] **OpenAPI/Swagger docs** — Auto-generated API documentation

### Phase 5: Testing & Quality (2-3 days)
- [ ] **Integration tests with testcontainers-go** — Real DB in tests
- [ ] **E2E tests with Playwright** — Browser automation
- [ ] **Load testing with k6** — Simulate 1000 concurrent users
- [ ] **Security scanning** — Snyk, Dependabot, Trivy
- [ ] **Code coverage reporting** — Coveralls/Codecov integration
- [ ] **Pre-commit hooks** — Husky + lint-staged
- [ ] **Dependency update automation** — Renovate bot

### Phase 6: DevOps Maturity (2-3 days)
- [ ] **Staging environment** — Separate Terraform workspace
- [ ] **Blue-green deployments** — Zero-downtime strategy
- [ ] **Database migration strategy** — golang-migrate in CI/CD
- [ ] **Rollback automation** — Automatic rollback on health check failure
- [ ] **Infrastructure drift detection** — Terraform plan in scheduled runs
- [ ] **Cost monitoring** — AWS Cost Explorer alerts
- [ ] **Disaster recovery plan** — Cross-region backup strategy

---

## ESTIMATED TIMELINE

| Phase | Duration | Effort |
|-------|----------|--------|
| Current (Phase 1) | ✅ Complete | ~40 hours |
| Phase 2: Frontend | 2-3 days | ~15 hours |
| Phase 3: AWS Hardening | 2-3 days | ~15 hours |
| Phase 4: Advanced Features | 3-5 days | ~25 hours |
| Phase 5: Testing | 2-3 days | ~15 hours |
| Phase 6: DevOps | 2-3 days | ~15 hours |
| **Total Remaining** | **~14-20 days** | **~85 hours** |

---

## PRIORITY ORDER FOR HENNGE APPLICATION

If you're building this for HENNGE internship application, focus on:

1. **Get local setup working** (Phase 1-3 above) — This proves you can run it
2. **Deploy to AWS** (Phase 4 above) — This proves cloud skills
3. **Add WebSocket multi-container sync** — This shows advanced Go concurrency
4. **Write integration tests** — This shows testing discipline
5. **Add audit logging** — This shows security awareness
6. **Document everything** — Already done ✅

The current codebase already demonstrates **every HENNGE requirement**. The remaining work is polish and production hardening.

---

## DAILY WORKFLOW (After Setup)

```bash
# Morning: Start local stack
make dev

# Work on code...
# Backend changes auto-reload with air (install: go install github.com/cosmtrek/air@latest)
# Frontend changes auto-reload with Vite HMR

# Test before commit
make test
make lint

# Commit
git add .
git commit -m "feat: add task detail page"
git push origin main

# Watch CI/CD deploy automatically
# Check: https://github.com/YOUR_USERNAME/taskmaster/actions

# Evening: Stop local stack
make dev-down
```

---

**You're now fully equipped to build, deploy, and extend TaskMaster.**
