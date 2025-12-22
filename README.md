# DAB AWS Go Service Worker

A production-ready microservices platform demonstrating cloud-native DevOps practices with API and Worker services sharing a PostgreSQL database, deployed on AWS ECS Fargate.

## Live Demo

**Production URL**: https://app.novaferi.net

| Endpoint | URL |
|----------|-----|
| Dashboard | https://app.novaferi.net |
| API Health | https://app.novaferi.net/api/health |
| List Jobs | https://app.novaferi.net/api/jobs |

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        Route 53 (DNS)                           │
│                      app.novaferi.net                           │
└─────────────────────────┬───────────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────────┐
│              Application Load Balancer (HTTPS)                  │
│                    ACM Certificate                              │
└────────┬────────────────────────────────────────┬───────────────┘
         │ /api/*                                 │ /*
┌────────▼────────┐                    ┌──────────▼──────────┐
│   API Service   │                    │  Frontend Service   │
│   (Go + Gin)    │                    │  (React + Nginx)    │
│   Port 8080     │                    │     Port 80         │
└────────┬────────┘                    └─────────────────────┘
         │
    ┌────▼────┐
    │   SQS   │◄──────────────────┐
    │  Queue  │                   │
    └────┬────┘                   │
         │                        │
┌────────▼────────┐      ┌────────┴────────┐
│ Worker Service  │      │    Scheduler    │
│  (Go + GORM)    │      │  (Cron Jobs)    │
└────────┬────────┘      └─────────────────┘
         │
┌────────▼────────┐
│   PostgreSQL    │
│   (RDS 16.6)    │
└─────────────────┘
```
<img width="2837" height="1738" alt="image" src="https://github.com/user-attachments/assets/e016b54b-36e4-4531-8ae2-266ea09e413b" />

**Components:**
- **API Service**: RESTful API built with Go 1.23 and Gin framework
- **Worker Service**: Background job processor with integrated scheduler
- **Frontend**: React dashboard with real-time updates (Vite + TypeScript)
- **Database**: PostgreSQL 16.6 on RDS
- **Message Queue**: Amazon SQS with Dead Letter Queue
- **Infrastructure**: OpenTofu/Terraform IaC

## Key Features

- **Production-Ready**: Health checks, graceful shutdown, structured logging
- **Security-First**: GORM ORM (SQL injection prevention), input validation, typed models
- **Local Development**: LocalStack for SQS simulation, Docker Compose
- **Real-Time Updates**: Frontend auto-refreshes job status every 2 seconds
- **Clean Architecture**: Dependency injection, repository pattern, interface-based design
- **CI/CD Pipeline**: GitHub Actions for automated testing and deployment

## Quick Start

### Prerequisites
- Docker & Docker Compose
- Go 1.23+
- AWS CLI configured
- OpenTofu/Terraform 1.5+

### Local Development
```bash
# Clone the repository
git clone git@github.com:amayabdaniel/dab-aws-go-service-worker.git
cd dab-aws-go-service-worker

# Run locally with Docker Compose
docker-compose up -d --build

# Create SQS queue in LocalStack
docker-compose exec localstack awslocal sqs create-queue --queue-name jobs-queue

# Services available at:
# - Frontend: http://localhost:3000
# - API: http://localhost:8080/api/health
# - PostgreSQL: localhost:5432
```

### Running Tests
```bash
go test -v ./...
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check (ALB) |
| GET | `/api/health` | API health check |
| POST | `/api/jobs` | Create a new job |
| GET | `/api/jobs/:id` | Get job by ID |
| GET | `/api/jobs` | List jobs (supports `?status=` filter) |

### Create Job Request
```bash
curl -X POST https://app.novaferi.net/api/jobs \
  -H "Content-Type: application/json" \
  -d '{"type": "data-processing", "data": "sample input"}'
```

### Job Types
| Type | Description |
|------|-------------|
| `data-processing` | General data processing tasks |
| `cleanup` | Remove completed jobs older than 7 days |
| `health-report` | Generate system health metrics |
| `data-aggregation` | Daily statistics aggregation |
| `batch-import` | Process bulk data imports |

### Scheduled Tasks
The worker service runs these automatically:
- **Every 5 minutes**: Cleanup old completed jobs
- **Every hour**: Generate health report
- **Daily at 2 AM UTC**: Perform data aggregation
- **Every 30 seconds**: Check for batch import jobs

## Project Structure
```
.
├── cmd/
│   ├── api/                 # API service entrypoint
│   └── worker/              # Worker service entrypoint
├── internal/
│   ├── api/
│   │   ├── handlers/        # HTTP request handlers
│   │   └── middleware/      # Request validation, error handling
│   ├── database/            # Database connection
│   ├── interfaces/          # Dependency injection interfaces
│   ├── models/              # Data models (Job, JobPayload, etc.)
│   ├── queue/               # SQS client
│   ├── repository/          # Data access layer
│   ├── scheduler/           # Cron-based job scheduler
│   └── worker/              # Job processing logic
├── pkg/
│   ├── config/              # Configuration management
│   └── logger/              # Structured logging (slog)
├── frontend/                # React dashboard (Vite + TypeScript)
├── deployments/docker/      # Dockerfiles
├── infrastructure/terraform/ # IaC configuration
├── tests/                   # Integration tests
└── .github/workflows/       # CI/CD pipelines
```

## Deployment

### AWS Infrastructure
```bash
cd infrastructure/terraform

# Initialize (first time)
tofu init

# Plan and apply
tofu plan
tofu apply
```

### Docker Images
```bash
# Login to ECR
aws ecr get-login-password --region us-east-2 | \
  docker login --username AWS --password-stdin \
  $(aws sts get-caller-identity --query Account --output text).dkr.ecr.us-east-2.amazonaws.com

# Build and push (replace ACCOUNT_ID)
ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
REPO=dab-job-platform-c1bfa7bb

# API
docker build -t $ACCOUNT_ID.dkr.ecr.us-east-2.amazonaws.com/$REPO-api:latest \
  -f deployments/docker/api.Dockerfile .
docker push $ACCOUNT_ID.dkr.ecr.us-east-2.amazonaws.com/$REPO-api:latest

# Worker
docker build -t $ACCOUNT_ID.dkr.ecr.us-east-2.amazonaws.com/$REPO-worker:latest \
  -f deployments/docker/worker.Dockerfile .
docker push $ACCOUNT_ID.dkr.ecr.us-east-2.amazonaws.com/$REPO-worker:latest

# Frontend
docker build -t $ACCOUNT_ID.dkr.ecr.us-east-2.amazonaws.com/$REPO-frontend:latest \
  -f frontend/Dockerfile frontend/
docker push $ACCOUNT_ID.dkr.ecr.us-east-2.amazonaws.com/$REPO-frontend:latest

# Deploy to ECS
aws ecs update-service --cluster $REPO-cluster --service $REPO-api \
  --force-new-deployment --region us-east-2
aws ecs update-service --cluster $REPO-cluster --service $REPO-worker \
  --force-new-deployment --region us-east-2
aws ecs update-service --cluster $REPO-cluster --service $REPO-frontend \
  --force-new-deployment --region us-east-2
```

### Deployed Infrastructure

| Component | Resource |
|-----------|----------|
| Region | us-east-2 (Ohio) |
| ECS Cluster | dab-job-platform-c1bfa7bb-cluster |
| API Service | dab-job-platform-c1bfa7bb-api |
| Worker Service | dab-job-platform-c1bfa7bb-worker |
| Frontend Service | dab-job-platform-c1bfa7bb-frontend |
| Database | RDS PostgreSQL 16.6 (db.t3.micro) |
| Queue | SQS with Dead Letter Queue |
| Load Balancer | Application Load Balancer |
| SSL Certificate | ACM (app.novaferi.net) |
| DNS | Route 53 |

## CI/CD

### Pull Requests (CI)
- Run Go tests with race detection
- Build Docker images (API, Worker, Frontend)
- Run Terraform plan

### Merge to Master (CD)
- Apply Terraform changes
- Build and push Docker images to ECR
- Deploy to ECS with rolling update
- Health check verification

### Required GitHub Secrets

| Secret | Description |
|--------|-------------|
| `AWS_ACCESS_KEY_ID` | AWS access key |
| `AWS_SECRET_ACCESS_KEY` | AWS secret key |
| `TF_VAR_RDS_PASSWORD` | RDS database password (must match deployed value) |

### IAM Permissions Required
- `ecr:*` - Docker image management
- `ecs:*` - Service deployment
- `s3:*` on terraform state bucket
- `dynamodb:*` on lock table
- `rds:*`, `sqs:*`, `ec2:*`, `iam:*`, `logs:*`

## Monitoring

- **CloudWatch Logs**: Application logs from all ECS services
- **CloudWatch Metrics**: CPU, memory, request counts
- **ECS Exec**: SSH into running containers for debugging

### View Logs
```bash
# API logs
aws logs tail /ecs/dab-job-platform-c1bfa7bb-api --follow --region us-east-2

# Worker logs
aws logs tail /ecs/dab-job-platform-c1bfa7bb-worker --follow --region us-east-2
```

### ECS Exec (Debug)
```bash
aws ecs execute-command \
  --cluster dab-job-platform-c1bfa7bb-cluster \
  --task <task-id> \
  --container api \
  --interactive \
  --command "/bin/sh"
```

## Security

- Services run in private subnets (NAT Gateway for outbound)
- RDS accessible only from ECS security group
- TLS/SSL termination at ALB (HTTPS via ACM certificate)
- IAM roles with least privilege
- Database credentials passed via ECS task definition environment variables
- Sensitive terraform variables stored in GitHub Secrets

## License

This project is licensed under the GNU Affero General Public License v3.0 - see the [LICENSE](LICENSE) file for details.

**COPYRIGHT 2025 Daniel Amaya Buitrago. ALL RIGHTS RESERVED.**

This code is submitted for technical assessment purposes only. For permissions, contact: daniel.amaya.buitrago@outlook.com
