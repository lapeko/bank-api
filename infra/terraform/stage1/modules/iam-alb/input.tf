variable "app_name" {
  type = string
}

variable "oidc_issuer" {
  type = string
}

variable "oidc_provider_arn" {
  type = string
}

variable "tags" {
  type    = map(string)
  default = {}
}