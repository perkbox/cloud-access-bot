locals {
  name = "request_access_bot"

  bot_config = {
    BOT_CONFIG_S3_BUCKET = "request-access-bot-config"
    BOT_CONFIG_S3_KEY    = "config.yml"
  }
  bot_secrets = {
    SLACK_APP_TOKEN = "${aws_secretsmanager_secret.slack_auth.arn}:appToken::"
    SLACK_BOT_TOKEN = "${aws_secretsmanager_secret.slack_auth.arn}:botToken::"
  }
}


variable "slack_app_token" {
  description = "Slack application token (Secret)"
  type        = string
}

variable "slack_bot_token" {
  description = "Slack Bot token (Secret)"
  type        = string
}

variable "image" {
  description = "The docker image to launch within Fargate"
  type        = string
}

variable "aws_subnet_ids" {
  description = "subnet ids needed for ECS FARGATE PLACEMENT, Subnets should be part of the same VPC_ID "
  type        = list(string)
}

variable "aws_vpc_id" {
  description = "VPC ID for Security groups."
  type        = string
}


variable "tags" {
  default = {
    Stack     = "CloudAccessBot"
    ManagedBy = "Terraform"
  }
}
