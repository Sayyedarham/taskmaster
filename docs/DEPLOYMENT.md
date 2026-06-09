# Deployment Guide

## Local Development

### Quick Start
```bash
# Clone and setup
git clone <your-repo>
cd taskmaster

# Option 1: Use setup script
./scripts/setup.sh

# Option 2: Manual setup
cp .env.example .env
docker-compose up -d

# Verify
make test
curl http://localhost:8080/health
```

### Access Points
- **API**: http://localhost:8080
- **Frontend**: http://localhost:3000
- **Health**: http://localhost:8080/health

### Demo Credentials
- Admin: `admin@taskmaster.com` / `password123`
- Manager: `manager@taskmaster.com` / `password123`
- Member: `member@taskmaster.com` / `password123`

---

## AWS Production Deployment

### Prerequisites
- AWS CLI configured with appropriate credentials
- Terraform >= 1.5
- Docker
- GitHub repository with Actions enabled

### Step 1: AWS Setup

```bash
# Create S3 bucket for Terraform state
aws s3 mb s3://taskmaster-terraform-state --region us-east-1

# Enable versioning on state bucket
aws s3api put-bucket-versioning   --bucket taskmaster-terraform-state   --versioning-configuration Status=Enabled
```

### Step 2: GitHub Secrets

Add these secrets to your GitHub repository:

| Secret | Description |
|--------|-------------|
| `AWS_ACCESS_KEY_ID` | IAM user access key |
| `AWS_SECRET_ACCESS_KEY` | IAM user secret key |

### Step 3: First Deployment

```bash
# Initialize Terraform
cd infrastructure/environments/dev
terraform init

# Plan infrastructure
terraform plan

# Apply infrastructure
terraform apply

# Note: ECS will fail first time because no image exists yet
```

### Step 4: Push Image to ECR

```bash
# Build and push manually first time
aws ecr get-login-password | docker login --username AWS --password-stdin <account>.dkr.ecr.us-east-1.amazonaws.com

docker build -t taskmaster/api:latest ./backend
docker tag taskmaster/api:latest <account>.dkr.ecr.us-east-1.amazonaws.com/taskmaster/api:latest
docker push <account>.dkr.ecr.us-east-1.amazonaws.com/taskmaster/api:latest
```

### Step 5: Update ECS Service

```bash
aws ecs update-service   --cluster taskmaster   --service taskmaster-api   --force-new-deployment
```

### Step 6: Verify Deployment

```bash
# Get ALB DNS
terraform output alb_dns

# Test health endpoint
curl http://<alb-dns>/health

# Test API
curl -X POST http://<alb-dns>/api/v1/auth/register   -H "Content-Type: application/json"   -d '{"email":"test@example.com","password":"password123","name":"Test"}'
```

---

## CI/CD Pipeline

The GitHub Actions workflow handles:

1. **Lint**: `golangci-lint` checks code quality
2. **Unit Tests**: Mocked tests with race detection
3. **Integration Tests**: Real database with testcontainers
4. **Build**: Docker image with multi-stage build
5. **Push**: Image to Amazon ECR
6. **Terraform**: Apply infrastructure changes
7. **Deploy**: Force ECS new deployment
8. **Smoke Tests**: Verify deployment health

### Pipeline Triggers
- Push to `main`: Full CI/CD
- Pull Request: Lint + Test only

---

## Database Migrations

### Local
```bash
# Run migrations
make migrate-up

# Rollback one
make migrate-down

# Create new migration
make migrate-create name=add_notifications
```

### Production
```bash
# Run from CI/CD or bastion host
docker run --rm   -e DATABASE_URL="postgres://..."   migrate/migrate:v4.17.0   -path /migrations   -database "$DATABASE_URL"   up
```

---

## Monitoring & Observability

### CloudWatch Logs
```bash
# View API logs
aws logs tail /ecs/taskmaster-api --follow
```

### CloudWatch Metrics
- ECS CPU/Memory utilization
- RDS connections, CPU, storage
- ALB request count, 4xx/5xx errors
- Custom: Goroutines, GC pauses

### Alarms
- CPU > 70% for 5 min → Scale up
- 5xx errors > 10/min → Alert
- RDS connections > 80% → Alert

---

## Rollback Strategy

### Application Rollback
```bash
# Revert to previous ECS task definition
aws ecs update-service   --cluster taskmaster   --service taskmaster-api   --task-definition taskmaster-api:<PREVIOUS_REVISION>
```

### Database Rollback
```bash
# Rollback one migration
migrate -path migrations -database "$DB_URL" down 1
```

### Infrastructure Rollback
```bash
# Revert Terraform state
cd infrastructure/environments/dev
git checkout <previous-commit>
terraform apply
```

---

## Environment Promotion

```
Dev → Staging → Production
  │       │         │
  ▼       ▼         ▼
PR    Merge    Manual
Test  to main  approval
```

### Promoting to Production
1. Merge PR to `main`
2. CI/CD deploys to dev automatically
3. Run smoke tests against dev
4. Create release tag: `git tag v1.0.0`
5. GitHub Actions deploys to prod
6. Verify production health

---

## Troubleshooting

### ECS Task Won't Start
```bash
# Check task status
aws ecs describe-tasks   --cluster taskmaster   --tasks $(aws ecs list-tasks --cluster taskmaster --query 'taskArns[0]' --output text)

# Check logs
aws logs get-log-events   --log-group-name /ecs/taskmaster-api   --log-stream-name <stream-name>
```

### Database Connection Issues
```bash
# Test from ECS task
aws ecs execute-command   --cluster taskmaster   --task <task-id>   --container api   --interactive   --command "sh -c 'psql $DATABASE_URL -c "SELECT 1"'"
```

### High Latency
- Check CloudWatch X-Ray traces
- Review RDS slow query logs
- Verify Redis cache hit rate
- Check ALB target health

---

## Security Checklist

- [ ] AWS Secrets Manager for all secrets
- [ ] IAM roles with least privilege
- [ ] Security groups restrict traffic
- [ ] RDS not publicly accessible
- [ ] CloudFront HTTPS only
- [ ] JWT secret rotated regularly
- [ ] Rate limiting enabled
- [ ] Input validation on all endpoints
- [ ] Dependency scanning (Snyk, Dependabot)
- [ ] Container scanning (ECR image scanning)
