# Operational Runbook

## Overview
This runbook provides procedures for common operational tasks and incident response for the Job Management Platform (API, Worker, and Frontend services).

## Production Environment

### Live URLs
- **Dashboard**: https://app.novaferi.net
- **API Health**: https://app.novaferi.net/api/health
- **API Jobs**: https://app.novaferi.net/api/jobs

### AWS Resources
- **Region**: us-east-2 (Ohio)
- **ECS Cluster**: `dab-job-platform-c1bfa7bb-cluster`
- **Services**:
  - `dab-job-platform-c1bfa7bb-api`
  - `dab-job-platform-c1bfa7bb-worker`
  - `dab-job-platform-c1bfa7bb-frontend`
- **RDS Instance**: `dab-job-platform-c1bfa7bb-db`
- **SQS Queue**: `dab-job-platform-c1bfa7bb-jobs-queue`
- **ALB**: `dab-job-platform-c1bfa7bb-alb`

## Local Development

### Starting All Services
```bash
# Start all services with Docker Compose
docker-compose up -d --build

# Create SQS queue in LocalStack
docker-compose exec localstack awslocal sqs create-queue --queue-name jobs-queue
```

### Accessing Services
- Frontend: http://localhost:3000
- API: http://localhost:8080
- PostgreSQL: localhost:5432 (user: devuser, pass: devpass, db: jobsdb)
- LocalStack: http://localhost:4566

## Service Health Checks

### API Health
```bash
# Production
curl https://app.novaferi.net/api/health

# Local
curl http://localhost:8080/health
```
Expected response:
```json
{"service":"api","status":"healthy"}
```

### Full Infrastructure Health Check
```bash
# Check all ECS services
aws ecs describe-services \
  --cluster dab-job-platform-c1bfa7bb-cluster \
  --services dab-job-platform-c1bfa7bb-api dab-job-platform-c1bfa7bb-worker dab-job-platform-c1bfa7bb-frontend \
  --region us-east-2 \
  --query 'services[].{Service:serviceName,Running:runningCount,Desired:desiredCount,Status:status}'

# Check target group health
aws elbv2 describe-target-health \
  --target-group-arn arn:aws:elasticloadbalancing:us-east-2:018491521563:targetgroup/dab-job-platform-api-tg/82192b1f7ae46ae5 \
  --region us-east-2

# Check RDS status
aws rds describe-db-instances \
  --db-instance-identifier dab-job-platform-c1bfa7bb-db \
  --region us-east-2 \
  --query 'DBInstances[0].{Status:DBInstanceStatus}'

# Check SQS queue depth
aws sqs get-queue-attributes \
  --queue-url https://sqs.us-east-2.amazonaws.com/018491521563/dab-job-platform-c1bfa7bb-jobs-queue \
  --attribute-names ApproximateNumberOfMessages ApproximateNumberOfMessagesNotVisible \
  --region us-east-2
```

### Database Health
```bash
aws rds describe-db-instances --db-instance-identifier dab-job-platform-c1bfa7bb-db --region us-east-2
```

## Common Operations

### 1. Scaling Services

#### Scale API Service
```bash
aws ecs update-service \
  --cluster dab-job-platform-c1bfa7bb-cluster \
  --service dab-job-platform-c1bfa7bb-api \
  --desired-count 3 \
  --region us-east-2
```

#### Scale Worker Service
```bash
aws ecs update-service \
  --cluster dab-job-platform-c1bfa7bb-cluster \
  --service dab-job-platform-c1bfa7bb-worker \
  --desired-count 5 \
  --region us-east-2
```

#### Scale Frontend Service
```bash
aws ecs update-service \
  --cluster dab-job-platform-c1bfa7bb-cluster \
  --service dab-job-platform-c1bfa7bb-frontend \
  --desired-count 2 \
  --region us-east-2
```

### 2. Deploy New Version

#### Manual Deployment (Current Process)
```bash
# Login to ECR
aws ecr get-login-password --region us-east-2 | docker login --username AWS --password-stdin 018491521563.dkr.ecr.us-east-2.amazonaws.com

# Build and push API
docker build -t 018491521563.dkr.ecr.us-east-2.amazonaws.com/dab-job-platform-c1bfa7bb-api:latest -f deployments/docker/api.Dockerfile .
docker push 018491521563.dkr.ecr.us-east-2.amazonaws.com/dab-job-platform-c1bfa7bb-api:latest

# Build and push Worker
docker build -t 018491521563.dkr.ecr.us-east-2.amazonaws.com/dab-job-platform-c1bfa7bb-worker:latest -f deployments/docker/worker.Dockerfile .
docker push 018491521563.dkr.ecr.us-east-2.amazonaws.com/dab-job-platform-c1bfa7bb-worker:latest

# Build and push Frontend
docker build -t 018491521563.dkr.ecr.us-east-2.amazonaws.com/dab-job-platform-c1bfa7bb-frontend:latest -f frontend/Dockerfile frontend/
docker push 018491521563.dkr.ecr.us-east-2.amazonaws.com/dab-job-platform-c1bfa7bb-frontend:latest

# Force ECS deployments
aws ecs update-service --cluster dab-job-platform-c1bfa7bb-cluster --service dab-job-platform-c1bfa7bb-api --force-new-deployment --region us-east-2
aws ecs update-service --cluster dab-job-platform-c1bfa7bb-cluster --service dab-job-platform-c1bfa7bb-worker --force-new-deployment --region us-east-2
aws ecs update-service --cluster dab-job-platform-c1bfa7bb-cluster --service dab-job-platform-c1bfa7bb-frontend --force-new-deployment --region us-east-2
```

### 3. Database Operations

#### Connect to Database
```bash
# Get RDS endpoint
aws rds describe-db-instances \
  --db-instance-identifier dab-job-platform-c1bfa7bb-db \
  --region us-east-2 \
  --query 'DBInstances[0].Endpoint.Address' \
  --output text

# Connect (requires bastion or VPN - RDS is in private subnet)
psql -h <rds-endpoint> -U dbadmin -d jobsdb
```

#### Backup Database
```bash
aws rds create-db-snapshot \
  --db-instance-identifier dab-job-platform-c1bfa7bb-db \
  --db-snapshot-identifier manual-snapshot-$(date +%Y%m%d%H%M%S) \
  --region us-east-2
```

#### Run Migrations
Migrations are run automatically by GORM AutoMigrate on service startup.

## Incident Response

### 1. Frontend Not Loading

**Symptoms**: Blank page, 404 errors, API connection failed

**Steps**:
1. Check frontend service status
   ```bash
   aws ecs describe-services \
     --cluster dab-job-platform-c1bfa7bb-cluster \
     --services dab-job-platform-c1bfa7bb-frontend \
     --region us-east-2
   ```
2. Check CloudWatch logs
   ```bash
   aws logs tail /ecs/dab-job-platform-c1bfa7bb-frontend --follow --region us-east-2
   ```
3. Verify target group health
   ```bash
   aws elbv2 describe-target-health \
     --target-group-arn arn:aws:elasticloadbalancing:us-east-2:018491521563:targetgroup/dab-job-platform-fe-tg/959adee3bd543767 \
     --region us-east-2
   ```
4. Force new deployment if needed
   ```bash
   aws ecs update-service --cluster dab-job-platform-c1bfa7bb-cluster --service dab-job-platform-c1bfa7bb-frontend --force-new-deployment --region us-east-2
   ```

### 2. API Not Responding

**Symptoms**: Health check fails, 5xx errors, /api/health returns error

**Steps**:
1. Check ECS task status
   ```bash
   aws ecs describe-services \
     --cluster dab-job-platform-c1bfa7bb-cluster \
     --services dab-job-platform-c1bfa7bb-api \
     --region us-east-2 \
     --query 'services[0].{running:runningCount,desired:desiredCount,events:events[:3]}'
   ```
2. Check CloudWatch logs
   ```bash
   aws logs tail /ecs/dab-job-platform-c1bfa7bb-api --follow --region us-east-2
   ```
3. Check ALB target health
   ```bash
   aws elbv2 describe-target-health \
     --target-group-arn arn:aws:elasticloadbalancing:us-east-2:018491521563:targetgroup/dab-job-platform-api-tg/82192b1f7ae46ae5 \
     --region us-east-2
   ```
4. Check stopped task reason
   ```bash
   TASK=$(aws ecs list-tasks --cluster dab-job-platform-c1bfa7bb-cluster --service-name dab-job-platform-c1bfa7bb-api --desired-status STOPPED --region us-east-2 --query 'taskArns[0]' --output text)
   aws ecs describe-tasks --cluster dab-job-platform-c1bfa7bb-cluster --tasks $TASK --region us-east-2 --query 'tasks[0].stoppedReason'
   ```
5. Force new deployment
   ```bash
   aws ecs update-service --cluster dab-job-platform-c1bfa7bb-cluster --service dab-job-platform-c1bfa7bb-api --force-new-deployment --region us-east-2
   ```

### 3. Worker Not Processing Jobs

**Symptoms**: Queue depth increasing, jobs stuck in pending

**Steps**:
1. Check worker logs
   ```bash
   aws logs tail /ecs/dab-job-platform-c1bfa7bb-worker --follow --region us-east-2
   ```
2. Check SQS queue depth
   ```bash
   aws sqs get-queue-attributes \
     --queue-url https://sqs.us-east-2.amazonaws.com/018491521563/dab-job-platform-c1bfa7bb-jobs-queue \
     --attribute-names ApproximateNumberOfMessages ApproximateNumberOfMessagesNotVisible \
     --region us-east-2
   ```
3. Check dead letter queue
   ```bash
   aws sqs get-queue-attributes \
     --queue-url https://sqs.us-east-2.amazonaws.com/018491521563/dab-job-platform-c1bfa7bb-jobs-dlq \
     --attribute-names ApproximateNumberOfMessages \
     --region us-east-2
   ```
4. Scale workers if needed
   ```bash
   aws ecs update-service --cluster dab-job-platform-c1bfa7bb-cluster --service dab-job-platform-c1bfa7bb-worker --desired-count 2 --region us-east-2
   ```

### 4. Database Connection Errors

**Symptoms**: "connection refused" or timeout errors

**Steps**:
1. Check RDS instance status
   ```bash
   aws rds describe-db-instances --db-instance-identifier production-db
   ```
2. Verify security group rules
3. Check connection pool exhaustion
4. Review CloudWatch metrics for CPU/connections
5. Consider scaling RDS instance

### 5. High Memory Usage

**Symptoms**: ECS tasks being killed, OOM errors

**Steps**:
1. Check CloudWatch Container Insights
2. Review recent code changes for memory leaks
3. Increase task memory allocation
4. Add memory profiling if persistent

## Monitoring and Alerts

### Key Metrics to Monitor

1. **API Metrics**
   - Request rate
   - Error rate (4xx, 5xx)
   - Response time (p50, p95, p99)
   - Active connections

2. **Worker Metrics**
   - Jobs processed/minute
   - Processing time
   - Error rate
   - Queue depth

3. **Infrastructure Metrics**
   - ECS CPU/Memory utilization
   - RDS CPU/connections/storage
   - ALB target health
   - SQS message age

### Alert Thresholds

- API error rate > 1% for 5 minutes
- Worker error rate > 5% for 5 minutes
- Queue depth > 1000 messages
- RDS CPU > 80% for 10 minutes
- ECS task count < desired count for 5 minutes
- Frontend response time > 2s for 5 minutes

## Maintenance Procedures

### Weekly Tasks
1. Review CloudWatch logs for errors
2. Check database backup status
3. Review security group rules
4. Update dependencies if needed

### Monthly Tasks
1. Review AWS costs
2. Analyze performance metrics
3. Update documentation
4. Test disaster recovery

### Quarterly Tasks
1. Security audit
2. Load testing
3. Dependency updates
4. Architecture review

## Emergency Contacts

- On-call Engineer: See PagerDuty
- AWS Support: 1-800-xxx-xxxx
- Database Admin: dba-team@company.com
- Security Team: security@company.com

## Rollback Procedures

### Application Rollback
```bash
# List recent task definitions
aws ecs list-task-definitions --family-prefix api

# Update service with previous version
aws ecs update-service \
  --cluster production-cluster \
  --service api-service \
  --task-definition api:previous-version
```

### Database Rollback
```bash
# Point-in-time recovery
aws rds restore-db-instance-to-point-in-time \
  --source-db-instance-identifier production-db \
  --target-db-instance-identifier production-db-restored \
  --restore-time 2023-12-20T03:00:00.000Z
```

## Useful Commands

```bash
# View real-time API logs
aws logs tail /ecs/dab-job-platform-c1bfa7bb-api --follow --region us-east-2

# View real-time Worker logs
aws logs tail /ecs/dab-job-platform-c1bfa7bb-worker --follow --region us-east-2

# Get all service statuses
aws ecs describe-services \
  --cluster dab-job-platform-c1bfa7bb-cluster \
  --services dab-job-platform-c1bfa7bb-api dab-job-platform-c1bfa7bb-worker dab-job-platform-c1bfa7bb-frontend \
  --region us-east-2 \
  --query 'services[].{Service:serviceName,Running:runningCount,Desired:desiredCount}'

# Check SQS queue
aws sqs get-queue-attributes \
  --queue-url https://sqs.us-east-2.amazonaws.com/018491521563/dab-job-platform-c1bfa7bb-jobs-queue \
  --attribute-names All \
  --region us-east-2

# Force service updates (all services)
aws ecs update-service --cluster dab-job-platform-c1bfa7bb-cluster --service dab-job-platform-c1bfa7bb-api --force-new-deployment --region us-east-2
aws ecs update-service --cluster dab-job-platform-c1bfa7bb-cluster --service dab-job-platform-c1bfa7bb-worker --force-new-deployment --region us-east-2
aws ecs update-service --cluster dab-job-platform-c1bfa7bb-cluster --service dab-job-platform-c1bfa7bb-frontend --force-new-deployment --region us-east-2

# Test endpoints
curl https://app.novaferi.net/api/health
curl https://app.novaferi.net/api/jobs

# Create a test job
curl -X POST https://app.novaferi.net/api/jobs \
  -H "Content-Type: application/json" \
  -d '{"type":"data-processing","data":"Test job from runbook"}'
```

## Infrastructure Teardown

**WARNING**: This will destroy all resources and data!

```bash
cd infrastructure/terraform
tofu destroy
```