module "infrastructure" {
  source = "../../"

  aws_region   = "us-east-1"
  project_name = "taskmaster"
  environment  = "prod"
}
