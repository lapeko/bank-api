output "url" {
  value     = local.db_url
  sensitive = true
}

output "db_name" {
  value = local.db_name
}