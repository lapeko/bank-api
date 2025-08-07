variable "oidc_issuer" {
  type = string
}

variable "cluster_name" {
  type = string
}

variable "name" {
  type        = string
  default     = "bank-api"
  description = "App name"
}

variable "domain" {
  type    = string
  default = "vlapeka.click"
}

variable "k8s_ingress_namespace" {
  type    = string
  default = "ingress"
}
