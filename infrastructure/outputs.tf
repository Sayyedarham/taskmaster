output "alb_dns" {
  value = module.alb.dns_name
}

output "db_endpoint" {
  value = module.rds.endpoint
}

output "redis_endpoint" {
  value = module.elasticache.endpoint
}

output "cloudfront_domain" {
  value = module.s3_cloudfront.cloudfront_domain
}
