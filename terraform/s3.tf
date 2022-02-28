module "s3_bucket" {
  source  = "terraform-aws-modules/s3-bucket/aws"
  version = "2.14.1"

  bucket = "cloud-access-bot-config"

  acl = "private"

  versioning = {
    enabled = true
  }

  server_side_encryption_configuration = {
    rule = {
      apply_server_side_encryption_by_default = {
        sse_algorithm = "AES256"
      }
    }
  }

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
  tags                    = var.tags
}
