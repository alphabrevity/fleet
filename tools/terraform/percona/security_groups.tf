resource "aws_security_group" "lb" {
  name        = "percona load balancer"
  description = "percona Load balancer security group"
  vpc_id      = var.vpc_id
}

resource "aws_security_group_rule" "lb-ingress" {
  description = "percona: allow traffic from public internet"
  type        = "ingress"

  from_port   = "443"
  to_port     = "443"
  protocol    = "tcp"
  cidr_blocks = ["0.0.0.0/0"]

  security_group_id = aws_security_group.lb.id
}

resource "aws_security_group_rule" "lb-http-ingress" {
  description = "percona: allow traffic from public internet"
  type        = "ingress"

  from_port   = "80"
  to_port     = "80"
  protocol    = "tcp"
  cidr_blocks = ["0.0.0.0/0"]

  security_group_id = aws_security_group.lb.id
}
resource "aws_security_group_rule" "backend-egress" {
  description = "percona: allow all outbound traffic"
  type        = "egress"

  from_port   = 0
  to_port     = 0
  protocol    = "-1"
  cidr_blocks = ["0.0.0.0/0"]

  security_group_id = aws_security_group.backend.id
}

resource "aws_security_group" "backend" {
  name        = "percona backend"
  description = "percona Backend security group"
  vpc_id      = var.vpc_id

}

resource "aws_security_group_rule" "lb-egress" {
  description = "percona: allow all outbound traffic"
  type        = "egress"

  from_port   = 0
  to_port     = 0
  protocol    = "-1"
  cidr_blocks = ["0.0.0.0/0"]

  security_group_id = aws_security_group.lb.id
}
resource "aws_security_group_rule" "backend-ingress" {
  description = "percona: allow traffic from load balancer"
  type        = "ingress"

  from_port                = "80"
  to_port                  = "80"
  protocol                 = "tcp"
  source_security_group_id = aws_security_group.lb.id
  security_group_id        = aws_security_group.backend.id
}
