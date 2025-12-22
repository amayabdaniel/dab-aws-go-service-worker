terraform {
  backend "s3" {
    bucket         = "dab-terraform-state-bucket"
    key            = "job-platform/terraform.tfstate"
    region         = "us-east-2"
    encrypt        = true
    dynamodb_table = "dab-terraform-state-lock"
  }
}