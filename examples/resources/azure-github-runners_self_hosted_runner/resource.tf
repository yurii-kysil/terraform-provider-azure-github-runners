data "azure-github-runners_network_configuration" "main" {
  name = "production-network-config"
}

# Create a runner group for production workloads
resource "azure-github-runners_runner_group" "production" {
  name                     = "production-runners"
  visibility               = "all"
  network_configuration_id = azure-github-runners_network_configuration.main.id
}

# Create a self-hosted runner with readonly and custom labels
resource "azure-github-runners_self_hosted_runner" "main" {
  name            = "runner-01"
  runner_group_id = azure-github-runners_runner_group.production.id

  readonly_labels = [
    "self-hosted",
    "X64",
    "Linux"
  ]

  labels = [
    "custom-label",
    "production"
  ]

  work_folder = "_work"
}
