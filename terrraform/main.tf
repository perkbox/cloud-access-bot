resource "aws_cloudwatch_log_group" "log_group" {

  name              = "/aws/ecs/app/cloud_access_bot_logs"
  retention_in_days = 30

  tags = var.tags
}

resource "aws_ecs_task_definition" "cloud_access_bot_task" {

  family = local.name
  container_definitions = jsonencode([
    {
      name        = local.name
      image       = var.docker_image
      essential   = true
      environment = [for k, v in local.bot_config : { name = k, value = v }]
      secrets     = [for k, v in local.bot_secrets : { name = k, valueFrom = v }]
      logConfiguration = {
        logDriver = "awslogs"
        options = {
          awslogs-group         = aws_cloudwatch_log_group.log_group.name
          awslogs-region        = data.aws_region.current.name
          awslogs-stream-prefix = local.name
        }
      }
    }
  ])

  task_role_arn            = module.request_access_task_role.arn
  execution_role_arn       = module.request_access_task_execution_role.arn
  network_mode             = "awsvpc"
  cpu                      = 256
  memory                   = 512
  requires_compatibilities = ["FARGATE"]

  tags = var.tags
}


resource "aws_ecs_service" "cloud_access_bot_service" {

  name                = local.name
  task_definition     = aws_ecs_task_definition.cloud_access_bot_task.arn
  desired_count       = 1
  cluster             = data.aws_ecs_cluster.request_bot_main.arn
  launch_type         = "FARGATE"
  scheduling_strategy = "REPLICA"

  network_configuration {
    subnets          = var.aws_subnet_ids
    security_groups  = [aws_security_group.request_access.id]
    assign_public_ip = false
  }


  lifecycle {
    ignore_changes = [desired_count]
  }

  tags = var.tags
}
