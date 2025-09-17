terraform {
  required_providers {
    github-runners = {
      source  = "yurii-kysil/azure-github-runners"
      version = "~> 1.0"
    }
  }
}

provider "github-runners" {
  organization = var.organization

  # Option 1: Use Personal Access Token
  # token = var.github_token

  # Option 2: Use GitHub App Authentication
  app_auth {
    id              = var.github_app_id
    installation_id = var.github_app_installation_id
    pem_file        = var.github_app_pem_file
  }
}

# Network Configuration
resource "github_runners_network_configuration" "main" {
  name            = "production-network-config"
  compute_service = "actions"
  network_settings_ids = [
    "23456789ABDCEF1",
    "3456789ABDCEF12"
  ]
}

# Runner Group
resource "github_runners_runner_group" "main" {
  name       = "production-runners"
  visibility = "selected"
  selected_repository_ids = [
    123456789,
    987654321
  ]
  allows_public_repositories = false
  restricted_to_workflows    = true
  selected_workflows = [
    "octo-org/octo-repo/.github/workflows/deploy.yaml@refs/heads/main"
  ]
  network_configuration_id = github_runners_network_configuration.main.id
}

# Self-hosted Runner
resource "github_runners_self_hosted_runner" "main" {
  name            = "runner-01"
  runner_group_id = github_runners_runner_group.main.id
  labels = [
    "self-hosted",
    "X64",
    "Linux",
    "custom-label"
  ]
  work_folder = "_work"
}

# Data sources
data "github_runners_network_configuration" "current" {
  name = "production-network-config"
}

data "github_runners_runner_group" "default" {
  name = "Default"
}

data "github_runners_self_hosted_runner" "runner" {
  name            = "runner-01"
  runner_group_id = github_runners_runner_group.main.id
}

# Runner Applications
data "github_runners_runner_applications" "apps" {}

# Registration Token
data "github_runners_registration_token" "token" {}

# Remove Token
data "github_runners_remove_token" "remove_token" {}
