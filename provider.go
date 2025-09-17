package main

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("GITHUB_TOKEN", nil),
				Description: "The GitHub personal access token",
			},
			"base_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("GITHUB_BASE_URL", "https://api.github.com"),
				Description: "The GitHub base URL",
			},
			"organization": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The GitHub organization name",
			},
			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to use insecure connections",
			},
			"app_auth": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "GitHub App authentication configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The GitHub App ID",
						},
						"installation_id": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The GitHub App installation ID",
						},
						"pem_file": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Path to the GitHub App private key PEM file",
						},
					},
				},
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"azure-github-runners_network_configuration": resourceNetworkConfiguration(),
			"azure-github-runners_network_settings":      resourceNetworkSettings(),
			"azure-github-runners_runner_group":          resourceRunnerGroup(),
			"azure-github-runners_self_hosted_runner":    resourceSelfHostedRunner(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"azure-github-runners_network_configuration": dataSourceNetworkConfiguration(),
			"azure-github-runners_network_settings":      dataSourceNetworkSettings(),
			"azure-github-runners_runner_group":          dataSourceRunnerGroup(),
			"azure-github-runners_self_hosted_runner":    dataSourceSelfHostedRunner(),
			"azure-github-runners_runner_applications":   dataSourceRunnerApplications(),
			"azure-github-runners_registration_token":    dataSourceRegistrationToken(),
			"azure-github-runners_remove_token":          dataSourceRemoveToken(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	token := d.Get("token").(string)
	baseURL := d.Get("base_url").(string)
	organization := d.Get("organization").(string)
	insecure := d.Get("insecure").(bool)

	if organization == "" {
		return nil, diag.Errorf("GitHub organization is required")
	}

	// Check if GitHub App authentication is configured
	appAuthList := d.Get("app_auth").([]interface{})
	var appAuth *AppAuth
	if len(appAuthList) > 0 {
		appAuthConfig := appAuthList[0].(map[string]interface{})
		appAuth = &AppAuth{
			ID:             appAuthConfig["id"].(int),
			InstallationID: appAuthConfig["installation_id"].(int),
			PEMFile:        appAuthConfig["pem_file"].(string),
		}
	}

	// Either token or app_auth must be provided
	if token == "" && appAuth == nil {
		return nil, diag.Errorf("Either GitHub token or app_auth configuration is required")
	}

	config := &Config{
		Token:        token,
		BaseURL:      baseURL,
		Organization: organization,
		Insecure:     insecure,
		AppAuth:      appAuth,
	}

	client, err := config.Client()
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return client, nil
}

type Config struct {
	Token        string
	BaseURL      string
	Organization string
	Insecure     bool
	AppAuth      *AppAuth
}

type AppAuth struct {
	ID             int
	InstallationID int
	PEMFile        string
}

func (c *Config) Client() (*Client, error) {
	if c.AppAuth != nil {
		return NewClientWithAppAuth(c.AppAuth, c.BaseURL, c.Organization, c.Insecure)
	}
	return NewClient(c.Token, c.BaseURL, c.Organization, c.Insecure)
}
