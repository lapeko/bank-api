resource "kubernetes_namespace" "this" {
  metadata {
    name = var.secret_namespace
  }
}

resource "kubernetes_secret" "this" {
  metadata {
    name      = var.secret_name
    namespace = var.secret_namespace
  }
  data = {
    DB_PASSWORD = base64encode(var.secret_password)
    JWT_SECRET  = base64encode(var.secret_jwt)
  }
  type = "Opaque"
}

resource "kubernetes_service_account" "alb_controller" {
  metadata {
    name      = "aws-load-balancer-controller"
    namespace = "kube-system"
    annotations = {
      "eks.amazonaws.com/role-arn" = var.alb_controller_arn
    }
  }
}

resource "helm_release" "alb_controller" {
  depends_on = [kubernetes_service_account.alb_controller]

  name       = "aws-load-balancer-controller"
  namespace  = "kube-system"
  repository = "https://aws.github.io/eks-charts"
  chart      = "aws-load-balancer-controller"
  version    = "1.7.1"

  timeout       = 600
  wait          = true
  recreate_pods = true

  set = [
    {
      name  = "clusterName"
      value = var.cluster_name
    },
    {
      name  = "vpcId"
      value = var.vpc_id
    },
    {
      name  = "serviceAccount.create"
      value = "false"
    },
    {
      name  = "serviceAccount.name"
      value = kubernetes_service_account.alb_controller.metadata[0].name
    }
  ]
}
