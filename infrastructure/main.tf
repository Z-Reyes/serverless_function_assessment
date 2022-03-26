##### High Level Configuration Options #####
terraform {
  required_providers {
    null = {
      source = "hashicorp/null"
      version = "3.1.1"
    }
    aws = {
      source = "hashicorp/aws"
      version = "~> 2.0"
    }
    archive = {
      version = "~> 1.3.0"
    }
  }
}

provider "aws" {
  profile = "default"
  region  = "us-west-2"
}


##### Resources to compile and zip executable for uploading to AWS lambda #####

resource "null_resource" "compile" {
  triggers = {
    build_number = "${timestamp()}"
  }
  provisioner "local-exec" {
    command = "GOOS=linux go build -o ../bin/main ../cmd/main.go"
  } 
}

data "archive_file" "zip" {
    type = "zip"
    source_file = "../bin/main"
    output_path = "../bin/main.zip"
    depends_on = [null_resource.compile]
}

##### Define inputs and specifications for AWS lambda instance #####
resource "aws_lambda_function" "zach_test_function" {
    function_name = "zach_test_function"
    handler = "main"
    runtime = "go1.x"
    filename = data.archive_file.zip.output_path
    source_code_hash = data.archive_file.zip.output_base64sha256
    role = "${aws_iam_role.iam_for_lambda.arn}"
    memory_size = 128
    timeout = 10
}


##### Create IAM role for lambda instance #####

resource "aws_iam_role" "iam_for_lambda" {
  name = "iam_for_lambda"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}


##### Create APIGateway API #####
resource "aws_api_gateway_rest_api" "api" {
  name = "zach_test_api"
  endpoint_configuration {types = ["REGIONAL"]}
}

##### Define new API Gateway resource (this is used for a path parameter) #####
resource "aws_api_gateway_resource" "resource" {
  path_part   = "{requestip}"
  parent_id   = "${aws_api_gateway_rest_api.api.root_resource_id}"
  rest_api_id = "${aws_api_gateway_rest_api.api.id}"
}

##### Define GET methods for both existing resources #####
resource "aws_api_gateway_method" "method" {
  rest_api_id   = "${aws_api_gateway_rest_api.api.id}"
  resource_id   = "${aws_api_gateway_resource.resource.id}"
  http_method   = "GET"
  authorization = "NONE"
}

resource "aws_api_gateway_method" "default_method" {
  rest_api_id   = "${aws_api_gateway_rest_api.api.id}"
  resource_id   = "${aws_api_gateway_rest_api.api.root_resource_id}"
  http_method   = "GET"
  authorization = "NONE"
}

##### Integrate methods into appropriate resources #####
resource "aws_api_gateway_integration" "integration" {
  rest_api_id             = "${aws_api_gateway_rest_api.api.id}"
  resource_id             = "${aws_api_gateway_resource.resource.id}"
  http_method             = "${aws_api_gateway_method.method.http_method}"
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = "${aws_lambda_function.zach_test_function.invoke_arn}"
}

resource "aws_api_gateway_integration" "default_integration" {
  rest_api_id             = "${aws_api_gateway_rest_api.api.id}"
  resource_id             = "${aws_api_gateway_rest_api.api.root_resource_id}"
  http_method             = "${aws_api_gateway_method.default_method.http_method}"
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = "${aws_lambda_function.zach_test_function.invoke_arn}"
}

resource "aws_lambda_permission" "apigw_lambda" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = "${aws_lambda_function.zach_test_function.function_name}"
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_api_gateway_rest_api.api.execution_arn}/*/*/*"
}


##### Deploy resources and methods to API Gateway #####
resource "aws_api_gateway_deployment" "zach_deploy" {
  depends_on = [aws_api_gateway_integration.integration, aws_api_gateway_integration.default_integration]

  rest_api_id = "${aws_api_gateway_rest_api.api.id}"
  stage_name  = "v1"
}

output "url" {
  value = "${aws_api_gateway_deployment.zach_deploy.invoke_url}"
}