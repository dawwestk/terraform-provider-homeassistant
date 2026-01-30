package homeassistant

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestDataSourceZone_Schema(t *testing.T) {
	s := dataSourceZone().Schema

	// Test required field
	if !s["entity_id"].Required {
		t.Error("expected entity_id to be required")
	}

	// Test computed fields
	computedFields := []string{
		"name", "state", "latitude", "longitude", "radius",
		"passive", "icon", "editable", "persons",
		"last_changed", "last_updated",
	}
	for _, field := range computedFields {
		if !s[field].Computed {
			t.Errorf("expected %s to be computed", field)
		}
	}
}

func TestDataSourceZone_PersonsIsListOfString(t *testing.T) {
	s := dataSourceZone().Schema["persons"]

	if s.Type.String() != "TypeList" {
		t.Errorf("expected persons to be TypeList, got %s", s.Type.String())
	}
}

func TestDataSourceZone_LatLongAreFloat(t *testing.T) {
	s := dataSourceZone().Schema

	if s["latitude"].Type.String() != "TypeFloat" {
		t.Errorf("expected latitude to be TypeFloat, got %s", s["latitude"].Type.String())
	}
	if s["longitude"].Type.String() != "TypeFloat" {
		t.Errorf("expected longitude to be TypeFloat, got %s", s["longitude"].Type.String())
	}
}

// Acceptance tests - require a real Home Assistant instance

func TestAccDataSourceZone_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceZoneConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.homeassistant_zone.test", "name"),
					resource.TestCheckResourceAttrSet("data.homeassistant_zone.test", "latitude"),
					resource.TestCheckResourceAttrSet("data.homeassistant_zone.test", "longitude"),
				),
			},
		},
	})
}

func TestAccDataSourceZone_home(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceZoneConfig_home(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.homeassistant_zone.home", "entity_id", "zone.home"),
					resource.TestCheckResourceAttrSet("data.homeassistant_zone.home", "name"),
				),
			},
		},
	})
}

func testAccDataSourceZoneConfig_basic() string {
	return `
data "homeassistant_zone" "test" {
  entity_id = "zone.living_room"
}
`
}

func testAccDataSourceZoneConfig_home() string {
	return `
data "homeassistant_zone" "home" {
  entity_id = "zone.home"
}
`
}
