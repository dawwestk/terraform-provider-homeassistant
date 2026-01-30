package client

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// GetStates retrieves the state of all entities.
func (c *Client) GetStates() ([]State, error) {
	body, err := c.doRequest("GET", "/states", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get states: %w", err)
	}

	var states []State
	if err := json.Unmarshal(body, &states); err != nil {
		return nil, fmt.Errorf("failed to parse states response: %w", err)
	}

	return states, nil
}

// GetState retrieves the state of a specific entity.
func (c *Client) GetState(entityID string) (*State, error) {
	endpoint := fmt.Sprintf("/states/%s", url.PathEscape(entityID))
	body, err := c.doRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get state for %s: %w", entityID, err)
	}

	var state State
	if err := json.Unmarshal(body, &state); err != nil {
		return nil, fmt.Errorf("failed to parse state response: %w", err)
	}

	return &state, nil
}

// SetState updates the state of an entity.
// This creates the entity if it doesn't exist.
func (c *Client) SetState(entityID string, req StateUpdateRequest) (*State, error) {
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal state update request: %w", err)
	}

	endpoint := fmt.Sprintf("/states/%s", url.PathEscape(entityID))
	body, err := c.doRequest("POST", endpoint, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to set state for %s: %w", entityID, err)
	}

	var state State
	if err := json.Unmarshal(body, &state); err != nil {
		return nil, fmt.Errorf("failed to parse state response: %w", err)
	}

	return &state, nil
}
