output "vpc_id" {
  description = "ID of the sovereign VPC"
  value       = aws_vpc.sovereign_vpc.id
}

output "merkle_bucket_arn" {
  description = "ARN of the WORM-immutable Merkle transparency bucket"
  value       = aws_s3_bucket.merkle_transparency_bucket.arn
}

output "kms_key_arn" {
  description = "ARN of the Canadian sovereign KMS encryption key"
  value       = aws_kms_key.sovereign_cmk.arn
}
