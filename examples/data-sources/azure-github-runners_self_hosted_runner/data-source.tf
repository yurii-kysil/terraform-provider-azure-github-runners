# Retrieve a self-hosted runner by name
data "azure-github-runners_self_hosted_runner" "example" {
  name = "my-runner"
}
