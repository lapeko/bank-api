module "alb" {
  source = "./modules/alb"

  ingress_name      = "${var.name}-ingress"
  ingress_namespace = var.k8s_ingress_namespace
  domain            = var.domain
  cluster_name      = var.cluster_name
}

module "alb-iam" {
  source = "./modules/alb-iam"

  name        = var.name
  oidc_issuer = var.oidc_issuer
}
