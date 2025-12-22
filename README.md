# DAB AWS Go Service Worker

A production-ready microservices platform demonstrating cloud-native DevOps practices with API and Worker services sharing a PostgreSQL database, deployed on AWS ECS Fargate.

## 🌐 Live Demo

**Production URL**: https://app.novaferi.net

| Endpoint | URL |
|----------|-----|
| Dashboard | https://app.novaferi.net |
| API Health | https://app.novaferi.net/api/health |
| List Jobs | https://app.novaferi.net/api/jobs |

## ⚖️ Legal Notice

**COPYRIGHT © 2025 Daniel Amaya Buitrago. ALL RIGHTS RESERVED.**

This repository is protected under the GNU Affero General Public License v3.0 (AGPL-3.0).

### Restrictions:
- ❌ NO commercial use without written permission
- ❌ NO integration into proprietary systems
- ❌ NO derivative works without attribution
- ❌ NO private/internal use without compliance

### Requirements if used:
- ✅ Must open-source entire application
- ✅ Must include this copyright notice
- ✅ Must disclose all modifications
- ✅ Must provide source code to all users

**This code is submitted for technical assessment purposes only. Any other use is a violation of copyright law.**

For permissions, contact: daniel.amaya.buitrago@outlook.com

## 🏗️ Architecture

- **API Service**: RESTful API built with Go and Gin framework
- **Worker Service**: Background job processor with integrated scheduler
- **Frontend**: React dashboard with real-time updates (Vite + TypeScript)
- **Scheduler**: Cron-based task scheduler for recurring jobs
- **Database**: PostgreSQL for persistent storage
- **Message Queue**: Amazon SQS for asynchronous job processing
- **Container Orchestration**: AWS ECS Fargate
- **Load Balancer**: Application Load Balancer
- **Infrastructure**: Terraform for IaC

## 🎯 Key Features

- **AI-Assisted Development**: Built using Claude Code for enhanced productivity
- **Production-Ready**: Health checks, graceful shutdown, structured logging
- **Security-First**: GORM ORM (no SQL injection), input validation, typed data models
- **Local AWS Development**: LocalStack for SQS simulation
- **Real-Time Updates**: Frontend auto-refreshes job status every 2 seconds
- **SOLID Principles**: Dependency injection, repository pattern, clean architecture

## 🚀 Quick Start

### Prerequisites
- Docker & Docker Compose
- Go 1.21+
- AWS CLI configured
- Terraform 1.5+

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
# - API: http://localhost:8080/health
# - PostgreSQL: localhost:5432
```

### API Endpoints
- `GET /health` - Health check
- `POST /jobs` - Create a new job
- `GET /jobs/{id}` - Get job status
- `GET /jobs` - List all jobs (supports `?status=` filter)

### Job Types
- **data-processing** - General data processing tasks
- **cleanup** - Remove completed jobs older than 7 days
- **health-report** - Generate system health metrics
- **data-aggregation** - Daily statistics aggregation
- **batch-import** - Process bulk data imports

### Scheduled Tasks
The worker service automatically runs:
- **Every 5 minutes**: Cleanup old completed jobs
- **Every hour**: Generate health report
- **Daily at 2 AM**: Perform data aggregation
- **Every 30 seconds**: Check for batch import jobs

## 📁 Project Structure
```
.
├── cmd/                    # Application entry points
│   ├── api/               # API service main
│   └── worker/            # Worker service main
├── internal/              # Private application code
│   ├── api/              # API handlers and middleware
│   ├── models/           # Data models
│   ├── database/         # Database connection and migrations
│   ├── queue/            # SQS client
│   ├── scheduler/        # Cron-based job scheduler
│   ├── repository/       # Data access layer
│   ├── interfaces/       # Dependency injection interfaces
│   └── worker/           # Job processing logic
├── frontend/             # React dashboard
│   ├── src/              # React components and services
│   ├── public/           # Static assets
│   └── nginx.conf        # Nginx configuration
├── pkg/                   # Public packages
│   ├── config/           # Configuration management
│   └── logger/           # Structured logging
├── deployments/          # Deployment configurations
│   └── docker/           # Dockerfiles
├── infrastructure/       # Infrastructure as Code
│   └── terraform/        # Terraform modules and environments
└── tests/                # Test suites
```

## 🔧 Development

### Building Services
```bash
# Build API service
go build -o bin/api cmd/api/main.go

# Build Worker service
go build -o bin/worker cmd/worker/main.go

# Run tests
go test ./...
```

### Environment Variables
See `.env.example` for required configuration.

## 🚢 Deployment

### AWS Infrastructure (OpenTofu/Terraform)

```bash
# Set up backend (first time only)
cd infrastructure/terraform
./setup-backend.sh

# Initialize and deploy
tofu init
tofu plan
tofu apply
```

### Build and Deploy Docker Images

```bash
# Login to ECR
aws ecr get-login-password --region us-east-2 | docker login --username AWS --password-stdin <account-id>.dkr.ecr.us-east-2.amazonaws.com

# Build and push API
docker build -t <account-id>.dkr.ecr.us-east-2.amazonaws.com/<repo>-api:latest -f deployments/docker/api.Dockerfile .
docker push <account-id>.dkr.ecr.us-east-2.amazonaws.com/<repo>-api:latest

# Build and push Worker
docker build -t <account-id>.dkr.ecr.us-east-2.amazonaws.com/<repo>-worker:latest -f deployments/docker/worker.Dockerfile .
docker push <account-id>.dkr.ecr.us-east-2.amazonaws.com/<repo>-worker:latest

# Build and push Frontend
docker build -t <account-id>.dkr.ecr.us-east-2.amazonaws.com/<repo>-frontend:latest -f frontend/Dockerfile frontend/
docker push <account-id>.dkr.ecr.us-east-2.amazonaws.com/<repo>-frontend:latest

# Force ECS deployments
aws ecs update-service --cluster <cluster-name> --service <service-name> --force-new-deployment --region us-east-2
```

### Deployed Infrastructure

| Component | Resource |
|-----------|----------|
| Region | us-east-2 (Ohio) |
| ECS Cluster | dab-job-platform-c1bfa7bb-cluster |
| Database | RDS PostgreSQL 16.6 (db.t3.micro) |
| Queue | SQS with Dead Letter Queue |
| Load Balancer | Application Load Balancer |
| SSL Certificate | ACM (app.novaferi.net) |
| DNS | Route 53 |

### CI/CD with GitHub Actions

**On Pull Requests (CI):**
- Run Go tests
- Build Docker images (API, Worker, Frontend)
- Run Terraform plan

**On Merge to Master (CD):**
- Apply Terraform changes
- Build and push Docker images to ECR
- Deploy to ECS with rolling update
- Health check verification

#### Required GitHub Secrets

Add these secrets in your repository settings (`Settings > Secrets and variables > Actions`):

| Secret | Description |
|--------|-------------|
| `AWS_ACCESS_KEY_ID` | AWS access key with permissions for ECR, ECS, Terraform |
| `AWS_SECRET_ACCESS_KEY` | AWS secret access key |

#### IAM Permissions Required

The AWS credentials need these permissions:
- `ecr:*` - Push/pull Docker images
- `ecs:*` - Update services, describe tasks
- `s3:*` on terraform state bucket - Terraform state
- `dynamodb:*` on lock table - Terraform state locking
- `rds:*`, `sqs:*`, `ec2:*`, `iam:*`, `logs:*` - Infrastructure management

## 📊 Monitoring

- CloudWatch Logs for application logs
- CloudWatch Metrics for system metrics
- X-Ray for distributed tracing
- Custom dashboards for service health

## 🔐 Security

- All services run in private subnets
- Secrets managed via AWS Secrets Manager
- IAM roles with least privilege
- TLS/SSL for all communications
- Regular security scanning in CI/CD

## 📝 License

This project is licensed under the GNU Affero General Public License v3.0 - see the [LICENSE](LICENSE) file for details.

---

© 2025 Daniel Amaya Buitrago. All rights reserved.
