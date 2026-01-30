package homeassistant

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// testAccProviderFactories is used for acceptance tests
var testAccProviderFactories map[string]func() (*schema.Provider, error)

// testAccProvider is the provider instance used in tests
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviderFactories = map[string]func() (*schema.Provider, error){
		"homeassistant": func() (*schema.Provider, error) {
			return testAccProvider, nil
		},
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func TestProvider_HasExpectedResources(t *testing.T) {
	expectedResources := []string{
		"homeassistant_light",
		"homeassistant_zone",
	}

	provider := Provider()
	for _, name := range expectedResources {
		if _, ok := provider.ResourcesMap[name]; !ok {
			t.Errorf("expected resource %s to be registered", name)
		}
	}
}

func TestProvider_HasExpectedDataSources(t *testing.T) {
	expectedDataSources := []string{
		"homeassistant_light",
		"homeassistant_zone",
	}

	provider := Provider()
	for _, name := range expectedDataSources {
		if _, ok := provider.DataSourcesMap[name]; !ok {
			t.Errorf("expected data source %s to be registered", name)
		}
	}
}

func TestProvider_SchemaHasExpectedAttributes(t *testing.T) {
	expectedAttrs := []string{
		"bearer_token",
		"host_name",
		"port",
	}

	provider := Provider()
	for _, attr := range expectedAttrs {
		if _, ok := provider.Schema[attr]; !ok {
			t.Errorf("expected schema attribute %s to be defined", attr)
		}
	}
}

func TestProvider_BearerTokenIsSensitive(t *testing.T) {
	provider := Provider()
	if !provider.Schema["bearer_token"].Sensitive {
		t.Error("expected bearer_token to be marked as sensitive")
	}
}

func TestProvider_PortHasDefault(t *testing.T) {
	provider := Provider()
	portSchema := provider.Schema["port"]

	// Clear env var to test default
	os.Unsetenv("HA_PORT")

	defaultVal, err := portSchema.DefaultFunc()
	if err != nil {
		t.Fatalf("error getting default: %v", err)
	}

	if defaultVal != "8123" {
		t.Errorf("expected default port '8123', got %v", defaultVal)
	}
}

// testAccPreCheck validates the necessary test API keys exist
// in the testing environment
func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("HA_BEARER_TOKEN"); v == "" {
		t.Skip("HA_BEARER_TOKEN must be set for acceptance tests")
	}
	if v := os.Getenv("HA_HOST_NAME"); v == "" {
		t.Skip("HA_HOST_NAME must be set for acceptance tests")
	}
}
