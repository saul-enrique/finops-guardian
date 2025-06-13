// Definimos el proveedor de nube que vamos a usar.
// En una configuración real, aquí pondrías la región y las credenciales.
// Para Infracost, no necesita ser funcional, solo necesita existir.
provider "aws" {
  region = "us-east-1"
}

// Definimos nuestra infraestructura base.
// Un simple bucket de S3 para guardar archivos.
resource "aws_s3_bucket" "website_bucket" {
  bucket = "my-unique-finops-guardian-test-bucket-12345" // Los nombres de bucket S3 deben ser únicos globalmente

  tags = {
    Name        = "My FinOps Guardian Test Bucket"
    Environment = "Dev"
  }
}
// --- NUEVO RECURSO AÑADIDO EN ESTA RAMA ---
// Una pequeña instancia de servidor EC2 de tipo t2.micro
// que normalmente está cubierta por la capa gratuita de AWS,
// pero que Infracost identificará como un costo añadido.
resource "aws_instance" "web_server" {
  ami           = "ami-0c55b159cbfafe1f0" // Un ejemplo de AMI de Amazon Linux 2
  instance_type = "t2.micro"

  tags = {
    Name = "WebApp Server"
  }
}
resource "aws_instance" "another_server" {
  ami           = "ami-0c55b159cbfafe1f0"
  instance_type = "t2.small" # Un tipo diferente para que el costo sea distinto
}
