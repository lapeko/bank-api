data "aws_lb" "ing" {
  name = "${var.app_name}-ing"
}

resource "aws_route53_record" "app" {
  zone_id = var.hosted_zone_id
  name    = ""
  type    = "A"
  alias {
    name                   = data.aws_lb.ing.dns_name
    zone_id                = data.aws_lb.ing.zone_id
    evaluate_target_health = true
  }
}
