# Terraform Provider for GitHub Runners

A Terraform provider for managing GitHub organization network configurations and self-hosted runners.

## Features

- **Network Configuration Management**: Configure allowed IP ranges for GitHub Actions
- **Runner Group Management**: Create and manage self-hosted runner groups
- **Self-hosted Runner Management**: Manage self-hosted runners and their group assignments

## Installation

### Using Terraform Registry

```hcl
terraform {
  required_providers {
    github-runners = {
      source  = "yurii-kysil/github-runners"
      version = "~> 1.0"
    }
  }
}
```

**Note:** This provider will be available on the Terraform Registry after the first release.

### Local Development

1. Clone this repository
2. Build the provider:

   ```bash
   go build -o terraform-provider-azure-github-runners
   ```

3. Install the provider in your Terraform plugins directory

## Configuration

```hcl
provider "github-runners" {
  organization = var.organization
  
  # Option 1: Use Personal Access Token
  token = var.github_token
  
  # Option 2: Use GitHub App Authentication (recommended for production)
  app_auth {
    id              = var.github_app_id
    installation_id = var.github_app_installation_id
    pem_file        = var.github_app_pem_file
  }
  
  # Optional
  base_url = "https://api.github.com"
  insecure = false
}
```

### Authentication

The provider supports the following authentication methods:

#### Personal Access Token

Use a GitHub Personal Access Token for authentication:

```hcl
provider "github-runners" {
  token        = var.github_token
  organization = var.organization
}
```

#### GitHub App Authentication (Recommended)

Use GitHub App authentication for better security and higher rate limits:

```hcl
provider "github-runners" {
  organization = var.organization
  
  app_auth {
    id              = var.github_app_id
    installation_id = var.github_app_installation_id
    pem_file        = var.github_app_pem_file
  }
}
```

**GitHub App Setup:**

1. Create a GitHub App in your organization
2. Install the app on your organization
3. Download the private key (PEM file)
4. Note the App ID and Installation ID from the app settings

**Benefits of GitHub App Authentication:**

- Higher API rate limits (15,000 requests/hour vs 5,000 for PAT)
- Better security (no long-lived tokens)
- Fine-grained permissions
- Automatic token rotation

**Environment Variables:**

- `GITHUB_TOKEN`: Personal Access Token (alternative to `token` parameter)
- `GITHUB_BASE_URL`: GitHub API base URL (alternative to `base_url` parameter)

## Resources

### github_runners_network_configuration

Manages GitHub organization network configurations for Actions.

```hcl
resource "github_runners_network_configuration" "main" {
  name               = "production-network-config"
  compute_service    = "actions"
  network_settings_ids = [
    "23456789ABDCEF1",
    "3456789ABDCEF12"
  ]
}
```

### github_runners_runner_group

Manages GitHub self-hosted runner groups.

```hcl
resource "github_runners_runner_group" "main" {
  name       = "production-runners"
  visibility = "selected"
  selected_repository_ids = [
    123456789,
    987654321
  ]
  allows_public_repositories = false
  restricted_to_workflows   = true
  selected_workflows = [
    "octo-org/octo-repo/.github/workflows/deploy.yaml@refs/heads/main"
  ]
  network_configuration_id = github_runners_network_configuration.main.id
}
```

### github_runners_self_hosted_runner

Manages GitHub self-hosted runners.

```hcl
resource "github_runners_self_hosted_runner" "main" {
  name           = "runner-01"
  runner_group_id = github_runners_runner_group.main.id
  labels = [
    "self-hosted",
    "X64",
    "Linux",
    "custom-label"
  ]
  work_folder = "_work"
}
```

## Data Sources

### github_runners_network_configuration

Retrieves the current network configuration.

```hcl
data "github_runners_network_configuration" "current" {
  name = "production-network-config"
}
```

### github_runners_runner_group

Retrieves a runner group by name.

```hcl
data "github_runners_runner_group" "default" {
  name = "Default"
}
```

### github_runners_self_hosted_runner

Retrieves a self-hosted runner by name.

```hcl
data "github_runners_self_hosted_runner" "runner" {
  name           = "runner-01"
  runner_group_id = 123
}
```

### github_runners_runner_applications

Retrieves available runner applications for download.

```hcl
data "github_runners_runner_applications" "apps" {}
```

### github_runners_registration_token

Retrieves a registration token for the organization.

```hcl
data "github_runners_registration_token" "token" {}
```

### github_runners_remove_token

Retrieves a remove token for the organization.

```hcl
data "github_runners_remove_token" "remove_token" {}
```

## API Coverage

This provider covers the following GitHub REST API endpoints:

### Network Configurations

- `GET /orgs/{org}/settings/network-configurations`
- `POST /orgs/{org}/settings/network-configurations`
- `GET /orgs/{org}/settings/network-configurations/{network_configuration_id}`
- `PATCH /orgs/{org}/settings/network-configurations/{network_configuration_id}`
- `DELETE /orgs/{org}/settings/network-configurations/{network_configuration_id}`
- `GET /orgs/{org}/settings/network-settings/{network_settings_id}`

### Self-hosted Runner Groups

- `GET /orgs/{org}/actions/runner-groups`
- `POST /orgs/{org}/actions/runner-groups`
- `GET /orgs/{org}/actions/runner-groups/{runner_group_id}`
- `PATCH /orgs/{org}/actions/runner-groups/{runner_group_id}`
- `DELETE /orgs/{org}/actions/runner-groups/{runner_group_id}`
- `GET /orgs/{org}/actions/runner-groups/{runner_group_id}/repositories`
- `PUT /orgs/{org}/actions/runner-groups/{runner_group_id}/repositories`
- `PUT /orgs/{org}/actions/runner-groups/{runner_group_id}/repositories/{repository_id}`
- `DELETE /orgs/{org}/actions/runner-groups/{runner_group_id}/repositories/{repository_id}`
- `GET /orgs/{org}/actions/runner-groups/{runner_group_id}/runners`
- `PUT /orgs/{org}/actions/runner-groups/{runner_group_id}/runners`
- `PUT /orgs/{org}/actions/runner-groups/{runner_group_id}/runners/{runner_id}`
- `DELETE /orgs/{org}/actions/runner-groups/{runner_group_id}/runners/{runner_id}`
- `GET /orgs/{org}/actions/runner-groups/{runner_group_id}/hosted-runners`

### Self-hosted Runners

- `GET /orgs/{org}/actions/runners`
- `GET /orgs/{org}/actions/runners/{runner_id}`
- `DELETE /orgs/{org}/actions/runners/{runner_id}`
- `GET /orgs/{org}/actions/runners/downloads`
- `POST /orgs/{org}/actions/runners/generate-jitconfig`
- `POST /orgs/{org}/actions/runners/registration-token`
- `POST /orgs/{org}/actions/runners/remove-token`
- `GET /orgs/{org}/actions/runners/{runner_id}/labels`
- `POST /orgs/{org}/actions/runners/{runner_id}/labels`
- `PUT /orgs/{org}/actions/runners/{runner_id}/labels`
- `DELETE /orgs/{org}/actions/runners/{runner_id}/labels`
- `DELETE /orgs/{org}/actions/runners/{runner_id}/labels/{name}`

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 0.13
- [Go](https://golang.org/doc/install) >= 1.19 (for building from source)

## Building

```bash
go build -o terraform-provider-azure-github-runners
```

## Testing

```bash
go test ./...
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
