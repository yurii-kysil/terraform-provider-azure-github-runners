# Get network configuration by name
data "github_runners_network_configuration" "current" {
  name = "production-network-config"
}
