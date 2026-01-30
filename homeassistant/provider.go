package homeassistant

import (
	"context"
	"fmt"
	"os"

	"github.com/dawwestk/terraform-provider-homeassistant/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider returns the Home Assistant Terraform provider.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"bearer_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("HA_BEARER_TOKEN", nil),
				Description: "Long-lived access token for Home Assistant. Can also be set via HA_BEARER_TOKEN env var.",
			},
			"host_name": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("HA_HOST_NAME", nil),
				Description: "IP address or hostname of the Home Assistant instance. Can also be set via HA_HOST_NAME env var.",
			},
			"port": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("HA_PORT", "8123"),
				Description: "Port of the Home Assistant instance. Defaults to 8123. Can also be set via HA_PORT env var.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"homeassistant_light": resourceLight(),
			"homeassistant_zone":  resourceZone(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"homeassistant_light": dataSourceLight(),
			"homeassistant_zone":  dataSourceZone(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	token := d.Get("bearer_token").(string)
	hostName := d.Get("host_name").(string)
	port := d.Get("port").(string)

	if token == "" {
		return nil, diag.FromErr(fmt.Errorf("bearer_token is required (set via provider config or HA_BEARER_TOKEN env var)"))
	}

	if hostName == "" {
		return nil, diag.FromErr(fmt.Errorf("host_name is required (set via provider config or HA_HOST_NAME env var)"))
	}

	// Set environment variables for the client
	os.Setenv("HA_BEARER_TOKEN", token)
	os.Setenv("HA_HOST_NAME", hostName)
	os.Setenv("HA_PORT", port)

	c, err := client.NewClient()
	if err != nil {
		return nil, diag.FromErr(fmt.Errorf("failed to create Home Assistant client: %w", err))
	}

	// Verify connectivity
	_, err = c.Health()
	if err != nil {
		return nil, diag.FromErr(fmt.Errorf("failed to connect to Home Assistant API: %w", err))
	}

	return c, diags
}
