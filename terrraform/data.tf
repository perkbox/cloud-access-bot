data "aws_region" "current" {
}


data "aws_caller_identity" "current" {
}


data "aws_iam_policy_document" "request_access_task" {

  statement {
    sid    = "S3BucketAccess"
    effect = "Allow"
    actions = [
      "s3:*"
    ]
    resources = [
      module.config_bucket.s3_bucket.arn,
      "${module.config_bucket.s3_bucket.arn}/*"
    ]
  }
}


data "aws_iam_policy_document" "request_access_task_execution" {

  statement {
    effect = "Allow"
    actions = [
      "secretsmanager:GetSecretValue",
      "kms:Decrypt"
    ]
    resources = [
      aws_secretsmanager_secret.slack_auth.arn,
      aws_kms_key.slack_auth.arn
    ]
  }
}
