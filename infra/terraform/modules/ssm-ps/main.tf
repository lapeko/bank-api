resource "aws_ssm_parameter" "db_url" {
  name        = "/${var.db_name}/db_url"
  type        = "SecureString"
  value       = var.db_url
  overwrite   = true

  tags = merge(var.tags, {
    Name = "${var.name}-db-url"
  })
}

resource "aws_ssm_parameter" "jwt_secret" {
  name        = "/${var.jwt_secret_key}/jwt-secret-key"
  type        = "SecureString"
  value       = var.jwt_secret
  overwrite   = true

  tags = merge(var.tags, {
    Name = "${var.name}-jwt-secret"
  })
}