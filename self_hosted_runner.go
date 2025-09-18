package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSelfHostedRunner() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages GitHub self-hosted runners with JIT configuration.",
		CreateContext: resourceSelfHostedRunnerCreate,
		ReadContext:   resourceSelfHostedRunnerRead,
		UpdateContext: resourceSelfHostedRunnerUpdate,
		DeleteContext: resourceSelfHostedRunnerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the self-hosted runner",
			},
			"runner_group_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the runner group to add the runner to",
			},
			"readonly_labels": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Read-only labels to be set during runner creation (cannot be modified after creation)",
			},
			"labels": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Custom labels to add to the runner",
			},
			"work_folder": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "_work",
				Description: "Working directory for job execution",
			},
			"os": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Operating system of the runner",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the runner",
			},
			"busy": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the runner is busy",
			},
			"ephemeral": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the runner is ephemeral",
			},
			"all_labels": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "All labels associated with the runner",
			},
			"encoded_jit_config": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Encoded JIT configuration for the runner",
			},
		},
	}
}

func dataSourceSelfHostedRunner() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieves a GitHub self-hosted runner by name.",
		ReadContext: dataSourceSelfHostedRunnerRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the self-hosted runner",
			},
			"runner_group_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the runner group to search in",
			},
			"os": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Operating system of the runner",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the runner",
			},
			"busy": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the runner is busy",
			},
			"ephemeral": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the runner is ephemeral",
			},
			"all_labels": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "All labels associated with the runner",
			},
		},
	}
}

func resourceSelfHostedRunnerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	name := d.Get("name").(string)
	runnerGroupID := d.Get("runner_group_id").(int)
	readonlyLabels := expandStringList(d.Get("readonly_labels").([]interface{}))
	labels := expandStringList(d.Get("labels").([]interface{}))
	workFolder := d.Get("work_folder").(string)

	// Validate readonly_labels are not empty
	if len(readonlyLabels) == 0 {
		return diag.Errorf("readonly_labels cannot be empty")
	}

	// Validate labels are not empty
	if len(labels) == 0 {
		return diag.Errorf("labels cannot be empty")
	}

	// Create JIT configuration for the runner with readonly_labels
	req := &JITConfigRequest{
		Name:           name,
		RunnerGroupID:  runnerGroupID,
		ReadOnlyLabels: readonlyLabels,
		WorkFolder:     workFolder,
	}

	var result JITConfigResponse
	err := client.Post(ctx, fmt.Sprintf("/orgs/%s/actions/runners/generate-jitconfig", client.organization), req, &result)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(result.Runner.ID))
	d.Set("encoded_jit_config", result.EncodedJITConfig)

	// Set custom labels manually after runner creation
	setReq := &SetLabelsRequest{
		Labels: labels,
	}

	err = client.Put(ctx, fmt.Sprintf("/orgs/%s/actions/runners/%s/labels", client.organization, strconv.Itoa(result.Runner.ID)), setReq, nil)
	if err != nil {
		return diag.Errorf("failed to set runner labels: %v", err)
	}

	return resourceSelfHostedRunnerRead(ctx, d, m)
}

func resourceSelfHostedRunnerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	runnerID := d.Id()
	var runner SelfHostedRunner
	err := client.Get(ctx, fmt.Sprintf("/orgs/%s/actions/runners/%s", client.organization, runnerID), &runner)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", runner.Name)
	d.Set("os", runner.OS)
	d.Set("status", runner.Status)
	d.Set("busy", runner.Busy)
	d.Set("ephemeral", runner.Ephemeral)

	// Extract all label names
	labelNames := make([]string, len(runner.Labels))
	for i, label := range runner.Labels {
		labelNames[i] = label.Name
	}
	d.Set("all_labels", labelNames)

	return nil
}

func resourceSelfHostedRunnerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	runnerID := d.Id()

	// Handle runner group changes
	if d.HasChange("runner_group_id") {
		oldGroupID, newGroupID := d.GetChange("runner_group_id")
		oldGroup := oldGroupID.(int)
		newGroup := newGroupID.(int)

		// Remove from old group
		if oldGroup > 0 {
			err := client.Delete(ctx, fmt.Sprintf("/orgs/%s/actions/runner-groups/%d/runners/%s", client.organization, oldGroup, runnerID), nil)
			if err != nil {
				d.Set("runner_group_id", oldGroup)
				return diag.FromErr(err)
			}
		}

		// Add to new group
		if newGroup > 0 {
			err := client.Put(ctx, fmt.Sprintf("/orgs/%s/actions/runner-groups/%d/runners/%s", client.organization, newGroup, runnerID), nil, nil)
			if err != nil {
				d.Set("runner_group_id", nil)
				return diag.FromErr(err)
			}
		}
	}

	// Get current runner to check read-only labels
	var currentRunner SelfHostedRunner
	err := client.Get(ctx, fmt.Sprintf("/orgs/%s/actions/runners/%s", client.organization, runnerID), &currentRunner)
	if err != nil {
		return diag.Errorf("failed to get current runner: %v", err)
	}

	// Identify read-only labels from the runner (labels that are not custom)
	readOnlyLabelsFromRunner := make(map[string]bool)
	for _, label := range currentRunner.Labels {
		// Read-only labels have type "read-only" or are system labels
		if label.Type == "read-only" {
			readOnlyLabelsFromRunner[label.Name] = true
		}
	}

	// Check if readonly_labels has changed and validate no read-only labels are being removed
	if d.HasChange("readonly_labels") {
		oldReadonlyLabels, newReadonlyLabels := d.GetChange("readonly_labels")
		oldReadonlyLabelList := expandStringList(oldReadonlyLabels.([]interface{}))
		newReadonlyLabelList := expandStringList(newReadonlyLabels.([]interface{}))

		// Validate readonly_labels are not empty
		if len(newReadonlyLabelList) == 0 {
			d.Set("readonly_labels", oldReadonlyLabelList)
			return diag.Errorf("readonly_labels cannot be empty")
		}

		// Check if any read-only labels from the runner are being removed
		for _, oldLabel := range oldReadonlyLabelList {
			if readOnlyLabelsFromRunner[oldLabel] {
				// Check if this read-only label is still in the new list
				found := false
				for _, newLabel := range newReadonlyLabelList {
					if newLabel == oldLabel {
						found = true
						break
					}
				}
				if !found {
					d.Set("readonly_labels", oldReadonlyLabelList)
					return diag.Errorf("cannot remove read-only label: %s", oldLabel)
				}
			}
		}
	}

	// Update labels if changed
	if d.HasChange("labels") {
		oldLabels, newLabels := d.GetChange("labels")
		oldLabelList := expandStringList(oldLabels.([]interface{}))
		newLabelList := expandStringList(newLabels.([]interface{}))

		// Validate labels are not empty
		if len(newLabelList) == 0 {
			d.Set("labels", oldLabelList)
			return diag.Errorf("labels cannot be empty")
		}

		// Filter out read-only labels from the request (only send custom labels)
		customLabels := make([]string, 0)
		for _, label := range newLabelList {
			if !readOnlyLabelsFromRunner[label] {
				customLabels = append(customLabels, label)
			}
		}

		// Set labels (only custom labels)
		setReq := &SetLabelsRequest{
			Labels: customLabels,
		}

		err = client.Put(ctx, fmt.Sprintf("/orgs/%s/actions/runners/%s/labels", client.organization, runnerID), setReq, nil)
		if err != nil {
			d.Set("labels", oldLabelList)
			return diag.Errorf("failed to update runner labels: %v", err)
		}
	}

	return resourceSelfHostedRunnerRead(ctx, d, m)
}

func resourceSelfHostedRunnerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	runnerID := d.Id()

	// Get current runner to check status
	var currentRunner SelfHostedRunner
	err := client.Get(ctx, fmt.Sprintf("/orgs/%s/actions/runners/%s", client.organization, runnerID), &currentRunner)
	if err != nil {
		return diag.Errorf("failed to get runner status: %v", err)
	}

	// Validate that runner is offline before deletion
	if currentRunner.Status != "offline" {
		return diag.Errorf("cannot delete runner: runner must be offline before deletion, current status: %s", currentRunner.Status)
	}

	// Delete the runner from GitHub
	err = client.Delete(ctx, fmt.Sprintf("/orgs/%s/actions/runners/%s", client.organization, runnerID), nil)
	if err != nil {
		return diag.Errorf("failed to delete runner: %v", err)
	}

	d.SetId("")
	return nil
}

func dataSourceSelfHostedRunnerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	name := d.Get("name").(string)
	runnerGroupID := d.Get("runner_group_id").(int)

	// Search for runner by name
	var runnerList SelfHostedRunnerList
	path := fmt.Sprintf("/orgs/%s/actions/runners", client.organization)
	if runnerGroupID > 0 {
		path = fmt.Sprintf("/orgs/%s/actions/runner-groups/%d/runners", client.organization, runnerGroupID)
	}

	err := client.Get(ctx, path, &runnerList)
	if err != nil {
		return diag.FromErr(err)
	}

	var foundRunner *SelfHostedRunner
	for _, runner := range runnerList.Runners {
		if runner.Name == name {
			foundRunner = &runner
			break
		}
	}

	if foundRunner == nil {
		return diag.Errorf("Self-hosted runner with name '%s' not found", name)
	}

	d.SetId(strconv.Itoa(foundRunner.ID))
	d.Set("name", foundRunner.Name)
	d.Set("os", foundRunner.OS)
	d.Set("status", foundRunner.Status)
	d.Set("busy", foundRunner.Busy)
	d.Set("ephemeral", foundRunner.Ephemeral)

	// Extract all label names
	labelNames := make([]string, len(foundRunner.Labels))
	for i, label := range foundRunner.Labels {
		labelNames[i] = label.Name
	}
	d.Set("all_labels", labelNames)

	return nil
}

func dataSourceRunnerApplications() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieves available runner applications for download.",
		ReadContext: dataSourceRunnerApplicationsRead,
		Schema: map[string]*schema.Schema{
			"applications": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of runner applications available for download",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"os": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Operating system",
						},
						"architecture": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Architecture",
						},
						"download_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Download URL for the runner application",
						},
						"filename": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Filename of the runner application",
						},
					},
				},
			},
		},
	}
}

func dataSourceRegistrationToken() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieves a registration token for the organization.",
		ReadContext: dataSourceRegistrationTokenRead,
		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Registration token for the organization",
			},
			"expires_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Expiration time of the token",
			},
		},
	}
}

func dataSourceRemoveToken() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieves a remove token for the organization.",
		ReadContext: dataSourceRemoveTokenRead,
		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Remove token for the organization",
			},
			"expires_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Expiration time of the token",
			},
		},
	}
}

func dataSourceRunnerApplicationsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	var applications []RunnerApplication
	err := client.Get(ctx, fmt.Sprintf("/orgs/%s/actions/runners/downloads", client.organization), &applications)
	if err != nil {
		return diag.FromErr(err)
	}

	applicationList := make([]map[string]interface{}, len(applications))
	for i, app := range applications {
		applicationList[i] = map[string]interface{}{
			"os":           app.OS,
			"architecture": app.Architecture,
			"download_url": app.DownloadURL,
			"filename":     app.Filename,
		}
	}

	d.SetId("runner-applications")
	d.Set("applications", applicationList)

	return nil
}

func dataSourceRegistrationTokenRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	var token RegistrationToken
	err := client.Post(ctx, fmt.Sprintf("/orgs/%s/actions/runners/registration-token", client.organization), nil, &token)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("registration-token")
	d.Set("token", token.Token)
	d.Set("expires_at", token.ExpiresAt)

	return nil
}

func dataSourceRemoveTokenRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	var token RemoveToken
	err := client.Post(ctx, fmt.Sprintf("/orgs/%s/actions/runners/remove-token", client.organization), nil, &token)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("remove-token")
	d.Set("token", token.Token)
	d.Set("expires_at", token.ExpiresAt)

	return nil
}
