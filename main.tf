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
