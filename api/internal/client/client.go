package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

// InternalClient is the HTTP client for making internal API calls
type InternalClient struct {
	BaseURL    string
	Secret     string
	HTTPClient *http.Client
}

// InternalError represents an error response from internal APIs
type InternalError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e *InternalError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// NewInternalClient creates a new internal HTTP client
func NewInternalClient() *InternalClient {
	baseURL := os.Getenv("INTERNAL_SERVICE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	return &InternalClient{
		BaseURL: baseURL,
		Secret:  os.Getenv("INTERNAL_API_SECRET"),
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Do performs an HTTP request to an internal endpoint
func (c *InternalClient) Do(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.BaseURL+path, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Internal-API-Key", c.Secret)

	// Copy request ID for tracing if available in context
	if requestID, ok := ctx.Value("request_id").(string); ok && requestID != "" {
		req.Header.Set("X-Request-ID", requestID)
	} else {
		// Generate a new request ID if not available
		req.Header.Set("X-Request-ID", uuid.New().String())
	}

	return c.HTTPClient.Do(req)
}

// Get performs a GET request to an internal endpoint
func (c *InternalClient) Get(ctx context.Context, path string) (*http.Response, error) {
	return c.Do(ctx, http.MethodGet, path, nil)
}

// Post performs a POST request to an internal endpoint
func (c *InternalClient) Post(ctx context.Context, path string, body interface{}) (*http.Response, error) {
	return c.Do(ctx, http.MethodPost, path, body)
}

// Put performs a PUT request to an internal endpoint
func (c *InternalClient) Put(ctx context.Context, path string, body interface{}) (*http.Response, error) {
	return c.Do(ctx, http.MethodPut, path, body)
}

// ParseResponse parses the response body into the given interface
func ParseResponse[T any](resp *http.Response) (*T, error) {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for error response
	if resp.StatusCode >= 400 {
		var errorResp struct {
			Error InternalError `json:"error"`
		}
		if err := json.Unmarshal(body, &errorResp); err == nil && errorResp.Error.Code != "" {
			return nil, &errorResp.Error
		}
		return nil, &InternalError{
			Code:    fmt.Sprintf("HTTP_%d", resp.StatusCode),
			Message: string(body),
		}
	}

	var result T
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// ParseResponseWithData parses a response with a "data" wrapper
func ParseResponseWithData[T any](resp *http.Response) (*T, error) {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for error response
	if resp.StatusCode >= 400 {
		var errorResp struct {
			Error InternalError `json:"error"`
		}
		if err := json.Unmarshal(body, &errorResp); err == nil && errorResp.Error.Code != "" {
			return nil, &errorResp.Error
		}
		return nil, &InternalError{
			Code:    fmt.Sprintf("HTTP_%d", resp.StatusCode),
			Message: string(body),
		}
	}

	var wrapper struct {
		Data T `json:"data"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &wrapper.Data, nil
}
