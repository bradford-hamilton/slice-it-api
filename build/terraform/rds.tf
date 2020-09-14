resource "aws_db_instance" "default" {
  publicly_accessible  = false
  allocated_storage    = 5
  engine               = "postgres"
  instance_class       = "db.t2.micro"
  name                 = "slice_it_api_storage"
  username             = var.db_username
  password             = var.db_password
  db_subnet_group_name = aws_db_subnet_group.public.name

  depends_on = [aws_db_subnet_group.public]
}
