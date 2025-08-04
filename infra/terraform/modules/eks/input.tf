variable "name" {
  type    = string
}

variable "vpc_id" {
  type = string
}

variable "private_subnets" {
  type = list(string)
}

variable name {
    type        = string
    description = "EKS name"
}

variable "node_group_size" {
  type    = number
  default = 2
}

variable tags {
  type        = map(string)
  default     = {}
}
