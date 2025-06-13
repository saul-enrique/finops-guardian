package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Uso: finops-guardian <rama_base> <ref_a_restaurar>")
		os.Exit(1)
	}
	baseBranch := os.Args[1]
	returnRef := os.Args[2] // El 'ref' (commit SHA o rama) al que volver
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
	if err := exec.Command("git", "checkout", returnRef).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error al volver a la referencia de trabajo: %v\n", err)
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
	fmt.Println(stdout.String())
}
