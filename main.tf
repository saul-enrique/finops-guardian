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

resource "aws_db_instance" "final_db" {
  allocated_storage    = 10
  engine               = "mysql"
  engine_version       = "8.0"
  instance_class       = "db.t3.micro"
  username             = "final_user"
  password             = "MustBeSecretNow"
  skip_final_snapshot  = true
}
