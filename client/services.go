package client

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// GetServices retrieves all available services grouped by domain.
func (c *Client) GetServices() ([]ServiceDomain, error) {
	body, err := c.doRequest("GET", "/services", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get services: %w", err)
	}

	var services []ServiceDomain
	if err := json.Unmarshal(body, &services); err != nil {
		return nil, fmt.Errorf("failed to parse services response: %w", err)
	}

	return services, nil
}

// CallService calls a service in a specific domain.
// The serviceData can contain entity_id and any additional service-specific parameters.
func (c *Client) CallService(domain, service string, serviceData map[string]interface{}) ([]State, error) {
	var payload []byte
	var err error

	if serviceData != nil {
		payload, err = json.Marshal(serviceData)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal service data: %w", err)
		}
	}

	endpoint := fmt.Sprintf("/services/%s/%s", url.PathEscape(domain), url.PathEscape(service))
	body, err := c.doRequest("POST", endpoint, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to call service %s.%s: %w", domain, service, err)
	}

	// Service calls return an array of affected states
	var states []State
	if err := json.Unmarshal(body, &states); err != nil {
		return nil, fmt.Errorf("failed to parse service call response: %w", err)
	}

	return states, nil
}
