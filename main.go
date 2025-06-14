// main.go
package main

import (
	"bytes"          // Para construir strings de manera eficiente
	"encoding/json"  // Paquete para codificar y decodificar JSON
	"flag"           // ¡NUEVO! Paquete para manejar los flags de la línea de comandos
	"fmt"            // Paquete para formatear e imprimir texto
	"os"             // Paquete para interactuar con el sistema operativo (como leer archivos)
	"text/template"  // ¡NUEVO! Paquete para crear plantillas de texto
)

// --- Las structs no cambian, son las mismas de antes ---

// ResourceDiff representa un único recurso que ha cambiado.
type ResourceDiff struct {
	Name         string `json:"name"`
	ResourceType string `json:"resourceType"`
	MonthlyCost  string `json:"monthlyCost"`
}

// DiffDetail contiene el resumen de los cambios en un proyecto.
type DiffDetail struct {
	TotalMonthlyCost string         `json:"totalMonthlyCost"`
	Resources        []ResourceDiff `json:"resources"`
}

// Project representa un único proyecto de Terraform analizado.
type Project struct {
	Name string     `json:"name"`
	Diff DiffDetail `json:"diff"`
}

// InfracostOutput es la estructura de nivel superior que representa todo el archivo JSON.
type InfracostOutput struct {
	DiffTotalMonthlyCost string    `json:"diffTotalMonthlyCost"`
	Projects           []Project `json:"projects"`
}

// ¡NUEVA FUNCIÓN! Esta función toma los datos de Infracost y genera el reporte
// como un string. Esto separa la "lógica" de la "presentación".
// Usamos plantillas para que sea más fácil cambiar el formato del reporte en el futuro.
func generateReport(output InfracostOutput) (string, error) {
	// Plantilla para el comentario. Usamos la sintaxis de Go templates.
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
	// `template.New` crea una nueva plantilla.
	// `template.Must` envuelve la creación y hace 'panic' si hay un error en la plantilla.
	tmpl, err := template.New("report").Parse(reportTemplate)
	if err != nil {
		// Este error es para el programador, si la plantilla está mal escrita.
		return "", fmt.Errorf("error al parsear la plantilla del reporte: %w", err)
	}

	// Usaremos un buffer para escribir el resultado de la plantilla.
	var report bytes.Buffer
	// `tmpl.Execute` aplica los datos (output) a la plantilla y escribe el resultado en el buffer.
	err = tmpl.Execute(&report, output)
	if err != nil {
		// Este error puede ocurrir si los datos no coinciden con la plantilla.
		return "", fmt.Errorf("error al ejecutar la plantilla: %w", err)
	}

	return report.String(), nil
}

// La función main ahora se enfoca en la orquestación.
func main() {
	// Paso 1: Definir y parsear los flags de la línea de comandos.
	// flag.String define un flag llamado "file", con un valor por defecto "infracost-output.json",
	// y un texto de ayuda.
	filePath := flag.String("file", "infracost-output.json", "Ruta al archivo JSON de salida de Infracost.")
	flag.Parse() // Procesa los argumentos de la línea de comandos.

	fmt.Println("🚀 Iniciando FinOps Guardian...")

	// Paso 2: Leer el archivo JSON usando la ruta del flag (*filePath).
	// ¡Ojo al asterisco! Se usa para obtener el valor del puntero que devuelve flag.String.
	jsonData, err := os.ReadFile(*filePath)
	if err != nil {
		fmt.Printf("Error: No se pudo leer el archivo '%s': %v\n", *filePath, err)
		os.Exit(1)
	}

	// Paso 3: Decodificar el JSON (sin cambios aquí).
	var output InfracostOutput
	err = json.Unmarshal(jsonData, &output)
	if err != nil {
		fmt.Printf("Error: El archivo '%s' no es un JSON de Infracost válido: %v\n", *filePath, err)
		os.Exit(1)
	}

	// Paso 4: Generar el reporte llamando a nuestra nueva función.
	report, err := generateReport(output)
	if err != nil {
		fmt.Printf("Error al generar el reporte: %v\n", err)
		os.Exit(1)
	}

	// Paso 5: Imprimir el reporte final.
	fmt.Println(report)
}
