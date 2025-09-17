package main

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceNetworkConfiguration() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages GitHub organization network configurations for Actions.",
		CreateContext: resourceNetworkConfigurationCreate,
		ReadContext:   resourceNetworkConfigurationRead,
		UpdateContext: resourceNetworkConfigurationUpdate,
		DeleteContext: resourceNetworkConfigurationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the network configuration",
			},
			"compute_service": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "actions",
				ValidateFunc: validation.StringInSlice([]string{"none", "actions"}, false),
				Description:  "The hosted compute service to use for the network configuration",
			},
			"network_settings_ids": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The identifier of the network settings to use for the network configuration",
			},
			"created_on": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation timestamp of the network configuration",
			},
		},
	}
}

func dataSourceNetworkConfiguration() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieves a GitHub organization network configuration by name.",
		ReadContext: dataSourceNetworkConfigurationRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the network configuration",
			},
			"compute_service": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The hosted compute service to use for the network configuration",
			},
			"network_settings_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The identifier of the network settings to use for the network configuration",
			},
			"created_on": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation timestamp of the network configuration",
			},
		},
	}
}

func resourceNetworkConfigurationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	req := &CreateNetworkConfigurationRequest{
		Name:               d.Get("name").(string),
		ComputeService:     d.Get("compute_service").(string),
		NetworkSettingsIDs: expandStringList(d.Get("network_settings_ids").([]interface{})),
	}

	var result NetworkConfiguration
	err := client.Post(ctx, fmt.Sprintf("/orgs/%s/settings/network-configurations", client.organization), req, &result)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(result.ID)
	return resourceNetworkConfigurationRead(ctx, d, m)
}

func resourceNetworkConfigurationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	networkConfigID := d.Id()
	var config NetworkConfiguration
	err := client.Get(ctx, fmt.Sprintf("/orgs/%s/settings/network-configurations/%s", client.organization, networkConfigID), &config)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", config.Name)
	d.Set("compute_service", config.ComputeService)
	d.Set("network_settings_ids", config.NetworkSettingsIDs)
	d.Set("created_on", config.CreatedOn)

	return nil
}

func resourceNetworkConfigurationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	networkConfigID := d.Id()
	req := &UpdateNetworkConfigurationRequest{
		Name:               d.Get("name").(string),
		ComputeService:     d.Get("compute_service").(string),
		NetworkSettingsIDs: expandStringList(d.Get("network_settings_ids").([]interface{})),
	}

	var result NetworkConfiguration
	err := client.Patch(ctx, fmt.Sprintf("/orgs/%s/settings/network-configurations/%s", client.organization, networkConfigID), req, &result)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNetworkConfigurationRead(ctx, d, m)
}

func resourceNetworkConfigurationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	networkConfigID := d.Id()
	err := client.Delete(ctx, fmt.Sprintf("/orgs/%s/settings/network-configurations/%s", client.organization, networkConfigID), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func dataSourceNetworkConfigurationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	name := d.Get("name").(string)

	// Get all network configurations and find the one with matching name
	var configList NetworkConfigurationList
	err := client.Get(ctx, fmt.Sprintf("/orgs/%s/settings/network-configurations", client.organization), &configList)
	if err != nil {
		return diag.FromErr(err)
	}

	var foundConfig *NetworkConfiguration
	for _, config := range configList.NetworkConfigurations {
		if config.Name == name {
			foundConfig = &config
			break
		}
	}

	if foundConfig == nil {
		return diag.Errorf("Network configuration with name '%s' not found", name)
	}

	d.SetId(foundConfig.ID)
	d.Set("name", foundConfig.Name)
	d.Set("compute_service", foundConfig.ComputeService)
	d.Set("network_settings_ids", foundConfig.NetworkSettingsIDs)
	d.Set("created_on", foundConfig.CreatedOn)

	return nil
}

func expandStringList(configured []interface{}) []string {
	vs := make([]string, 0, len(configured))
	for _, v := range configured {
		vs = append(vs, v.(string))
	}
	return vs
}
