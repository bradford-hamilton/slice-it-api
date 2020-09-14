# Create a cloudwatch log group for streaming logs from the slice it server
resource "aws_cloudwatch_log_group" "slice_it_api_log_group" {
  name              = "/ecs/slice-it-api"
  retention_in_days = 30

  tags = {
    Name = "slice-it-api-log-group"
  }
}

# The slice-it-api log stream
resource "aws_cloudwatch_log_stream" "slice_it_api_log_stream" {
  name           = "slice-it-api-log-stream"
  log_group_name = aws_cloudwatch_log_group.slice_it_api_log_group.name
}
