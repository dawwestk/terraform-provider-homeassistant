package client

// State represents the state of an entity in Home Assistant.
type State struct {
	EntityID    string                 `json:"entity_id"`
	State       string                 `json:"state"`
	Attributes  map[string]interface{} `json:"attributes,omitempty"`
	LastChanged string                 `json:"last_changed,omitempty"`
	LastUpdated string                 `json:"last_updated,omitempty"`
	Context     *Context               `json:"context,omitempty"`
}

// Context represents the context of a state change.
type Context struct {
	ID       string `json:"id"`
	ParentID string `json:"parent_id,omitempty"`
	UserID   string `json:"user_id,omitempty"`
}

// Config represents the Home Assistant configuration.
type Config struct {
	Latitude        float64  `json:"latitude"`
	Longitude       float64  `json:"longitude"`
	Elevation       int      `json:"elevation"`
	UnitSystem      UnitSystem `json:"unit_system"`
	LocationName    string   `json:"location_name"`
	TimeZone        string   `json:"time_zone"`
	Components      []string `json:"components"`
	ConfigDir       string   `json:"config_dir"`
	WhitelistExternalDirs []string `json:"whitelist_external_dirs"`
	AllowlistExternalDirs []string `json:"allowlist_external_dirs"`
	AllowlistExternalURLs []string `json:"allowlist_external_urls"`
	Version         string   `json:"version"`
	ConfigSource    string   `json:"config_source"`
	SafeMode        bool     `json:"safe_mode"`
	State           string   `json:"state"`
	ExternalURL     string   `json:"external_url,omitempty"`
	InternalURL     string   `json:"internal_url,omitempty"`
	Currency        string   `json:"currency"`
	Country         string   `json:"country,omitempty"`
	Language        string   `json:"language"`
}

// UnitSystem represents the unit system configuration.
type UnitSystem struct {
	Length            string `json:"length"`
	AccumulatedPrecipitation string `json:"accumulated_precipitation"`
	Mass              string `json:"mass"`
	Pressure          string `json:"pressure"`
	Temperature       string `json:"temperature"`
	Volume            string `json:"volume"`
	WindSpeed         string `json:"wind_speed"`
}

// ServiceDomain represents a domain with its available services.
type ServiceDomain struct {
	Domain   string                `json:"domain"`
	Services map[string]ServiceDef `json:"services"`
}

// ServiceDef represents the definition of a service.
type ServiceDef struct {
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	Fields      map[string]ServiceField `json:"fields,omitempty"`
	Target      *ServiceTarget         `json:"target,omitempty"`
}

// ServiceField represents a field in a service definition.
type ServiceField struct {
	Name        string      `json:"name,omitempty"`
	Description string      `json:"description,omitempty"`
	Required    bool        `json:"required,omitempty"`
	Example     interface{} `json:"example,omitempty"`
	Selector    interface{} `json:"selector,omitempty"`
}

// ServiceTarget represents the target specification for a service.
type ServiceTarget struct {
	Entity []TargetEntity `json:"entity,omitempty"`
	Device []TargetDevice `json:"device,omitempty"`
	Area   []TargetArea   `json:"area,omitempty"`
}

// TargetEntity represents entity targeting criteria.
type TargetEntity struct {
	Domain      string `json:"domain,omitempty"`
	Integration string `json:"integration,omitempty"`
}

// TargetDevice represents device targeting criteria.
type TargetDevice struct {
	Integration string `json:"integration,omitempty"`
}

// TargetArea represents area targeting criteria.
type TargetArea struct{}

// Event represents an event type in Home Assistant.
type Event struct {
	Event         string `json:"event"`
	ListenerCount int    `json:"listener_count"`
}

// APIStatus represents the API status response.
type APIStatus struct {
	Message string `json:"message"`
}

// ServiceCallRequest represents a request to call a service.
type ServiceCallRequest struct {
	EntityID string                 `json:"entity_id,omitempty"`
	Data     map[string]interface{} `json:"-"`
}

// StateUpdateRequest represents a request to update an entity state.
type StateUpdateRequest struct {
	State      string                 `json:"state"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
}
