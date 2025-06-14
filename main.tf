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
