output vpc_id {
  value       = aws_vpc.this.id
}

variable "vpc_cidr" {
  type        = string
  description = "VPC CIDR block"
}

output "public_subnets" {
  value = aws_subnet.public[*].id
}

output "private_subnets" {
  value = aws_subnet.private[*].id
}