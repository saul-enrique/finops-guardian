// main.go
package main

import (
	"encoding/json" // Paquete para codificar y decodificar JSON
	"fmt"            // Paquete para formatear e imprimir texto
	"os"             // Paquete para interactuar con el sistema operativo (como leer archivos)
)

// Estas son las "structs" de Go. Definen la estructura de los datos que esperamos
// leer del JSON. Mapean los campos del JSON a campos en nuestro programa Go.
// La parte `json:"..."` se llama "tag" y le dice a Go exactamente qué campo
// del JSON corresponde a cada campo de la struct.

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
	TotalMonthlyCost   string    `json:"totalMonthlyCost"`
	DiffTotalMonthlyCost string    `json:"diffTotalMonthlyCost"`
	Projects           []Project `json:"projects"`
}


// La función main es el punto de entrada de cualquier programa en Go.
func main() {
	fmt.Println("🚀 Iniciando FinOps Guardian...")

	// Paso 1: Leer el contenido del archivo JSON de muestra.
	// os.ReadFile devuelve los datos del archivo y un posible error.
	jsonData, err := os.ReadFile("infracost-output.json")
	if err != nil {
		// Si hay un error (ej. el archivo no existe), el programa se detendrá
		// y mostrará el error. Esto es crucial para un código robusto.
		fmt.Printf("Error leyendo el archivo JSON: %v\n", err)
		os.Exit(1)
	}

	// Paso 2: Preparar una variable para almacenar los datos decodificados.
	// Creamos una variable llamada 'output' que tiene la estructura de InfracostOutput.
	var output InfracostOutput

	// Paso 3: Decodificar (o "Unmarshal") el JSON.
	// Tomamos los datos del archivo (jsonData) y los vertimos en nuestra
	// variable 'output'. Go automáticamente mapeará los campos gracias a las structs.
	err = json.Unmarshal(jsonData, &output)
	if err != nil {
		fmt.Printf("Error decodificando el JSON: %v\n", err)
		os.Exit(1)
	}

	// Paso 4: ¡Éxito! Imprimir los resultados que nos interesan.
	// Ahora podemos acceder a los datos de forma nativa en Go, ej: output.DiffTotalMonthlyCost
	fmt.Println("--------------------------------------------------")
	fmt.Printf("📊 Reporte de FinOps Guardian 📈\n")
	fmt.Printf(
		"Este cambio aumentará los costos mensuales estimados en $%s.\n",
		output.DiffTotalMonthlyCost,
	)

	fmt.Println("Recursos añadidos/modificados:")
	// Iteramos a través de cada recurso en el desglose de la diferencia.
	for _, project := range output.Projects {
		if len(project.Diff.Resources) > 0 {
			for _, resource := range project.Diff.Resources {
				fmt.Printf(
					"  ■ %s (%s): +$%s\n",
					resource.Name,
					resource.ResourceType,
					resource.MonthlyCost,
				)
			}
		}
	}
	fmt.Println("--------------------------------------------------")
}
