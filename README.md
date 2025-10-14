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
go get github.com/idmission/nubarium-go
```

## Configuration

The module uses Viper to load configuration from a `.env` file or environment variables.

Create a `.env` file in the project root:

```bash
NUBARIUM_ENDPOINT=https://ocr.nubarium.com/ocr/v2/comprobante_domicilio
NUBARIUM_USERNAME=your-username
NUBARIUM_PASSWORD=your-password
```

## Usage

### Basic Example

```go
package main

import (
    "fmt"
    "log"

    nubarium "github.com/idmission/nubarium-go"
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
    endpoint := viper.GetString("NUBARIUM_ENDPOINT")
    username := viper.GetString("NUBARIUM_USERNAME")
    password := viper.GetString("NUBARIUM_PASSWORD")

    // Create a new client using option pattern
    client := nubarium.NewClient(
        nubarium.WithBaseURL(endpoint),
        nubarium.WithCredentials(username, password),
    )

    // Option 1: Send raw JSON string
    jsonRequest := `{"field1": "value1", "field2": "value2"}`
    response, err := client.SendRequest(jsonRequest)
    if err != nil {
        log.Fatalf("Error: %v", err)
    }
    fmt.Printf("Status: %d\n", response.StatusCode)
    fmt.Println("Response:", response.JSONData)

    // Option 2: Send Go struct (automatically marshaled to JSON)
    type Payload struct {
        Field1 string `json:"field1"`
        Field2 string `json:"field2"`
    }
    payload := Payload{Field1: "value1", Field2: "value2"}
    response2, err := client.SendRequestWithPayload(payload)
    if err != nil {
        log.Fatalf("Error: %v", err)
    }

    // Option 3: Parse response into struct
    type ResponseData struct {
        Status string `json:"status"`
    }
    var data ResponseData
    if err := response2.ParseResponse(&data); err != nil {
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
    nubarium "github.com/idmission/nubarium-go"
)

// Create custom retryable client
retryClient := retryablehttp.NewClient()
retryClient.RetryMax = 5                      // Max number of retries
retryClient.RetryWaitMin = 1 * time.Second   // Min wait between retries
retryClient.RetryWaitMax = 30 * time.Second  // Max wait between retries
retryClient.Logger = someLogger               // Custom logger

// Inject the custom client
client := nubarium.NewClient(
    nubarium.WithBaseURL("https://api.nubarium.com/endpoint"),
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
func (c *Client) SendRequest(jsonRequest string) (*Response, error)
```

Sends a JSON string request to Nubarium API and returns the response.

### SendRequestWithPayload

```go
func (c *Client) SendRequestWithPayload(payload interface{}) (*Response, error)
```

Marshals a Go struct to JSON and sends it to Nubarium API.

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

## Dependencies

- `github.com/spf13/viper` - Configuration management
- `github.com/hashicorp/go-retryablehttp` - Automatic HTTP retries with exponential backoff

## License

Apache License 2.0
