locals {
  issuer = aws_eks_cluster.this.identity[0].oidc[0].issuer
}

data "aws_iam_openid_connect_provider" "this" {
  url = local.issuer
}

output "cluster_name" {
  value = aws_eks_cluster.this.name
}

output "cluster_endpoint" {
  value = aws_eks_cluster.this.endpoint
}

output "cluster_ca" {
  value = aws_eks_cluster.this.certificate_authority[0].data
}

output "oidc_issuer" {
  value = replace(local.issuer, "https://", "")
}

output "oidc_provider_arn" {
  value = data.aws_iam_openid_connect_provider.this.arn
}
