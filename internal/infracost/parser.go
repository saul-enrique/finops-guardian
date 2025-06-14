// internal/infracost/parser.go
package infracost

// El nombre del paquete es 'infracost', el nombre del directorio.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"text/template"
)

// Las structs son las mismas, pero ahora viven en este paquete.
// Nota que las hemos hecho "Públicas" (con mayúscula inicial) para poder
// usarlas desde main.go. Ej: InfracostOutput en lugar de infracostOutput.
type ResourceDiff struct {
	Name         string `json:"name"`
	ResourceType string `json:"resourceType"`
	MonthlyCost  string `json:"monthlyCost"`
}

type DiffDetail struct {
	TotalMonthlyCost string         `json:"totalMonthlyCost"`
	Resources        []ResourceDiff `json:"resources"`
}

type Project struct {
	Name string     `json:"name"`
	Diff DiffDetail `json:"diff"`
}

type InfracostOutput struct {
	DiffTotalMonthlyCost string    `json:"diffTotalMonthlyCost"`
	Projects           []Project `json:"projects"`
}

// Parse y GenerateReport ahora son funciones públicas de este paquete.

// ParseJSONFile lee y decodifica el archivo JSON de Infracost.
// Devuelve la estructura de datos poblada.
func ParseJSONFile(filePath string) (InfracostOutput, error) {
	var output InfracostOutput

	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		return output, fmt.Errorf("no se pudo leer el archivo '%s': %w", filePath, err)
	}

	err = json.Unmarshal(jsonData, &output)
	if err != nil {
		return output, fmt.Errorf("el archivo '%s' no es un JSON de Infracost válido: %w", filePath, err)
	}

	return output, nil
}


// GenerateReport toma los datos de Infracost y produce el string del reporte.
func GenerateReport(output InfracostOutput) (string, error) {
	const reportTemplate = `
--------------------------------------------------
📊 **Reporte de FinOps Guardian** 📈

Este cambio aumentará los costos mensuales estimados en **${{ .DiffTotalMonthlyCost }}**.

**Recursos añadidos/modificados:**
{{- range .Projects }}
{{- range .Diff.Resources }}
  - ` + "`{{ .Name }}`" + ` ({{ .ResourceType }}): +${{ .MonthlyCost }}
{{- end }}
{{- end }}
--------------------------------------------------
`
	tmpl, err := template.New("report").Parse(reportTemplate)
	if err != nil {
		return "", fmt.Errorf("error al parsear la plantilla del reporte: %w", err)
	}

	var report bytes.Buffer
	err = tmpl.Execute(&report, output)
	if err != nil {
		return "", fmt.Errorf("error al ejecutar la plantilla: %w", err)
	}

	return report.String(), nil
}
