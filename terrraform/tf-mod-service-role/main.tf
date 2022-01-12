terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 3.44.0"
    }
  }

}

resource "aws_iam_role" "role" {
  name = var.role_name
  path = var.iam_path

  assume_role_policy = data.aws_iam_policy_document.assume_role_policy.json

  tags = {
    Owner     = var.tag_owner
    Stack     = var.tag_stack
    ManagedBy = "Terraform"
  }
}

resource "aws_iam_role_policy" "policy" {
  for_each = { for policy in var.policies : policy.name => policy }

  role   = aws_iam_role.role.id
  name   = each.value.name
  policy = each.value.policy
}

resource "aws_iam_role_policy_attachment" "policy_attachment" {
  for_each = { for attachement in var.policy_attachments : attachement => attachement }

  role       = aws_iam_role.role.id
  policy_arn = each.key
}
