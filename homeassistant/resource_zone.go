package homeassistant

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/dawwestk/terraform-provider-homeassistant/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceZone() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceZoneCreate,
		ReadContext:   resourceZoneRead,
		UpdateContext: resourceZoneUpdate,
		DeleteContext: resourceZoneDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Friendly name of the zone.",
			},
			"latitude": {
				Type:         schema.TypeFloat,
				Required:     true,
				ValidateFunc: validation.FloatBetween(-90, 90),
				Description:  "Latitude coordinate of the zone center.",
			},
			"longitude": {
				Type:         schema.TypeFloat,
				Required:     true,
				ValidateFunc: validation.FloatBetween(-180, 180),
				Description:  "Longitude coordinate of the zone center.",
			},
			"radius": {
				Type:         schema.TypeFloat,
				Optional:     true,
				Default:      100.0,
				ValidateFunc: validation.FloatAtLeast(0),
				Description:  "Zone radius in meters. Defaults to 100.",
			},
			"passive": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, the zone will not trigger automations. Defaults to false.",
			},
			"icon": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "mdi:map-marker",
				Description: "MDI icon for the zone (e.g., mdi:home). Defaults to mdi:map-marker.",
			},
			// Computed attributes
			"entity_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The entity ID of the zone (e.g., zone.home).",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Number of persons currently in the zone.",
			},
			"editable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the zone is editable.",
			},
		},
	}
}

// slugify converts a name to a slug suitable for entity_id
func slugify(name string) string {
	// Convert to lowercase
	slug := strings.ToLower(name)
	// Replace spaces and special chars with underscores
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	slug = reg.ReplaceAllString(slug, "_")
	// Trim underscores from ends
	slug = strings.Trim(slug, "_")
	return slug
}

func resourceZoneCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	name := d.Get("name").(string)
	entityID := fmt.Sprintf("zone.%s", slugify(name))

	// Build state update request
	attributes := map[string]interface{}{
		"friendly_name": name,
		"latitude":      d.Get("latitude").(float64),
		"longitude":     d.Get("longitude").(float64),
		"radius":        d.Get("radius").(float64),
		"passive":       d.Get("passive").(bool),
		"icon":          d.Get("icon").(string),
		"editable":      true,
		"persons":       []string{},
	}

	stateReq := client.StateUpdateRequest{
		State:      "0",
		Attributes: attributes,
	}

	_, err := c.SetState(entityID, stateReq)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create zone: %w", err))
	}

	d.SetId(entityID)

	return resourceZoneRead(ctx, d, m)
}

func resourceZoneRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	var diags diag.Diagnostics

	entityID := d.Id()

	state, err := c.GetState(entityID)
	if err != nil {
		// If the zone doesn't exist, remove it from state
		d.SetId("")
		return diags
	}

	d.Set("entity_id", entityID)
	d.Set("state", state.State)

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

	return diags
}

func resourceZoneUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	entityID := d.Id()
	name := d.Get("name").(string)

	// Build state update request with current values
	attributes := map[string]interface{}{
		"friendly_name": name,
		"latitude":      d.Get("latitude").(float64),
		"longitude":     d.Get("longitude").(float64),
		"radius":        d.Get("radius").(float64),
		"passive":       d.Get("passive").(bool),
		"icon":          d.Get("icon").(string),
		"editable":      true,
	}

	// Preserve existing persons list
	currentState, err := c.GetState(entityID)
	if err == nil {
		if persons, ok := currentState.Attributes["persons"]; ok {
			attributes["persons"] = persons
		}
	}

	stateReq := client.StateUpdateRequest{
		State:      d.Get("state").(string),
		Attributes: attributes,
	}

	// If state is empty, default to "0"
	if stateReq.State == "" {
		stateReq.State = "0"
	}

	_, err = c.SetState(entityID, stateReq)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to update zone: %w", err))
	}

	return resourceZoneRead(ctx, d, m)
}

func resourceZoneDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	var diags diag.Diagnostics

	entityID := d.Id()

	// Set state to "unavailable" to effectively remove the zone
	// Note: HA doesn't have a direct delete API for states, 
	// but setting to unavailable marks it as inactive
	stateReq := client.StateUpdateRequest{
		State:      "unavailable",
		Attributes: map[string]interface{}{},
	}

	_, err := c.SetState(entityID, stateReq)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete zone: %w", err))
	}

	d.SetId("")

	return diags
}
