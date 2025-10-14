package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"

	nubarium "github.com/idmission/nubarium-go"
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

	// Create request payload with base64 encoded image
	type RequestPayload struct {
		Comprobante string `json:"comprobante"`
	}

	payload := RequestPayload{
		Comprobante: base64Image,
	}

	// Send request
	response, err := client.SendRequestWithPayload(payload)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}

	// Display results
	fmt.Printf("Status Code: %d\n", response.StatusCode)
	fmt.Println("\nJSON Response:")

	// Pretty print the JSON response
	var prettyJSON interface{}
	if err := json.Unmarshal([]byte(response.JSONData), &prettyJSON); err == nil {
		prettyBytes, _ := json.MarshalIndent(prettyJSON, "", "  ")
		fmt.Println(string(prettyBytes))
	} else {
		fmt.Println(response.JSONData)
	}
}
