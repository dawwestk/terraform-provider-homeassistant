package homeassistant

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// getTestLightEntityID returns the light entity ID to use for tests.
// Set HA_TEST_LIGHT_ENTITY env var to specify a real light entity.
func getTestLightEntityID() string {
	if v := os.Getenv("HA_TEST_LIGHT_ENTITY"); v != "" {
		return v
	}
	return "light.test_light"
}

func TestResourceLight_Schema(t *testing.T) {
	s := resourceLight().Schema

	// Test required fields
	requiredFields := []string{"entity_id"}
	for _, field := range requiredFields {
		if !s[field].Required {
			t.Errorf("expected %s to be required", field)
		}
	}

	// Test optional fields
	optionalFields := []string{"state", "brightness", "brightness_pct", "rgb_color", "color_temp_kelvin", "transition", "effect"}
	for _, field := range optionalFields {
		if s[field].Required {
			t.Errorf("expected %s to be optional", field)
		}
	}
}

func TestResourceLight_HasImporter(t *testing.T) {
	r := resourceLight()
	if r.Importer == nil {
		t.Error("expected resource to have an importer")
	}
}

func TestResourceLight_BrightnessValidation(t *testing.T) {
	s := resourceLight().Schema["brightness"]

	if s.ValidateFunc == nil {
		t.Fatal("expected brightness to have validation")
	}

	// Test valid values
	validValues := []int{0, 128, 255}
	for _, v := range validValues {
		warns, errs := s.ValidateFunc(v, "brightness")
		if len(errs) > 0 {
			t.Errorf("expected %d to be valid, got errors: %v", v, errs)
		}
		if len(warns) > 0 {
			t.Errorf("expected %d to have no warnings, got: %v", v, warns)
		}
	}
}

func TestResourceLight_BrightnessPctValidation(t *testing.T) {
	s := resourceLight().Schema["brightness_pct"]

	if s.ValidateFunc == nil {
		t.Fatal("expected brightness_pct to have validation")
	}

	// Test valid values
	validValues := []int{0, 50, 100}
	for _, v := range validValues {
		warns, errs := s.ValidateFunc(v, "brightness_pct")
		if len(errs) > 0 {
			t.Errorf("expected %d to be valid, got errors: %v", v, errs)
		}
		if len(warns) > 0 {
			t.Errorf("expected %d to have no warnings, got: %v", v, warns)
		}
	}
}

func TestResourceLight_StateValidation(t *testing.T) {
	s := resourceLight().Schema["state"]

	if s.ValidateFunc == nil {
		t.Fatal("expected state to have validation")
	}

	// Test valid values
	validValues := []string{"on", "off"}
	for _, v := range validValues {
		warns, errs := s.ValidateFunc(v, "state")
		if len(errs) > 0 {
			t.Errorf("expected '%s' to be valid, got errors: %v", v, errs)
		}
		if len(warns) > 0 {
			t.Errorf("expected '%s' to have no warnings, got: %v", v, warns)
		}
	}

	// Test invalid value
	warns, errs := s.ValidateFunc("invalid", "state")
	if len(errs) == 0 {
		t.Error("expected 'invalid' to fail validation")
	}
	_ = warns // warnings may or may not be present
}

func TestResourceLight_RGBColorMaxItems(t *testing.T) {
	s := resourceLight().Schema["rgb_color"]

	if s.MaxItems != 3 {
		t.Errorf("expected rgb_color MaxItems to be 3, got %d", s.MaxItems)
	}
	if s.MinItems != 3 {
		t.Errorf("expected rgb_color MinItems to be 3, got %d", s.MinItems)
	}
}

// Acceptance tests - require a real Home Assistant instance
// Run with: TF_ACC=1 HA_TEST_LIGHT_ENTITY=light.your_light go test -v ./homeassistant/

func testAccLightPreCheck(t *testing.T) {
	testAccPreCheck(t)
	if v := os.Getenv("HA_TEST_LIGHT_ENTITY"); v == "" {
		t.Skip("HA_TEST_LIGHT_ENTITY must be set for light acceptance tests")
	}
}

func TestAccResourceLight_basic(t *testing.T) {
	entityID := getTestLightEntityID()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccLightPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceLightConfig_basic(entityID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("homeassistant_light.test", "entity_id", entityID),
					// Just check entity_id is set, state may vary due to timing
					resource.TestCheckResourceAttrSet("homeassistant_light.test", "state"),
				),
			},
		},
	})
}

func TestAccResourceLight_withBrightness(t *testing.T) {
	entityID := getTestLightEntityID()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccLightPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceLightConfig_withBrightness(entityID, 200),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("homeassistant_light.test", "entity_id", entityID),
					resource.TestCheckResourceAttrSet("homeassistant_light.test", "state"),
					resource.TestCheckResourceAttrSet("homeassistant_light.test", "brightness"),
				),
			},
		},
	})
}

func TestAccResourceLight_update(t *testing.T) {
	entityID := getTestLightEntityID()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccLightPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceLightConfig_withBrightness(entityID, 100),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("homeassistant_light.test", "entity_id", entityID),
				),
			},
			{
				Config: testAccResourceLightConfig_withBrightness(entityID, 255),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("homeassistant_light.test", "entity_id", entityID),
				),
			},
		},
	})
}

func TestAccResourceLight_import(t *testing.T) {
	entityID := getTestLightEntityID()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccLightPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceLightConfig_basic(entityID),
			},
			{
				ResourceName:            "homeassistant_light.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"state", "brightness", "rgb_color"},
			},
		},
	})
}

func testAccResourceLightConfig_basic(entityID string) string {
	return fmt.Sprintf(`
resource "homeassistant_light" "test" {
  entity_id = %q
  state     = "on"
}
`, entityID)
}

func testAccResourceLightConfig_withBrightness(entityID string, brightness int) string {
	return fmt.Sprintf(`
resource "homeassistant_light" "test" {
  entity_id  = %q
  state      = "on"
  brightness = %d
}
`, entityID, brightness)
}

// Note: buildLightServiceData cannot be easily unit tested because it relies on
// d.GetRawConfig() which requires a full Terraform config context that isn't
// available via schema.TestResourceDataRaw. This function is tested indirectly
// through the acceptance tests.
