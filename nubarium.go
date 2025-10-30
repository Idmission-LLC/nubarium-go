package nubarium

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
	dateparser "github.com/markusmobius/go-dateparser"
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
	retryClient.HTTPClient = &http.Client{
		Transport: http.DefaultTransport.(*http.Transport).Clone(),
		Timeout:   5 * time.Minute, // Nubarium API timeout
	}

	client := &Client{RetryableClient: retryClient}

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
func (c *Client) SendRequest(ctx context.Context, endpoint string, jsonRequest string) (*Response, error) {
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
	resp, err := c.RetryableClient.Do(req.WithContext(ctx))
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
	var testJSON any
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
func (c *Client) SendRequestWithPayload(ctx context.Context, endpoint string, payload any) (*Response, error) {
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling payload to JSON: %w", err)
	}
	return c.SendRequest(ctx, endpoint, string(jsonBytes))
}

// ParseResponse unmarshals the JSON response into the provided struct
func (r *Response) ParseResponse(v any) error {
	return json.Unmarshal([]byte(r.JSONData), v)
}

// Endpoint constants for Nubarium OCR API
const (
	EndpointComprobanteDomicilio = "/ocr/v2/comprobante_domicilio"
)

// ComprobanteDomicilioRequest represents the request payload for the comprobante_domicilio endpoint
type ComprobanteDomicilioRequest struct {
	Comprobante string `json:"comprobante"` // URL or base64-encoded document
}

// ComprobanteDomicilioValidaciones represents the validation results in the comprobante_domicilio response
type ComprobanteDomicilioValidaciones struct {
	CodigoNumerico string         `json:"codigoNumerico"`
	Fecha          string         `json:"fecha"`
	NumeroServicio string         `json:"numeroServicio"`
	Tarifa         string         `json:"tarifa"`
	TotalPagar     StringOrObject `json:"totalPagar"`
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
	TotalPagar       StringOrObject                   `json:"totalPagar"`
	TotalPagar2      StringOrObject                   `json:"totalPagar2"`
	Validaciones     ComprobanteDomicilioValidaciones `json:"validaciones"`

	ParsedDate time.Time `json:"parsedDate"`
	DateError  error     `json:"dateError"`
}

// StringOrObject unmarshals JSON values that may be a string or an object.
// The raw JSON is preserved; helpers are provided to access common forms.
type StringOrObject struct {
	raw json.RawMessage
}

// UnmarshalJSON implements json.Unmarshaler to accept string or object.
func (so *StringOrObject) UnmarshalJSON(data []byte) error {
	if so == nil {
		return errors.New("StringOrObject: nil receiver")
	}
	if data == nil {
		so.raw = nil
		return nil
	}
	so.raw = append(so.raw[:0], data...)
	return nil
}

// MarshalJSON returns the raw JSON unchanged to preserve original shape.
func (so StringOrObject) MarshalJSON() ([]byte, error) {
	if so.raw == nil {
		return []byte("null"), nil
	}
	return so.raw, nil
}

// IsString reports whether the underlying JSON value is a string.
func (so StringOrObject) IsString() bool {
	return len(so.raw) >= 2 && so.raw[0] == '"' && so.raw[len(so.raw)-1] == '"'
}

// String returns the string value if it is a string; otherwise a compact JSON representation of the object.
func (so StringOrObject) String() string {
	if so.IsString() {
		var s string
		if err := json.Unmarshal(so.raw, &s); err == nil {
			return s
		}
	}
	var out any
	if err := json.Unmarshal(so.raw, &out); err == nil {
		b, err := json.Marshal(out)
		if err == nil {
			return string(b)
		}
	}
	return string(so.raw)
}

// UnmarshalObject unmarshals the underlying value into dst if it is an object.
func (so StringOrObject) UnmarshalObject(dst any) error {
	if so.IsString() {
		return fmt.Errorf("value is string, not object")
	}
	return json.Unmarshal(so.raw, dst)
}

// SendComprobanteDomicilio is a convenience method for sending a comprobante_domicilio request
// It automatically parses the response into a ComprobanteDomicilioResponse struct
// The documentSource parameter can be either a URL or a base64-encoded document string
func (c *Client) SendComprobanteDomicilio(ctx context.Context, documentSource string) (result *ComprobanteDomicilioResponse, err error) {
	payload := ComprobanteDomicilioRequest{
		Comprobante: documentSource,
	}
	response, err := c.SendRequestWithPayload(ctx, EndpointComprobanteDomicilio, payload)
	if err != nil {
		return nil, err
	}

	result = &ComprobanteDomicilioResponse{}
	if err := response.ParseResponse(result); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	if dt, err := dateparser.Parse(nil, result.Fecha); err == nil {
		result.ParsedDate = dt.Time
	} else {
		result.DateError = err
	}

	return result, nil
}
