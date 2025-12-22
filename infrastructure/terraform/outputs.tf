output "alb_dns_name" {
  description = "DNS name of the load balancer"
  value       = aws_lb.main.dns_name
}

output "frontend_url" {
  description = "URL for the frontend"
  value       = "http://${aws_lb.main.dns_name}"
}

output "api_url" {
  description = "URL for the API"
  value       = "http://${aws_lb.main.dns_name}/api"
}

output "health_check_url" {
  description = "URL for health check"
  value       = "http://${aws_lb.main.dns_name}/health"
}

output "ecr_repositories" {
  description = "ECR repository URLs"
  value = {
    api      = aws_ecr_repository.api.repository_url
    worker   = aws_ecr_repository.worker.repository_url
    frontend = aws_ecr_repository.frontend.repository_url
  }
}

output "rds_endpoint" {
  description = "RDS instance endpoint"
  value       = aws_db_instance.main.endpoint
  sensitive   = true
}

output "sqs_queue_url" {
  description = "SQS queue URL"
  value       = aws_sqs_queue.jobs.url
}

output "ecs_cluster_name" {
  description = "ECS cluster name"
  value       = aws_ecs_cluster.main.name
}

output "vpc_id" {
  description = "VPC ID"
  value       = aws_vpc.main.id
}

output "deployment_commands" {
  description = "Commands to deploy the application"
  value       = <<-EOT
    # Build and push Docker images:
    aws ecr get-login-password --region ${var.aws_region} | docker login --username AWS --password-stdin ${split("/", aws_ecr_repository.api.repository_url)[0]}
    
    # API
    docker build -t ${aws_ecr_repository.api.repository_url}:latest -f deployments/docker/api.Dockerfile .
    docker push ${aws_ecr_repository.api.repository_url}:latest
    
    # Worker
    docker build -t ${aws_ecr_repository.worker.repository_url}:latest -f deployments/docker/worker.Dockerfile .
    docker push ${aws_ecr_repository.worker.repository_url}:latest
    
    # Frontend
    docker build -t ${aws_ecr_repository.frontend.repository_url}:latest -f frontend/Dockerfile frontend/
    docker push ${aws_ecr_repository.frontend.repository_url}:latest
    
    # Force new deployment:
    aws ecs update-service --cluster ${aws_ecs_cluster.main.name} --service ${local.name_suffix}-api --force-new-deployment
    aws ecs update-service --cluster ${aws_ecs_cluster.main.name} --service ${local.name_suffix}-worker --force-new-deployment
    aws ecs update-service --cluster ${aws_ecs_cluster.main.name} --service ${local.name_suffix}-frontend --force-new-deployment
  EOT
}