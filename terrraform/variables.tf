locals {
  name  = "request_access_bot"
  owner = "Devops"
  image = "491933125842.dkr.ecr.eu-west-1.amazonaws.com/infrastructure/request-access-bot:v0.0.5"

  default_tags = {
    Owner     = local.owner
    ManagedBy = "Terraform"
  }

  bot_config = {
    BOT_CONFIG_S3_BUCKET = "request-access-bot-config"
    BOT_CONFIG_S3_KEY    = "config.yml"
  }
  bot_secrets = {
    SLACK_APP_TOKEN = "${aws_secretsmanager_secret.slack_auth.arn}:appToken::"
    SLACK_BOT_TOKEN = "${aws_secretsmanager_secret.slack_auth.arn}:botToken::"
  }
}


variable "vpc_tag_name" {
  description = "VPC name Tag used to locate the VPC ID"
}

variable "subnet_tag_name" {
  description = "Subnet name Tag used to locate the SUBNET IDs, must be in the VPC provided above."
}

variable "slack_app_token" {
  description = "Slack application token (Secret)"
}

variable "slack_bot_token" {
  description = "Slack Bot token (Secret)"
}
