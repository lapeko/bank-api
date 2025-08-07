resource "aws_db_subnet_group" "this" {
  name       = "${var.name}-db-subnet-group"
  subnet_ids = var.private_subnets

  tags = merge(var.tags, {
    Name = "${var.name}-db-subnet-group"
  })
}

resource "aws_db_instance" "postgres" {
  identifier              = "${var.name}-postgres"
  engine                  = "postgres"
  engine_version          = "17"
  instance_class          = var.instance_class
  allocated_storage       = 20
  db_subnet_group_name    = aws_db_subnet_group.this.name
  vpc_security_group_ids  = [aws_security_group.db.id]
  username                = var.db_username
  password                = var.db_password
  db_name                 = replace(var.name, "-", "_")
  skip_final_snapshot     = true
  deletion_protection     = true
  publicly_accessible     = false
  multi_az                = false
  backup_retention_period = 0
  apply_immediately       = true

  tags = merge(var.tags, {
    Name = "${var.name}-db"
  })
}

resource "aws_security_group" "db" {
  name        = "${var.name}-db-sg"
  vpc_id      = var.vpc_id
  description = "Allow DB access from private subnets"

  ingress {
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = [var.vpc_cidr]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(var.tags, {
    Name = "${var.name}-db-sg"
  })
}

locals {
  db_url = "postgres://${var.db_username}:${var.db_password}@${aws_db_instance.postgres.address}:${aws_db_instance.postgres.port}/${var.name}"
}
