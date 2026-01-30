package homeassistant

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestDataSourceLight_Schema(t *testing.T) {
	s := dataSourceLight().Schema

	// Test required field
	if !s["entity_id"].Required {
		t.Error("expected entity_id to be required")
	}

	// Test computed fields
	computedFields := []string{
		"state", "brightness", "rgb_color", "color_temp_kelvin",
		"effect", "friendly_name", "color_mode", "supported_color_modes",
		"last_changed", "last_updated",
	}
	for _, field := range computedFields {
		if !s[field].Computed {
			t.Errorf("expected %s to be computed", field)
		}
	}
}

func TestDataSourceLight_RGBColorIsListOfInt(t *testing.T) {
	s := dataSourceLight().Schema["rgb_color"]

	if s.Type.String() != "TypeList" {
		t.Errorf("expected rgb_color to be TypeList, got %s", s.Type.String())
	}
}

func TestDataSourceLight_SupportedColorModesIsListOfString(t *testing.T) {
	s := dataSourceLight().Schema["supported_color_modes"]

	if s.Type.String() != "TypeList" {
		t.Errorf("expected supported_color_modes to be TypeList, got %s", s.Type.String())
	}
}

// Acceptance tests - require a real Home Assistant instance

func TestAccDataSourceLight_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceLightConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.homeassistant_light.test", "state"),
					resource.TestCheckResourceAttrSet("data.homeassistant_light.test", "friendly_name"),
				),
			},
		},
	})
}

func testAccDataSourceLightConfig_basic() string {
	return `
data "homeassistant_light" "test" {
  entity_id = "light.desk_light"
}
`
}
