variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = ""
}

variable "environment" {
  description = "Environment name"
  type        = string
  default     = ""
}

variable "project_name" {
  description = "Project name"
  type        = string
  default     = ""
}

variable "vpc_cidr" {
  description = "CIDR block for VPC"
  type        = string
  default     = ""
}

variable "availability_zones" {
  description = "List of availability zones"
  type        = list(string)
  default     = []
}

variable "rds_instance_class" {
  description = "RDS instance class"
  type        = string
  default     = ""
}

variable "rds_allocated_storage" {
  description = "RDS allocated storage in GB"
  type        = number
  default     = 0
}

variable "rds_username" {
  description = "RDS master username"
  type        = string
  default     = ""
  sensitive   = true
}

variable "rds_password" {
  description = "RDS master password"
  type        = string
  default     = ""
  sensitive   = true
}

variable "ecs_task_cpu" {
  description = "ECS task CPU units"
  type        = string
  default     = ""
}

variable "ecs_task_memory" {
  description = "ECS task memory in MB"
  type        = string
  default     = ""
}

variable "api_desired_count" {
  description = "Desired count for API service"
  type        = number
  default     = 0
}

variable "worker_desired_count" {
  description = "Desired count for Worker service"
  type        = number
  default     = 0
}

variable "frontend_desired_count" {
  description = "Desired count for Frontend service"
  type        = number
  default     = 0
}

variable "enable_multi_az" {
  description = "Enable Multi-AZ for RDS"
  type        = bool
  default     = false
}

variable "enable_container_insights" {
  description = "Enable Container Insights for ECS"
  type        = bool
  default     = false
}

variable "enable_nat_gateway" {
  description = "Enable NAT Gateway"
  type        = bool
  default     = true
}

variable "log_retention_days" {
  description = "CloudWatch log retention in days"
  type        = number
  default     = 0
}