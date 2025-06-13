package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type InfracostOutput struct {
	PastTotalMonthlyCost string `json:"pastTotalMonthlyCost"`
	TotalMonthlyCost     string `json:"totalMonthlyCost"`
	DiffTotalMonthlyCost string `json:"diffTotalMonthlyCost"`
}

func main() {
    // AHORA ESPERAMOS 3 ARGUMENTOS
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Uso: finops-guardian <rama_base> <rama_a_restaurar>")
		os.Exit(1)
	}
	baseBranch := os.Args[1]
	returnBranch := os.Args[2] // La rama a la que volver
	baselineFile := "infracost-base.json"

	if err := exec.Command("git", "checkout", baseBranch).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error al cambiar a la rama base: %v\n", err)
		os.Exit(1)
	}

	breakdownCmd := exec.Command("infracost", "breakdown", "--path", ".", "--format", "json", "--out-file", baselineFile)
	if output, err := breakdownCmd.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "Error al crear la línea base de Infracost: %s\n", string(output))
		os.Exit(1)
	}

    // --- ¡LA CORRECCIÓN ESTÁ AQUÍ! ---
    // Usamos el nombre de la rama explícito en lugar de '-'.
	if err := exec.Command("git", "checkout", returnBranch).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error al volver a la rama de trabajo: %v\n", err)
		os.Exit(1)
	}

	diffCmd := exec.Command("infracost", "diff", "--path", ".", "--compare-to", baselineFile, "--format", "json", "--no-color")
	var stdout, stderr bytes.Buffer
	diffCmd.Stdout = &stdout
	diffCmd.Stderr = &stderr
	if err := diffCmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error al ejecutar 'infracost diff': %s\n", stderr.String())
		os.Exit(1)
	}
	jsonOutput := stdout.Bytes()

	var outputData InfracostOutput
	if err := json.Unmarshal(jsonOutput, &outputData); err != nil {
		fmt.Fprintf(os.Stderr, "Error al analizar el JSON de Infracost: %v\n", err)
		os.Exit(1)
	}

	var reportBuilder strings.Builder
	reportBuilder.WriteString("### 🛡️ Reporte de FinOps Guardian (Nuestro Bot) 🛡️\n\n")
	reportBuilder.WriteString("| Descripción | Costo Mensual |\n")
	reportBuilder.WriteString("| :--- | :--- |\n")
	reportBuilder.WriteString(fmt.Sprintf("| Costo Anterior | **$%s** |\n", outputData.PastTotalMonthlyCost))
	reportBuilder.WriteString(fmt.Sprintf("| Nuevo Costo | **$%s** |\n", outputData.TotalMonthlyCost))
	reportBuilder.WriteString(fmt.Sprintf("| **Impacto del Cambio** | **$%s** |\n", outputData.DiffTotalMonthlyCost))

	// Esta es la sintaxis especial para crear un "output" para la GitHub Action.
	// La sintaxis ha cambiado ligeramente en las nuevas versiones de Actions.
    // `>> $GITHUB_OUTPUT` es la forma más moderna y robusta.
	outputFile := os.Getenv("GITHUB_OUTPUT")
	file, err := os.OpenFile(outputFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "No se pudo abrir el archivo de salida de GitHub: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	fmt.Fprintf(file, "report=%s\n", reportBuilder.String())
}
