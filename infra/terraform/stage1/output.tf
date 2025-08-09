output "cluster_name" {
  value = module.eks.cluster_name
}

output "vpc_id" {
  value = module.vpc.id
}

output "alb_controller_role_arn" {
  value = module.iam-alb.alb_controller_role_arn
}
