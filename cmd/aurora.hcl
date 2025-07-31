env "aws" {
  migration {
    dir = "file://database/migration"
  }

  url = "postgres://${var.aws_dsql_username}:${urlescape(data.aws_dsql_token.this)}@${var.aws_dsql_host}/example-api"
}

data "aws_dsql_token" "this" {
  username = var.aws_dsql_username
  endpoint = var.aws_dsql_host
  region   = var.aws_region
}

variable "aws_dsql_username" {
  type    = string
  default = "example-api"
}

variable "aws_dsql_host" {
  type    = string
  default = getenv("EXAMPLE_API_DATABASE_HOST")
}

variable "aws_region" {
  type    = string
  default = getenv("AWS_REGION")
}
