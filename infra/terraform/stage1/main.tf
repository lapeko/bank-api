data "aws_availability_zones" "available" {}

resource "random_password" "db" {
  length           = 16
  override_special = "!#%^&*()_+-=[]{}<>:"
}

resource "random_password" "jwt_secret" {
  length = 16
}

locals {
  azs      = data.aws_availability_zones.available.names
  app_name = "${var.app_name}-db"
}

module "db" {
  source = "./modules/db"

  app_name        = var.app_name
  db_username     = var.db_user
  db_password     = random_password.db.result
  private_subnets = module.vpc.private_subnets
  instance_class  = "db.t3.micro"
  vpc_id          = module.vpc.id
  vpc_cidr        = var.vpc_cidr
}

module "eks" {
  source = "./modules/eks"

  app_name        = var.app_name
  vpc_id          = module.vpc.id
  private_subnets = module.vpc.private_subnets
  node_group_size = 2
}

module "iam-alb" {
  source = "./modules/iam-alb"

  app_name          = var.app_name
  oidc_issuer       = module.eks.oidc_issuer
  oidc_provider_arn = module.eks.oidc_provider_arn
}

module "iam-external-secret" {
  source = "./modules/iam-external-secret"

  app_name          = var.app_name
  oidc_issuer       = module.eks.oidc_issuer
  oidc_provider_arn = module.eks.oidc_provider_arn
}

module "ssm-ps" {
  source = "./modules/ssm-ps"

  app_name   = var.app_name
  db_url     = module.db.url
  jwt_secret = random_password.jwt_secret.result
}

module "vpc" {
  source = "./modules/vpc"

  app_name = var.app_name
  azs      = slice(local.azs, 0, 2)
  vpc_cidr = var.vpc_cidr
}