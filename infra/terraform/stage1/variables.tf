variable "name" {
  type        = string
  default     = "bank-api"
  description = "App name"
}

variable "db_user" {
  type    = string
  default = "bank-api"
}

variable "vpc_cidr" {
  type    = string
  default = "10.0.0.0/16"
}
