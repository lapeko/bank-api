output "external-secrets-role-arn" {
  value = module.eks.oidc_provider_arn
}