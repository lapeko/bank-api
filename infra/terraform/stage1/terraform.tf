terraform {
  required_version = "~> 1.12"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.0"
    }
  }
  backend "s3" {
    bucket         = "bank-api-tfstate"
    key            = "stage1/terraform.tfstate"
    region         = "eu-central-1"
    dynamodb_table = "bank-api-tf-locks"
    encrypt        = true
  }
}
