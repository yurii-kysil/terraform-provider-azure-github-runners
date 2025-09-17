# Create a network configuration for GitHub Actions
resource "azure-github-runners_network_configuration" "main" {
  name            = "production-network-config"
  compute_service = "actions"
  network_settings_ids = [
    "23456789ABDCEF1",
    "3456789ABDCEF12"
  ]
}
