variable "db_username" {
  type        = string
  description = "DB username"
}

variable "db_password" {
  type        = string
  description = "DB password"
  sensitive   = true
}

variable "private_subnets" {
  type        = list(string)
}

variable "name" {
  type        = string
  description = "RDS name"
}

variable tags {
  type        = map(string)
  default     = {}
}