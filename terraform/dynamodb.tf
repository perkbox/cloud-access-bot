module "dynamodb_table" {
  source  = "terraform-aws-modules/dynamodb-table/aws"
  version = "1.1.0"

  name      = "cloud_access_bot"
  hash_key  = "PK"
  range_key = "SK"

  attributes = [
    {
      name = "PK"
      type = "S"
    },
    {
      name = "SK"
      type = "S"
    }
  ]

  tags = var.tags
}
