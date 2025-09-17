# Get self-hosted runner by name
data "github_runners_self_hosted_runner" "runner" {
  name            = "runner-01"
  runner_group_id = 123
}
