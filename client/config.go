package client

import (
	"encoding/json"
	"fmt"
)

// Health checks if the Home Assistant API is running.
// Returns the API status message on success.
func (c *Client) Health() (*APIStatus, error) {
	body, err := c.doRequest("GET", "/", nil)
	if err != nil {
		return nil, fmt.Errorf("health check failed: %w", err)
	}

	var status APIStatus
	if err := json.Unmarshal(body, &status); err != nil {
		return nil, fmt.Errorf("failed to parse health response: %w", err)
	}

	return &status, nil
}

// GetConfig retrieves the current Home Assistant configuration.
func (c *Client) GetConfig() (*Config, error) {
	body, err := c.doRequest("GET", "/config", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	var config Config
	if err := json.Unmarshal(body, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config response: %w", err)
	}

	return &config, nil
}
