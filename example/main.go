package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"

	nubarium "github.com/Idmission-LLC/nubarium-go"
	"github.com/spf13/viper"
)

func init() {
	// Set the file name of the configurations file
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")

	// Enable reading from environment variables
	viper.AutomaticEnv()

	// Read the configuration file
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Could not read config file: %v", err)
		log.Println("Using environment variables or defaults")
	}
}

func main() {
	// Check for command line argument
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <image_file_path>")
	}

	imagePath := os.Args[1]

	// Load configuration from .env file or environment variables
	endpoint := viper.GetString("NUBARIUM_ENDPOINT")
	username := viper.GetString("NUBARIUM_USERNAME")
	password := viper.GetString("NUBARIUM_PASSWORD")

	// Validate required configuration
	if endpoint == "" {
		log.Fatal("NUBARIUM_ENDPOINT is required")
	}
	if username == "" {
		log.Fatal("NUBARIUM_USERNAME is required")
	}
	if password == "" {
		log.Fatal("NUBARIUM_PASSWORD is required")
	}

	// Read the image file
	imageData, err := os.ReadFile(imagePath)
	if err != nil {
		log.Fatalf("Error reading image file: %v", err)
	}

	// Convert image to base64
	base64Image := base64.StdEncoding.EncodeToString(imageData)

	// Initialize the Nubarium client
	client := nubarium.NewClient(
		nubarium.WithBaseURL(endpoint),
		nubarium.WithCredentials(username, password),
	)

	// Send request using the convenience method - response is automatically parsed
	result, err := client.SendComprobanteDomicilio(base64Image)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}

	// Display results
	fmt.Println("=== Comprobante Domicilio OCR Results ===")
	fmt.Printf("\nStatus: %s\n", result.Status)
	fmt.Printf("Tipo: %s\n", result.Tipo)
	fmt.Printf("Nombre: %s\n", result.Nombre)
	fmt.Printf("Número de Servicio: %s\n", result.NumeroServicio)
	fmt.Printf("Total a Pagar: $%s\n", result.TotalPagar)
	fmt.Printf("Fecha Límite de Pago: %s\n", result.FechaLimitePago)

	fmt.Printf("\nDirección:\n")
	fmt.Printf("  Calle: %s\n", result.Calle)
	fmt.Printf("  Colonia: %s\n", result.Colonia)
	fmt.Printf("  Ciudad: %s\n", result.Ciudad)
	fmt.Printf("  CP: %s\n", result.CP)

	fmt.Printf("\nValidaciones:\n")
	fmt.Printf("  Código Numérico: %s\n", result.Validaciones.CodigoNumerico)
	fmt.Printf("  Fecha: %s\n", result.Validaciones.Fecha)
	fmt.Printf("  Número Servicio: %s\n", result.Validaciones.NumeroServicio)
	fmt.Printf("  Tarifa: %s\n", result.Validaciones.Tarifa)
	fmt.Printf("  Total a Pagar: %s\n", result.Validaciones.TotalPagar)

	// Pretty print the full JSON
	fmt.Println("\n=== Full JSON Response ===")
	prettyBytes, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(prettyBytes))
}
