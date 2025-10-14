package nubarium

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
)

// Client represents a Nubarium API client
type Client struct {
	BaseURL         string
	Username        string
	Password        string
	RetryableClient *retryablehttp.Client
}

// ClientOption is a function that configures a Client
type ClientOption func(*Client)

// WithBaseURL sets the base URL for the Nubarium API
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.BaseURL = baseURL
	}
}

// WithCredentials sets the username and password for Basic Auth
func WithCredentials(username, password string) ClientOption {
	return func(c *Client) {
		c.Username = username
		c.Password = password
	}
}

// WithRetryableClient allows injecting a custom retryable HTTP client
func WithRetryableClient(client *retryablehttp.Client) ClientOption {
	return func(c *Client) {
		c.RetryableClient = client
	}
}

// NewClient creates a new Nubarium client with the provided options
func NewClient(opts ...ClientOption) *Client {
	// Create default retryable client
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 3
	retryClient.Logger = nil // Disable default logging

	client := &Client{
		RetryableClient: retryClient,
	}

	// Apply all options
	for _, opt := range opts {
		opt(client)
	}

	return client
}

// Response represents the response from Nubarium
type Response struct {
	JSONData   string
	StatusCode int
	Headers    http.Header
}

// SendRequest sends a JSON request to a specific Nubarium API endpoint with automatic retries
func (c *Client) SendRequest(endpoint string, jsonRequest string) (*Response, error) {
	// Construct full URL
	fullURL := c.BaseURL + endpoint

	// Step 1: Prepare retryable HTTP request
	req, err := retryablehttp.NewRequest(http.MethodPost, fullURL, bytes.NewBufferString(jsonRequest))
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %w", err)
	}

	// Step 2: Set headers
	req.Header.Set("Content-Type", "application/json")

	// Add Basic Auth if credentials are provided
	if c.Username != "" && c.Password != "" {
		auth := c.Username + ":" + c.Password
		encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
		req.Header.Set("Authorization", "Basic "+encodedAuth)
	}

	// Step 3: Send request with automatic retries
	resp, err := c.RetryableClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Step 4: Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	jsonResponse := string(body)

	// Check if response is actually JSON by trying to validate it
	// If it's not JSON (e.g., error page, HTML), return error
	var testJSON interface{}
	if err := json.Unmarshal(body, &testJSON); err != nil {
		return &Response{
			JSONData:   jsonResponse,
			StatusCode: resp.StatusCode,
			Headers:    resp.Header,
		}, fmt.Errorf("API returned non-JSON response (status %d): %s", resp.StatusCode, jsonResponse)
	}

	return &Response{
		JSONData:   jsonResponse,
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
	}, nil
}

// SendRequestWithPayload sends a request with a struct payload that will be marshaled to JSON
func (c *Client) SendRequestWithPayload(endpoint string, payload interface{}) (*Response, error) {
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling payload to JSON: %w", err)
	}
	return c.SendRequest(endpoint, string(jsonBytes))
}

// ParseResponse unmarshals the JSON response into the provided struct
func (r *Response) ParseResponse(v interface{}) error {
	return json.Unmarshal([]byte(r.JSONData), v)
}

// Endpoint constants for Nubarium OCR API
const (
	EndpointComprobanteDomicilio = "/ocr/v2/comprobante_domicilio"
)

// ComprobanteDomicilioRequest represents the request payload for the comprobante_domicilio endpoint
type ComprobanteDomicilioRequest struct {
	Comprobante string `json:"comprobante"` // Base64 encoded image
}

// ComprobanteDomicilioValidaciones represents the validation results in the comprobante_domicilio response
type ComprobanteDomicilioValidaciones struct {
	CodigoNumerico string `json:"codigoNumerico"`
	Fecha          string `json:"fecha"`
	NumeroServicio string `json:"numeroServicio"`
	Tarifa         string `json:"tarifa"`
	TotalPagar     string `json:"totalPagar"`
}

// ComprobanteDomicilioResponse represents the response from the comprobante_domicilio endpoint
type ComprobanteDomicilioResponse struct {
	QR               string                           `json:"QR"`
	Calle            string                           `json:"calle"`
	Ciudad           string                           `json:"ciudad"`
	ClaveMensaje     string                           `json:"claveMensaje"`
	CodigoBarras     string                           `json:"codigoBarras"`
	CodigoNumerico   string                           `json:"codigoNumerico"`
	CodigoValidacion string                           `json:"codigoValidacion"`
	Colonia          string                           `json:"colonia"`
	CP               string                           `json:"cp"`
	Fecha            string                           `json:"fecha"`
	FechaLimitePago  string                           `json:"fechaLimitePago"`
	Multiplicador    string                           `json:"multiplicador"`
	Nombre           string                           `json:"nombre"`
	NumeroMedidor    string                           `json:"numeroMedidor"`
	NumeroServicio   string                           `json:"numeroServicio"`
	PeriodoFacturado string                           `json:"periodoFacturado"`
	Referencia       string                           `json:"referencia"`
	RMU2             string                           `json:"rmu2"`
	Status           string                           `json:"status"`
	Tarifa           string                           `json:"tarifa"`
	Tipo             string                           `json:"tipo"`
	TotalPagar       string                           `json:"totalPagar"`
	TotalPagar2      string                           `json:"totalPagar2"`
	Validaciones     ComprobanteDomicilioValidaciones `json:"validaciones"`
}

// SendComprobanteDomicilio is a convenience method for sending a comprobante_domicilio request
// It automatically parses the response into a ComprobanteDomicilioResponse struct
func (c *Client) SendComprobanteDomicilio(base64Image string) (*ComprobanteDomicilioResponse, error) {
	payload := ComprobanteDomicilioRequest{
		Comprobante: base64Image,
	}
	response, err := c.SendRequestWithPayload(EndpointComprobanteDomicilio, payload)
	if err != nil {
		return nil, err
	}

	var result ComprobanteDomicilioResponse
	if err := response.ParseResponse(&result); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}
	return &result, nil
}
