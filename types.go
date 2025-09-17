package main

// NetworkConfiguration represents a GitHub organization network configuration
type NetworkConfiguration struct {
	ID                 string   `json:"id,omitempty"`
	Name               string   `json:"name"`
	ComputeService     string   `json:"compute_service"`
	NetworkSettingsIDs []string `json:"network_settings_ids"`
	CreatedOn          string   `json:"created_on,omitempty"`
}

// NetworkConfigurationList represents the response for listing network configurations
type NetworkConfigurationList struct {
	TotalCount            int                    `json:"total_count"`
	NetworkConfigurations []NetworkConfiguration `json:"network_configurations"`
}

// NetworkSettings represents a GitHub organization network settings
type NetworkSettings struct {
	ID                     string `json:"id,omitempty"`
	NetworkConfigurationID string `json:"network_configuration_id,omitempty"`
	Name                   string `json:"name"`
	SubnetID               string `json:"subnet_id"`
	Region                 string `json:"region"`
}

// CreateNetworkConfigurationRequest represents the request to create a network configuration
type CreateNetworkConfigurationRequest struct {
	Name               string   `json:"name"`
	ComputeService     string   `json:"compute_service,omitempty"`
	NetworkSettingsIDs []string `json:"network_settings_ids"`
}

// UpdateNetworkConfigurationRequest represents the request to update a network configuration
type UpdateNetworkConfigurationRequest struct {
	Name               string   `json:"name,omitempty"`
	ComputeService     string   `json:"compute_service,omitempty"`
	NetworkSettingsIDs []string `json:"network_settings_ids,omitempty"`
}

// RunnerGroup represents a GitHub self-hosted runner group
type RunnerGroup struct {
	ID                           int      `json:"id,omitempty"`
	Name                         string   `json:"name"`
	Visibility                   string   `json:"visibility"`
	Default                      bool     `json:"default"`
	SelectedRepositoriesURL      string   `json:"selected_repositories_url,omitempty"`
	RunnersURL                   string   `json:"runners_url,omitempty"`
	HostedRunnersURL             string   `json:"hosted_runners_url,omitempty"`
	NetworkConfigurationID       string   `json:"network_configuration_id,omitempty"`
	Inherited                    bool     `json:"inherited,omitempty"`
	AllowsPublicRepositories     bool     `json:"allows_public_repositories,omitempty"`
	RestrictedToWorkflows        bool     `json:"restricted_to_workflows,omitempty"`
	SelectedWorkflows            []string `json:"selected_workflows,omitempty"`
	WorkflowRestrictionsReadOnly bool     `json:"workflow_restrictions_read_only,omitempty"`
}

// RunnerGroupList represents the response for listing runner groups
type RunnerGroupList struct {
	TotalCount   int           `json:"total_count"`
	RunnerGroups []RunnerGroup `json:"runner_groups"`
}

// SelfHostedRunner represents a GitHub self-hosted runner
type SelfHostedRunner struct {
	ID        int    `json:"id,omitempty"`
	Name      string `json:"name"`
	OS        string `json:"os,omitempty"`
	Status    string `json:"status,omitempty"`
	Busy      bool   `json:"busy,omitempty"`
	Ephemeral bool   `json:"ephemeral,omitempty"`
	Labels    []struct {
		ID   int    `json:"id,omitempty"`
		Name string `json:"name"`
		Type string `json:"type,omitempty"`
	} `json:"labels,omitempty"`
}

// RunnerApplication represents a GitHub runner application download
type RunnerApplication struct {
	OS           string `json:"os,omitempty"`
	Architecture string `json:"architecture,omitempty"`
	DownloadURL  string `json:"download_url,omitempty"`
	Filename     string `json:"filename,omitempty"`
}

// JITConfigRequest represents the request to create a JIT configuration
type JITConfigRequest struct {
	Name          string   `json:"name"`
	RunnerGroupID int      `json:"runner_group_id"`
	Labels        []string `json:"labels"`
	WorkFolder    string   `json:"work_folder,omitempty"`
}

// JITConfigResponse represents the response for JIT configuration
type JITConfigResponse struct {
	Runner           SelfHostedRunner `json:"runner"`
	EncodedJITConfig string           `json:"encoded_jit_config"`
}

// RegistrationToken represents a GitHub registration token
type RegistrationToken struct {
	Token     string `json:"token,omitempty"`
	ExpiresAt string `json:"expires_at,omitempty"`
}

// RemoveToken represents a GitHub remove token
type RemoveToken struct {
	Token     string `json:"token,omitempty"`
	ExpiresAt string `json:"expires_at,omitempty"`
}

// RunnerLabel represents a GitHub runner label
type RunnerLabel struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name"`
	Type string `json:"type,omitempty"`
}

// RunnerLabelList represents the response for listing runner labels
type RunnerLabelList struct {
	TotalCount int           `json:"total_count"`
	Labels     []RunnerLabel `json:"labels"`
}

// AddLabelsRequest represents the request to add labels to a runner
type AddLabelsRequest struct {
	Labels []string `json:"labels"`
}

// SetLabelsRequest represents the request to set labels for a runner
type SetLabelsRequest struct {
	Labels []string `json:"labels"`
}

// CreateRunnerGroupRequest represents the request to create a runner group
type CreateRunnerGroupRequest struct {
	Name                     string   `json:"name"`
	Visibility               string   `json:"visibility,omitempty"`
	SelectedRepositoryIDs    []int    `json:"selected_repository_ids,omitempty"`
	Runners                  []int    `json:"runners,omitempty"`
	AllowsPublicRepositories bool     `json:"allows_public_repositories,omitempty"`
	RestrictedToWorkflows    bool     `json:"restricted_to_workflows,omitempty"`
	SelectedWorkflows        []string `json:"selected_workflows,omitempty"`
	NetworkConfigurationID   string   `json:"network_configuration_id,omitempty"`
}

// UpdateRunnerGroupRequest represents the request to update a runner group
type UpdateRunnerGroupRequest struct {
	Name                     string   `json:"name,omitempty"`
	Visibility               string   `json:"visibility,omitempty"`
	AllowsPublicRepositories *bool    `json:"allows_public_repositories,omitempty"`
	RestrictedToWorkflows    *bool    `json:"restricted_to_workflows,omitempty"`
	SelectedWorkflows        []string `json:"selected_workflows,omitempty"`
	NetworkConfigurationID   *string  `json:"network_configuration_id,omitempty"`
}

// SetRepositoriesForRunnerGroupRequest represents the request to set repositories for a runner group
type SetRepositoriesForRunnerGroupRequest struct {
	SelectedRepositoryIDs []int `json:"selected_repository_ids"`
}

// SetRunnersForRunnerGroupRequest represents the request to set runners for a runner group
type SetRunnersForRunnerGroupRequest struct {
	Runners []int `json:"runners"`
}

// SelfHostedRunnerList represents the response for listing self-hosted runners
type SelfHostedRunnerList struct {
	TotalCount int                `json:"total_count"`
	Runners    []SelfHostedRunner `json:"runners"`
}

// HostedRunner represents a GitHub-hosted runner
type HostedRunner struct {
	ID            int    `json:"id,omitempty"`
	Name          string `json:"name"`
	RunnerGroupID int    `json:"runner_group_id,omitempty"`
	Platform      string `json:"platform,omitempty"`
	Image         struct {
		ID   string `json:"id,omitempty"`
		Size int    `json:"size,omitempty"`
	} `json:"image,omitempty"`
	MachineSizeDetails struct {
		ID        string `json:"id,omitempty"`
		CPUCores  int    `json:"cpu_cores,omitempty"`
		MemoryGB  int    `json:"memory_gb,omitempty"`
		StorageGB int    `json:"storage_gb,omitempty"`
	} `json:"machine_size_details,omitempty"`
	Status          string `json:"status,omitempty"`
	MaximumRunners  int    `json:"maximum_runners,omitempty"`
	PublicIPEnabled bool   `json:"public_ip_enabled,omitempty"`
	PublicIPs       []struct {
		Enabled bool   `json:"enabled,omitempty"`
		Prefix  string `json:"prefix,omitempty"`
		Length  int    `json:"length,omitempty"`
	} `json:"public_ips,omitempty"`
	LastActiveOn string `json:"last_active_on,omitempty"`
}

// HostedRunnerList represents the response for listing hosted runners
type HostedRunnerList struct {
	TotalCount int            `json:"total_count"`
	Runners    []HostedRunner `json:"runners"`
}

// Repository represents a GitHub repository
type Repository struct {
	ID          int    `json:"id,omitempty"`
	NodeID      string `json:"node_id,omitempty"`
	Name        string `json:"name,omitempty"`
	FullName    string `json:"full_name,omitempty"`
	Private     bool   `json:"private,omitempty"`
	HTMLURL     string `json:"html_url,omitempty"`
	Description string `json:"description,omitempty"`
	Fork        bool   `json:"fork,omitempty"`
	URL         string `json:"url,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
	PushedAt    string `json:"pushed_at,omitempty"`
	Visibility  string `json:"visibility,omitempty"`
}

// RepositoryList represents the response for listing repositories
type RepositoryList struct {
	TotalCount   int          `json:"total_count"`
	Repositories []Repository `json:"repositories"`
}
