terraform {
  backend "s3" {
    bucket = var.bucket
    key    = var.key
    region = var.region
  }
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "3.0.2"
    }
    aws = {
      source  = "hashicorp/aws"
      version = ">=3.0.0"
    }
  }
  required_version = ">= 0.14.0"
}

# AWS provider
provider "aws" {
  region = var.region
}

# Docker provider
provider "docker" {
  registry_auth {
    address  = data.aws_ecr_authorization_token.token.proxy_endpoint
    username = data.aws_ecr_authorization_token.token.user_name
    password = data.aws_ecr_authorization_token.token.password
  }
}

# Get authorization credentials to push to ecr
data "aws_ecr_authorization_token" "token" {}
# Get availability zones for this region
data "aws_availability_zones" "this" {}

locals {
  # Tags
  tags = {
    environment = var.environment
    app_name    = var.app_name
    region      = var.region
  }
  # Elastic container registry data
  registry_endpoint      = replace(data.aws_ecr_authorization_token.token.proxy_endpoint, "https://", "")
  ecr_name               = format("ecr-%s-%s-%s", var.app_name, var.environment, var.region)
  task_definition_prefix = format("ecs-td-%s-%s-%s", var.app_name, var.environment, var.region)
  containers = [
    {
      name                 = "bff"
      repository_name      = format("%s-%s", local.ecr_name, "bff")
      endpoint             = format("%s/%s:latest", local.registry_endpoint, format("%s-%s", local.ecr_name, "bff"))
      task_definition_name = format("%s-%s", local.task_definition_prefix, "bff")
      port                 = 8080
      host_port            = 8080
      environment = {
        "BFF_ENABLEPLAYGROUND" = "true"
      }
      assign_public_ip = true
      public           = true
    },
    {
      name                 = "item"
      repository_name      = format("%s-%s", local.ecr_name, "item")
      endpoint             = format("%s/%s:latest", local.registry_endpoint, format("%s-%s", local.ecr_name, "item"))
      task_definition_name = format("%s-%s", local.task_definition_prefix, "item")
      port                 = 8081
      host_port            = 8081
      environment          = {}
      assign_public_ip     = true
      public               = true
    },
    {
      name                 = "tournament"
      repository_name      = format("%s-%s", local.ecr_name, "tournament")
      endpoint             = format("%s/%s:latest", local.registry_endpoint, format("%s-%s", local.ecr_name, "tournament"))
      task_definition_name = format("%s-%s", local.task_definition_prefix, "tournament")
      port                 = 8082
      host_port            = 8082
      environment          = {}
      assign_public_ip     = true
      public               = true
    },
    {
      name                 = "team"
      repository_name      = format("%s-%s", local.ecr_name, "team")
      endpoint             = format("%s/%s:latest", local.registry_endpoint, format("%s-%s", local.ecr_name, "team"))
      task_definition_name = format("%s-%s", local.task_definition_prefix, "team")
      port                 = 8083
      host_port            = 8083
      environment          = {}
      assign_public_ip     = true
      public               = true
    },
    {
      name                 = "record"
      repository_name      = format("%s-%s", local.ecr_name, "record")
      endpoint             = format("%s/%s:latest", local.registry_endpoint, format("%s-%s", local.ecr_name, "record"))
      task_definition_name = format("%s-%s", local.task_definition_prefix, "record")
      port                 = 8084
      host_port            = 8084
      environment          = {}
      assign_public_ip     = true
      public               = true
    }
  ]
}

# VPC
module "vpc" {
  source             = "./vpc"
  name               = format("vpc-%s-%s-%s", var.app_name, var.environment, var.region)
  availability_zones = data.aws_availability_zones.this.names
  tags               = local.tags
}

# Elastic container registry
module "ecr" {
  source     = "./ecr"
  containers = local.containers
  name       = local.ecr_name
  tags       = local.tags
}


# Elastic container service
module "ecs" {
  source          = "./ecs"
  name            = format("ecs-%s-%s-%s", var.app_name, var.environment, var.region)
  tags            = local.tags
  containers      = local.containers
  repository_arn  = module.ecr.arn
  vpc_id          = module.vpc.vpc_id
  public_subnets  = module.vpc.public_subnet_ids
  private_subnets = module.vpc.private_subnet_ids
  security_groups = [module.vpc.security_group_id]
}
