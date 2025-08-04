provider "aws" {
  region = "eu-central-1"
}

data "aws_availability_zones" "available" {}

resource "random_password" "db" {
  length           = 16
  upper            = true
  lower            = true
  numeric          = true
  special          = true
  override_special = "!@#$%^&*()-_=+[]{}<>:?"
}

resource "random_password" "jwt_secret" {
  length  = 16
  upper   = true
  lower   = true
  numeric = true
  special = true
}

locals {
  azs     = data.aws_availability_zones.available.names
  db_name = "${var.name}-db"
}

module "eks" {
  source = "./modules/eks"

  name            = "${var.name}-eks"
  vpc_id          = module.vpc.vpc_id
  private_subnets = module.vpc.private_subnets
  node_group_size = 2
}

module "db" {
  source = "./modules/db"

  name            = "${var.name}-db"
  db_username     = var.db_user
  db_password     = random_password.db.result
  private_subnets = module.vpc.private_subnets
  instance_class  = "db.t3.micro"
  vpc_id          = module.vpc.vpc_id
  vpc_cidr        = var.vpc_cidr
}

module "ssm-ps" {
  source = "./modules/ssm-ps"

  name           = "${var.name}-ssm"
  db_name        = local.db_name
  db_url         = module.db.url
  jwt_secret_key = "${var.name}-jwt-secret"
  jwt_secret     = random_password.jwt_secret.result
}

module "vpc" {
  source = "./modules/vpc"

  name     = "${var.name}-vpc"
  azs      = slice(local.azs, 0, 2)
  vpc_cidr = var.vpc_cidr
}