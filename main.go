package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	// Ejecutamos el comando de comparación directa. Es el método ideal para CI/CD.
	cmd := exec.Command("infracost", "diff", "--path", ".", "--compare-to", "git::origin/main", "--format", "json", "--no-color")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		// Si hay un error, lo imprimimos al error estándar para que los logs de la Action lo muestren.
		fmt.Fprintln(os.Stderr, "Error al ejecutar Infracost diff:", stderr.String())
		os.Exit(1)
	}

	// Si todo va bien, imprimimos el JSON puro a la salida estándar.
	fmt.Println(stdout.String())
}
