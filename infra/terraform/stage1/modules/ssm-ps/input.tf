variable "db_name" {
  type        = string
  description = "Will be used in SSM parameter name"
}

variable "db_url" {
  type      = string
  sensitive = true
}

variable "jwt_secret_key" {
  type        = string
  description = "Will be used in SSM parameter name"
}

variable "jwt_secret" {
  type      = string
  sensitive = true
}

variable "name" {
  type        = string
  description = "SSM name"
}

variable "tags" {
  type    = map(string)
  default = {}
}