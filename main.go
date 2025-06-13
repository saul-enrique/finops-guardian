package main

import (
	"fmt"
	"os"
	"os/exec" // <-- ¡NUEVO! Importamos el paquete para ejecutar comandos
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Uso: finops-guardian <rama_base> <rama_actual>")
		fmt.Println("Ejemplo: finops-guardian main feature/new-infra")
		os.Exit(1)
	}

	baseBranch := os.Args[1]
	// La rama actual es la que ya está "checked out", por lo que no la necesitamos para el comando diff.
	// La mantenemos por si la usamos en el futuro para los comentarios.
	// currentBranch := os.Args[2]

	fmt.Printf("Analizando diferencias de costo entre la rama actual y '%s'...\n", baseBranch)

	// --- LÓGICA CENTRAL ---
	// Construimos el argumento de comparación para Infracost.
	// El formato es git::[nombre_de_la_rama]
	compareTo := fmt.Sprintf("git::%s", baseBranch)

	// Preparamos el comando que vamos a ejecutar.
	// infracost diff --path=. --compare-to=git::main --format=json
	// Lo separamos en partes para que sea más claro.
	cmd := exec.Command("infracost", "diff", "--path=.", "--compare-to", compareTo, "--format=json")

	// Ejecutamos el comando y capturamos su salida (stdout) y cualquier error (stderr).
	output, err := cmd.CombinedOutput() // Usamos CombinedOutput para capturar stdout y stderr juntos

	// Manejamos el error. Si el comando falla, 'err' no será 'nil'.
	if err != nil {
		fmt.Println("Error al ejecutar Infracost:")
		fmt.Println(string(output)) // Imprimimos la salida del error para poder depurar
		os.Exit(1)
	}

	// Si todo fue bien, imprimimos la salida JSON de Infracost.
	fmt.Println("¡Análisis de Infracost completado!")
	fmt.Println("Salida JSON:")
	fmt.Println(string(output))
}
