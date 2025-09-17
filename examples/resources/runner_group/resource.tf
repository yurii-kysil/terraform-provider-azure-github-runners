# Create a runner group for production workloads
resource "azure-github-runners_runner_group" "production" {
  name       = "production-runners"
  visibility = "selected"

  selected_repository_ids = [
    123456789,
    987654321
  ]

  allows_public_repositories = false
  restricted_to_workflows    = true
  selected_workflows = [
    "my-org/my-repo/.github/workflows/deploy.yaml@refs/heads/main"
  ]

  network_configuration_id = azure-github-runners_network_configuration.main.id
}
