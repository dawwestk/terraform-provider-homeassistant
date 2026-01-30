package homeassistant

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestResourceZone_Schema(t *testing.T) {
	s := resourceZone().Schema

	// Test required fields
	requiredFields := []string{"name", "latitude", "longitude"}
	for _, field := range requiredFields {
		if !s[field].Required {
			t.Errorf("expected %s to be required", field)
		}
	}

	// Test optional fields
	optionalFields := []string{"radius", "passive", "icon"}
	for _, field := range optionalFields {
		if s[field].Required {
			t.Errorf("expected %s to be optional", field)
		}
	}

	// Test computed fields
	computedFields := []string{"entity_id", "state", "editable"}
	for _, field := range computedFields {
		if !s[field].Computed {
			t.Errorf("expected %s to be computed", field)
		}
	}
}

func TestResourceZone_HasImporter(t *testing.T) {
	r := resourceZone()
	if r.Importer == nil {
		t.Error("expected resource to have an importer")
	}
}

func TestResourceZone_LatitudeValidation(t *testing.T) {
	s := resourceZone().Schema["latitude"]

	if s.ValidateFunc == nil {
		t.Fatal("expected latitude to have validation")
	}

	// Test valid values
	validValues := []float64{-90.0, 0.0, 45.5, 90.0}
	for _, v := range validValues {
		warns, errs := s.ValidateFunc(v, "latitude")
		if len(errs) > 0 {
			t.Errorf("expected %f to be valid, got errors: %v", v, errs)
		}
		if len(warns) > 0 {
			t.Errorf("expected %f to have no warnings, got: %v", v, warns)
		}
	}
}

func TestResourceZone_LongitudeValidation(t *testing.T) {
	s := resourceZone().Schema["longitude"]

	if s.ValidateFunc == nil {
		t.Fatal("expected longitude to have validation")
	}

	// Test valid values
	validValues := []float64{-180.0, 0.0, 90.5, 180.0}
	for _, v := range validValues {
		warns, errs := s.ValidateFunc(v, "longitude")
		if len(errs) > 0 {
			t.Errorf("expected %f to be valid, got errors: %v", v, errs)
		}
		if len(warns) > 0 {
			t.Errorf("expected %f to have no warnings, got: %v", v, warns)
		}
	}
}

func TestResourceZone_RadiusDefault(t *testing.T) {
	s := resourceZone().Schema["radius"]

	if s.Default != 100.0 {
		t.Errorf("expected radius default to be 100.0, got %v", s.Default)
	}
}

func TestResourceZone_PassiveDefault(t *testing.T) {
	s := resourceZone().Schema["passive"]

	if s.Default != false {
		t.Errorf("expected passive default to be false, got %v", s.Default)
	}
}

func TestResourceZone_IconDefault(t *testing.T) {
	s := resourceZone().Schema["icon"]

	if s.Default != "mdi:map-marker" {
		t.Errorf("expected icon default to be 'mdi:map-marker', got %v", s.Default)
	}
}

func TestSlugify(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Home", "home"},
		{"Living Room", "living_room"},
		{"Office 2", "office_2"},
		{"Test-Zone", "test_zone"},
		{"  Spaces  ", "spaces"},
		{"Special!@#Chars", "special_chars"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := slugify(tt.input)
			if result != tt.expected {
				t.Errorf("slugify(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

// Acceptance tests - require a real Home Assistant instance
// Run with: TF_ACC=1 go test -v ./homeassistant/

func TestAccResourceZone_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceZoneConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("homeassistant_zone.test", "name", "Test Zone"),
					resource.TestCheckResourceAttr("homeassistant_zone.test", "latitude", "51.5074"),
					resource.TestCheckResourceAttr("homeassistant_zone.test", "longitude", "-0.1278"),
				),
			},
		},
	})
}

func TestAccResourceZone_complete(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceZoneConfig_complete(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("homeassistant_zone.test", "name", "Office"),
					resource.TestCheckResourceAttr("homeassistant_zone.test", "latitude", "51.5074"),
					resource.TestCheckResourceAttr("homeassistant_zone.test", "longitude", "-0.1278"),
					resource.TestCheckResourceAttr("homeassistant_zone.test", "radius", "50"),
					resource.TestCheckResourceAttr("homeassistant_zone.test", "passive", "true"),
					resource.TestCheckResourceAttr("homeassistant_zone.test", "icon", "mdi:office-building"),
				),
			},
		},
	})
}

func TestAccResourceZone_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceZoneConfig_withRadius(50),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("homeassistant_zone.test", "radius", "50"),
				),
			},
			{
				Config: testAccResourceZoneConfig_withRadius(100),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("homeassistant_zone.test", "radius", "100"),
				),
			},
		},
	})
}

func TestAccResourceZone_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceZoneConfig_basic(),
			},
			{
				ResourceName:      "homeassistant_zone.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourceZoneConfig_basic() string {
	return `
resource "homeassistant_zone" "test" {
  name      = "Test Zone"
  latitude  = 51.5074
  longitude = -0.1278
}
`
}

func testAccResourceZoneConfig_complete() string {
	return `
resource "homeassistant_zone" "test" {
  name      = "Office"
  latitude  = 51.5074
  longitude = -0.1278
  radius    = 50
  passive   = true
  icon      = "mdi:office-building"
}
`
}

func testAccResourceZoneConfig_withRadius(radius int) string {
	return fmt.Sprintf(`
resource "homeassistant_zone" "test" {
  name      = "Test Zone"
  latitude  = 51.5074
  longitude = -0.1278
  radius    = %d
}
`, radius)
}
