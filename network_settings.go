package main

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNetworkSettings() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkSettingsCreate,
		ReadContext:   resourceNetworkSettingsRead,
		UpdateContext: resourceNetworkSettingsUpdate,
		DeleteContext: resourceNetworkSettingsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the network settings",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The subnet ID for the network settings",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The region for the network settings",
			},
			"network_configuration_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The network configuration ID this settings belongs to",
			},
		},
	}
}

func dataSourceNetworkSettings() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieves GitHub organization network settings by ID.",
		ReadContext: dataSourceNetworkSettingsRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the network settings",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The subnet ID for the network settings",
			},
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The region for the network settings",
			},
			"network_configuration_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The network configuration ID this settings belongs to",
			},
		},
	}
}

func resourceNetworkSettingsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Note: GitHub doesn't provide an API to create network settings directly
	// This resource assumes the network settings are created through other means
	// and we're just managing their configuration
	return diag.Errorf("Network settings cannot be created directly through the API. They must be created through the GitHub UI or other means.")
}

func resourceNetworkSettingsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	networkSettingsID := d.Id()
	var settings NetworkSettings
	err := client.Get(ctx, fmt.Sprintf("/orgs/%s/settings/network-settings/%s", client.organization, networkSettingsID), &settings)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", settings.Name)
	d.Set("subnet_id", settings.SubnetID)
	d.Set("region", settings.Region)
	d.Set("network_configuration_id", settings.NetworkConfigurationID)

	return nil
}

func resourceNetworkSettingsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Note: GitHub doesn't provide an API to update network settings directly
	return diag.Errorf("Network settings cannot be updated directly through the API")
}

func resourceNetworkSettingsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Note: GitHub doesn't provide an API to delete network settings directly
	return diag.Errorf("Network settings cannot be deleted directly through the API")
}

func dataSourceNetworkSettingsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Note: GitHub doesn't provide an API to list network settings
	// This would need to be implemented based on the actual API if available
	return diag.Errorf("Network settings lookup by name is not supported by the GitHub API")
}
