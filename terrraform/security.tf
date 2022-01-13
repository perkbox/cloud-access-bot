# KMS encryption key for Google OAuth Client ID/secret
resource "aws_kms_key" "slack_auth" {


  tags = var.tags
}

resource "aws_secretsmanager_secret" "slack_auth" {

  name       = "cloud_access_bot/slack_auth"
  kms_key_id = aws_kms_key.slack_auth.id
  tags       = var.tags
}


resource "aws_secretsmanager_secret_version" "secrets_store" {
  secret_id = aws_secretsmanager_secret.slack_auth.id
  secret_string = jsonencode({
    appToken = var.slack_app_token
    botToken = var.slack_bot_token
  })
}

resource "aws_security_group" "request_access" {
  name        = "ecs-${local.name}"
  description = "Security group for Request Accessn Bot in ECS"
  vpc_id      = var.aws_vpc_id

  tags = merge(var.tags, { Name = "ecs-${local.name}" })
}


resource "aws_security_group_rule" "request_access_egress" {
  security_group_id = aws_security_group.request_access.id
  type              = "egress"
  description       = "Allow egress"
  from_port         = 0
  to_port           = 0
  protocol          = "-1"
  cidr_blocks       = ["0.0.0.0/0"]
}
