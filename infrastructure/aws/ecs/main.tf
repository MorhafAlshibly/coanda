resource "aws_ecs_cluster" "this" {
  name = var.name
  tags = var.tags
}

resource "aws_ecs_cluster_capacity_providers" "this" {
  cluster_name = aws_ecs_cluster.this.name
  capacity_providers = [
    "FARGATE",
    "FARGATE_SPOT",
  ]
  default_capacity_provider_strategy {
    base              = 1
    weight            = 100
    capacity_provider = "FARGATE_SPOT"
  }
}

resource "aws_iam_policy" "this" {
  name        = var.name
  description = "Allow ECS to pull from ECR for specific repositories"
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "ecr:GetAuthorizationToken",
          "ecr:BatchCheckLayerAvailability",
          "ecr:GetDownloadUrlForLayer",
          "ecr:BatchGetImage",
          "logs:CreateLogStream",
          "logs:PutLogEvents"
        ]
        Effect   = "Allow"
        Resource = "*"
      }
    ]
  })
  tags = var.tags
}

resource "aws_iam_role" "this" {
  name = var.name
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "ecs-tasks.amazonaws.com"
        }
      }
    ]
  })
  tags = var.tags
}

resource "aws_iam_role_policy_attachment" "this" {
  role       = aws_iam_role.this.name
  policy_arn = aws_iam_policy.this.arn
}

resource "aws_ecs_task_definition" "this" {
  # Depends on the IAM role policy attachment
  depends_on               = [aws_iam_role_policy_attachment.this]
  for_each                 = { for container in var.containers : container.name => container }
  family                   = each.value.task_definition_name
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = 256
  memory                   = 512
  execution_role_arn       = aws_iam_role.this.arn
  container_definitions = jsonencode([
    {
      name      = each.value.name
      image     = each.value.endpoint
      cpu       = 256
      memory    = 512
      essential = true
      // Environment variables, check if the map is not empty
      environment = keys(each.value.environment) != [] ? [
        for key, value in each.value.environment : {
          name  = key
          value = value
        }
      ] : []
      portMappings = [
        {
          containerPort = each.value.port
          hostPort      = each.value.host_port
        }
      ]
    },
  ])
  runtime_platform {
    operating_system_family = "LINUX"
    cpu_architecture        = "X86_64"
  }
  tags = var.tags
}

resource "aws_ecs_service" "this" {
  depends_on      = [aws_ecs_cluster_capacity_providers.this]
  for_each        = { for container in var.containers : container.name => container }
  name            = each.value.name
  cluster         = aws_ecs_cluster.this.name
  task_definition = aws_ecs_task_definition.this[each.value.name].arn
  capacity_provider_strategy {
    capacity_provider = "FARGATE_SPOT"
    base              = 0
    weight            = 1
  }
  desired_count = 1
  network_configuration {
    subnets          = each.value.public ? var.public_subnets : var.private_subnets
    security_groups  = var.security_groups
    assign_public_ip = each.value.assign_public_ip
  }
  tags = var.tags
}
