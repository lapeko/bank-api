variable "app_name" {
  type = string
}

variable "azs" {
  type        = list(string)
  description = "Availability zones"
}

variable "vpc_cidr" {
  type        = string
  description = "VPC CIDR block"
}

variable "tags" {
  type    = map(string)
  default = {}
}