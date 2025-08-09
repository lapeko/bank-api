data "aws_iam_openid_connect_provider" "eks" {
  arn = var.oidc_provider_arn
}

resource "aws_iam_role" "alb_controller" {
  name = "${var.app_name}-alb-controller-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Principal = {
        Federated = data.aws_iam_openid_connect_provider.eks.arn
      }
      Action = "sts:AssumeRoleWithWebIdentity"
      Condition = {
        StringEquals = {
          "${var.oidc_issuer}:sub" = "system:serviceaccount:kube-system:aws-load-balancer-controller"
          "${var.oidc_issuer}:aud" = "sts.amazonaws.com"
        }
      }
    }]
  })

  tags = merge(var.tags, {
    Name = "${var.app_name}-alb-controller-role"
  })
}

resource "aws_iam_policy" "alb_controller" {
  name   = "${var.app_name}-alb-controller-policy"
  policy = data.http.alb_iam_policy.response_body

  tags = merge(var.tags, {
    Name = "${var.app_name}-alb-controller-policy"
  })
}

resource "aws_iam_role_policy_attachment" "alb_controller_policy" {
  role       = aws_iam_role.alb_controller.name
  policy_arn = aws_iam_policy.alb_controller.arn
}

data "http" "alb_iam_policy" {
  url = "https://raw.githubusercontent.com/kubernetes-sigs/aws-load-balancer-controller/main/docs/install/iam_policy.json"
}