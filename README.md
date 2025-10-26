# Nubarium Go Client

A Go module that replicates the functionality of the Java Nubarium adapter for calling the Nubarium API.

## Features

- Direct JSON request/response handling
- HTTP Basic Authentication support
- Automatic retries with exponential backoff (using HashiCorp retryablehttp)
- Configurable retry settings
- Send raw JSON strings or Go structs
- Parse JSON responses into Go structs
- Clean, simple API

## Installation

```bash
go get github.com/Idmission-LLC/nubarium-go
```

## Configuration

The module uses Viper to load configuration from a `.env` file or environment variables.

Create a `.env` file in the project root:

```bash
NUBARIUM_ENDPOINT=https://ocr.nubarium.com
NUBARIUM_USERNAME=your-username
NUBARIUM_PASSWORD=your-password
```

Note: The `NUBARIUM_ENDPOINT` should be set to the base URL only. Specific endpoints are handled by the client methods.

## Usage

### Basic Example

```go
package main

import (
    "fmt"
    "log"

    nubarium "github.com/Idmission-LLC/nubarium-go"
    "github.com/spf13/viper"
)

func init() {
    viper.SetConfigName(".env")
    viper.SetConfigType("env")
    viper.AddConfigPath(".")
    viper.AutomaticEnv()

    if err := viper.ReadInConfig(); err != nil {
        log.Printf("Warning: %v", err)
    }
}

func main() {
    // Load configuration
    baseURL := viper.GetString("NUBARIUM_ENDPOINT")
    username := viper.GetString("NUBARIUM_USERNAME")
    password := viper.GetString("NUBARIUM_PASSWORD")

    // Create a new client using option pattern
    client := nubarium.NewClient(
        nubarium.WithBaseURL(baseURL),
        nubarium.WithCredentials(username, password),
    )

    // Option 1: Use endpoint-specific convenience method
    documentSource := "https://example.com/document.pdf" // or base64-encoded document
    response, err := client.SendComprobanteDomicilio(documentSource)
    if err != nil {
        log.Fatalf("Error: %v", err)
    }
    fmt.Printf("Status: %d\n", response.StatusCode)
    fmt.Println("Response:", response.JSONData)

    // Option 2: Use generic method with endpoint constant
    type Payload struct {
        Field1 string `json:"field1"`
        Field2 string `json:"field2"`
    }
    payload := Payload{Field1: "value1", Field2: "value2"}
    response2, err := client.SendRequestWithPayload("/ocr/v2/some_endpoint", payload)
    if err != nil {
        log.Fatalf("Error: %v", err)
    }

    // Option 3: Send raw JSON string to a specific endpoint
    jsonRequest := `{"field1": "value1", "field2": "value2"}`
    response3, err := client.SendRequest("/ocr/v2/some_endpoint", jsonRequest)
    if err != nil {
        log.Fatalf("Error: %v", err)
    }

    // Option 4: Parse response into struct
    type ResponseData struct {
        Status string `json:"status"`
    }
    var data ResponseData
    if err := response3.ParseResponse(&data); err != nil {
        log.Printf("Error parsing: %v", err)
    }
    fmt.Printf("Status: %s\n", data.Status)
}
```

### Custom Retry Configuration

You can customize the retry behavior by injecting your own retryable HTTP client:

```go
import (
    "time"
    retryablehttp "github.com/hashicorp/go-retryablehttp"
    nubarium "github.com/Idmission-LLC/nubarium-go"
)

// Create custom retryable client
retryClient := retryablehttp.NewClient()
retryClient.RetryMax = 5                      // Max number of retries
retryClient.RetryWaitMin = 1 * time.Second   // Min wait between retries
retryClient.RetryWaitMax = 30 * time.Second  // Max wait between retries
retryClient.Logger = someLogger               // Custom logger

// Inject the custom client
client := nubarium.NewClient(
    nubarium.WithBaseURL("https://ocr.nubarium.com"),
    nubarium.WithCredentials("username", "password"),
    nubarium.WithRetryableClient(retryClient),
)
```

## How It Works

The Go module sends JSON directly to Nubarium API:

1. **Request Preparation**

   - Accepts raw JSON string or Go struct
   - Marshals struct to JSON if needed
   - Adds `Content-Type: application/json` header
   - Adds Basic Authentication header if credentials provided

2. **Request Execution**

   - Sends POST request to Nubarium API
   - Uses HTTP client with configurable timeout

3. **Response Handling**
   - Returns JSON response as string
   - Validates response is valid JSON
   - Provides helper to parse into Go structs

## API Reference

### Client

```go
type Client struct {
    BaseURL         string
    Username        string
    Password        string
    RetryableClient *retryablehttp.Client
}
```

### NewClient

```go
func NewClient(opts ...ClientOption) *Client
```

Creates a new Nubarium client with the provided options.

### Client Options

```go
func WithBaseURL(baseURL string) ClientOption
func WithCredentials(username, password string) ClientOption
func WithRetryableClient(client *retryablehttp.Client) ClientOption
```

- `WithBaseURL`: Sets the base URL for the Nubarium API
- `WithCredentials`: Sets the username and password for Basic Authentication
- `WithRetryableClient`: Injects a custom retryable HTTP client (overrides default)

### SendRequest

```go
func (c *Client) SendRequest(ctx context.Context, endpoint string, jsonRequest string) (*Response, error)
```

Sends a JSON string request to a specific Nubarium API endpoint and returns the response.

### SendRequestWithPayload

```go
func (c *Client) SendRequestWithPayload(ctx context.Context, endpoint string, payload interface{}) (*Response, error)
```

Marshals a Go struct to JSON and sends it to a specific Nubarium API endpoint.

### SendComprobanteDomicilio

```go
func (c *Client) SendComprobanteDomicilio(ctx context.Context, documentSource string) (*ComprobanteDomicilioResponse, error)
```

Convenience method for sending a comprobante_domicilio OCR request. The `documentSource` parameter accepts either a URL or a base64-encoded document string. Returns the parsed response automatically.

### Response

```go
type Response struct {
    XMLData    string      // Deprecated, not used
    JSONData   string      // JSON response from API
    StatusCode int         // HTTP status code
    Headers    http.Header // Response headers
}
```

### ParseResponse

```go
func (r *Response) ParseResponse(v interface{}) error
```

Unmarshals the JSON response into the provided Go struct.

## Available Endpoints

The package currently supports the following endpoints:

### Comprobante Domicilio OCR

**Endpoint constant:** `nubarium.EndpointComprobanteDomicilio`

**Convenience method:** `client.SendComprobanteDomicilio(ctx context.Context, documentSource string) (*ComprobanteDomicilioResponse, error)`

**Response struct:** `ComprobanteDomicilioResponse`

**Example:**

```go
client := nubarium.NewClient(
    nubarium.WithBaseURL("https://ocr.nubarium.com"),
    nubarium.WithCredentials("username", "password"),
)

// Option 1: Using a URL
documentSource := "https://example.com/comprobante.pdf"
result, err := client.SendComprobanteDomicilio(context.Background(), documentSource)

// Option 2: Using base64-encoded document
// base64Document := "your-base64-encoded-document"
// result, err := client.SendComprobanteDomicilio(context.Background(), base64Document)
if err != nil {
    log.Fatalf("Error: %v", err)
}

// Response is automatically parsed into typed struct
fmt.Printf("Status: %s\n", result.Status)
fmt.Printf("Nombre: %s\n", result.Nombre)
fmt.Printf("Total a Pagar: $%s\n", result.TotalPagar)
fmt.Printf("Validación de Fecha: %s\n", result.Validaciones.Fecha)
```

**Response Fields:**

The `ComprobanteDomicilioResponse` struct contains:

- Basic Info: `QR`, `Tipo`, `Status`, `ClaveMensaje`
- Personal: `Nombre`, `Calle`, `Colonia`, `Ciudad`, `CP`
- Service: `NumeroServicio`, `NumeroMedidor`, `Tarifa`, `Referencia`
- Billing: `TotalPagar`, `TotalPagar2`, `Fecha`, `FechaLimitePago`, `PeriodoFacturado`
- Codes: `CodigoBarras`, `CodigoNumerico`, `CodigoValidacion`, `RMU2`, `Multiplicador`
- Validaciones: Nested object with validation results for `CodigoNumerico`, `Fecha`, `NumeroServicio`, `Tarifa`, `TotalPagar`

## Adding New Endpoints

To add support for a new Nubarium endpoint:

1. **Add an endpoint constant** in `nubarium.go`:

```go
const (
    EndpointComprobanteDomicilio = "/ocr/v2/comprobante_domicilio"
    EndpointNewEndpoint          = "/ocr/v2/new_endpoint" // Add here
)
```

2. **(Optional) Define a request struct** if the endpoint has specific payload requirements:

```go
type NewEndpointRequest struct {
    Field1 string `json:"field1"`
    Field2 string `json:"field2"`
}
```

3. **(Optional) Create a convenience method**:

```go
func (c *Client) SendNewEndpoint(ctx context.Context, field1, field2 string) (*Response, error) {
    payload := NewEndpointRequest{
        Field1: field1,
        Field2: field2,
    }
    return c.SendRequestWithPayload(ctx, EndpointNewEndpoint, payload)
}
```

4. **Use the new endpoint**:

```go
// Using convenience method (if created)
response, err := client.SendNewEndpoint(ctx, "value1", "value2")

// Or using generic method
payload := NewEndpointRequest{Field1: "value1", Field2: "value2"}
response, err := client.SendRequestWithPayload(ctx, EndpointNewEndpoint, payload)
```

## Dependencies

- `github.com/spf13/viper` - Configuration management
- `github.com/hashicorp/go-retryablehttp` - Automatic HTTP retries with exponential backoff

## License

Apache License 2.0
