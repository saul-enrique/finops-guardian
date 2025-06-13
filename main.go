package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings" // ¡NUEVO! Para construir nuestro string de reporte.
)

type InfracostOutput struct {
	PastTotalMonthlyCost string `json:"pastTotalMonthlyCost"`
	TotalMonthlyCost     string `json:"totalMonthlyCost"`
	DiffTotalMonthlyCost string `json:"diffTotalMonthlyCost"`
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Uso: finops-guardian <rama_base>")
		os.Exit(1)
	}
	baseBranch := os.Args[1]
	baselineFile := "infracost-base.json"

	// ... (Toda la lógica de checkout y diff se mantiene igual) ...
	if err := exec.Command("git", "checkout", baseBranch).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error al cambiar a la rama base: %v\n", err)
		os.Exit(1)
	}

	breakdownCmd := exec.Command("infracost", "breakdown", "--path", ".", "--format", "json", "--out-file", baselineFile)
	if output, err := breakdownCmd.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "Error al crear la línea base de Infracost: %s\n", string(output))
		os.Exit(1)
	}

	if err := exec.Command("git", "checkout", "-").Run(); err != nil {
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

	// --- ¡EL GRAN CAMBIO ESTÁ AQUÍ! ---
	// En lugar de imprimir, construimos un comentario en formato Markdown.
	var reportBuilder strings.Builder
	reportBuilder.WriteString("### 🛡️ Reporte de FinOps Guardian 🛡️\n\n")
	reportBuilder.WriteString("| Descripción | Costo Mensual |\n")
	reportBuilder.WriteString("| :--- | :--- |\n")
	reportBuilder.WriteString(fmt.Sprintf("| Costo Anterior | **$%s** |\n", outputData.PastTotalMonthlyCost))
	reportBuilder.WriteString(fmt.Sprintf("| Nuevo Costo | **$%s** |\n", outputData.TotalMonthlyCost))
	reportBuilder.WriteString(fmt.Sprintf("| **Impacto del Cambio** | **$%s** |\n", outputData.DiffTotalMonthlyCost))

	// Esta es la sintaxis especial para crear un "output" para la GitHub Action.
	// Le decimos: "crea una variable de salida llamada 'report' con el contenido de nuestro comentario".
	fmt.Printf("::set-output name=report::%s\n", reportBuilder.String())
}
