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

# get authorization credentials to push to ecr
data "aws_ecr_authorization_token" "token" {}

# Tags
locals {
  tags = {
    environment = var.environment
    app_name    = var.app_name
    region      = var.region
  }
}

# Container registry
module "registry" {
  source = "./registry"
  containers = [
    {
      name = "bff"
    },
    {
      name = "item"
    },
    {
      name = "tournament"
    },
    {
      name = "team"
    },
    {
      name = "record"
    }
  ]
  name     = format("ecr-%s-%s-%s", var.app_name, var.environment, var.region)
  endpoint = data.aws_ecr_authorization_token.token.proxy_endpoint
  tags     = local.tags
}