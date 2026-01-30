package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestNewClient_Success(t *testing.T) {
	// Set required environment variables
	os.Setenv("HA_BEARER_TOKEN", "test-token")
	os.Setenv("HA_HOST_NAME", "localhost")
	os.Setenv("HA_PORT", "8123")
	defer func() {
		os.Unsetenv("HA_BEARER_TOKEN")
		os.Unsetenv("HA_HOST_NAME")
		os.Unsetenv("HA_PORT")
	}()

	client, err := NewClient()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if client == nil {
		t.Fatal("expected client to not be nil")
	}

	expectedBaseURL := "http://localhost:8123/api"
	if client.BaseURL != expectedBaseURL {
		t.Errorf("expected BaseURL %s, got %s", expectedBaseURL, client.BaseURL)
	}

	if client.Token != "test-token" {
		t.Errorf("expected Token 'test-token', got %s", client.Token)
	}
}

func TestNewClient_DefaultPort(t *testing.T) {
	os.Setenv("HA_BEARER_TOKEN", "test-token")
	os.Setenv("HA_HOST_NAME", "192.168.1.100")
	os.Unsetenv("HA_PORT")
	defer func() {
		os.Unsetenv("HA_BEARER_TOKEN")
		os.Unsetenv("HA_HOST_NAME")
	}()

	client, err := NewClient()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedBaseURL := "http://192.168.1.100:8123/api"
	if client.BaseURL != expectedBaseURL {
		t.Errorf("expected BaseURL %s, got %s", expectedBaseURL, client.BaseURL)
	}
}

func TestNewClient_MissingToken(t *testing.T) {
	os.Unsetenv("HA_BEARER_TOKEN")
	os.Setenv("HA_HOST_NAME", "localhost")
	defer os.Unsetenv("HA_HOST_NAME")

	_, err := NewClient()
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestNewClient_MissingIPAddress(t *testing.T) {
	os.Setenv("HA_BEARER_TOKEN", "test-token")
	os.Unsetenv("HA_HOST_NAME")
	defer os.Unsetenv("HA_BEARER_TOKEN")

	_, err := NewClient()
	if err == nil {
		t.Fatal("expected error for missing IP address")
	}
}

// createTestClient creates a client pointing to a test server
func createTestClient(server *httptest.Server) *Client {
	return &Client{
		BaseURL:    server.URL,
		Token:      "test-token",
		HTTPClient: server.Client(),
	}
}

func TestClient_Health(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.URL.Path != "/" {
			t.Errorf("expected path '/', got %s", r.URL.Path)
		}
		if r.Method != "GET" {
			t.Errorf("expected GET method, got %s", r.Method)
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("expected Bearer token header, got %s", r.Header.Get("Authorization"))
		}

		// Return response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(APIStatus{Message: "API running."})
	}))
	defer server.Close()

	client := createTestClient(server)
	status, err := client.Health()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if status.Message != "API running." {
		t.Errorf("expected message 'API running.', got %s", status.Message)
	}
}

func TestClient_GetConfig(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/config" {
			t.Errorf("expected path '/config', got %s", r.URL.Path)
		}

		config := Config{
			LocationName: "Home",
			TimeZone:     "Europe/London",
			Version:      "2024.1.0",
			Latitude:     51.5074,
			Longitude:    -0.1278,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(config)
	}))
	defer server.Close()

	client := createTestClient(server)
	config, err := client.GetConfig()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if config.LocationName != "Home" {
		t.Errorf("expected LocationName 'Home', got %s", config.LocationName)
	}
	if config.Version != "2024.1.0" {
		t.Errorf("expected Version '2024.1.0', got %s", config.Version)
	}
}

func TestClient_GetStates(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/states" {
			t.Errorf("expected path '/states', got %s", r.URL.Path)
		}

		states := []State{
			{EntityID: "light.living_room", State: "on"},
			{EntityID: "light.bedroom", State: "off"},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(states)
	}))
	defer server.Close()

	client := createTestClient(server)
	states, err := client.GetStates()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(states) != 2 {
		t.Fatalf("expected 2 states, got %d", len(states))
	}
	if states[0].EntityID != "light.living_room" {
		t.Errorf("expected entity_id 'light.living_room', got %s", states[0].EntityID)
	}
}

func TestClient_GetState(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/states/light.living_room" {
			t.Errorf("expected path '/states/light.living_room', got %s", r.URL.Path)
		}

		state := State{
			EntityID: "light.living_room",
			State:    "on",
			Attributes: map[string]interface{}{
				"brightness":    float64(255),
				"friendly_name": "Living Room Light",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(state)
	}))
	defer server.Close()

	client := createTestClient(server)
	state, err := client.GetState("light.living_room")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if state.EntityID != "light.living_room" {
		t.Errorf("expected entity_id 'light.living_room', got %s", state.EntityID)
	}
	if state.State != "on" {
		t.Errorf("expected state 'on', got %s", state.State)
	}
}

func TestClient_SetState(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/states/light.living_room" {
			t.Errorf("expected path '/states/light.living_room', got %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("expected POST method, got %s", r.Method)
		}

		// Decode request body
		var req StateUpdateRequest
		json.NewDecoder(r.Body).Decode(&req)
		if req.State != "on" {
			t.Errorf("expected state 'on', got %s", req.State)
		}

		// Return updated state
		state := State{
			EntityID: "light.living_room",
			State:    "on",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(state)
	}))
	defer server.Close()

	client := createTestClient(server)
	state, err := client.SetState("light.living_room", StateUpdateRequest{State: "on"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if state.State != "on" {
		t.Errorf("expected state 'on', got %s", state.State)
	}
}

func TestClient_CallService(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/services/light/turn_on" {
			t.Errorf("expected path '/services/light/turn_on', got %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("expected POST method, got %s", r.Method)
		}

		// Decode request body
		var req map[string]interface{}
		json.NewDecoder(r.Body).Decode(&req)
		if req["entity_id"] != "light.living_room" {
			t.Errorf("expected entity_id 'light.living_room', got %v", req["entity_id"])
		}

		// Return affected states
		states := []State{
			{EntityID: "light.living_room", State: "on"},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(states)
	}))
	defer server.Close()

	client := createTestClient(server)
	states, err := client.CallService("light", "turn_on", map[string]interface{}{
		"entity_id": "light.living_room",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(states) != 1 {
		t.Fatalf("expected 1 state, got %d", len(states))
	}
	if states[0].State != "on" {
		t.Errorf("expected state 'on', got %s", states[0].State)
	}
}

func TestClient_GetEvents(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/events" {
			t.Errorf("expected path '/events', got %s", r.URL.Path)
		}

		events := []Event{
			{Event: "state_changed", ListenerCount: 10},
			{Event: "call_service", ListenerCount: 5},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(events)
	}))
	defer server.Close()

	client := createTestClient(server)
	events, err := client.GetEvents()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
}

func TestClient_FireEvent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/events/custom_event" {
			t.Errorf("expected path '/events/custom_event', got %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("expected POST method, got %s", r.Method)
		}

		response := FireEventResponse{Message: "Event custom_event fired."}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := createTestClient(server)
	resp, err := client.FireEvent("custom_event", map[string]interface{}{"data": "test"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.Message != "Event custom_event fired." {
		t.Errorf("expected message 'Event custom_event fired.', got %s", resp.Message)
	}
}

func TestClient_ErrorResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Entity not found"))
	}))
	defer server.Close()

	client := createTestClient(server)
	_, err := client.GetState("light.nonexistent")
	if err == nil {
		t.Fatal("expected error for 404 response")
	}
}

func TestClient_GetServices(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/services" {
			t.Errorf("expected path '/services', got %s", r.URL.Path)
		}

		services := []ServiceDomain{
			{
				Domain: "light",
				Services: map[string]ServiceDef{
					"turn_on":  {Name: "Turn on", Description: "Turn on a light"},
					"turn_off": {Name: "Turn off", Description: "Turn off a light"},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(services)
	}))
	defer server.Close()

	client := createTestClient(server)
	services, err := client.GetServices()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(services) != 1 {
		t.Fatalf("expected 1 service domain, got %d", len(services))
	}
	if services[0].Domain != "light" {
		t.Errorf("expected domain 'light', got %s", services[0].Domain)
	}
	if len(services[0].Services) != 2 {
		t.Errorf("expected 2 services, got %d", len(services[0].Services))
	}
}
