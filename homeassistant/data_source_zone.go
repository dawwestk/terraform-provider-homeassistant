package homeassistant

import (
	"context"
	"fmt"

	"github.com/dawwestk/terraform-provider-homeassistant/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceZone() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceZoneRead,

		Schema: map[string]*schema.Schema{
			"entity_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The entity ID of the zone (e.g., zone.home).",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Friendly name of the zone.",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Number of persons currently in the zone.",
			},
			"latitude": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "Latitude coordinate of the zone center.",
			},
			"longitude": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "Longitude coordinate of the zone center.",
			},
			"radius": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "Zone radius in meters.",
			},
			"passive": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If true, the zone will not trigger automations.",
			},
			"icon": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "MDI icon for the zone.",
			},
			"editable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the zone is editable.",
			},
			"persons": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of person entity IDs currently in the zone.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"last_changed": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timestamp of the last state change.",
			},
			"last_updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timestamp of the last update.",
			},
		},
	}
}

func dataSourceZoneRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	var diags diag.Diagnostics

	entityID := d.Get("entity_id").(string)

	state, err := c.GetState(entityID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read zone state: %w", err))
	}

	d.SetId(entityID)
	d.Set("state", state.State)
	d.Set("last_changed", state.LastChanged)
	d.Set("last_updated", state.LastUpdated)

	// Extract attributes
	if friendlyName, ok := state.Attributes["friendly_name"]; ok {
		if fn, ok := friendlyName.(string); ok {
			d.Set("name", fn)
		}
	}

	if latitude, ok := state.Attributes["latitude"]; ok {
		if lat, ok := latitude.(float64); ok {
			d.Set("latitude", lat)
		}
	}

	if longitude, ok := state.Attributes["longitude"]; ok {
		if lon, ok := longitude.(float64); ok {
			d.Set("longitude", lon)
		}
	}

	if radius, ok := state.Attributes["radius"]; ok {
		if r, ok := radius.(float64); ok {
			d.Set("radius", r)
		}
	}

	if passive, ok := state.Attributes["passive"]; ok {
		if p, ok := passive.(bool); ok {
			d.Set("passive", p)
		}
	}

	if icon, ok := state.Attributes["icon"]; ok {
		if i, ok := icon.(string); ok {
			d.Set("icon", i)
		}
	}

	if editable, ok := state.Attributes["editable"]; ok {
		if e, ok := editable.(bool); ok {
			d.Set("editable", e)
		}
	}

	if persons, ok := state.Attributes["persons"]; ok {
		if pList, ok := persons.([]interface{}); ok {
			personStrings := make([]string, len(pList))
			for i, p := range pList {
				if ps, ok := p.(string); ok {
					personStrings[i] = ps
				}
			}
			d.Set("persons", personStrings)
		}
	}

	return diags
}
