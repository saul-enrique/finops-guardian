# main.tf

# Un recurso que ya existía (no cambia)
resource "aws_s3_bucket" "existente" {
  bucket = "mi-bucket-de-reportes-financieros"
}

# Un recurso nuevo añadido en esta PR
resource "aws_instance" "servidor_web" {
  ami           = "ami-0c55b159cbfafe1f0" # Amazon Linux 2
  instance_type = "t2.micro" # Nivel gratuito de AWS

  tags = {
    Name = "ServidorWeb-Nuevo"
  }
}

# Otro recurso nuevo añadido en esta PR
resource "aws_db_instance" "base_de_datos" {
  allocated_storage    = 20
  engine               = "mysql"
  engine_version       = "8.0"
  instance_class       = "db.t2.small" # Un tamaño pequeño, no gratuito
  db_name              = "mi_app_db"
  username             = "admin"
  password             = "una-contraseña-muy-segura"
  skip_final_snapshot  = true
}

# Este es el cambio que activará la action
# Segundo comentario para probar el workflow mejorado
