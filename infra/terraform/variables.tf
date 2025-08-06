variable "name" {
  type        = string
  default     = "bank-api"
  description = "App name"
}

variable "db_user" {
  type    = string
  default = "bank_admin"
}

variable "vpc_cidr" {
  type    = string
  default = "10.0.0.0/16"
}

variable "domain" {
  type    = string
  default = "vlapeka.click"
}

variable "k8s_ingress_namespace" {
  type    = string
  default = "ingress"
}
