terraform {
  required_version = ">= 1.6.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.40"
    }
  }
  backend "s3" {
    bucket         = "cog-sovereign-tfstate-ca-central-1"
    key            = "production/terraform.tfstate"
    region         = "ca-central-1"
    encrypt        = true
    dynamodb_table = "cog-tfstate-locks"
  }
}

provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      Project     = "CanadaEconomicOpportunityGraph"
      Environment = var.environment
      Classification = "Protected-B"
      DataSovereignty = "Canada-Only"
    }
  }
}

# Sovereign Isolated VPC in Canada Central (Montreal)
resource "aws_vpc" "sovereign_vpc" {
  cidr_block           = var.vpc_cidr
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = {
    Name = "cog-sovereign-vpc-ca-central-1"
  }
}

# S3 Immutable Merkle Transparency Log Bucket with Object Lock (WORM compliance)
resource "aws_s3_bucket" "merkle_transparency_bucket" {
  bucket        = "cog-sovereign-merkle-log-${var.environment}"
  force_destroy = false

  object_lock_enabled = true

  tags = {
    Name        = "Sovereign Merkle Transparency Archive"
    Compliance  = "CEGS-WORM-Immutable"
  }
}

resource "aws_s3_bucket_versioning" "merkle_versioning" {
  bucket = aws_s3_bucket.merkle_transparency_bucket.id
  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_server_side_encryption_configuration" "merkle_encryption" {
  bucket = aws_s3_bucket.merkle_transparency_bucket.id

  rule {
    apply_server_side_encryption_by_default {
      kms_master_key_id = aws_kms_key.sovereign_cmk.arn
      sse_algorithm     = "aws:kms"
    }
  }
}

# Dedicated Canadian Sovereign KMS Customer Managed Key
resource "aws_kms_key" "sovereign_cmk" {
  description             = "Canadian Sovereign Master Encryption Key for CEO-G Lakehouse and Merkle Logs"
  deletion_window_in_days = 30
  enable_key_rotation     = true

  tags = {
    Name = "cog-sovereign-kms-key"
  }
}
