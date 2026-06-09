terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
  backend "s3" {
    bucket = "taskmaster-terraform-state"
    key    = "infrastructure/terraform.tfstate"
    region = "us-east-1"
  }
}

provider "aws" {
  region = var.aws_region
}

module "vpc" {
  source = "./modules/vpc"
  name   = var.project_name
}

module "rds" {
  source     = "./modules/rds"
  name       = var.project_name
  vpc_id     = module.vpc.vpc_id
  subnet_ids = module.vpc.private_subnet_ids
}

module "elasticache" {
  source     = "./modules/elasticache"
  name       = var.project_name
  vpc_id     = module.vpc.vpc_id
  subnet_ids = module.vpc.private_subnet_ids
}

module "alb" {
  source     = "./modules/alb"
  name       = var.project_name
  vpc_id     = module.vpc.vpc_id
  subnet_ids = module.vpc.public_subnet_ids
}

module "ecs" {
  source          = "./modules/ecs"
  name            = var.project_name
  vpc_id          = module.vpc.vpc_id
  private_subnets = module.vpc.private_subnet_ids
  alb_sg_id       = module.alb.sg_id
  alb_target_arn  = module.alb.target_group_arn
  db_endpoint     = module.rds.endpoint
  redis_endpoint  = module.elasticache.endpoint
}

module "s3_cloudfront" {
  source = "./modules/s3_cloudfront"
  name   = var.project_name
}
