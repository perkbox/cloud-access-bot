# Cloud Access Bot Terraform Module

This Terraform Module will help get you started quickly with the Access bot in your own AWS account.

## Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider\_aws) | n/a |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_aws_subnet_ids"></a> [aws\_subnet\_ids](#input\_aws\_subnet\_ids) | subnet ids needed for ECS FARGATE PLACEMENT, Subnets should be part of the same VPC\_ID | `list(string)` | n/a | yes |
| <a name="input_aws_vpc_id"></a> [aws\_vpc\_id](#input\_aws\_vpc\_id) | VPC ID for Security groups. | `string` | n/a | yes |
| <a name="input_docker_image"></a> [docker\_image](#input\_docker\_image) | The docker image to launch within Fargate | `string` | n/a | yes |
| <a name="input_slack_app_token"></a> [slack\_app\_token](#input\_slack\_app\_token) | Slack application token (Secret) | `string` | n/a | yes |
| <a name="input_slack_bot_token"></a> [slack\_bot\_token](#input\_slack\_bot\_token) | Slack Bot token (Secret) | `string` | n/a | yes |
| <a name="input_tags"></a> [tags](#input\_tags) | n/a | `map` | <pre>{<br>  "ManagedBy": "Terraform",<br>  "Stack": "CloudAccessBot"<br>}</pre> | no |


## Usage 
This is just an example of how the module can be used within your own terraform code.

**Example of module Usage.**
```js
module "cloud_access_bot_fargate" {
  source = "git::ssh://git@github.com/perkbox/cloud-access-bot//terraform?ref=master"
  slack_app_token = "xapp-1-A02K..."
  slack_bot_token = "xoxb-2557..."
  docker_image = "alpine:latest"
  aws_subnet_ids = ["subnet-1110022","subnet-12345aeee"]
  aws_vpc_id = "vpc-812b69fe4"
  tags  = {
    Stack     = "CloudAccessBot"
    ManagedBy = "Terraform"
  }
}
```