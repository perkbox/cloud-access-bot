# ECS task execution role for grafana ECS service. Allows ECS to write the service logs to 
# CloudWatch logs.
module "request_access_task_execution_role" {
  source = "./tf-mod-service-role"

  role_name = "ecs-request-access-task-execution"
  services  = ["ecs-tasks"]
  policies = concat(
    [
      {
        name   = "SecretsManagerAccess"
        policy = data.aws_iam_policy_document.request_access_task_execution.json
      }
    ]
  )
  policy_attachments = [
    "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
  ]

  tag_stack = var.stack
}

# ECS task role for grafana ECS service. This determines the AWS permissions the container running in ECS
# has. In this case, mainly used to allow grafana to assume IAM roles in the other AWS accounts in order to 
# access CloudWatch metrics in those accounts
module "request_access_task_role" {
  source = "./tf-mod-service-role"

  role_name = "ecs-request-access-task"
  services  = ["ecs-tasks"]
  policies = concat(
    [
      {
        name   = "RequestAccess"
        policy = data.aws_iam_policy_document.request_access_task.json
      }
    ]
  )
  policy_attachments = []

  tag_stack = var.stack
}
