# Create a self-hosted runner with custom labels
resource "github_runners_self_hosted_runner" "main" {
  name            = "runner-01"
  runner_group_id = github_runners_runner_group.production.id

  labels = [
    "self-hosted",
    "X64",
    "Linux",
    "custom-label"
  ]

  work_folder = "_work"
}
