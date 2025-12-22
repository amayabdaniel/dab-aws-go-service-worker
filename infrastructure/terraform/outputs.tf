output "ecr_api_repository_url" {
  description = "URL of the API ECR repository"
  value       = aws_ecr_repository.api.repository_url
}

output "ecr_worker_repository_url" {
  description = "URL of the Worker ECR repository"
  value       = aws_ecr_repository.worker.repository_url
}

output "ecr_frontend_repository_url" {
  description = "URL of the Frontend ECR repository"
  value       = aws_ecr_repository.frontend.repository_url
}

output "ecs_cluster_name" {
  description = "Name of the ECS cluster"
  value       = aws_ecs_cluster.main.name
}

output "alb_dns_name" {
  description = "DNS name of the Application Load Balancer"
  value       = aws_lb.main.dns_name
}

output "rds_endpoint" {
  description = "RDS instance endpoint"
  value       = aws_db_instance.main.endpoint
  sensitive   = true
}

output "bastion_instance_id" {
  description = "Bastion host instance ID for SSM port forwarding"
  value       = aws_instance.bastion.id
}

output "bastion_public_ip" {
  description = "Bastion host Elastic IP"
  value       = aws_eip.bastion.public_ip
}
