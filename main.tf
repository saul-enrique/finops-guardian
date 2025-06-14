provider "aws" {
  region = "us-east-1"
}

resource "aws_s3_bucket" "website_bucket" {
  bucket = "my-unique-finops-guardian-test-bucket-12345"
  tags = {
    Name        = "My FinOps Guardian Test Bucket"
    Environment = "Dev"
  }
}

resource "aws_db_instance" "showcase_db" {
  allocated_storage    = 10
  engine               = "mysql"
  engine_version       = "8.0"
  instance_class       = "db.t3.micro"
  username             = "showcase_user"
  password             = "MustBeVerySecret"
  skip_final_snapshot  = true
}
# Demostración final para el portafolio
resource "aws_alb" "application_load_balancer" {
  name               = "final-demo-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = []
  subnets            = ["subnet-123456", "subnet-abcdef"]
}
