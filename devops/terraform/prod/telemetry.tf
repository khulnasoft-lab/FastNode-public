resource "aws_kinesis_firehose_delivery_stream" "fastnode_status" {
    provider = aws.east
    name = "fastnode_status"
    destination = "extended_s3"

    extended_s3_configuration {
        role_arn   = "arn:aws:iam::XXXXXXX:role/firehose_delivery_role"
        bucket_arn = "arn:aws:s3:::fastnode-metrics"
        prefix     = "firehose/fastnode_status/"

        buffer_size        = 5
        buffer_interval    = 600
        compression_format = "GZIP"

        error_output_prefix = "firehose/failures/fastnode_status/"

        cloudwatch_logging_options {
            enabled         = true
            log_group_name  = "/aws/kinesisfirehose/fastnode_status"
            log_stream_name = "S3Delivery"
        }
    }
}

locals {
    telemetry_response_template = <<EOF
#set($inputRoot = $input.path('$'))
{ }
EOF
}

# ~~ fastnode_status

resource "aws_api_gateway_rest_api" "Telemetry" {
    provider = aws.east
    name     = "Telemetry"
}
resource "aws_api_gateway_resource" "Telemetry_fastnode_status" {
    provider    = aws.east
    rest_api_id = aws_api_gateway_rest_api.Telemetry.id
    parent_id   = aws_api_gateway_rest_api.Telemetry.root_resource_id
    path_part   = "fastnode_status"
}
resource "aws_api_gateway_method" "Telemetry_fastnode_status" {
    provider    = aws.east
    rest_api_id = aws_api_gateway_rest_api.Telemetry.id
    resource_id = aws_api_gateway_resource.Telemetry_fastnode_status.id
    http_method = "POST"

    authorization     = "NONE"
    api_key_required = true
}
resource "aws_api_gateway_integration" "Telemetry_fastnode_status" {
    provider    = aws.east
    rest_api_id = aws_api_gateway_rest_api.Telemetry.id
    resource_id = aws_api_gateway_resource.Telemetry_fastnode_status.id
    http_method = aws_api_gateway_method.Telemetry_fastnode_status.http_method

    type                    = "AWS"
    integration_http_method = "POST"
    uri                     = "arn:aws:XXXXXXX:us-east-1:firehose:action/PutRecord"
    credentials             = "arn:aws:iam::XXXXXXX:role/APIGatewayPushToKinesis"

    request_templates = {
        "application/json" = templatefile("${path.module}/templates/telemetry_request_mapping.tmpl",
                                          { stream_name = "fastnode_status" })
    }
}
resource "aws_api_gateway_integration_response" "Telemetry_fastnode_status" {
    provider    = aws.east
    rest_api_id = aws_api_gateway_rest_api.Telemetry.id
    resource_id = aws_api_gateway_resource.Telemetry_fastnode_status.id
    http_method = aws_api_gateway_method.Telemetry_fastnode_status.http_method
    status_code = "200"

    response_templates = {
        "application/json" = local.telemetry_response_template
    }
}
resource "aws_api_gateway_method_response" "Telemetry_fastnode_status" {
    provider    = aws.east
    rest_api_id = aws_api_gateway_rest_api.Telemetry.id
    resource_id = aws_api_gateway_resource.Telemetry_fastnode_status.id
    http_method = aws_api_gateway_method.Telemetry_fastnode_status.http_method
    status_code = "200"
}

# ~~ client_events

resource "aws_api_gateway_resource" "Telemetry_client_events" {
    provider    = aws.east
    rest_api_id = aws_api_gateway_rest_api.Telemetry.id
    parent_id   = aws_api_gateway_rest_api.Telemetry.root_resource_id
    path_part   = "client_events"
}
resource "aws_api_gateway_method" "Telemetry_client_events" {
    provider    = aws.east
    rest_api_id = aws_api_gateway_rest_api.Telemetry.id
    resource_id = aws_api_gateway_resource.Telemetry_client_events.id
    http_method = "POST"

    authorization     = "NONE"
    api_key_required = true
}
resource "aws_api_gateway_integration" "Telemetry_client_events" {
    provider    = aws.east
    rest_api_id = aws_api_gateway_rest_api.Telemetry.id
    resource_id = aws_api_gateway_resource.Telemetry_client_events.id
    http_method = aws_api_gateway_method.Telemetry_client_events.http_method

    type                    = "AWS"
    integration_http_method = "POST"
    uri                     = "arn:aws:XXXXXXX:us-east-1:firehose:action/PutRecord"
    credentials             = "arn:aws:iam::XXXXXXX:role/APIGatewayPushToKinesis"

    request_templates = {
        "application/json" = templatefile("${path.module}/templates/telemetry_request_mapping_old.tmpl",
                                          { stream_name = "client_events" })
    }
}
resource "aws_api_gateway_integration_response" "Telemetry_client_events" {
    provider    = aws.east
    rest_api_id = aws_api_gateway_rest_api.Telemetry.id
    resource_id = aws_api_gateway_resource.Telemetry_client_events.id
    http_method = aws_api_gateway_method.Telemetry_client_events.http_method
    status_code = "200"

    response_templates = {
        "application/json" = local.telemetry_response_template
    }
}
resource "aws_api_gateway_method_response" "Telemetry_client_events" {
    provider    = aws.east
    rest_api_id = aws_api_gateway_rest_api.Telemetry.id
    resource_id = aws_api_gateway_resource.Telemetry_client_events.id
    http_method = aws_api_gateway_method.Telemetry_client_events.http_method
    status_code = "200"
}

# ~~ fastnode_service

resource "aws_api_gateway_resource" "Telemetry_fastnode_service" {
    provider    = aws.east
    rest_api_id = aws_api_gateway_rest_api.Telemetry.id
    parent_id   = aws_api_gateway_rest_api.Telemetry.root_resource_id
    path_part   = "fastnode_service"
}
resource "aws_api_gateway_method" "Telemetry_fastnode_service" {
    provider    = aws.east
    rest_api_id = aws_api_gateway_rest_api.Telemetry.id
    resource_id = aws_api_gateway_resource.Telemetry_fastnode_service.id
    http_method = "POST"

    authorization     = "NONE"
    api_key_required = true
}
resource "aws_api_gateway_integration" "Telemetry_fastnode_service" {
    provider    = aws.east
    rest_api_id = aws_api_gateway_rest_api.Telemetry.id
    resource_id = aws_api_gateway_resource.Telemetry_fastnode_service.id
    http_method = aws_api_gateway_method.Telemetry_fastnode_service.http_method

    type                    = "AWS"
    integration_http_method = "POST"
    uri                     = "arn:aws:XXXXXXX:us-east-1:firehose:action/PutRecord"
    credentials             = "arn:aws:iam::XXXXXXX:role/APIGatewayPushToKinesis"

    request_templates = {
        "application/json" = templatefile("${path.module}/templates/telemetry_request_mapping_old.tmpl",
                                          { stream_name = "fastnode_service" })
    }
}
resource "aws_api_gateway_integration_response" "Telemetry_fastnode_service" {
    provider    = aws.east
    rest_api_id = aws_api_gateway_rest_api.Telemetry.id
    resource_id = aws_api_gateway_resource.Telemetry_fastnode_service.id
    http_method = aws_api_gateway_method.Telemetry_fastnode_service.http_method
    status_code = "200"

    response_templates = {
        "application/json" = local.telemetry_response_template
    }
}
resource "aws_api_gateway_method_response" "Telemetry_fastnode_service" {
    provider    = aws.east
    rest_api_id = aws_api_gateway_rest_api.Telemetry.id
    resource_id = aws_api_gateway_resource.Telemetry_fastnode_service.id
    http_method = aws_api_gateway_method.Telemetry_fastnode_service.http_method
    status_code = "200"
}
