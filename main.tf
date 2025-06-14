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

resource "aws_instance" "web_server" {
    ami           = "ami-0c55b159cbfafe1f0" // Un ejemplo de AMI de Amazon Linux 2
    instance_type = "t2.micro"

    tags = {
      Name = "WebApp Server"
    }
}

resource "aws_instance" "another_server" {
    ami           = "ami-0c55b159cbfafe1f0"
    instance_type = "t2.small"
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