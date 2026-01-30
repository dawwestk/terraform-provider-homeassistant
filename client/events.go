package client

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// GetEvents retrieves a list of all event types that can be fired.
func (c *Client) GetEvents() ([]Event, error) {
	body, err := c.doRequest("GET", "/events", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}

	var events []Event
	if err := json.Unmarshal(body, &events); err != nil {
		return nil, fmt.Errorf("failed to parse events response: %w", err)
	}

	return events, nil
}

// FireEventResponse represents the response from firing an event.
type FireEventResponse struct {
	Message string `json:"message"`
}

// FireEvent fires an event with the specified type and optional data.
func (c *Client) FireEvent(eventType string, eventData map[string]interface{}) (*FireEventResponse, error) {
	var payload []byte
	var err error

	if eventData != nil {
		payload, err = json.Marshal(eventData)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal event data: %w", err)
		}
	}

	endpoint := fmt.Sprintf("/events/%s", url.PathEscape(eventType))
	body, err := c.doRequest("POST", endpoint, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to fire event %s: %w", eventType, err)
	}

	var response FireEventResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse fire event response: %w", err)
	}

	return &response, nil
}
