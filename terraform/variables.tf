variable "aws_region" {
  description = "Primary AWS sovereign region (Canada Central)"
  type        = string
  default     = "ca-central-1"
}

variable "environment" {
  description = "Target deployment environment (development, staging, production)"
  type        = string
  default     = "production"
}

variable "vpc_cidr" {
  description = "CIDR block for the sovereign VPC"
  type        = string
  default     = "10.100.0.0/16"
}
