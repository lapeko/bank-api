data "aws_availability_zones" "available" {}

resource "random_password" "db" {
  length           = 16
  special          = false
  override_special = "!#%^&*()_+-=[]{}<>:"
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

resource "time_sleep" "wait_for_alb" {
  depends_on      = [module.k8s]
  create_duration = "120s"
}

module "alb" {
  depends_on = [time_sleep.wait_for_alb]

  source = "./modules/alb"

  ingress_name      = "${var.name}-ingress"
  ingress_namespace = var.k8s_ingress_namespace
  domain            = var.domain
  cluster_name      = module.eks.cluster_name
}

module "alb-iam" {
  source = "./modules/alb-iam"

  name        = var.name
  oidc_issuer = module.eks.oidc_issuer
}

module "db" {
  source = "./modules/db"

  name            = var.name
  db_username     = var.db_user
  db_password     = random_password.db.result
  private_subnets = module.vpc.private_subnets
  instance_class  = "db.t3.micro"
  vpc_id          = module.vpc.id
  vpc_cidr        = var.vpc_cidr
}

module "eks" {
  source = "./modules/eks"

  name            = "${var.name}-eks"
  vpc_id          = module.vpc.id
  private_subnets = module.vpc.private_subnets
  node_group_size = 2
}

module "k8s" {
  source = "./modules/k8s"

  secret_name        = "${var.name}-secret"
  secret_namespace   = var.name
  secret_password    = random_password.db.result
  secret_jwt         = random_password.jwt_secret.result
  cluster_name       = module.eks.cluster_name
  alb_controller_arn = module.alb-iam.alb_controller_arn
  vpc_id             = module.vpc.id
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