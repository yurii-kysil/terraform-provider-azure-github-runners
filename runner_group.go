package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceRunnerGroup() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages GitHub self-hosted runner groups.",
		CreateContext: resourceRunnerGroupCreate,
		ReadContext:   resourceRunnerGroupRead,
		UpdateContext: resourceRunnerGroupUpdate,
		DeleteContext: resourceRunnerGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the runner group",
			},
			"visibility": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "all",
				ValidateFunc: validation.StringInSlice([]string{"all", "selected", "private"}, false),
				Description:  "Visibility of the runner group",
			},
			"selected_repository_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: "List of repository IDs that can access the runner group",
			},
			"runners": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: "List of runner IDs in the group",
			},
			"allows_public_repositories": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether public repositories can use the runner group",
			},
			"restricted_to_workflows": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether the runner group is restricted to specific workflows",
			},
			"selected_workflows": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of workflows that can use the runner group",
			},
			"network_configuration_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The identifier of a hosted compute network configuration",
			},
			"default": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether this is the default runner group",
			},
			"inherited": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the runner group is inherited",
			},
			"selected_repositories_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL for selected repositories",
			},
			"runners_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL for runners",
			},
			"hosted_runners_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL for hosted runners",
			},
			"workflow_restrictions_read_only": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether workflow restrictions are read-only",
			},
		},
	}
}

func dataSourceRunnerGroup() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieves a GitHub self-hosted runner group by name.",
		ReadContext: dataSourceRunnerGroupRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the runner group",
			},
			"visibility": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Visibility of the runner group",
			},
			"selected_repository_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: "List of repository IDs that can access the runner group",
			},
			"runners": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: "List of runner IDs in the group",
			},
			"allows_public_repositories": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether public repositories can use the runner group",
			},
			"restricted_to_workflows": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the runner group is restricted to specific workflows",
			},
			"selected_workflows": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of workflows that can use the runner group",
			},
			"network_configuration_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The identifier of a hosted compute network configuration",
			},
			"default": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether this is the default runner group",
			},
			"inherited": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the runner group is inherited",
			},
			"selected_repositories_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL for selected repositories",
			},
			"runners_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL for runners",
			},
			"hosted_runners_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL for hosted runners",
			},
			"workflow_restrictions_read_only": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether workflow restrictions are read-only",
			},
		},
	}
}

func resourceRunnerGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	visibility := d.Get("visibility").(string)
	selectedRepositoryIDs := expandIntList(d.Get("selected_repository_ids").([]interface{}))

	// Validate visibility and selected_repository_ids
	if visibility == "all" && len(selectedRepositoryIDs) > 0 {
		return diag.Errorf("selected_repository_ids cannot be set when visibility is 'all'")
	}
	if visibility == "selected" && len(selectedRepositoryIDs) == 0 {
		return diag.Errorf("selected_repository_ids cannot be empty when visibility is 'selected'")
	}

	req := &CreateRunnerGroupRequest{
		Name:                     d.Get("name").(string),
		Visibility:               visibility,
		SelectedRepositoryIDs:    selectedRepositoryIDs,
		Runners:                  expandIntList(d.Get("runners").([]interface{})),
		AllowsPublicRepositories: d.Get("allows_public_repositories").(bool),
		RestrictedToWorkflows:    d.Get("restricted_to_workflows").(bool),
		SelectedWorkflows:        expandStringList(d.Get("selected_workflows").([]interface{})),
		NetworkConfigurationID:   d.Get("network_configuration_id").(string),
	}

	var result RunnerGroup
	err := client.Post(ctx, fmt.Sprintf("/orgs/%s/actions/runner-groups", client.organization), req, &result)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(result.ID))
	return resourceRunnerGroupRead(ctx, d, m)
}

func resourceRunnerGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	runnerGroupID := d.Id()
	var runnerGroup RunnerGroup
	err := client.Get(ctx, fmt.Sprintf("/orgs/%s/actions/runner-groups/%s", client.organization, runnerGroupID), &runnerGroup)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", runnerGroup.Name)
	d.Set("visibility", runnerGroup.Visibility)
	d.Set("allows_public_repositories", runnerGroup.AllowsPublicRepositories)
	d.Set("restricted_to_workflows", runnerGroup.RestrictedToWorkflows)
	d.Set("selected_workflows", runnerGroup.SelectedWorkflows)
	d.Set("network_configuration_id", runnerGroup.NetworkConfigurationID)
	d.Set("default", runnerGroup.Default)
	d.Set("inherited", runnerGroup.Inherited)
	d.Set("selected_repositories_url", runnerGroup.SelectedRepositoriesURL)
	d.Set("runners_url", runnerGroup.RunnersURL)
	d.Set("hosted_runners_url", runnerGroup.HostedRunnersURL)
	d.Set("workflow_restrictions_read_only", runnerGroup.WorkflowRestrictionsReadOnly)

	return nil
}

func resourceRunnerGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	runnerGroupID := d.Id()
	visibility := d.Get("visibility").(string)
	selectedRepositoryIDs := expandIntList(d.Get("selected_repository_ids").([]interface{}))

	// Validate visibility and selected_repository_ids
	if visibility == "all" && len(selectedRepositoryIDs) > 0 {
		return diag.Errorf("selected_repository_ids cannot be set when visibility is 'all'")
	}
	if visibility == "selected" && len(selectedRepositoryIDs) == 0 {
		return diag.Errorf("selected_repository_ids cannot be empty when visibility is 'selected'")
	}

	networkConfigID := d.Get("network_configuration_id").(string)
	req := &UpdateRunnerGroupRequest{
		Name:                     d.Get("name").(string),
		Visibility:               visibility,
		AllowsPublicRepositories: boolPtr(d.Get("allows_public_repositories").(bool)),
		RestrictedToWorkflows:    boolPtr(d.Get("restricted_to_workflows").(bool)),
		SelectedWorkflows:        expandStringList(d.Get("selected_workflows").([]interface{})),
		NetworkConfigurationID:   &networkConfigID,
	}

	var result RunnerGroup
	err := client.Patch(ctx, fmt.Sprintf("/orgs/%s/actions/runner-groups/%s", client.organization, runnerGroupID), req, &result)
	if err != nil {
		return diag.FromErr(err)
	}

	// Update repositories if changed
	if d.HasChange("selected_repository_ids") {
		_, newRepos := d.GetChange("selected_repository_ids")
		newRepoList := expandIntList(newRepos.([]interface{}))

		// Set repositories (replaces the entire list)
		setReq := &SetRepositoriesForRunnerGroupRequest{
			SelectedRepositoryIDs: newRepoList,
		}
		err := client.Put(ctx, fmt.Sprintf("/orgs/%s/actions/runner-groups/%s/repositories", client.organization, runnerGroupID), setReq, nil)
		if err != nil {
			return diag.Errorf("failed to update runner group repositories: %v", err)
		}
	}

	// Update runners if changed
	if d.HasChange("runners") {
		_, newRunners := d.GetChange("runners")
		newRunnerList := expandIntList(newRunners.([]interface{}))

		// Set runners (replaces the entire list)
		setReq := &SetRunnersForRunnerGroupRequest{
			Runners: newRunnerList,
		}
		err := client.Put(ctx, fmt.Sprintf("/orgs/%s/actions/runner-groups/%s/runners", client.organization, runnerGroupID), setReq, nil)
		if err != nil {
			return diag.Errorf("failed to update runner group runners: %v", err)
		}
	}

	return resourceRunnerGroupRead(ctx, d, m)
}

func resourceRunnerGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	runnerGroupID := d.Id()
	err := client.Delete(ctx, fmt.Sprintf("/orgs/%s/actions/runner-groups/%s", client.organization, runnerGroupID), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func dataSourceRunnerGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	name := d.Get("name").(string)

	// Get all runner groups and find the one with matching name
	var runnerGroupList RunnerGroupList
	err := client.Get(ctx, fmt.Sprintf("/orgs/%s/actions/runner-groups", client.organization), &runnerGroupList)
	if err != nil {
		return diag.FromErr(err)
	}

	var foundRunnerGroup *RunnerGroup
	for _, rg := range runnerGroupList.RunnerGroups {
		if rg.Name == name {
			foundRunnerGroup = &rg
			break
		}
	}

	if foundRunnerGroup == nil {
		return diag.Errorf("Runner group with name '%s' not found", name)
	}

	d.SetId(strconv.Itoa(foundRunnerGroup.ID))
	d.Set("name", foundRunnerGroup.Name)
	d.Set("visibility", foundRunnerGroup.Visibility)
	d.Set("allows_public_repositories", foundRunnerGroup.AllowsPublicRepositories)
	d.Set("restricted_to_workflows", foundRunnerGroup.RestrictedToWorkflows)
	d.Set("selected_workflows", foundRunnerGroup.SelectedWorkflows)
	d.Set("network_configuration_id", foundRunnerGroup.NetworkConfigurationID)
	d.Set("default", foundRunnerGroup.Default)
	d.Set("inherited", foundRunnerGroup.Inherited)
	d.Set("selected_repositories_url", foundRunnerGroup.SelectedRepositoriesURL)
	d.Set("runners_url", foundRunnerGroup.RunnersURL)
	d.Set("hosted_runners_url", foundRunnerGroup.HostedRunnersURL)
	d.Set("workflow_restrictions_read_only", foundRunnerGroup.WorkflowRestrictionsReadOnly)

	return nil
}

func expandIntList(configured []interface{}) []int {
	vs := make([]int, 0, len(configured))
	for _, v := range configured {
		vs = append(vs, v.(int))
	}
	return vs
}

func boolPtr(b bool) *bool {
	return &b
}
