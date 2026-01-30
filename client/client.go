package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Client holds the configuration for connecting to a Home Assistant instance.
type Client struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
}

// NewClient creates a new Home Assistant API client.
// It reads configuration from environment variables:
//   - HA_BEARER_TOKEN (required): Long-lived access token
//   - HA_HOST_NAME (required): IP address or hostname of the HA instance
//   - HA_PORT (optional): Port number, defaults to 8123
func NewClient() (*Client, error) {
	token := os.Getenv("HA_BEARER_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("HA_BEARER_TOKEN environment variable is required")
	}

	hostName := os.Getenv("HA_HOST_NAME")
	if hostName == "" {
		return nil, fmt.Errorf("HA_HOST_NAME environment variable is required")
	}

	port := os.Getenv("HA_PORT")
	if port == "" {
		port = "8123"
	}

	baseURL := fmt.Sprintf("http://%s:%s/api", hostName, port)

	return &Client{
		BaseURL: baseURL,
		Token:   token,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// doRequest executes an HTTP request with proper authentication headers.
func (c *Client) doRequest(method, endpoint string, body []byte) ([]byte, error) {
	url := fmt.Sprintf("%s%s", c.BaseURL, endpoint)

	var req *http.Request
	var err error

	if body != nil {
		req, err = http.NewRequest(method, url, bytes.NewReader(body))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}
