package homeassistant

import (
	"context"
	"fmt"

	"github.com/dawwestk/terraform-provider-homeassistant/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceLight() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLightRead,

		Schema: map[string]*schema.Schema{
			"entity_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The entity ID of the light (e.g., light.living_room).",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current state of the light ('on' or 'off').",
			},
			"brightness": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Current brightness level (0-255).",
			},
			"rgb_color": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Current RGB color as a list of 3 integers [R, G, B].",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"color_temp_kelvin": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Current color temperature in Kelvin.",
			},
			"effect": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current light effect name.",
			},
			"friendly_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Friendly name of the light.",
			},
			"supported_color_modes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of supported color modes.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"color_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current color mode.",
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

func dataSourceLightRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	var diags diag.Diagnostics

	entityID := d.Get("entity_id").(string)

	state, err := c.GetState(entityID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read light state: %w", err))
	}

	d.SetId(entityID)
	d.Set("state", state.State)
	d.Set("last_changed", state.LastChanged)
	d.Set("last_updated", state.LastUpdated)

	// Extract attributes
	if friendlyName, ok := state.Attributes["friendly_name"]; ok {
		if fn, ok := friendlyName.(string); ok {
			d.Set("friendly_name", fn)
		}
	}

	if brightness, ok := state.Attributes["brightness"]; ok {
		if b, ok := brightness.(float64); ok {
			d.Set("brightness", int(b))
		}
	}

	if rgbColor, ok := state.Attributes["rgb_color"]; ok {
		if rgb, ok := rgbColor.([]interface{}); ok && len(rgb) == 3 {
			rgbInts := make([]int, 3)
			for i, v := range rgb {
				if f, ok := v.(float64); ok {
					rgbInts[i] = int(f)
				}
			}
			d.Set("rgb_color", rgbInts)
		}
	}

	if colorTemp, ok := state.Attributes["color_temp_kelvin"]; ok {
		if ct, ok := colorTemp.(float64); ok {
			d.Set("color_temp_kelvin", int(ct))
		}
	}

	if effect, ok := state.Attributes["effect"]; ok {
		if e, ok := effect.(string); ok {
			d.Set("effect", e)
		}
	}

	if colorMode, ok := state.Attributes["color_mode"]; ok {
		if cm, ok := colorMode.(string); ok {
			d.Set("color_mode", cm)
		}
	}

	if supportedModes, ok := state.Attributes["supported_color_modes"]; ok {
		if modes, ok := supportedModes.([]interface{}); ok {
			modeStrings := make([]string, len(modes))
			for i, m := range modes {
				if s, ok := m.(string); ok {
					modeStrings[i] = s
				}
			}
			d.Set("supported_color_modes", modeStrings)
		}
	}

	return diags
}
