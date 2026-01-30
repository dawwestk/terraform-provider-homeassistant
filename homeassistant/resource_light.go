package homeassistant

import (
	"context"
	"fmt"

	"github.com/dawwestk/terraform-provider-homeassistant/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceLight() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLightCreate,
		ReadContext:   resourceLightRead,
		UpdateContext: resourceLightUpdate,
		DeleteContext: resourceLightDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"entity_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The entity ID of the light (e.g., light.living_room).",
			},
			"state": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"on", "off"}, false),
				Description:  "Desired state of the light: 'on' or 'off'. If not specified, reads current state from Home Assistant.",
			},
			"brightness": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 255),
				Description:  "Brightness level (0-255).",
			},
			"brightness_pct": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 100),
				Description:  "Brightness percentage (0-100).",
			},
			"rgb_color": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    3,
				MinItems:    3,
				Description: "RGB color as a list of 3 integers [R, G, B], each 0-255.",
				Elem: &schema.Schema{
					Type:         schema.TypeInt,
					ValidateFunc: validation.IntBetween(0, 255),
				},
			},
			"color_temp_kelvin": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(2000, 6500),
				Description:  "Color temperature in Kelvin (2000-6500).",
			},
			"transition": {
				Type:         schema.TypeFloat,
				Optional:     true,
				ValidateFunc: validation.FloatBetween(0, 300),
				Description:  "Transition time in seconds (0-300).",
			},
			"effect": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Light effect name.",
			},
		},
	}
}

func resourceLightCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	entityID := d.Get("entity_id").(string)
	state := d.Get("state").(string)

	// If state is not specified, default to "on"
	if state == "" {
		state = "on"
	}

	// Build service data
	serviceData := buildLightServiceData(d)

	// Call the appropriate service
	var service string
	if state == "on" {
		service = "turn_on"
	} else {
		service = "turn_off"
	}

	_, err := c.CallService("light", service, serviceData)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to set light state: %w", err))
	}

	d.SetId(entityID)

	return resourceLightRead(ctx, d, m)
}

func resourceLightRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	var diags diag.Diagnostics

	entityID := d.Id()

	haState, err := c.GetState(entityID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read light state: %w", err))
	}

	// Set entity_id if not already set (happens during import)
	if d.Get("entity_id").(string) == "" {
		d.Set("entity_id", entityID)
	}

	// Set the state from Home Assistant
	d.Set("state", haState.State)

	// Update brightness from attributes if available
	if brightness, ok := haState.Attributes["brightness"]; ok {
		if b, ok := brightness.(float64); ok {
			d.Set("brightness", int(b))
		}
	}

	// Update rgb_color from attributes if available
	if rgbColor, ok := haState.Attributes["rgb_color"]; ok {
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

	// Update color_temp_kelvin from attributes if available
	if colorTemp, ok := haState.Attributes["color_temp_kelvin"]; ok {
		if ct, ok := colorTemp.(float64); ok {
			d.Set("color_temp_kelvin", int(ct))
		}
	}

	// Update effect from attributes if available
	if effect, ok := haState.Attributes["effect"]; ok {
		if e, ok := effect.(string); ok {
			d.Set("effect", e)
		}
	}

	return diags
}

func resourceLightUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	state := d.Get("state").(string)

	// If state is not specified, default to "on"
	if state == "" {
		state = "on"
	}

	// Build service data
	serviceData := buildLightServiceData(d)

	// Call the appropriate service
	var service string
	if state == "on" {
		service = "turn_on"
	} else {
		service = "turn_off"
	}

	_, err := c.CallService("light", service, serviceData)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update light state: %w", err))
	}

	return resourceLightRead(ctx, d, m)
}

func resourceLightDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	var diags diag.Diagnostics

	entityID := d.Get("entity_id").(string)

	// Turn off the light when the resource is deleted
	serviceData := map[string]interface{}{
		"entity_id": entityID,
	}

	_, err := c.CallService("light", "turn_off", serviceData)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to turn off light: %w", err))
	}

	d.SetId("")

	return diags
}

// buildLightServiceData builds the service data map from the resource data.
// It only includes values that are explicitly configured by the user (not computed).
func buildLightServiceData(d *schema.ResourceData) map[string]interface{} {
	entityID := d.Get("entity_id").(string)

	serviceData := map[string]interface{}{
		"entity_id": entityID,
	}

	// Get the raw config to check what's actually configured vs computed
	rawConfig := d.GetRawConfig()

	// brightness - check if explicitly configured
	if !rawConfig.GetAttr("brightness").IsNull() {
		serviceData["brightness"] = d.Get("brightness").(int)
	}

	// brightness_pct - check if explicitly configured
	if !rawConfig.GetAttr("brightness_pct").IsNull() {
		serviceData["brightness_pct"] = d.Get("brightness_pct").(int)
	}

	// rgb_color - check if explicitly configured
	if !rawConfig.GetAttr("rgb_color").IsNull() {
		if v, ok := d.GetOk("rgb_color"); ok {
			rgbList := v.([]interface{})
			if len(rgbList) == 3 {
				rgb := make([]int, 3)
				for i, val := range rgbList {
					rgb[i] = val.(int)
				}
				serviceData["rgb_color"] = rgb
			}
		}
	}

	// color_temp_kelvin - check if explicitly configured
	if !rawConfig.GetAttr("color_temp_kelvin").IsNull() {
		serviceData["color_temp_kelvin"] = d.Get("color_temp_kelvin").(int)
	}

	// transition - check if explicitly configured
	if !rawConfig.GetAttr("transition").IsNull() {
		serviceData["transition"] = d.Get("transition").(float64)
	}

	// effect - check if explicitly configured
	if !rawConfig.GetAttr("effect").IsNull() {
		serviceData["effect"] = d.Get("effect").(string)
	}

	return serviceData
}
