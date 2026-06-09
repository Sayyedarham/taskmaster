resource "aws_secretsmanager_secret" "app" {
  name = "${var.name}/app-secrets"
}

resource "aws_secretsmanager_secret_version" "app" {
  secret_id = aws_secretsmanager_secret.app.id
  secret_string = jsonencode({
    database_url   = "postgres://taskmaster:${random_password.db_password.result}@${var.db_endpoint}/taskmaster"
    redis_password = random_password.redis_password.result
    jwt_secret     = random_password.jwt_secret.result
  })
}

resource "random_password" "db_password" {
  length  = 32
  special = false
}

resource "random_password" "redis_password" {
  length  = 32
  special = false
}

resource "random_password" "jwt_secret" {
  length  = 64
  special = true
}

output "secret_arn" {
  value = aws_secretsmanager_secret.app.arn
}
