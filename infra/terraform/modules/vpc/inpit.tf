variable "name" {
  type        = string
  description = "Project/cluster name prefix"
}

variable "azs" {
  type        = list(string)
  description = "Availability zones"
}

variable "vpc_cidr" {
  type        = string
  description = "VPC CIDR block"
}