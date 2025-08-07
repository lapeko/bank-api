variable "name" {
  type        = string
  description = "EKS name"
}

variable "vpc_id" {
  type = string
}

variable "private_subnets" {
  type = list(string)
}

variable "node_group_size" {
  type    = number
  default = 2
}

variable "tags" {
  type    = map(string)
  default = {}
}
