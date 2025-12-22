# Infrastructure as Code with OpenTofu

This directory contains the OpenTofu configuration for deploying the Job Management Platform to AWS.

## Architecture

- **VPC**: Multi-AZ setup with public and private subnets
- **ECS Fargate**: Running API, Worker, and Frontend services
- **RDS PostgreSQL**: Multi-AZ database
- **SQS**: Message queue for job processing
- **ALB**: Application Load Balancer with path-based routing
- **ECR**: Container registries for all services

## Prerequisites

1. AWS CLI configured with appropriate credentials
2. OpenTofu installed (v1.11.0+)
3. Docker installed for building images

## Deployment Steps

### 1. Initialize OpenTofu
```bash
tofu init
```

### 2. Review the plan
```bash
tofu plan
```

### 3. Apply infrastructure
```bash
tofu apply
```

### 4. Build and push Docker images
After the infrastructure is created, use the commands from the output:
```bash
# The apply command will output the exact commands to run
tofu output deployment_commands
```

### 5. Access the application
```bash
# Get the application URL
tofu output frontend_url
```

## Cost Optimization

The infrastructure is configured for minimal costs by default:
- Single NAT Gateway (configurable via `enable_nat_gateway`)
- Single-AZ RDS deployment (Multi-AZ can be enabled via `enable_multi_az`)
- Minimal instance sizes (t3.micro for RDS, 256 CPU/512 MB for ECS)
- Container Insights disabled by default (enable via `enable_container_insights`)
- Short log retention (7 days by default)
- Worker service uses Fargate Spot for 70% cost savings
- Single instance of each service by default

## Destroy Infrastructure

To tear down all resources:
```bash
tofu destroy
```

## Configuration

Edit `terraform.tfvars` to customize:
- Instance sizes
- Number of service replicas
- Database configuration
- Region settings

## Security Notes

- All services run in private subnets
- RDS is encrypted at rest
- Security groups follow least privilege
- Change default RDS password before production use