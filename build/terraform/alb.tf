# Application load balancer for delivering traffic to the target group
resource "aws_alb" "main" {
  name            = "slice-it-api-load-balancer"
  subnets         = flatten(["${aws_subnet.public.*.id}"])
  security_groups = ["${aws_security_group.lb.id}"]
}

# Target group definition - http over port 80, main vpc, health checks to /ping
resource "aws_alb_target_group" "app" {
  name        = "slice-it-api-target-group"
  port        = 80
  protocol    = "HTTP"
  vpc_id      = aws_vpc.main.id
  target_type = "ip"

  health_check {
    healthy_threshold   = "3"
    interval            = "30"
    protocol            = "HTTP"
    matcher             = "200"
    timeout             = "3"
    path                = var.health_check_path
    unhealthy_threshold = "2"
  }
}

# Redirect all traffic from the ALB to the target group
resource "aws_alb_listener" "front_end" {
  load_balancer_arn = aws_alb.main.id
  port              = "80"
  protocol          = "HTTP"

  default_action {
    target_group_arn = aws_alb_target_group.app.id
    type             = "forward"
  }
}
