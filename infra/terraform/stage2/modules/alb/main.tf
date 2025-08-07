data "aws_lb" "api_alb" {
  tags = {
    "kubernetes.io/ingress-name"      = var.ingress_name
    "kubernetes.io/ingress-namespace" = var.ingress_namespace
  }
}

resource "aws_route53_zone" "main" {
  name = var.domain
}

resource "aws_route53_record" "api" {
  zone_id = aws_route53_zone.main.zone_id
  name    = var.domain
  type    = "A"

  alias {
    name                   = data.aws_lb.api_alb.dns_name
    zone_id                = data.aws_lb.api_alb.zone_id
    evaluate_target_health = true
  }
}
