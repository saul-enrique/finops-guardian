package main

import (
	"bytes" // ¡NUEVO! Necesario para nuestros "buffers" de captura.
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

type InfracostOutput struct {
	PastTotalMonthlyCost string `json:"pastTotalMonthlyCost"`
	TotalMonthlyCost     string `json:"totalMonthlyCost"`
	DiffTotalMonthlyCost string `json:"diffTotalMonthlyCost"`
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Uso: finops-guardian <rama_base>")
		fmt.Println("Ejemplo: finops-guardian main")
		os.Exit(1)
	}
	baseBranch := os.Args[1]
	baselineFile := "infracost-base.json"

	fmt.Printf("Analizando la rama base '%s' para crear la línea base...\n", baseBranch)
	if err := exec.Command("git", "checkout", baseBranch).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error al cambiar a la rama base: %v\n", err)
		os.Exit(1)
	}

	breakdownCmd := exec.Command("infracost", "breakdown", "--path", ".", "--format", "json", "--out-file", baselineFile)
	if output, err := breakdownCmd.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "Error al crear la línea base de Infracost: %s\n", string(output))
		os.Exit(1)
	}
	fmt.Printf("Línea base guardada en '%s'\n", baselineFile)

	if err := exec.Command("git", "checkout", "-").Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error al volver a la rama de trabajo: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Comparando la rama actual con la línea base...")
	diffCmd := exec.Command("infracost", "diff", "--path", ".", "--compare-to", baselineFile, "--format", "json", "--no-color")
	
    // --- LA CORRECCIÓN DEFINITIVA ESTÁ AQUÍ ---
	// Creamos buffers separados para stdout y stderr.
    var stdout, stderr bytes.Buffer
    diffCmd.Stdout = &stdout // Le decimos al comando que escriba la salida estándar aquí.
    diffCmd.Stderr = &stderr // Y la salida de error aquí.

    // Ejecutamos el comando solo con Run(), que no combina las salidas.
    err := diffCmd.Run()
	if err != nil {
        // Si hay un error, la razón estará en stderr.
		fmt.Fprintf(os.Stderr, "Error al ejecutar 'infracost diff': %s\n", stderr.String())
		os.Exit(1)
	}

	// El JSON puro está ahora en stdout. Lo pasamos al analizador.
    jsonOutput := stdout.Bytes()

	fmt.Println("Análisis completado. Extrayendo resumen de costos...")
	var outputData InfracostOutput
	if err := json.Unmarshal(jsonOutput, &outputData); err != nil {
		fmt.Fprintf(os.Stderr, "Error al analizar el JSON de Infracost: %v\n", err)
        // Para depurar, imprimimos lo que recibimos que no es JSON válido.
        fmt.Fprintf(os.Stderr, "Salida recibida:\n%s\n", string(jsonOutput))
		os.Exit(1)
	}

	fmt.Println("\n--- 🛡️ Reporte de FinOps Guardian 🛡️ ---")
	fmt.Printf("Costo Mensual Anterior: $%s\n", outputData.PastTotalMonthlyCost)
	fmt.Printf("Nuevo Costo Mensual:    $%s\n", outputData.TotalMonthlyCost)
	fmt.Printf("Impacto de este cambio:  $%s\n", outputData.DiffTotalMonthlyCost)
	fmt.Println("------------------------------------")
}
